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
)

// TODO use datastore.Key.Encode() instead of UUID
type Note struct {
	Title    string
	Body     string `datastore:",noindex"`
	UserKeys []*datastore.Key
	TagKeys  []*datastore.Key  // sorted by Tag.Name
	context  appengine.Context `datastore:",noindex"`
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
