package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
	"github.com/sirupsen/logrus"
)

// CreateUser creates a new user.
func (d *Database) CreateUser(u *model.User) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists([]byte(usersBucket))
		if err != nil {
			return err
		}

		username := []byte(u.Username)
		user, err := json.Marshal(u)
		if err != nil {
			return err
		}

		if users.Get(username) != nil {
			return fmt.Errorf("bolt: user already exists")
		}
		if err := users.Put(username, user); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		d.Logger.WithFields(logrus.Fields{
			"username": u.Username,
		}).Info("bolt: created user")
	} else {
		d.Logger.WithFields(logrus.Fields{
			"username": u.Username,
			"err":      err,
		}).Warn("bolt: failed to create user")
	}
	return err
}

// FetchUser returns the user with the given username.
func (d *Database) FetchUser(username string) (*model.User, error) {
	var u *model.User
	err := d.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(usersBucket))
		if users == nil {
			return fmt.Errorf("bolt: users not found")
		}

		user := users.Get([]byte(username))
		if user == nil {
			return fmt.Errorf("bolt: user not found")
		}

		if err := json.Unmarshal(user, &u); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		d.Logger.WithFields(logrus.Fields{
			"username": username,
		}).Info("bolt: fetched user")
	} else {
		d.Logger.WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Warn("bolt: failed to fetch user")
	}
	return u, err
}
