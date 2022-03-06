package repo

import (
	"github.com/quantosnetwork/Quantosp2p/config"
	"github.com/quantosnetwork/Quantosp2p/filestore"
	"github.com/quantosnetwork/Quantosp2p/keystore"
	"github.com/quantosnetwork/Quantosp2p/quantos"
	"context"
	"errors"
	"github.com/ipfs/go-datastore/query"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"time"
)

var (
	ErrApiNotRunning = errors.New("api not running")
)

type Repo interface {
	Config() (*config.Config, error)
	BackupConfig(prefix string) (string, error)
	SetConfig(*config.Config) error
	SetConfigKey(key string, value interface{})
	GetConfigKey(key string) (interface{}, error)
	Datastore() Datastore
	GetStoreSize(context.Context) (uint64, error)
	Keystore() keystore.Keystore
	FileManager() *filestore.FileManager
	ChainManager(networkId string) *quantos.BlockchainManager
	SetAPIAddr(addr ma.Multiaddr) error
	SwarmKey() ([]byte, error)
	io.Closer
}

type Datastore interface {
	Read
	Write
	Sync(ctx context.Context, prefix Key) error
	io.Closer
}

type Write interface {
	Put(ctx context.Context, key Key, value []byte) error
	Delete(ctx context.Context, key Key) error
}

type Read interface {
	Get(ctx context.Context, key Key) (value []byte, err error)
	Has(ctx context.Context, key Key) (exists bool, err error)
	GetSize(ctx context.Context, key Key) (size int, err error)
	Query(ctx context.Context, q query.Query) (query.Results, error)
}

type Batching interface {
	Datastore
	Batch(ctx context.Context) (Batch, error)
}

type CheckedDatastore interface {
	Datastore

	Check(ctx context.Context) error
}

// ScrubbedDatastore is an interface that should be implemented by datastores
// which want to provide a mechanism to check data integrity and/or
// error correction.
type ScrubbedDatastore interface {
	Datastore

	Scrub(ctx context.Context) error
}

// GCDatastore is an interface that should be implemented by datastores which
// don't free disk space by just removing data from them.
type GCDatastore interface {
	Datastore

	CollectGarbage(ctx context.Context) error
}

// PersistentDatastore is an interface that should be implemented by datastores
// which can report disk usage.
type PersistentDatastore interface {
	Datastore

	// DiskUsage returns the space used by a datastore, in bytes.
	DiskUsage(ctx context.Context) (uint64, error)
}

// DiskUsage checks if a Datastore is a
// PersistentDatastore and returns its DiskUsage(),
// otherwise returns 0.
func DiskUsage(ctx context.Context, d Datastore) (uint64, error) {
	persDs, ok := d.(PersistentDatastore)
	if !ok {
		return 0, nil
	}
	return persDs.DiskUsage(ctx)
}

// TTLDatastore is an interface that should be implemented by datastores that
// support expiring entries.
type TTLDatastore interface {
	Datastore
	TTL
}

// TTL encapulates the methods that deal with entries with time-to-live.
type TTL interface {
	PutWithTTL(ctx context.Context, key Key, value []byte, ttl time.Duration) error
	SetTTL(ctx context.Context, key Key, ttl time.Duration) error
	GetExpiration(ctx context.Context, key Key) (time.Time, error)
}

// Txn extends the Datastore type. Txns allow users to batch queries and
// mutations to the Datastore into atomic groups, or transactions. Actions
// performed on a transaction will not take hold until a successful call to
// Commit has been made. Likewise, transactions can be aborted by calling
// Discard before a successful Commit has been made.
type Txn interface {
	Read
	Write

	// Commit finalizes a transaction, attempting to commit it to the Datastore.
	// May return an error if the transaction has gone stale. The presence of an
	// error is an indication that the data was not committed to the Datastore.
	Commit(ctx context.Context) error
	// Discard throws away changes recorded in a transaction without committing
	// them to the underlying Datastore. Any calls made to Discard after Commit
	// has been successfully called will have no effect on the transaction and
	// state of the Datastore, making it safe to defer.
	Discard(ctx context.Context)
}

// TxnDatastore is an interface that should be implemented by datastores that
// support transactions.
type TxnDatastore interface {
	Datastore

	NewTransaction(ctx context.Context, readOnly bool) (Txn, error)
}

// Errors

type dsError struct {
	error
	isNotFound bool
}

func (e *dsError) NotFound() bool {
	return e.isNotFound
}

// ErrNotFound is returned by Get and GetSize when a datastore does not map the
// given key to a value.
var ErrNotFound error = &dsError{error: errors.New("datastore: key not found"), isNotFound: true}

// GetBackedHas provides a default Datastore.Has implementation.
// It exists so Datastore.Has implementations can use it, like so:
//
// func (*d SomeDatastore) Has(key Key) (exists bool, err error) {
//   return GetBackedHas(d, key)
// }
func GetBackedHas(ctx context.Context, ds Read, key Key) (bool, error) {
	_, err := ds.Get(ctx, key)
	switch err {
	case nil:
		return true, nil
	case ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

// GetBackedSize provides a default Datastore.GetSize implementation.
// It exists so Datastore.GetSize implementations can use it, like so:
//
// func (*d SomeDatastore) GetSize(key Key) (size int, err error) {
//   return GetBackedSize(d, key)
// }
func GetBackedSize(ctx context.Context, ds Read, key Key) (int, error) {
	value, err := ds.Get(ctx, key)
	if err == nil {
		return len(value), nil
	}
	return -1, err
}

type Batch interface {
	Write

	Commit(ctx context.Context) error
}