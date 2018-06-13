package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
)

// CreateUser creates a new user.
func (d *Database) CreateUser(u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return d.create(tx, "user", usersBucket, []byte(u.Username), u)
	})
}

// UpdateUser updates the given user.
func (d *Database) UpdateUser(u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
}

// UpdateUserUsername updates the given user, including changes to the username.
func (d *Database) UpdateUserUsername(username string, u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if username == u.Username {
			return d.update(tx, "user", usersBucket, []byte(u.Username), u)
		}

		if err := d.delete(tx, "user", usersBucket, []byte(username)); err != nil {
			return err
		}
		return d.create(tx, "user", usersBucket, []byte(u.Username), u)
	})
}

// GetUser returns the user with the given username.
func (d *Database) GetUser(username string) (*model.User, error) {
	var u *model.User
	err := d.db.View(func(tx *bolt.Tx) error {
		return d.get(tx, "user", usersBucket, []byte(username), &u)
	})
	return u, err
}
