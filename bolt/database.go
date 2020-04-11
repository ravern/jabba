package bolt

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ravern/jabba/errors"
)

// Bucket names.
const (
	linksBucket = "Links"
	usersBucket = "Users"
	authsBucket = "Auths"
)

// Database represents the database connection.
type Database struct {
	Path               string
	VisitCountInterval time.Duration // interval to write visit count caches to the database

	// Underlying bolt database instance
	db *bolt.DB

	// Usage counters for each slug
	visitCounts   map[string]int
	visitCountsMu sync.Mutex
}

// Open opens up a connection to the database.
func (d *Database) Open() error {
	db, err := bolt.Open(d.Path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}
	d.db = db

	d.visitCounts = make(map[string]int)
	go func() {
		for {
			time.Sleep(d.VisitCountInterval)
			d.updateLinkCounts()
		}
	}()

	return nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// bucket returns the given bucket.
//
// Panics if the given bucket does not exist.
func (d *Database) bucket(tx *bolt.Tx, bucket string, name string) *bolt.Bucket {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		panic(fmt.Sprintf("bolt: %s bucket not found", name))
	}
	return b
}

// create puts the given value with the given key into the given transaction, but
// only if it doesn't exist.
func (d *Database) create(tx *bolt.Tx, name string, bucket string, key []byte, value interface{}) error {
	err := d.get(tx, name, bucket, key, &struct{}{})
	if err == nil {
		return errors.Error{
			Type:    errors.AlreadyExists,
			Message: fmt.Sprintf("bolt: %s already exists", name),
		}
	} else if err.(errors.Error).Type != errors.NotFound {
		return err
	}
	return d.put(tx, name, bucket, key, value)
}

// update puts the given value with the given key into the given transaction, but
// only if it exists.
func (d *Database) update(tx *bolt.Tx, name string, bucket string, key []byte, value interface{}) error {
	if err := d.get(tx, name, bucket, key, &struct{}{}); err != nil {
		return err
	}
	return d.put(tx, name, bucket, key, value)
}

// put puts the given value with the given key into the given transaction.
func (d *Database) put(tx *bolt.Tx, name string, bucket string, key []byte, value interface{}) error {
	b := d.bucket(tx, bucket, name)

	v, err := json.Marshal(value)
	if err != nil {
		return errors.Error{
			Type:    errors.FailedMarshal,
			Message: fmt.Sprintf("bolt: failed to marshal %s: %v", name, err),
		}
	}

	if err := b.Put(key, v); err != nil {
		return errors.Error{
			Type:    errors.NotPut,
			Message: fmt.Sprintf("bolt: failed to put %s", name),
		}
	}

	return nil
}

// get gets the value with the given key from the given transaction.
func (d *Database) get(tx *bolt.Tx, name string, bucket string, key []byte, value interface{}) error {
	b := d.bucket(tx, bucket, name)

	v := b.Get(key)
	if v == nil {
		return errors.Error{
			Type:    errors.NotFound,
			Message: fmt.Sprintf("bolt: failed to find %s", name),
		}
	}

	if err := json.Unmarshal(v, value); err != nil {
		return errors.Error{
			Type:    errors.FailedMarshal,
			Message: fmt.Sprintf("bolt: failed to unmarshal %s: %v", name, err),
		}
	}

	return nil
}

// delete deletes the value with the given key from the given transaction.
func (d *Database) delete(tx *bolt.Tx, name string, bucket string, key []byte) error {
	b := d.bucket(tx, bucket, name)

	if err := b.Delete(key); err != nil {
		return errors.Error{
			Type:    errors.NotDeleted,
			Message: fmt.Sprintf("bolt: failed to delete %s", name),
		}
	}

	return nil
}
