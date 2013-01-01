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

package tessernote

import (
	"appengine"
	"appengine/datastore"
	"github.com/oschmid/cachestore"
	"github.com/oschmid/tessernote/hashtag"
	"strings"
	"unicode"
)

type Tag struct {
	Name         string // unique per Notebook
	NotebookKeys []*datastore.Key
	NoteKeys     []*datastore.Key
	ChildKeys    []*datastore.Key
}

// Notebooks returns all Notebooks that share this Tag.
func (tag Tag) Notebooks(c appengine.Context) ([]Notebook, error) {
	notebooks := make([]Notebook, len(tag.NotebookKeys))
	err := cachestore.GetMulti(c, tag.NotebookKeys, notebooks)
	if err != nil {
		c.Errorf("getting tag notebooks: %s", err)
	}
	return notebooks, err
}

// Notes returns all Notes that share this Tag.
func (tag Tag) Notes(c appengine.Context) ([]Note, error) {
	notes := make([]Note, len(tag.NoteKeys))
	err := cachestore.GetMulti(c, tag.NoteKeys, notes)
	if err != nil {
		c.Errorf("getting tag notes: %s", err)
		return notes, err
	}
	for i, note := range notes {
		note.ID = tag.NoteKeys[i].Encode()
	}
	return notes, nil
}

// Children returns all Tags this Note is parent to.
func (tag Tag) Children(c appengine.Context) ([]Tag, error) {
	children := make([]Tag, len(tag.ChildKeys))
	err := cachestore.GetMulti(c, tag.ChildKeys, children)
	if err != nil {
		c.Errorf("getting tag children: %s", err)
	}
	return children, err
}

// RelatedNotes returns the union of all Notes referred to by a set of Tags.
func RelatedNotes(tags []Tag, c appengine.Context) ([]Note, error) {
	if len(tags) == 0 {
		return *new([]Note), nil
	}

	noteKeys := tags[0].NoteKeys
	for i := 1; i < len(tags); i++ {
		noteKeys = unionKeys(noteKeys, tags[i].NoteKeys)
	}

	notes, err := make([]Note, len(noteKeys)), *new(error)
	if len(noteKeys) > 0 {
		err = cachestore.GetMulti(c, noteKeys, notes)
		if err != nil {
			c.Errorf("getting related notes: %s", err)
		}
	}
	return notes, err
}

// ParseTagNames parses a string for hashtags.
func ParseTagNames(text string) []string {
	var names []string
	matches := hashtag.Regex.FindAllString(text, len(text))
	for _, match := range matches {
		name := strings.TrimLeftFunc(match, isHashtagDecoration)
		names = append(names, name)
	}
	return names
}

// isHashtagDecoration returns true for hash characters (#) and whitespace
func isHashtagDecoration(r rune) bool {
	return r == '#' || r == '\uFF03' || unicode.IsSpace(r)
}

// NewTag creates a new Tag for a Note in a Notebook
func NewTag(name string, note Note, notebook Notebook, c appengine.Context) *Tag {
	tag := new(Tag)
	tag.Name = name
	tag.NotebookKeys = []*datastore.Key{notebook.Key(c)}
	tag.NoteKeys = []*datastore.Key{note.Key(c)}
	return tag
}

// Name returns the names of Tags in tag.
func Name(tag []Tag) []string {
	name := make([]string, len(tag))
	for i, t := range tag {
		name[i] = t.Name
	}
	return name
}
