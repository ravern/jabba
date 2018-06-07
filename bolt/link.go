package bolt

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
	"github.com/sirupsen/logrus"
)

// CreateLink creates a new link.
func (d *Database) CreateLink(l *model.Link) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists([]byte(linksBucket))
		if err != nil {
			return err
		}

		slug := []byte(l.Slug)

		var b bytes.Buffer
		if err := gob.NewEncoder(&b).Encode(l); err != nil {
			return err
		}
		link := b.Bytes()

		if links.Get(link) != nil {
			return fmt.Errorf("bolt: link already exists")
		}
		if err := links.Put(slug, link); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		d.Logger.WithFields(logrus.Fields{
			"slug": l.Slug,
		}).Info("bolt: created link")
	} else {
		d.Logger.WithFields(logrus.Fields{
			"slug": l.Slug,
			"err":  err,
		}).Warn("bolt: failed to create link")
	}
	return err
}

// FetchLinks returns the links created by the given user.
func (d *Database) FetchLinks(u *model.User) ([]*model.Link, error) {
	var ll []*model.Link
	err := d.db.View(func(tx *bolt.Tx) error {
		links := tx.Bucket([]byte(linksBucket))
		if links == nil {
			return fmt.Errorf("bolt: links not found")
		}

		for _, slug := range u.LinkSlugs {
			var l *model.Link

			link := links.Get([]byte(slug))
			if link == nil {
				d.Logger.WithFields(logrus.Fields{
					"slug": slug,
				}).Warn("bolt: couldn't find link")
				continue
			}

			b := bytes.NewBuffer(link)
			if err := gob.NewDecoder(b).Decode(&l); err != nil {
				d.Logger.WithFields(logrus.Fields{
					"err": err,
				}).Warn("bolt: failed to decode link")
				continue
			}

			ll = append(ll, l)
		}

		return nil
	})
	if err == nil {
		d.Logger.Infof("bolt: fetched %d links", len(ll))
	} else {
		d.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Info("bolt: failed to fetch links")
	}
	return ll, err
}

// FetchLink returns the link with the given slug.
func (d *Database) FetchLink(slug string) (*model.Link, error) {
	var l *model.Link
	err := d.db.View(func(tx *bolt.Tx) error {
		links := tx.Bucket([]byte(linksBucket))
		if links == nil {
			return fmt.Errorf("bolt: links not found")
		}

		link := links.Get([]byte(slug))
		if link == nil {
			return fmt.Errorf("bolt: link not found")
		}

		b := bytes.NewBuffer(link)
		if err := gob.NewDecoder(b).Decode(&l); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		d.Logger.WithFields(logrus.Fields{
			"slug": slug,
		}).Info("bolt: fetched link")
	} else {
		d.Logger.WithFields(logrus.Fields{
			"slug": slug,
			"err":  err,
		}).Warn("bolt: failed to fetch link")
	}
	return l, err
}
