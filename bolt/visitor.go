package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/model"
)

// PutVisitor creates a new visitor or overwrites the existing one.
func (d *Database) PutVisitor(v *model.Visitor) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		visitors, err := tx.CreateBucketIfNotExists([]byte(visitorsBucket))
		if err != nil {
			return errors.Error{
				Type:    errors.NotFound,
				Message: "bolt: failed to find visitors bucket",
			}
		}

		token := []byte(v.Token)
		visitor, err := json.Marshal(v)
		if err != nil {
			return errors.Error{
				Type:    errors.FailedMarshal,
				Message: fmt.Sprintf("bolt: failed to marshal visitor: %v", err),
			}
		}

		if err := visitors.Put(token, visitor); err != nil {
			return errors.Error{
				Type:    errors.NotPut,
				Message: "bolt: failed to update visitor",
			}
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
			return errors.Error{
				Type:    errors.NotFound,
				Message: "bolt: failed to find visitors bucket",
			}
		}

		visitor := visitors.Get([]byte(token))
		if visitor == nil {
			return errors.Error{
				Type:    errors.NotFound,
				Message: "bolt: failed to find visitor",
			}
		}

		if err := json.Unmarshal(visitor, &v); err != nil {
			return errors.Error{
				Type:    errors.FailedMarshal,
				Message: fmt.Sprintf("bolt: failed to unmarshal visitor: %v", err),
			}
		}

		return nil
	})
	return v, err
}
