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
	"strings"
	"time"
)

type Note struct {
	ID           string `datastore:"-"` // datastore.Key.Encode()
	Body         string
	Created      time.Time
	LastModified time.Time
	TagKeys      []*datastore.Key
	NotebookKeys []*datastore.Key
	context      appengine.Context `datastore:"-"`
}

func (note Note) Tags() ([]Tag, error) {
	var tags []Tag
	err := datastore.GetMulti(note.context, note.TagKeys, tags)
	if err != nil {
		return tags, err
	}
	for _, tag := range tags {
		tag.context = note.context
	}
	return tags, nil
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

func (note *Note) SetBody(body string) {
	note.Body = body
	note.LastModified = time.Now()
}

// parses tag names from body
func (note *Note) ParseTagNames() []string {
	var names []string
	matches := Hashtag.FindAllString(note.Body, len(note.Body))
	for _, match := range matches {
		name := strings.TrimFunc(match, isHashtagDecoration)
		names = append(names, name)
	}
	return names
}
