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
	"appengine/user"
)
type User struct {
	Id       string
	TagKeys  []*datastore.Key
	NoteKeys []*datastore.Key
	context  appengine.Context `datastore:",noindex"`
}

func (u User) Tags() []Tag {
	var tags []Tag
	datastore.GetMulti(u.context, u.TagKeys, tags)
	for _, t := range tags {
		t.context = u.context
	}
	return tags
}

func (u User) Notes() []Note {
	var notes []Note
	datastore.GetMulti(u.context, u.NoteKeys, notes)
	for _, n := range notes {
		n.context = u.context
	}
	return notes
}

func GetUser(c appengine.Context) *User {
	g := new(User)
	u := user.Current(c)
	k := datastore.NewKey(c, "User", u.FederatedIdentity, 0, nil)
	if err := datastore.Get(c, k, u); err != nil {
		// TODO create new user
	}
	g.context = c
	return g
}

func PutUser(u User) {
	k := datastore.NewKey(u.context, "User", u.Id, 0, nil)
	if _, err := datastore.Put(u.context, k, &u); err != nil {
		// TODO handle failed put
	}
}
