package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
)

// PutVisitor creates a new visitor or overwrites the existing one.
func (d *Database) PutVisitor(v *model.Visitor) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return d.put(tx, "visitor", visitorsBucket, []byte(v.Token), v)
	})
}

// GetVisitor returns the visitor with the given token.
func (d *Database) GetVisitor(token string) (*model.Visitor, error) {
	var v *model.Visitor
	err := d.db.View(func(tx *bolt.Tx) error {
		return d.get(tx, "visitor", visitorsBucket, []byte(token), &v)
	})
	return v, err
}
