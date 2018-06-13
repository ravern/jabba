package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/model"
)

// CreateLink creates a new link and adds that link to the user.
func (d *Database) CreateLink(l *model.Link, u *model.User) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		if err := d.create(tx, "link", linksBucket, []byte(l.Slug), l); err != nil {
			return err
		}
		u.LinkSlugs = append(u.LinkSlugs, l.Slug)
		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
	if err != nil {
		if i, ok := u.FindLinkSlug(l.Slug); ok {
			u.LinkSlugs = append(u.LinkSlugs[:i], u.LinkSlugs[i+1:]...)
		}
	}
	return err
}

// DeleteLink deletes the given link and removes that link from the user.
func (d *Database) DeleteLink(l *model.Link, u *model.User) error {
	var (
		i  int
		ok bool
	)
	err := d.db.Update(func(tx *bolt.Tx) error {
		if err := d.delete(tx, "links", linksBucket, []byte(l.Slug)); err != nil {
			return err
		}

		i, ok = u.FindLinkSlug(l.Slug)
		if !ok {
			return errors.Error{
				Type:    errors.Unauthorized,
				Message: "bolt: failed to find link in user",
			}
		}
		u.LinkSlugs = append(u.LinkSlugs[:i], u.LinkSlugs[i+1:]...)

		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
	if err != nil && ok {
		u.LinkSlugs = append(u.LinkSlugs[:i], append([]string{l.Slug}, u.LinkSlugs[i:]...)...)
	}
	return err
}

// IncrementLinkCount increments the usage count of the given link.
func (d *Database) IncrementLinkCount(l *model.Link) {
	d.countsMu.Lock()
	defer d.countsMu.Unlock()

	_, ok := d.counts[l.Slug]
	if !ok {
		d.counts[l.Slug] = 0
	}

	d.counts[l.Slug]++
	l.Count++
}

// updateLinkCounts writes the usage counts cache into the database.
func (d *Database) updateLinkCounts() {
	d.countsMu.Lock()
	defer d.countsMu.Unlock()

	d.db.Update(func(tx *bolt.Tx) error {
		for slug, count := range d.counts {
			var l *model.Link
			if err := d.get(tx, "link", linksBucket, []byte(slug), &l); err != nil {
				continue
			}

			l.Count += count

			d.put(tx, "link", linksBucket, []byte(slug), l)
		}

		return nil
	})

	d.counts = make(map[string]int)
}

// UpdateLinkSlug updates the given link, including changes to the slug.
func (d *Database) UpdateLinkSlug(slug string, l *model.Link, u *model.User) error {
	var (
		i  int
		ok bool
	)
	err := d.db.Update(func(tx *bolt.Tx) error {
		if slug == l.Slug {
			return d.update(tx, "link", linksBucket, []byte(l.Slug), l)
		}

		if err := d.delete(tx, "link", linksBucket, []byte(slug)); err != nil {
			return err
		}
		if err := d.create(tx, "link", linksBucket, []byte(l.Slug), l); err != nil {
			return err
		}

		var ok bool
		i, ok := u.FindLinkSlug(slug)
		if !ok {
			return errors.Error{
				Type:    errors.Unauthorized,
				Message: "bolt: failed to find link in user",
			}
		}
		u.LinkSlugs[i] = l.Slug

		return d.update(tx, "user", usersBucket, []byte(u.Username), u)
	})
	if err != nil && ok {
		u.LinkSlugs[i] = slug
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
