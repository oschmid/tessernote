/*
This file is part of Tessernote.

Tessernote is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Tessernote is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Tessernote.  If not, see <http://www.gnu.org/licenses/>.
*/

package note

import (
	"appengine"
	"appengine/datastore"
	"log"
)

type Tag struct {
	Name         string // unique per user
	NotebookKeys []*datastore.Key
	NoteKeys     []*datastore.Key
	ChildKeys    []*datastore.Key
}

func (tag Tag) Notebooks(c appengine.Context) ([]Notebook, error) {
	notebooks := make([]Notebook, len(tag.NotebookKeys))
	err := datastore.GetMulti(c, tag.NotebookKeys, notebooks)
	if err != nil {
		log.Println("getMulti:notebooks", err)
	}
	return notebooks, err
}

func (tag Tag) Notes(c appengine.Context) ([]Note, error) {
	notes := make([]Note, len(tag.NoteKeys))
	err := datastore.GetMulti(c, tag.NoteKeys, notes)
	if err != nil {
		log.Println("getMulti:notes", err)
		return notes, err
	}
	for i, note := range notes {
		note.ID = tag.NoteKeys[i].Encode()
	}
	return notes, nil
}

func (tag Tag) Children(c appengine.Context) ([]Tag, error) {
	children := make([]Tag, len(tag.ChildKeys))
	err := datastore.GetMulti(c, tag.ChildKeys, children)
	if err != nil {
		log.Println("getMulti:children", err)
	}
	return children, err
}
