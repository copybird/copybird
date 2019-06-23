package operator

import (
	"fmt"
	"time"

	backupv1 "github.com/copybird/copybird/operator/pkg/apis/backup/v1"
	backupclientset "github.com/copybird/copybird/operator/pkg/client/clientset/versioned"
	backupinformer "github.com/copybird/copybird/operator/pkg/client/informers/externalversions/backup/v1"
	listers "github.com/copybird/copybird/operator/pkg/client/listers/backup/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/batch/v1"
	v1_beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	batch "k8s.io/client-go/informers/batch"
	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/batch/v1"
	listersv1_beta1 "k8s.io/client-go/listers/batch/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "backup-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Backup is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Backup fails
	// to sync due to a Job of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Job already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Backup"
	// MessageResourceSynced is the message used for an Event fired when a Backup
	// is synced successfully
	MessageResourceSynced = "Backup synced successfully"
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	logger          *log.Entry
	kubeclientset   kubernetes.Interface
	backupclientset backupclientset.Interface
	jobsLister      listersv1.JobLister
	jobsSynced      cache.InformerSynced
	cronjobsLister  listersv1_beta1.CronJobLister
	cronjobsSynced  cache.InformerSynced
	backupsLister   listers.BackupLister
	backupsSynced   cache.InformerSynced
	workqueue       workqueue.RateLimitingInterface
}

//NewController creates controller instance
func NewController(kubeclientset kubernetes.Interface,
	backupclientset *backupclientset.Clientset,
	batchInterface batch.Interface,
	backupInformer backupinformer.BackupInformer) *Controller {

	controller := &Controller{
		logger:          log.NewEntry(log.New()),
		kubeclientset:   kubeclientset,
		backupclientset: backupclientset,
		jobsLister:      batchInterface.V1().Jobs().Lister(),
		jobsSynced:      batchInterface.V1().Jobs().Informer().HasSynced,
		cronjobsLister:  batchInterface.V1beta1().CronJobs().Lister(),
		cronjobsSynced:  batchInterface.V1beta1().CronJobs().Informer().HasSynced,
		backupsLister:   backupInformer.Lister(),
		backupsSynced:   backupInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Backups"),
	}

	// Set up an event handler for when Backup resources change
	backupInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueBackup,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueBackup(new)
		},
	})

	batchInterface.V1().Jobs().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newJob := new.(*v1.Job)
			oldJob := old.(*v1.Job)
			if newJob.ResourceVersion == oldJob.ResourceVersion {
				// Periodic resync will send update events for all known Jobs.
				// Two different versions of the same Job will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	batchInterface.V1beta1().CronJobs().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newCronJob := new.(*v1_beta1.CronJob)
			oldCronJob := old.(*v1_beta1.CronJob)
			if newCronJob.ResourceVersion == oldCronJob.ResourceVersion {
				// Periodic resync will send update events for all known CronJobs.
				// Two different versions of the same Job will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting Backup controller")

	// Wait for the caches to be synced before starting workers
	log.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.jobsSynced, c.cronjobsSynced, c.backupsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Info("Starting workers")
	// Launch two workers to process Backup resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Info("Started workers")
	<-stopCh
	log.Info("Shutting down workers")

	return nil
}

// runWorker executes the loop to process new items added to the queue
func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
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
		log.Infof("Successfully synced '%s'", key)
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
	backup, err := c.backupsLister.Backups(namespace).Get(name)
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
		err := c.proceedToCronJob(backup)
		return err
	}

	err = c.proceedToJob(backup)
	return err
}

// enqueueBackup takes a Backup resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Backup.
func (c *Controller) enqueueBackup(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Backup resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Backup resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		log.Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	log.Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Backup, we should not do anything more
		// with it.
		if ownerRef.Kind != "Backup" {
			return
		}

		backup, err := c.backupsLister.Backups(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			log.Infof("ignoring orphaned object '%s' of backup '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueBackup(backup)
		return
	}
}

func (c *Controller) proceedToJob(backup *backupv1.Backup) error {
	// Get the job with the name specified in Backup.spec
	job, err := c.jobsLister.Jobs(backup.Namespace).Get(backup.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		job, err = c.kubeclientset.BatchV1().Jobs(backup.Namespace).Create(newJob(backup))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Job is not controlled by this Backup resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(job, backup) {
		msg := fmt.Sprintf(MessageResourceExists, job.Name)
		return fmt.Errorf(msg)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) proceedToCronJob(backup *backupv1.Backup) error {
	// Get the job with the name specified in Backup.spec
	job, err := c.cronjobsLister.CronJobs(backup.Namespace).Get(backup.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		job, err = c.kubeclientset.BatchV1beta1().CronJobs(backup.Namespace).Create(newCronJob(backup))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Job is not controlled by this Backup resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(job, backup) {
		msg := fmt.Sprintf(MessageResourceExists, job.Name)
		return fmt.Errorf(msg)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	return nil
}

func newJob(backup *backupv1.Backup) *v1.Job {
	return &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backup.Spec.Name,
			Namespace: backup.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(backup, backupv1.SchemeGroupVersion.WithKind("Backup")),
			},
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

func newCronJob(backup *backupv1.Backup) *v1_beta1.CronJob {
	return &v1_beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backup.Spec.Name,
			Namespace: backup.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(backup, backupv1.SchemeGroupVersion.WithKind("Backup")),
			},
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
