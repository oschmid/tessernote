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
	"string/collections/maps"
	"string/collections/sets"
	"testing"
)

func TestAddNote(t *testing.T) {
	tags := *sets.New("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)

	actual, err := notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	compareNote(note, *actual, t)
}

func TestNonExistentNote(t *testing.T) {
	notebook := NewNoteBook()
	note, err := notebook.Note("this uuid doesn't refer to anything")
	if note != nil {
		t.Fatal("note should be nil")
	}
	if err == nil {
		t.Fatal("fetching a non-existent note should result in error")
	}
}

func TestDeleteNote(t *testing.T) {
	tags := *sets.New("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)
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
	tags := *sets.New("tag1", "tag2")
	note := *NewNote(title, "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)

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

func TestDeleteNoteTags(t *testing.T) {
	note := *NewNote("title", "body", *sets.New("tag1", "tag2"))
	notebook := NewNoteBook()
	notebook.Set(note)
	notebook.Set(*NewNote("title", "body", *sets.New("tag1", "tag3")))

	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 1}
	actual := *notebook.Tags()
	compareMaps(expected, actual, t)

	notebook.Delete(note.Id)
	expected = map[string]int{"tag1": 1, "tag3": 1}
	actual = *notebook.Tags()
	compareMaps(expected, actual, t)
}

func TestSetBody(t *testing.T) {
	tags := *sets.New("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)

	note.Body = "body2"
	actual, err := notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if actual.Body == note.Body {
		t.Fatal("NoteBook storage updated before call to update")
	}

	notebook.Set(note)
	actual, err = notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Body)
	}
}

func TestSetTags(t *testing.T) {
	tags := *sets.New("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)

	note.Tags = *sets.New("tag3", "tag4", "tag5")
	actual, err := notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if sets.Equal(actual.Tags, note.Tags) {
		t.Fatal("NoteBook storage updated before call to update")
	}

	notebook.Set(note)
	actual, err = notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if !sets.Equal(actual.Tags, note.Tags) {
		t.Fatal("expected=%v actual=%v", note.Tags, actual.Tags)
	}
}

func TestSetNew(t *testing.T) {
	tags := *sets.New("tag1", "tag2")
	note := *NewNote("title", "body", tags)
	notebook := NewNoteBook()
	notebook.Set(note)

	actual, err := notebook.Note(note.Id)
	if actual == nil {
		t.Fatal("note was not added")
	}
	if err != nil {
		t.Fatal("note was not added")
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
		if titles[i] != title {
			t.Fatalf("expected=%v actual=%v", title, titles[i])
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
	if titles[0] != "title1" {
		t.Fatalf("expected=title1 actual=%v", titles[0])
	}
	if titles[1] != "title2" {
		t.Fatalf("expected=title2 actual=%v", titles[1])
	}
}

func TestDistinguishingTagsFromTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.Set(*NewNote("title", "body", *sets.New("tag1", "tag2")))
	notebook.Set(*NewNote("title", "body", *sets.New("tag2", "tag3")))
	notebook.Set(*NewNote("title", "body", *sets.New("tag1", "tag2", "tag3")))

	expected := map[string]int{"tag1": 1, "tag2": 2, "tag3": 2}
	actual := *notebook.Tags("tag2", "tag3")
	if !maps.Equal(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestRenameTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.Set(Note{"uuid1", "title1", "body1", *sets.New("tag1", "tag2")})
	notebook.Set(Note{"uuid2", "title2", "body2", *sets.New("tag2", "tag3", "tag4")})
	notebook.Set(Note{"uuid3", "title3", "body3", *sets.New("tag1", "tag4")})

	notebook.RenameTag("tag1", "tag4")
	expectedTags := map[string]int{"tag2": 2, "tag3": 1, "tag4": 3}
	actualTags := *notebook.Tags()
	if !maps.Equal(expectedTags, actualTags) {
		t.Fatalf("expected=%v actual=%v", expectedTags, actualTags)
	}
	checkNoteTags(notebook, "uuid1", *sets.New("tag2", "tag4"), t)
	checkNoteTags(notebook, "uuid2", *sets.New("tag2", "tag3", "tag4"), t)
	checkNoteTags(notebook, "uuid3", *sets.New("tag4"), t)
}

func TestDeleteTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.Set(Note{"uuid1", "title1", "body1", *sets.New("tag1", "tag2")})
	notebook.Set(Note{"uuid2", "title2", "body2", *sets.New("tag2", "tag3", "tag4")})
	notebook.Set(Note{"uuid3", "title3", "body3", *sets.New("tag4")})

	notebook.DeleteTag("tag4")
	expectedTags := map[string]int{"tag1": 1, "tag2": 2, "tag3": 1}
	actualTags := *notebook.Tags()
	if !maps.Equal(expectedTags, actualTags) {
		t.Fatalf("expected=%v actual=%v", expectedTags, actualTags)
	}
	checkNoteTags(notebook, "uuid1", *sets.New("tag1", "tag2"), t)
	checkNoteTags(notebook, "uuid2", *sets.New("tag2", "tag3"), t)
	checkNoteTags(notebook, "uuid3", *sets.New(), t)
}

// helper functions

func newFullNoteBook(num int) NoteBook {
	notebook := *NewNoteBook()
	for i := 0; i < num; i++ {
		tags := *sets.New(number("tag", i), number("tag", i+1))
		notebook.Set(*NewNote(number("title", i), "body", tags))
	}
	return notebook
}

func number(text string, num int) string {
	return fmt.Sprint(text, num)
}

func compareNote(expected Note, actual Note, t *testing.T) {
	if !expected.Equal(actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func compareMaps(expected map[string]int, actual map[string]int, t *testing.T) {
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func checkNoteTags(notebook *NoteBook, uuid string, expected map[string]bool, t *testing.T) {
	note, err := notebook.Note(uuid)
	if err != nil {
		t.Fatal(err)
		return
	}

	actual := note.Tags
	if !sets.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}
