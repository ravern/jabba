package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/model"
)

// CreateUserLink creates a new link and adds that link to the user.
func (d *Database) CreateUserLink(l *model.Link, u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := d.create(tx, "link", linksBucket, []byte(l.Slug), l); err != nil {
			return err
		}
		u.LinkSlugs = append(u.LinkSlugs, l.Slug)
		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
}

// CreateVisitorLink creates a new link and adds that link to the visitor.
func (d *Database) CreateVisitorLink(l *model.Link, v *model.Visitor) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := d.create(tx, "link", linksBucket, []byte(l.Slug), l); err != nil {
			return err
		}
		v.LinkSlugs = append(v.LinkSlugs, l.Slug)
		return d.update(tx, "visitor", visitorsBucket, []byte(v.Token), v)
	})
}

// DeleteUserLink deletes the link with the given slug and removes that link
// from the user.
func (d *Database) DeleteUserLink(slug string, u *model.User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := d.delete(tx, "links", linksBucket, []byte(slug)); err != nil {
			return err
		}

		i, ok := u.FindLinkSlug(slug)
		if !ok {
			return errors.Error{
				Type:    errors.Unauthorized,
				Message: "bolt: failed to find link in user",
			}
		}
		u.LinkSlugs = append(u.LinkSlugs[:i], u.LinkSlugs[i+1:]...)

		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
}

// DeleteVisitorLink deletes the link with the given slug and removes that link
// from the visitor.
func (d *Database) DeleteVisitorLink(slug string, v *model.Visitor) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := d.delete(tx, "links", linksBucket, []byte(slug)); err != nil {
			return err
		}

		i, ok := v.FindLinkSlug(slug)
		if !ok {
			return errors.Error{
				Type:    errors.Unauthorized,
				Message: "bolt: failed to find link in visitor",
			}
		}
		v.LinkSlugs = append(v.LinkSlugs[:i], v.LinkSlugs[i+1:]...)

		return d.update(tx, "visitor", visitorsBucket, []byte(v.Token), v)
	})
}

// IncrementLinkCount increments the visit count of the given link.
func (d *Database) IncrementLinkCount(l *model.Link) error {
	count := l.Count
	err := d.db.Update(func(tx *bolt.Tx) error {
		l.Count++
		return d.update(tx, "link", linksBucket, []byte(l.Slug), l)
	})
	if err != nil {
		l.Count = count
	}
	return err
}

// GetLinks returns the links with the given slugs.
func (d *Database) GetLinks(slugs []string) ([]*model.Link, error) {
	var ll []*model.Link
	err := d.db.View(func(tx *bolt.Tx) error {
		for _, slug := range slugs {
			var l *model.Link
			if err := d.get(tx, "link", linksBucket, []byte(slug), &l); err != nil {
				continue
			}
			ll = append(ll, l)
		}

		return nil
	})
	return ll, err
}

// GetLink returns the link with the given slug.
func (d *Database) GetLink(slug string) (*model.Link, error) {
	var l *model.Link
	err := d.db.View(func(tx *bolt.Tx) error {
		return d.get(tx, "link", linksBucket, []byte(slug), &l)
	})
	return l, err
}
