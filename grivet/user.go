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
}

func (u User) Tags(c appengine.Context) []Tag {
	var tags []Tag
	datastore.GetMulti(c, u.TagKeys, tags)
	return tags
}

func (u User) Notes(c appengine.Context) []Note {
	var notes []Note
	datastore.GetMulti(c, u.NoteKeys, notes)
	return notes
}

func GetUser(c appengine.Context, u user.User) *User {
	k := datastore.NewKey(c, "User", u.FederatedIdentity, 0, nil)
	g := new(User)
	if err := datastore.Get(c, k, u); err != nil {
		// TODO create new user
	}
	return g
}

func PutUser(c appengine.Context, u User) {
	k := datastore.NewKey(c, "User", u.Id, 0, nil)
	if _, err := datastore.Put(c, k, &u); err != nil {
		// TODO handle failed put
	}
}
