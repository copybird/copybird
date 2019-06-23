package mongo

import (
	"context"
	"fmt"
	input2 "github.com/copybird/copybird/modules/backup/input"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MODULE_NAME is name of module
const MODULE_NAME = "mongodb"

type (
	// BackupInputMongodb represents struct for dumping mongo
	BackupInputMongodb struct {
		input2.Input
		config *MongoConfig
		reader io.Reader
		writer io.Writer
		conn   *mongo.Client
	}
)

// GetName returns name
func (m *BackupInputMongodb) GetName() string {
	return MODULE_NAME
}

// GetConfig returns config
func (m *BackupInputMongodb) GetConfig() interface{} {
	return &MongoConfig{}
}

// InitPipe initializes pipes
func (m *BackupInputMongodb) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *BackupInputMongodb) InitModule(cfg interface{}) error {
	m.config = cfg.(*MongoConfig)
	cO := options.Client().ApplyURI(DSN)
	conn, err := mongo.Connect(context.TODO(), cO)
	if err != nil {
		return err
	}
	m.conn = conn
	return nil
}

// Run runs module
func (m *BackupInputMongodb) Run() error {
	return nil
}

// Close closes
func (m *BackupInputMongodb) Close() error {
	m.conn.Disconnect(context.TODO())
	return nil
}
func (m *BackupInputMongodb) getDatabases() ([]string, error) {
	return m.conn.ListDatabaseNames(context.TODO(), bsonx.Doc{})

}
func (m *BackupInputMongodb) getCollections(dbName string) ([]string, error) {
	var colls []string
	collections, err := m.conn.Database(dbName).ListCollections(nil, bson.M{})
	if err != nil {
		return colls, err
	}
	for collections.Next(nil) {
		colNameRaw := collections.Current.Lookup("name")
		colName, ok := colNameRaw.StringValueOK()
		if !ok {
			return colls, fmt.Errorf("invalid collection name: %v", colNameRaw)
		}
		colls = append(colls, colName)
	}
	return colls, nil

}
