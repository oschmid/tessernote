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

type Tag struct {
	Name      string
	UserKeys  []*datastore.Key
	NoteKeys  []*datastore.Key
	ParentKey *datastore.Key
	ChildKeys []*datastore.Key
}

func (t Tag) Users(c appengine.Context) []User {
	var users []User
	datastore.GetMulti(c, t.UserKeys, users)
	return users
}

func (t Tag) Notes(c appengine.Context) []Note {
	var notes []Note
	datastore.GetMulti(c, t.NoteKeys, notes)
	return notes
}

func (t Tag) Parent(c appengine.Context) Tag {
	var parent Tag
	datastore.Get(c, t.ParentKey, parent)
	return parent
}

func (t Tag) Children(c appengine.Context) []Tag {
	var children []Tag
	datastore.GetMulti(c, t.ChildKeys, children)
	return children
}
