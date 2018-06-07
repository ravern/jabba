package bolt

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// Bucket names.
const (
	usersBucket = "Users"
	linksBucket = "Links"
)

// Database represents the database connection.
type Database struct {
	Path   string
	Logger logrus.FieldLogger
	db     *bolt.DB
}

// Open opens up a connection to the database.
func (d *Database) Open() error {
	db, err := bolt.Open(d.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	d.db = db

	return nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}
