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

func (tag Tag) Notebooks() ([]Notebook, error) {
	var notebooks []Notebook
	err := datastore.GetMulti(tag.context, tag.NotebookKeys, notebooks)
	if err != nil {
		return notebooks, err
	}
	for _, notebook := range notebooks {
		notebook.context = tag.context
	}
	return notebooks, nil
}

func (tag Tag) Notes() ([]Note, error) {
	var notes []Note
	err := datastore.GetMulti(tag.context, tag.NoteKeys, notes)
	if err != nil {
		return notes, err
	}
	for i, note := range notes {
		note.ID = tag.NoteKeys[i].Encode()
		note.context = tag.context
	}
	return notes, nil
}

func (tag Tag) Children() ([]Tag, error) {
	var children []Tag
	err := datastore.GetMulti(tag.context, tag.ChildKeys, children)
	if err != nil {
		return children, err
	}
	for _, child := range children {
		child.context = tag.context
	}
	return children, err
}
