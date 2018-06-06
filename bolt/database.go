package bolt

import "github.com/boltdb/bolt"

// Database represents the database connection.
type Database struct {
	Path string
	db   *bolt.DB
}

// Open opens up a connection to the Bolt database.
func (d *Database) Open() error {
	db, err := bolt.Open(d.Path, 0644, nil)
	if err != nil {
		return err
	}
	d.db = db

	return nil
}
