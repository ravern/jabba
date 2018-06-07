package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
)

// CreateLink creates a new link.
func (d *Database) CreateLink(l *model.Link) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists([]byte(linksBucket))
		if err != nil {
			return err
		}

		slug := []byte(l.Slug)
		link, err := json.Marshal(l)
		if err != nil {
			return err
		}

		if links.Get(link) != nil {
			return fmt.Errorf("bolt: link already exists")
		}
		if err := links.Put(slug, link); err != nil {
			return err
		}

		return nil
	})
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
				continue
			}

			if err := json.Unmarshal(link, &l); err != nil {
				continue
			}

			ll = append(ll, l)
		}

		return nil
	})
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

		if err := json.Unmarshal(link, &l); err != nil {
			return err
		}

		return nil
	})
	return l, err
}
