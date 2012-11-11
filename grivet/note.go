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

type Note struct {
	Body     string `datastore:",noindex"`
	UserKeys []*datastore.Key
	TagKeys  []*datastore.Key
}

func (n Note) Users(c appengine.Context) []User {
	var users []User
	datastore.GetMulti(c, n.UserKeys, users)
	return users
}

func (n Note) Tags(c appengine.Context) []Tag {
	var tags []Tag
	datastore.GetMulti(c, n.TagKeys, tags)
	return tags
}
