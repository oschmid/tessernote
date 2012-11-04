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
	"sort"
	"string/collections/sets"
	"strings"
)

type NoteBook struct {
	Notes map[string]Note            // note ID -> Note
	Tags  map[string]map[string]bool // tag name -> note IDs -> true if note has tag
}

func NewNoteBook() *NoteBook {
	noteBook := new(NoteBook)
	noteBook.Notes = make(map[string]Note)
	noteBook.Tags = make(map[string]map[string]bool)
	return noteBook
}

// Returns the UUIDs of all notes in the subset of notes specified by "tags"
// If no tags are specified, returns all note UUIDs
func (n NoteBook) UUIDs(tags ...string) map[string]bool {
	if tags == nil || len(tags) == 0 {
		return n.allUUIDs()
	}

	// strip non-existent tags
	strippedTags := *new([]string)
	for _, tag := range tags {
		_, contained := n.Tags[tag]
		if contained {
			strippedTags = append(strippedTags, tag)
		}
	}
	tags = strippedTags
	if len(tags) == 0 {
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
	for id, tagged := range n.Tags[tag] {
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
		title := n.Notes[id].Title
		titles = append(titles, []string{title, id})
	}
	sort.Sort(StringSliceSlice(titles))
	return titles
}

// Given a set of tags T, returns all the tags that refer to all the notes
// referred to by T and just how many notes each tag refers to
// If "tags" is empty, returns all tags (and how many notes each tag refers to
func (n NoteBook) RelatedTags(tags ...string) *map[string]int {
	notes := n.UUIDs(tags...)
	super := make(map[string]int)
	for id, _ := range notes {
		note := n.Notes[id]
		super = *union(super, note.Tags())
	}
	return &super
}

func union(a map[string]int, b map[string]bool) *map[string]int {
	for tag, _ := range b {
		count, _ := a[tag]
		a[tag] = count + 1
	}
	return &a
}

func (n *NoteBook) Delete(id string) {
	// delete body
	delete(n.Notes, id)

	// delete note from tags
	note := n.Notes[id]
	for tag, _ := range note.Tags() {
		delete(n.Tags[tag], id)
	}
}

// Adds note if it didn't exist before, updates all information if it did.
func (n *NoteBook) SetNote(note Note) {
	oldNote, contained := n.Notes[note.Id]
	oldTags := *sets.New()
	if contained {
		oldTags = oldNote.Tags()
	}

	// set body
	n.Notes[note.Id] = note

	// set tags
	if !sets.Equal(oldTags, note.Tags()) {
		// remove note from tags it no longer has
		remove := *sets.Difference(oldTags, note.Tags())
		for tag, _ := range remove {
			delete(n.Tags[tag], note.Id)
		}

		// add note to tags it has gained
		add := *sets.Difference(note.Tags(), oldTags)
		for tag, _ := range add {
			n.addTag(tag, note.Id)
		}
	}
}

func (n *NoteBook) addTag(tag string, noteId string) {
	if n.Tags[tag] == nil {
		n.Tags[tag] = make(map[string]bool)
	}
	n.Tags[tag][noteId] = true
}

func (n *NoteBook) RenameTag(old string, new string) {
	// replace tag in note bodies
	notes := n.Tags[old]
	for id, _ := range notes {
		note, contained := n.Notes[id]
		if contained {
			body := strings.Replace(note.Body, "#"+old+" ", "#"+new+" ", -1)
			body = strings.Replace(body, "#"+old, "#"+new, -1)
			note.SetBody(body)
			n.SetNote(note)
		}
	}

	delete(n.Tags, old)
	n.Tags[new] = notes
}

func (n *NoteBook) DeleteTag(tag string) {
	// remove tag from note bodies
	notes := n.Tags[tag]
	for id, _ := range notes {
		note, contained := n.Notes[id]
		if contained {
			body := strings.Replace(note.Body, "#"+tag+" ", "", -1)
			body = strings.Replace(body, "#"+tag, "", -1)
			note.SetBody(body)
			n.SetNote(note)
		}
	}

	delete(n.Tags, tag)
}
