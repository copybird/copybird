package operator

import (
	"os"
	"time"

	"github.com/knative/pkg/signals"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	//need this line to parse gcp credentials
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	backupclientset "github.com/copybird/copybird/operator/pkg/client/clientset/versioned"
	backupinformers "github.com/copybird/copybird/operator/pkg/client/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
)

//Run starts the operator
func Run() {

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	backupClient, err := backupclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(client, time.Second*30)
	backupInformerFactory := backupinformers.NewSharedInformerFactory(backupClient, time.Second*30)

	controller := NewController(client, backupClient,
		kubeInformerFactory.Batch(),
		backupInformerFactory.Backups().V1().Backups())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	backupInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		log.Fatalf("Error running controller: %s", err.Error())
	}
}
