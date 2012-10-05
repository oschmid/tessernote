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
	"testing"
)

func TestAddNote(t *testing.T) {
	note := Note{Title:"title", Body:"body", Tags:[]string{"tag1","tag2"}}
	notebook := new(NoteBook)
	notebook.Add(note)

	actual := notebook.Note(note.Title)
	if actual == nil ||
		actual.Title != note.Title ||
		actual.Body != note.Body ||
		len(actual.Tags) != len(note.Tags) ||
		actual.Tags[0] != note.Tags[0] ||
		actual.Tags[1] != note.Tags[1] {
		t.Fail()
	}
}

func TestRemoveNote(t *testing.T) {
	t.Fail()
}

func TestEditNote(t *testing.T) {
	t.Fail()
}

func TestAllTitles(t *testing.T) {
	t.Fail()
}

func TestAllTitlesOfTag(t *testing.T) {
	t.Fail()
}
