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
	"testing"
)

func TestAddNote(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{"title", "body", tags}
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
	if !contains(actual.Tags, note.Tags[0]) {
		t.Fatalf("actual does not contain=%s", note.Tags[0])
	}
	if !contains(actual.Tags, note.Tags[1]) {
		t.Fatalf("actual does not contain=%s", note.Tags[1])
	}
}

func TestDeleteNote(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{"title", "body", tags}
	notebook := NewNoteBook()
	notebook.Add(note)
	titles := notebook.Titles()
	if len(titles) != 1 {
		t.Fatalf("note not added")
	}

	notebook.Delete(note.Title)
	titles = notebook.Titles()
	if len(titles) > 0 {
		t.Fatalf("expected=0 actual=%d", len(titles))
	}
}

func TestEditNoteBody(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{"title", "body", tags}
	notebook := NewNoteBook()
	notebook.Add(note)

	note.Body = "body2"
	notebook.Update(note)
	actual := notebook.Note(note.Title)
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Body)
	}
}

func TestUpdateTags(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{"title", "body", tags}
	notebook := NewNoteBook()
	notebook.Add(note)

	note.Tags = []string{"tag3", "tag4", "tag5"}
	notebook.Update(note)
	actual := notebook.Note(note.Title)
	if !equals(actual.Tags, note.Tags) {
		t.Fatal("expected=%v actual=%v", note.Tags, actual.Tags)
	}
}

func TestAllTitles(t *testing.T) {
	num := 10
	notebook := newFullNoteBook(num)

	titles := notebook.Titles()
	if len(titles) != num {
		t.Fatalf("expected=%d actual=%d", num, len(titles))
	}
	for i := 0; i < num; i++ {
		title := number("title", i)
		if !contains(titles, title) {
			t.Fatalf("%v does not contain %s", titles, title)
		}
	}
}

func TestAllTitlesOfTag(t *testing.T) {
	num := 10
	notebook := newFullNoteBook(num)

	titles := notebook.Titles(number("tag", 2))
	if len(titles) != 2 {
		t.Fatalf("expected=%d actual=%d %v", 2, len(titles), titles)
	}
}

// helper functions

func newFullNoteBook(num int) NoteBook {
	notebook := *NewNoteBook()
	for i := 0; i < num; i++ {
		tags := []string{number("tag", i), number("tag", i+1)}
		notebook.Add(Note{number("title", i), number("body", 0), tags})
	}
	return notebook
}

func number(text string, num int) string {
	return fmt.Sprint(text, num)
}
