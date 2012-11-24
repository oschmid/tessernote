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
	Name         string // unique per user
	NotebookKeys []*datastore.Key
	NoteKeys     []*datastore.Key
	ChildKeys    []*datastore.Key
	context      appengine.Context
}

func (t Tag) Notebooks() []Notebook {
	var users []Notebook
	datastore.GetMulti(t.context, t.NotebookKeys, users)
	for _, u := range users {
		u.context = t.context
	}
	return users
}

func (t Tag) Notes() []Note {
	var notes []Note
	datastore.GetMulti(t.context, t.NoteKeys, notes)
	for i, n := range notes {
		n.ID = t.NoteKeys[i].Encode()
		n.context = t.context
	}
	return notes
}

func (t Tag) Children() []Tag {
	var children []Tag
	datastore.GetMulti(t.context, t.ChildKeys, children)
	for _, tag := range children {
		tag.context = t.context
	}
	return children
}
