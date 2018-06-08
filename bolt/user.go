package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/model"
)

// CreateUser creates a new user.
func (d *Database) CreateUser(u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists([]byte(usersBucket))
		if err != nil {
			return errors.Error{
				Type:    errors.FailedPut,
				Message: "bolt: failed to create users bucket",
			}
		}

		username := []byte(u.Username)
		user, err := json.Marshal(u)
		if err != nil {
			return errors.Error{
				Type:    errors.FailedMarshal,
				Message: fmt.Sprintf("bolt: failed to marshal user: %v", err),
			}
		}

		if users.Get(username) != nil {
			return errors.Error{
				Type:    errors.AlreadyExists,
				Message: "bolt: user already exists",
			}
		}
		if err := users.Put(username, user); err != nil {
			return errors.Error{
				Type:    errors.FailedPut,
				Message: "bolt: failed to create user",
			}
		}

		return nil
	})
}

// FetchUser returns the user with the given username.
func (d *Database) FetchUser(username string) (*model.User, error) {
	var u *model.User
	err := d.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(usersBucket))
		if users == nil {
			return errors.Error{
				Type:    errors.NotFound,
				Message: "bolt: failed to find users bucket",
			}
		}

		user := users.Get([]byte(username))
		if user == nil {
			return errors.Error{
				Type:    errors.NotFound,
				Message: "bolt: failed to find user",
			}
		}

		if err := json.Unmarshal(user, &u); err != nil {
			return errors.Error{
				Type:    errors.FailedMarshal,
				Message: fmt.Sprintf("bolt: failed to unmarshal user: %v", err),
			}
		}

		return nil
	})
	return u, err
}
