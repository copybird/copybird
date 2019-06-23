package operator

import (
	"fmt"
	"time"

	backupv1 "github.com/copybird/copybird/operator/pkg/apis/backup/v1"
	listers "github.com/copybird/copybird/operator/pkg/client/listers/backup/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/batch/v1"
	v1_beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"k8s.io/klog"
)

const controllerAgentName = "backup-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Job of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Job already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Backup"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Backup synced successfully"
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	logger         *log.Entry
	clientset      kubernetes.Interface
	workqueue      workqueue.RateLimitingInterface
	backupLister   listers.BackupLister
	backupInformer cache.SharedIndexInformer
	handler        Handler
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// handle a panic with logging and exiting
	defer utilruntime.HandleCrash()
	// ignore new items in the queue but when all goroutines
	// have completed existing items then shutdown
	defer c.workqueue.ShutDown()

	c.logger.Info("Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.backupInformer.Run(stopCh)

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.logger.Info("Controller.Run: cache sync complete")

	// run the runWorker method every second with a stop channel
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.backupInformer.HasSynced()
}

// runWorker executes the loop to process new items added to the queue
func (c *Controller) runWorker() {
	log.Info("Controller.runWorker: starting")

	// invoke processNextItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextItem() {
		log.Info("Controller.runWorker: processing next item")
	}

	log.Info("Controller.runWorker: completed")
}

// processNextItem retrieves each queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Backup resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Backup resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Backup resource with this namespace/name
	backup, err := c.backupLister.Backups(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("backup '%s' in work queue no longer exists", key))
			return nil
		}

		return nil
	}

	if backup.Spec.Name == "" {
		utilruntime.HandleError(fmt.Errorf("%s: job name must be specified", key))
		return nil
	}

	if backup.Spec.Cron != "" {
		err := c.ProceedToCronJob(backup)
		return err
	}

	err = c.ProceedToJob(backup)
	return err
}

func (c *Controller) ProceedToJob(backup *backupv1.Backup) error {
	log.Info("Proceed to Job")
	// Get the Job with the name specified in Backup.spec
	job, err := c.clientset.BatchV1().Jobs(backup.Namespace).Get(backup.Spec.Name, metav1.GetOptions{})

	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		_, err = c.clientset.BatchV1().Jobs(backup.Namespace).Create(NewJob(backup))
		if err != nil {
			return err
		}
		return nil
	}

	job, err = c.clientset.BatchV1().Jobs(backup.Namespace).Update(job)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) ProceedToCronJob(backup *backupv1.Backup) error {
	log.Info("Proceed to CronJob")

	// Get the Job with the name specified in Backup.spec
	job, err := c.clientset.BatchV1beta1().CronJobs(backup.Namespace).Get(backup.Spec.Name, metav1.GetOptions{})

	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		_, err = c.clientset.BatchV1beta1().CronJobs(backup.Namespace).Create(NewCronJob(backup))
		if err != nil {
			return err
		}
		return nil
	}

	job, err = c.clientset.BatchV1beta1().CronJobs(backup.Namespace).Update(job)
	if err != nil {
		return err
	}
	return nil
}

func NewJob(backup *backupv1.Backup) *v1.Job {
	return &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: backup.Spec.Name,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: backup.Name,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "hello-world",
							Image: "hello-world:latest",
						},
					},
				},
			},
		},
	}
}

func NewCronJob(backup *backupv1.Backup) *v1_beta1.CronJob {
	return &v1_beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: backup.Spec.Name,
		},
		Spec: v1_beta1.CronJobSpec{
			Schedule: backup.Spec.Cron,
			JobTemplate: v1_beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: backup.Name,
				},
				Spec: v1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: backup.Name,
						},
						Spec: corev1.PodSpec{
							RestartPolicy: "OnFailure",
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "hello-world",
									Image: "hello-world:latest",
								},
							},
						},
					},
				},
			},
		},
	}
}
