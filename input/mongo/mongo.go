package mongo

import (
	"io"

	"github.com/copybird/copybird/input"
	"github.com/globalsign/mgo"
)

const MODULE_NAME = "mongo"

type (
	// MongoDumper represents struct for dumping mongo
	MongoDumper struct {
		input.Input
		config *MongoConfig
		reader io.Reader
		writer io.Writer
		conn   *mgo.Session
	}
)

// GetName returns name
func (d *MongoDumper) GetName() string {
	return MODULE_NAME
}

// GetConfig returns config
func (d *MongoDumper) GetConfig() interface{} {
	return &MongoConfig{}
}

// InitPipe initializes pipes
func (d *MongoDumper) InitPipe(w io.Writer, r io.Reader) error {
	d.reader = r
	d.writer = w
	return nil
}

// InitModule initializes module
func (d *MongoDumper) InitModule(cfg interface{}) error {
	d.config = cfg.(*MongoConfig)
	conn, err := mgo.Dial(d.config.Host)
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

// Run runs module
func (d *MongoDumper) Run() error {
	return nil
}

// Close closes
func (d *MongoDumper) Close() error {
	d.conn.Close()
	return nil
}
