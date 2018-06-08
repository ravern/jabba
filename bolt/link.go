package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/model"
)

// CreateUserLink creates a new link and adds that link to the user.
func (d *Database) CreateUserLink(l *model.Link, u *model.User) error {
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

		users := tx.Bucket([]byte(usersBucket))
		if users == nil {
			return fmt.Errorf("bolt: users not found")
		}

		u.LinkSlugs = append(u.LinkSlugs, l.Slug)
		user, err := json.Marshal(u)
		if err != nil {
			return err
		}

		if err := users.Put([]byte(u.Username), user); err != nil {
			return err
		}

		return nil
	})
}

// CreateVisitorLink creates a new link and adds that link to the visitor.
func (d *Database) CreateVisitorLink(l *model.Link, v *model.Visitor) error {
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

		visitors := tx.Bucket([]byte(visitorsBucket))
		if visitors == nil {
			return fmt.Errorf("bolt: visitors not found")
		}

		v.LinkSlugs = append(v.LinkSlugs, l.Slug)
		visitor, err := json.Marshal(v)
		if err != nil {
			return err
		}

		if err := visitors.Put([]byte(v.Token), visitor); err != nil {
			return err
		}

		return nil
	})
}

// IncrementLinkCount increments the visit count of the given link.
func (d *Database) IncrementLinkCount(l *model.Link) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists([]byte(linksBucket))
		if err != nil {
			return err
		}

		l.Count++

		slug := []byte(l.Slug)
		link, err := json.Marshal(l)
		if err != nil {
			return err
		}

		if err := links.Put(slug, link); err != nil {
			return err
		}

		return nil
	})
}

// FetchLinks returns the links with the given slugs.
func (d *Database) FetchLinks(slugs []string) ([]*model.Link, error) {
	var ll []*model.Link
	err := d.db.View(func(tx *bolt.Tx) error {
		links := tx.Bucket([]byte(linksBucket))
		if links == nil {
			return fmt.Errorf("bolt: links not found")
		}

		for _, slug := range slugs {
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
