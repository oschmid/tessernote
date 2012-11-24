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
)

type Note struct {
	ID           string `datastore:"-"` // datastore.Key.Encode()
	Body         string
	Created      time.Time
	LastModified time.Time
	NotebookKeys []*datastore.Key
	context      appengine.Context `datastore:"-"`
}

func (note Note) Notebooks() ([]Notebook, error) {
	var notebooks []Notebook
	err := datastore.GetMulti(note.context, note.NotebookKeys, notebooks)
	if err != nil {
		return notebooks, err
	}
	for _, notebook := range notebooks {
		notebook.context = note.context
	}
	return notebooks, nil
}

// TODO get tags by parsing body for hashtags

func (note *Note) SetBody(body string) error {
	note.Body = body
	note.LastModified = time.Now()
	key, err := datastore.DecodeKey(note.ID)
	if err != nil {
		return err
	}
	_, err = datastore.Put(note.context, key, note)
	return err
}
