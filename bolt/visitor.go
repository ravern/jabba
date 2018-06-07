package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
)

// PutVisitor creates a new visitor or overwrites the existing one.
func (d *Database) PutVisitor(v *model.Visitor) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		visitors, err := tx.CreateBucketIfNotExists([]byte(visitorsBucket))
		if err != nil {
			return err
		}

		token := []byte(v.Token)
		visitor, err := json.Marshal(v)
		if err != nil {
			return err
		}

		if err := visitors.Put(token, visitor); err != nil {
			return err
		}

		return nil
	})
}

// FetchVisitor returns the visitor with the given token.
func (d *Database) FetchVisitor(token string) (*model.Visitor, error) {
	var v *model.Visitor
	err := d.db.View(func(tx *bolt.Tx) error {
		visitors := tx.Bucket([]byte(visitorsBucket))
		if visitors == nil {
			return fmt.Errorf("bolt: users not found")
		}

		visitor := visitors.Get([]byte(token))
		if visitor == nil {
			return fmt.Errorf("bolt: user not found")
		}

		if err := json.Unmarshal(visitor, &v); err != nil {
			return err
		}

		return nil
	})
	return v, err
}
