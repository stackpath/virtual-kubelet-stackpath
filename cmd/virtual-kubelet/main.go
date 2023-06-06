package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client"
	"github.com/stackpath/vk-stackpath-provider/internal/auth"
	"github.com/stackpath/vk-stackpath-provider/internal/config"
	spprovider "github.com/stackpath/vk-stackpath-provider/internal/provider"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	logruslogger "github.com/virtual-kubelet/virtual-kubelet/log/logrus"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	v1 "k8s.io/api/core/v1"
)

type inputVars struct {
	nodeName         string
	providerConfig   string
	startupTimeout   time.Duration
	clientCACert     string
	kubeConfig       string
	disableTaint     bool
	logLevel         string
	podSyncWorkers   int
	fullResyncPeriod time.Duration
	taintKey         string
	taintEffect      string
	taintValue       string
	noVerifyClients  bool
	serverCertPath   string
	serverKeyPath    string
}

var inputs = inputVars{
	nodeName:         "stackpath-edge-provider",
	providerConfig:   "",
	startupTimeout:   0,
	clientCACert:     "",
	kubeConfig:       os.Getenv("KUBECONFIG"),
	disableTaint:     false,
	logLevel:         "info",
	podSyncWorkers:   50,
	fullResyncPeriod: 0,
	taintKey:         "virtual-kubelet.io/provider",
	taintEffect:      string(v1.TaintEffectNoSchedule),
	taintValue:       "stackpath",
	noVerifyClients:  false,
	serverCertPath:   os.Getenv("APISERVER_CERT_LOCATION"),
	serverKeyPath:    os.Getenv("APISERVER_KEY_LOCATION"),
}

var (
	buildVersion   = "N/A"
	k8sVersion     = "v1.25.0" // This should follow the version of k8s.io we are importing
	binaryFilename = filepath.Base(os.Args[0])
	description    = fmt.Sprintf("%s implements a node on a Kubernetes cluster using StackPath Workload API to run pods.", binaryFilename)
	listenPort     = int32(10250)
)

func init() {
	home, _ := homedir.Dir()
	if home != "" {
		inputs.kubeConfig = filepath.Join(home, ".kube", "config")
	}

	virtualKubeletCommand.Flags().StringVar(&inputs.nodeName, "nodename", inputs.nodeName, "kubernetes node name")
	virtualKubeletCommand.Flags().StringVar(&inputs.kubeConfig, "kube-config", inputs.kubeConfig, "kubeconfig file")
	virtualKubeletCommand.Flags().StringVar(&inputs.providerConfig, "provider-config", inputs.providerConfig, "provider configuration file")
	virtualKubeletCommand.Flags().DurationVar(&inputs.startupTimeout, "startup-timeout", inputs.startupTimeout, "how long to wait for the virtual-kubelet to start")
	virtualKubeletCommand.Flags().StringVar(&inputs.clientCACert, "client-verify-ca", inputs.clientCACert, "CA cert to use to verify client requests")
	virtualKubeletCommand.Flags().BoolVar(&inputs.noVerifyClients, "no-verify-clients", inputs.noVerifyClients, "do not require client certificate validation")
	virtualKubeletCommand.Flags().StringVar(&inputs.logLevel, "log-level", inputs.logLevel, "log level")
	virtualKubeletCommand.Flags().IntVar(&inputs.podSyncWorkers, "pod-sync-workers", inputs.podSyncWorkers, `set the number of pod synchronization workers`)
	virtualKubeletCommand.Flags().DurationVar(&inputs.fullResyncPeriod, "full-resync-period", inputs.fullResyncPeriod, "how often to perform a full resync of pods between kubernetes and the provider")
	virtualKubeletCommand.Flags().BoolVar(&inputs.disableTaint, "disable-taint", inputs.disableTaint, "disable the node taint")
	virtualKubeletCommand.Flags().StringVar(&inputs.taintKey, "taint-key", inputs.taintKey, "a string identifier used to mark a node in Kubernetes with a specific characteristic, influencing the scheduling of pods on that node")
	virtualKubeletCommand.Flags().StringVar(&inputs.taintEffect, "taint-effect", inputs.taintEffect, "a string representing the desired effect of a taint on a Kubernetes node, determining how pods should be scheduled or evicted in relation to the node")
	virtualKubeletCommand.Flags().StringVar(&inputs.taintValue, "taint-value", inputs.taintValue, "a string that provides additional context or details about the taintKey, helping differentiate between different taints with the same key on a Kubernetes node")
	virtualKubeletCommand.Flags().StringVar(&inputs.serverCertPath, "api-server-cert", inputs.serverCertPath, "the API server's public certificate")
	virtualKubeletCommand.Flags().StringVar(&inputs.serverKeyPath, "api-server-key", inputs.serverKeyPath, "the API server's private key in a Kubernetes cluster")
}

// virtualKubeletCommand is the main command that runs the virtual kubelet
var virtualKubeletCommand = &cobra.Command{
	Use:   binaryFilename,
	Short: description,
	Long:  description,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.StandardLogger()
		logLevel, err := logrus.ParseLevel(inputs.logLevel)

		if err != nil {
			logrus.WithError(err).Fatal("Error parsing log level")
		}
		logger.SetLevel(logLevel)

		ctx := log.WithLogger(cmd.Context(), logruslogger.FromLogrus(logrus.NewEntry(logger)))

		if err := runNode(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.G(ctx).Fatal(err)
			} else {
				log.G(ctx).Debug(err)
			}
		}
	},
}

// runNode creates and runs a virtual-kubelet node
func runNode(ctx context.Context) error {
	// Create API config and runtime
	apiConfig, err := config.NewConfig(ctx)
	if err != nil {
		log.G(ctx).Fatal(err)
	}

	// Add edge location to the node name
	inputs.nodeName = fmt.Sprintf("%s-%s", inputs.nodeName, strings.ToLower(apiConfig.CityCode))

	runtime, err := auth.NewRuntime(ctx, apiConfig.ClientID, apiConfig.ClientSecret, apiConfig.ApiHost, buildVersion)
	if err != nil {
		log.G(ctx).Fatal(err)
	}

	// Create StackPath client
	stackpathClient := workload_client.New(runtime, nil)

	// Create and run node
	node, err := nodeutil.NewNode(inputs.nodeName,
		func(cfg nodeutil.ProviderConfig) (nodeutil.Provider, node.NodeProvider, error) {
			if port := os.Getenv("KUBELET_PORT"); port != "" {
				p, err := strconv.Atoi(port)
				if err != nil {
					return nil, nil, err
				}
				listenPort = int32(p)
			}
			p, err := spprovider.NewStackpathProvider(ctx, stackpathClient, apiConfig, cfg, os.Getenv("VKUBELET_POD_IP"), listenPort)
			p.ConfigureNode(ctx, cfg.Node)
			return p, nil, err
		},
		withClient,
		withTaint,
		withVersion,
		withTLSConfig,
		withWebhookAuth,
		configureRoutes,
		func(cfg *nodeutil.NodeConfig) error {
			cfg.InformerResyncPeriod = inputs.fullResyncPeriod
			cfg.NumWorkers = inputs.podSyncWorkers
			cfg.HTTPListenAddr = fmt.Sprintf(":%d", listenPort)
			return nil
		},
	)
	if err != nil {
		return err
	}

	go func() error {
		err = node.Run(ctx)
		if err != nil {
			return fmt.Errorf("error running the node: %w", err)
		}
		return nil
	}()

	if err := node.WaitReady(ctx, inputs.startupTimeout); err != nil {
		return fmt.Errorf("error waiting for node to be ready: %w", err)
	}

	<-node.Done()
	return node.Err()
}

// withClient sets up the Kubernetes client for the node
func withClient(cfg *nodeutil.NodeConfig) error {
	client, err := nodeutil.ClientsetFromEnv(inputs.kubeConfig)
	if err != nil {
		return err
	}
	return nodeutil.WithClient(client)(cfg)
}

// withVersion sets the Kubelet Version reported by the node
func withVersion(cfg *nodeutil.NodeConfig) error {
	cfg.NodeSpec.Status.NodeInfo.KubeletVersion = strings.Join([]string{k8sVersion, "vk-stackpath", buildVersion}, "-")
	return nil
}

func withTLSConfig(cfg *nodeutil.NodeConfig) error {
	if inputs.serverCertPath == "" || inputs.serverKeyPath == "" {
		return nil
	}
	return nodeutil.WithTLSConfig(nodeutil.WithKeyPairFromPath(inputs.serverCertPath, inputs.serverKeyPath), withCA)(cfg)
}

func withWebhookAuth(cfg *nodeutil.NodeConfig) error {
	cfg.Handler = api.InstrumentHandler(nodeutil.WithAuth(nodeutil.NoAuth(), cfg.Handler))
	return nil
}

// withCA sets up the client CA for the node
func withCA(cfg *tls.Config) error {
	if inputs.clientCACert == "" {
		return nil
	}
	if err := nodeutil.WithCAFromPath(inputs.clientCACert)(cfg); err != nil {
		return fmt.Errorf("error getting CA from path: %w", err)
	}
	if inputs.noVerifyClients {
		cfg.ClientAuth = tls.NoClientCert
	}
	return nil
}

// configureRoutes sets up the HTTP routes for the virtual kubelet
func configureRoutes(cfg *nodeutil.NodeConfig) error {
	mux := http.NewServeMux()
	cfg.Handler = mux
	return nodeutil.AttachProviderRoutes(mux)(cfg)
}

// withTaint sets up the taint for the node
func withTaint(cfg *nodeutil.NodeConfig) error {
	if inputs.disableTaint {
		return nil
	}

	taint := v1.Taint{
		Key:   inputs.taintKey,
		Value: inputs.taintValue,
	}
	switch inputs.taintEffect {
	case string(v1.TaintEffectNoSchedule):
		taint.Effect = v1.TaintEffectNoSchedule
	case string(v1.TaintEffectNoExecute):
		taint.Effect = v1.TaintEffectNoExecute
	case string(v1.TaintEffectPreferNoSchedule):
		taint.Effect = v1.TaintEffectPreferNoSchedule
	default:
		return errdefs.InvalidInputf("taint effect %q is not supported", inputs.taintEffect)
	}
	cfg.NodeSpec.Spec.Taints = append(cfg.NodeSpec.Spec.Taints, taint)
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := virtualKubeletCommand.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Fatal("error running command")
		}
	}
}
