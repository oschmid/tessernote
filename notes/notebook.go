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

// TODO store notes by GUID

type NoteBook struct {
	notes map[string]string          // note title -> note body
	tags  map[string]map[string]bool // tag name -> note titles -> true if note has tag
}

func NewNoteBook() *NoteBook {
	noteBook := new(NoteBook)
	noteBook.notes = make(map[string]string)
	noteBook.tags = make(map[string]map[string]bool)
	return noteBook
}

// Returns the titles of all notes in the subset of notes specified by "tags"
// If no tags are specified, returns all notes
func (n NoteBook) Titles(tags ...string) []string {
	if tags == nil || len(tags) == 0 {
		return n.allTitles()
	}

	// get the notes of the first tag
	notes := n.allTitlesOfTag(tags[0])
	if len(notes) == 0 {
		return notes
	}

	// intersect with each following subset of notes
	for _, tag := range tags[1:] {
		titles := n.allTitlesOfTag(tag)
		notes = intersection(notes, titles)
		if len(notes) == 0 {
			break
		}
	}

	return notes
}

func (n NoteBook) allTitles() []string {
	titles := []string{}
	for title, _ := range n.notes {
		titles = append(titles, title)
	}
	return titles
}

func (n NoteBook) allTitlesOfTag(tag string) []string {
	titles := []string{}
	for title, tagged := range n.tags[tag] {
		if tagged {
			titles = append(titles, title)
		}
	}
	return titles
}

func (n NoteBook) Add(note Note) {
	n.notes[note.Title] = note.Body
	for _, tag := range note.Tags {
		n.addTag(tag, note.Title)
	}
}

func (n NoteBook) addTag(tag string, note string) {
	if n.tags[tag] == nil {
		n.tags[tag] = make(map[string]bool)
	}
	n.tags[tag][note] = true
}

func (n NoteBook) Note(title string) Note {
	body := n.notes[title]
	tags := n.TagsOfNote(title)
	return Note{title, body, tags}
}

func (n NoteBook) TagsOfNote(title string) []string {
	tags := *new([]string)
	for tag, notes := range n.tags {
		if notes[title] {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (n NoteBook) Delete(title string) {
	// delete body
	delete(n.notes, title)

	// delete note from tags
	for _, tag := range n.TagsOfNote(title) {
		delete(n.tags[tag], title)
	}
}

func (n NoteBook) Update(note Note) {
	// update body
	n.notes[note.Title] = note.Body

	// update tags
	oldTags := n.TagsOfNote(note.Title)
	if !equals(oldTags, note.Tags) {
		// remove note from tags it no longer has
		remove := difference(oldTags, note.Tags)
		for _, tag := range remove {
			delete(n.tags[tag], note.Title)
		}

		// add note to tags it has gained
		add := difference(note.Tags, oldTags)
		for _, tag := range add {
			n.addTag(tag, note.Title)
		}
	}
}
