package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/ravern/jabba/model"
)

// UpdateAuths updates the auths of the given link.
//
// This will create and delete auths from the database where necessary, as well
// as update existing ones. If an error occurs performing any of the updates,
// it will be returned.
func (d *Database) UpdateAuths(aa []*model.Auth, l *model.Link) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		for _, id := range l.AuthIDs {
			var found bool
			for _, a := range aa {
				if a.ID == id {
					found = true
					break
				}
			}
			if !found {
				if err := d.delete(tx, "auth", authsBucket, []byte(id)); err != nil {
					return err
				}
			}
		}

		var ids []string
		for _, a := range aa {
			if err := d.put(tx, "auth", authsBucket, []byte(a.ID), a); err != nil {
				return err
			}
			ids = append(ids, a.ID)
		}
		l.AuthIDs = ids

		if err := d.update(tx, "link", linksBucket, []byte(l.Slug), l); err != nil {
			return err
		}

		return nil
	})
}

// GetAuths returns the auths with the given IDs.
func (d *Database) GetAuths(ids []string) ([]*model.Auth, error) {
	var aa []*model.Auth
	err := d.db.View(func(tx *bolt.Tx) error {
		for _, id := range ids {
			var a *model.Auth
			if err := d.get(tx, "auth", authsBucket, []byte(id), &a); err != nil {
				continue
			}
			aa = append(aa, a)
		}

		return nil
	})
	return aa, err
}
