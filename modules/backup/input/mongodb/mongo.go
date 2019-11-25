package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/copybird/copybird/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Module Constants
const (
	groupName  = "backup"
	typeName   = "input"
	moduleName = "mongodb"
)

var (
	timeout = time.Second * 2
	header  = `{"database":"%s","collection":"%s","time":"%d"}`
)

type (
	// BackupInputMongodb represents struct for dumping mongo
	BackupInputMongodb struct {
		core.Module
		config         *MongoConfig
		reader         io.Reader
		writer         io.Writer
		conn           *mongo.Client
		rootCtx        context.Context
		cancelFn       func()
		dbCount        int
		collCount      int
		docCount       int64
		bytesOut       uint64
		recCount       uint64
		startTimestamp time.Time
	}
)

// SetDefaultTimeout sets default timeout used for operations
func SetDefaultTimeout(t time.Duration) {
	timeout = t
}

// InitPipe initializes pipes
func (m *BackupInputMongodb) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w

	return nil
}

// InitModule initializes module
func (m *BackupInputMongodb) InitModule(cfg interface{}) error {
	conf, ok := cfg.(*MongoConfig)
	if !ok {
		return fmt.Errorf("config type mismatch, expected: %T actual: %T", m.config, cfg)
	}
	m.config = conf

	rCtx, cancel := context.WithCancel(context.Background())
	m.cancelFn = cancel
	m.rootCtx = rCtx

	ctx, cancel := context.WithTimeout(m.rootCtx, timeout)
	defer cancel()

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(m.config.DSN))
	if err != nil {
		return err
	}

	m.conn = conn
	return nil
}

// Run runs module
func (m *BackupInputMongodb) Run(ctx context.Context) error {
	m.startTimestamp = time.Now()

	dbs, err := m.getDatabases(m.rootCtx)
	if err != nil {
		return fmt.Errorf("unable to fetch databases: %v", err)
	}

	for _, db := range dbs {
		colls, err := m.getCollections(m.rootCtx, db)
		if err != nil {
			return fmt.Errorf("unable to fetch collectioins: %v", err)
		}

		for _, coll := range colls {
			if err := m.exportCollection(db, coll); err != nil {
				return fmt.Errorf("unable to export collection: %v", err)
			}
			m.collCount++
		}
		m.dbCount++
	}
	_ = fmt.Sprintf("exported [databases: %d collections: %d documents: %d bytes out: %d duration: %v]", m.dbCount, m.collCount, m.bytesOut, m.docCount, time.Since(m.startTimestamp).Seconds())
	m.cancelFn()
	return nil
}

// Close disconnects from the server
func (m *BackupInputMongodb) Close() error {
	m.cancelFn()
	return m.conn.Disconnect(m.rootCtx)
}

func (m *BackupInputMongodb) getDatabases(ctx context.Context) ([]string, error) {
	return m.conn.ListDatabaseNames(ctx, bsonx.Doc{})
}

func (m *BackupInputMongodb) getCollections(ctx context.Context, dbName string) ([]string, error) {
	var colls []string

	collections, err := m.conn.Database(dbName).ListCollections(m.rootCtx, bson.M{})
	if err != nil {
		return colls, err
	}

	for collections.Next(m.rootCtx) {
		colNameRaw := collections.Current.Lookup("name")
		colName, ok := colNameRaw.StringValueOK()
		if !ok {
			return colls, fmt.Errorf("invalid collection name: %v", colNameRaw)
		}
		colls = append(colls, colName)
	}
	return colls, nil

}

func (m *BackupInputMongodb) exportCollection(dbName, collName string) error {
	curr, err := m.conn.Database(dbName).Collection(collName).Find(m.rootCtx, bson.D{})
	if err != nil {
		return fmt.Errorf("unable to fetch documents :%v", err)
	}
	defer curr.Close(m.rootCtx)

	if _, err = m.writer.Write([]byte(fmt.Sprintf(header, dbName, collName, time.Now().UnixNano()) + "\n")); err != nil {
		return fmt.Errorf("unable to write header: %v", err)
	}

	for curr.Next(m.rootCtx) {
		var res bson.M
		if err := curr.Decode(&res); err != nil {
			return fmt.Errorf("unable to decode document: %v", err)
		}

		data, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("unable to marshal document: %v", err)
		}
		data = append(data, "\n"...)
		n, err := m.writer.Write(data)
		if err != nil {
			return fmt.Errorf("unable to write document data: %v", err)
		}
		m.bytesOut += uint64(n)
		m.docCount++
		if n != len(data) {
			return fmt.Errorf("expected write: %d actual: %d", len(data), n)
		}
	}

	return nil
}

// GetGroup returns group
func (m *BackupInputMongodb) GetGroup() core.ModuleGroup { return groupName }

// GetType returns type
func (m *BackupInputMongodb) GetType() core.ModuleType { return typeName }

// GetName returns name
func (m *BackupInputMongodb) GetName() string { return moduleName }

// GetConfig returns config
func (m *BackupInputMongodb) GetConfig() interface{} { return &MongoConfig{} }
