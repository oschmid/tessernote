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

type NoteBook struct {
	notes map[string] string // note title -> note body
	tags map[string] map[string]bool // tag name -> note titles -> true if note has tag
}

// Returns the titles of all notes in the subset of notes specified by "tags"
// If no tags are specified, returns all notes
func (n NoteBook) Notes(tags ...string) []string {
	if tags == nil || len(tags) == 0 {
		return n.allTitles()
	}

	// get the notes of the first tag
	notes := n.allTitlesOfTag(tags[0])

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
	titles := make([]string, len(n.tags))
	for title, tagged := range n.tags[tag] {
		if tagged {
			titles = append(titles, title)
		}
	}
	return titles
}

func (n NoteBook) Add(note Note) {
	n.notes[note.Title]=note.Body
	for _, tag := range note.Tags {
		n.tags[tag][note.Title]=true
	}
}

func (n NoteBook) Note(title string) Note {
	body := n.notes[title]
	tags := n.Tags(title)
	return Note{title, body, tags}
}

func (n NoteBook) Tags(title string) []string {
	tags := make([]string, len(n.tags))
	for tag, notes := range n.tags {
		if notes[title] {
			tags = append(tags, tag)
		}
	}
	return tags
}
