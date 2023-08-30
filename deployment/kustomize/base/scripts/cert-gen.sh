#!/bin/bash

#install deps
apt update -y
apt install jq curl net-tools vim openssl original-awk gettext-base -y

#EKS issue https://github.com/aws/containers-roadmap/issues/1604

export NODE_NAME=$(hostname)
export INTERNAL_IP=$(ifconfig eth0 | grep inet | head -n 1 | awk '{print $2}')

openssl genrsa -out /etc/virtual-kubelet/key.pem 2048
openssl req -new -key /etc/virtual-kubelet/key.pem -out  /etc/virtual-kubelet/vk-sp.csr -config <(envsubst < /tmp/cert/csr.conf)

CSR=$(cat /etc/virtual-kubelet/vk-sp.csr | base64 | tr -d "\n")

 curl -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
          --data '{
            "kind": "CertificateSigningRequest",
            "apiVersion": "certificates.k8s.io/v1",
            "metadata": {
                "name": "vk-sp"
            },
            "spec": {
              "request": "'${CSR}'",
              "signerName": "kubernetes.io/kube-apiserver-client",
              "expirationSeconds": 86400,
            "usages": ["digital signature", "key encipherment", "server auth"]
            }
          }' \
        "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/apis/certificates.k8s.io/v1/certificatesigningrequests?fieldManager=kubectl-client-side-apply&fieldValidation=Strict" -k | jq

sleep 10

 curl -X PUT \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
          --data '{
          "kind": "CertificateSigningRequest",
          "apiVersion": "certificates.k8s.io/v1",
          "metadata": {
            "name": "vk-sp"
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
        }' \
        "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/apis/certificates.k8s.io/v1/certificatesigningrequests/vk-sp/approval" -k | jq

sleep 10

curl -X GET  \
    -H "Accept: application/json" \
    -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
    "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/apis/certificates.k8s.io/v1/certificatesigningrequests/vk-sp" -k | jq .status.certificate -r | base64 -d > /etc/virtual-kubelet/cert.pem
