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
package notes

import (
	"fmt"
	"sort"
	"string/collections/sets"
	"strings"
)

const TitleBodySeparator string = "\n"

// TODO export notes and tags for gob?
type NoteBook struct {
	Notes map[string]string          // note ID -> note title and body
	tags  map[string]map[string]bool // tag name -> note IDs -> true if note has tag
}

func NewNoteBook() *NoteBook {
	noteBook := new(NoteBook)
	noteBook.Notes = make(map[string]string)
	noteBook.tags = make(map[string]map[string]bool)
	return noteBook
}

// Returns the UUIDs of all notes in the subset of notes specified by "tags"
// If no tags are specified, returns all note UUIDs
func (n NoteBook) UUIDs(tags ...string) map[string]bool {
	if tags == nil || len(tags) == 0 {
		return n.allUUIDs()
	}

	// get the note UUIDs associated with the first tag
	uuids := n.allUUIDsOfTag(tags[0])
	if len(uuids) == 0 {
		return uuids
	}

	// intersect with each following subset of note UUIDs
	for _, tag := range tags[1:] {
		tagUUIDs := n.allUUIDsOfTag(tag)
		uuids = *sets.Intersection(uuids, tagUUIDs)
		if len(uuids) == 0 {
			break
		}
	}

	return uuids
}

func (n NoteBook) allUUIDs() map[string]bool {
	uuids := make(map[string]bool)
	for id, _ := range n.Notes {
		uuids[id] = true
	}
	return uuids
}

func (n NoteBook) allUUIDsOfTag(tag string) map[string]bool {
	uuids := make(map[string]bool)
	for id, tagged := range n.tags[tag] {
		if tagged {
			uuids[id] = true
		}
	}
	return uuids
}

// Returns the titles and IDs of all notes in the subset of notes specified by "tags"
// If no tags are specified, returns the titles and IDs of all notes.
func (n NoteBook) Titles(tags ...string) [][]string {
	titles := [][]string{}
	for id, _ := range n.UUIDs(tags...) {
		title := strings.SplitN(n.Notes[id], TitleBodySeparator, 2)[0]
		titles = append(titles, []string{title, id})
	}
	sort.Sort(StringSliceSlice(titles))
	return titles
}

func (n NoteBook) Note(id string) (*Note, error) {
	note := strings.SplitN(n.Notes[id], TitleBodySeparator, 2)
	if len(note) != 2 {
		return nil, fmt.Errorf("note \"%s\" does not exist", id)
	}
	title, body := note[0], note[1]
	tags := n.TagsOfNote(id)
	return &Note{id, title, body, tags}, nil
}

func (n NoteBook) TagsOfNote(id string) map[string]bool {
	tags := make(map[string]bool)
	for tag, notes := range n.tags {
		if notes[id] {
			tags[tag] = true
		}
	}
	return tags
}

// Given a set of tags T, returns all the tags that refer to all the notes
// referred to by T and just how many notes each tag refers to
// If "tags" is empty, returns all tags (and how many notes each tag refers to
func (n NoteBook) Tags(tags ...string) *map[string]int {
	notes := n.UUIDs(tags...)
	super := make(map[string]int)
	for id, _ := range notes {
		super = *union(super, n.TagsOfNote(id))
	}
	return &super
}

func union(a map[string]int, b map[string]bool) *map[string]int {
	for tag, _ := range b {
		count, contained := a[tag]
		if contained {
			a[tag] = count + 1
		} else {
			a[tag] = 1
		}
	}
	return &a
}

func (n NoteBook) Delete(id string) {
	// delete body
	delete(n.Notes, id)

	// delete note from tags
	for tag, _ := range n.TagsOfNote(id) {
		delete(n.tags[tag], id)
	}
}

// Adds note if it didn't exist before, updates all information if it did.
func (n NoteBook) Set(note Note) {
	// set body
	n.Notes[note.Id] = note.Title + TitleBodySeparator + note.Body

	// set tags
	oldTags := n.TagsOfNote(note.Id)
	if !sets.Equal(oldTags, note.Tags) {
		// remove note from tags it no longer has
		remove := *sets.Difference(oldTags, note.Tags)
		for tag, _ := range remove {
			delete(n.tags[tag], note.Id)
		}

		// add note to tags it has gained
		add := *sets.Difference(note.Tags, oldTags)
		for tag, _ := range add {
			n.addTag(tag, note.Id)
		}
	}
}

func (n NoteBook) addTag(tag string, noteId string) {
	if n.tags[tag] == nil {
		n.tags[tag] = make(map[string]bool)
	}
	n.tags[tag][noteId] = true
}

func (n NoteBook) RenameTag(old string, new string) {
	notes, contained := n.tags[old]
	if contained {
		previous, contained := n.tags[new]
		if contained {
			n.tags[new] = *sets.Union(previous, notes)
		} else {
			n.tags[new] = notes
		}
		delete(n.tags, old)
	}
}

func (n NoteBook) DeleteTag(tag string) {
	delete(n.tags, tag)
}
