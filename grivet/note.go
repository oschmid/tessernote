/*
This file is part of Grivet.

Grivet is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Grivet is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Grivet.  If not, see <http://www.gnu.org/licenses/>.
*/
package grivet

import (
	"appengine"
	"appengine/datastore"
	"time"
	"strings"
)

type Note struct {
	ID           *datastore.Key `datastore:"-"`
	Title        string
	Body         string
	Created      time.Time
	LastModified time.Time
	UserKeys     []*datastore.Key
	TagKeys      []*datastore.Key  // sorted by Tag.Name
	context      appengine.Context `datastore:"-"`
}

func (n Note) Users() []User {
	var users []User
	datastore.GetMulti(n.context, n.UserKeys, users)
	for _, u := range users {
		u.context = n.context
	}
	return users
}

func (n Note) Tags() []Tag {
	var tags []Tag
	datastore.GetMulti(n.context, n.TagKeys, tags)
	for _, t := range tags {
		t.context = n.context
	}
	return tags
}

func (n *Note) Update(text string) error {
	n.Title, n.Body = titleAndBodyFrom(text)
	n.LastModified = time.Now()
	_, err := datastore.Put(n.context, n.ID, n)
	return err
}

func titleAndBodyFrom(text string) (title, body string) {
	split := strings.SplitN(text, "\n", 2)
	return split[0], split[1]
}
