#!/bin/bash

set -e

# Function to make a Kubernetes API request
function kubectl_request {
    local method="$1"
    local path="$2"
    local data="$3"
    
    curl -s -X "$method" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        -d "$data" \
        "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT$path" -k | jq
}

# Install dependencies
apt update -y
apt install jq curl net-tools vim openssl original-awk gettext-base -y

export NODE_NAME=$(hostname)
export INTERNAL_IP=$(ifconfig eth0 | awk '/inet / {print $2}')

# Retrieve providerID
PROVIDER_ID=$(kubectl_request "GET" "/api/v1/nodes" "" | jq ".items | .[].spec.providerID")
SIGNER_NAME="kubernetes.io/kubelet-serving"
USAGES="[\"digital signature\", \"key encipherment\", \"server auth\"]"

# Check if providerID contains "eks"
# EKS issue https://github.com/aws/containers-roadmap/issues/1604
if [[ "$PROVIDER_ID" == *"aws"* ]]; then
    SIGNER_NAME="beta.eks.amazonaws.com/app-serving"
    USAGES="[\"server auth\"]"
fi

# Generate key and CSR
openssl genrsa -out /etc/virtual-kubelet/key.pem 2048
openssl req -new -key /etc/virtual-kubelet/key.pem -out /etc/virtual-kubelet/vk-sp.csr -config <(envsubst < /tmp/cert/csr.conf)

CSR=$(cat /etc/virtual-kubelet/vk-sp.csr | base64 | tr -d "\n")

CERT_NAME=vk-sp-$(date | md5sum  | awk '{print $1}')

# Create and approve CSR
body='{
    "kind": "CertificateSigningRequest",
    "apiVersion": "certificates.k8s.io/v1",
    "metadata": {
        "name": "'${CERT_NAME}'"
    },
    "spec": {
        "request": "'${CSR}'",
        "signerName": "'${SIGNER_NAME}'",
        "expirationSeconds": 315360000,
        "usages": '${USAGES}'
    }
}'
kubectl_request "POST" "/apis/certificates.k8s.io/v1/certificatesigningrequests?fieldManager=kubectl-client-side-apply&fieldValidation=Strict" "${body}"

sleep 10

kubectl_request "PUT" "/apis/certificates.k8s.io/v1/certificatesigningrequests/${CERT_NAME}/approval" '{
    "kind": "CertificateSigningRequest",
    "apiVersion": "certificates.k8s.io/v1",
    "metadata": {
        "name": "'${CERT_NAME}'"
    },
    "status": {
        "conditions": [
            {
                "type": "Approved",
                "status": "True",
                "reason": "KubectlApprove",
                "message": "This CSR was approved by kubectl certificate approve."
            }
        ]
    }
}'

sleep 10

# Get and save the certificate
kubectl_request "GET" "/apis/certificates.k8s.io/v1/certificatesigningrequests/${CERT_NAME}" "" | jq -r '.status.certificate' | base64 -d > /etc/virtual-kubelet/cert.pem

# Delete signing request
kubectl_request "DELETE" "/apis/certificates.k8s.io/v1/certificatesigningrequests/${CERT_NAME}" ""

# Check if the certificate is valid
if openssl x509 -noout -in /etc/virtual-kubelet/cert.pem; then 
    echo "Certificate successfully generated and signed"
else 
    echo "Error during certificate generation. Falling back to self-signed certificate."
    openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 -subj "/CN=sp-vk" -keyout /etc/virtual-kubelet/key.pem -out /etc/virtual-kubelet/cert.pem
fi
