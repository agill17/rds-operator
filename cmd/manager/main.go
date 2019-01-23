package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/agill17/rds-operator/pkg/apis"
	"github.com/agill17/rds-operator/pkg/controller"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/ready"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("operator-sdk Version: %v", sdkVersion.Version))
}
func getAwsCreds() (string, string) {
	return os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")
}

func ensureCredsAreDefined() {
	access, secret := getAwsCreds()
	if access == "" || secret == "" {
		logrus.Warnf("Operator cannot find aws creds, please create a secret and deploy the secret with aws credentials in the same namespace as the operator")
		logrus.Warnf("Secret name must be rds-access-creds and it must contain AWS_ACCESS and AWS_SECRET")
		logrus.Warnf("Assuming your node has an IAM role with enough access to create/update/delete rds instances")
	}
}

func main() {
	flag.Parse()

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(logf.ZapLogger(false))

	printVersion()

	// namespace, err := k8sutil.GetWatchNamespace()
	// if err != nil {
	// 	log.Error(err, "failed to get watch namespace")
	// 	os.Exit(1)
	// }

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Become the leader before proceeding
	leader.Become(context.TODO(), "rds-operator-lock")

	r := ready.NewFileReady()
	err = r.Set()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	defer r.Unset()

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: ""})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "manager exited non-zero")
		os.Exit(1)
	}
}
