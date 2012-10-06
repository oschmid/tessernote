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
	tags := []string{"tag1", "tag2"}
	note := Note{Title: "title", Body: "body", Tags: tags}
	notebook := NewNoteBook()
	notebook.Add(note)

	actual := notebook.Note(note.Title)
	if actual.Title != note.Title {
		t.Fatalf("expected=%s actual=%s", note.Title, actual.Title)
	}
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Title)
	}
	if len(actual.Tags) != len(note.Tags) {
		t.Fatalf("expected=%d actual=%d", len(note.Tags), len(actual.Tags))
	}
	if actual.Tags[0] != note.Tags[0] && actual.Tags[1] != note.Tags[0] {
		t.Fatalf("actual does not contain=%s", note.Tags[0])
	}
	if actual.Tags[0] != note.Tags[1] && actual.Tags[1] != note.Tags[1] {
		t.Fatalf("actual does not contain=%s", note.Tags[1])
	}
}

func TestRemoveNote(t *testing.T) {
	t.Fatal("Not yet implemented")
}

func TestEditNote(t *testing.T) {
	t.Fatal("Not yet implemented")
}

func TestAllTitles(t *testing.T) {
	t.Fatal("Not yet implemented")
}

func TestAllTitlesOfTag(t *testing.T) {
	t.Fatal("Not yet implemented")
}
