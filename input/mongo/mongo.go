package mongo

import (
	"context"
	"io"

	"github.com/copybird/copybird/input"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MODULE_NAME is name of module
const MODULE_NAME = "mongo"

type (
	// MongoDumper represents struct for dumping mongo
	MongoDumper struct {
		input.Input
		config *MongoConfig
		reader io.Reader
		writer io.Writer
		conn   *mongo.Client
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
	cO := options.Client().ApplyURI(d.config.DSN)
	conn, err := mongo.Connect(context.TODO(), cO)
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
	d.conn.Disconnect(context.TODO())
	return nil
}
func (d *MongoDumper) getDatabases() ([]string, error) {
	return d.conn.ListDatabaseNames(context.TODO(), bsonx.Doc{})

}
