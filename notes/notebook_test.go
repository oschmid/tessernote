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
	tags := set("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Add(note)

	actual := notebook.Note(note.Id)
	if actual.Title != note.Title {
		t.Fatalf("expected=%s actual=%s", note.Title, actual.Title)
	}
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Title)
	}
	if len(actual.Tags) != len(note.Tags) {
		t.Fatalf("expected=%d actual=%d", len(note.Tags), len(actual.Tags))
	}
	if !equals(actual.Tags, note.Tags) {
		t.Fatalf("expected=%v actual=%v", note.Tags, actual.Tags)
	}
}

func TestDeleteNote(t *testing.T) {
	tags := set("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Add(note)
	titles := notebook.Titles()
	if len(titles) != 1 {
		t.Fatalf("note not added")
	}

	notebook.Delete(note.Id)
	titles = notebook.Titles()
	if len(titles) > 0 {
		t.Fatalf("expected=0 actual=%d", len(titles))
	}
}

func TestDeleteNonExistentNote(t *testing.T) {
	title := "title"
	tags := set("tag1", "tag2")
	note := *NewNote(title, "body", tags)
	notebook := NewNoteBook()
	notebook.Add(note)

	note2 := *NewNote(title, "body", tags)
	if note2.Id == note.Id {
		t.Fatalf("note IDs are not unique")
	}

	notebook.Delete(note2.Id)
	titles := notebook.Titles()
	if len(titles) != 1 {
		t.Fatalf("expected=1 actual=%d %v", len(titles), titles)
	}
}

func TestEditNoteBody(t *testing.T) {
	tags := set("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Add(note)

	note.Body = "body2"
	actual := notebook.Note(note.Id)
	if actual.Body == note.Body {
		t.Fatal("NoteBook storage updated before call to update")
	}

	notebook.Update(note)
	actual = notebook.Note(note.Id)
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Body)
	}
}

func TestUpdateTags(t *testing.T) {
	tags := set("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Add(note)

	note.Tags = set("tag3", "tag4", "tag5")
	actual := notebook.Note(note.Id)
	if equals(actual.Tags, note.Tags) {
		t.Fatal("NoteBook storage updated before call to update")
	}

	notebook.Update(note)
	actual = notebook.Note(note.Id)
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
		tags := set(number("tag", i), number("tag", i+1))
		notebook.Add(*NewNote(number("title", i), "body", tags))
	}
	return notebook
}

func number(text string, num int) string {
	return fmt.Sprint(text, num)
}

func contains(slice []string, elem string) bool {
	for _, value := range slice {
		if value == elem {
			return true
		}
	}
	return false
}
