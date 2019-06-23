package operator

import (
	log "github.com/sirupsen/logrus"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// BackupHandler is a sample implementation of Handler
type BackupHandler struct{}

// Init handles any handler initialization
func (t *BackupHandler) Init() error {
	log.Info("BackupHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *BackupHandler) ObjectCreated(obj interface{}) {
	log.Info("BackupHandler.ObjectCreated")
}

// ObjectDeleted is called when an object is deleted
func (t *BackupHandler) ObjectDeleted(obj interface{}) {
	log.Info("BackupHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *BackupHandler) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("BackupHandler.ObjectUpdated")
}
