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
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)

	actual, err := notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	compareNote(note, *actual, t)
}

func TestGetNonExistentNote(t *testing.T) {
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
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)
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
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)

	note2 := *NewNote(title, "body #tag1 #tag2")
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
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)
	notebook.SetNote(*NewNote("title", "body #tag1 #tag3"))

	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 1}
	actual := *notebook.RelatedTags()
	compareMaps(expected, actual, t)

	notebook.Delete(note.Id)
	expected = map[string]int{"tag1": 1, "tag3": 1}
	actual = *notebook.RelatedTags()
	compareMaps(expected, actual, t)
}

func TestSetNoteBody(t *testing.T) {
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)

	note.Body = "body2"
	actual, err := notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if actual.Body == note.Body {
		t.Fatal("NoteBook storage updated before call to update")
	}

	notebook.SetNote(note)
	actual, err = notebook.Note(note.Id)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if actual.Body != note.Body {
		t.Fatalf("expected=%s actual=%s", note.Body, actual.Body)
	}
}

func TestSetNewNote(t *testing.T) {
	note := *NewNote("title", "body #tag1 #tag2")
	notebook := NewNoteBook()
	notebook.SetNote(note)

	actual, err := notebook.Note(note.Id)
	if actual == nil {
		t.Fatal("note was not added")
	}
	if err != nil {
		t.Fatal("note was not added")
	}
}

func TestGetAllTitles(t *testing.T) {
	num := 10
	notebook := newFullNoteBook(num)

	titles := notebook.Titles()
	if len(titles) != num {
		t.Fatalf("expected=%d actual=%d", num, len(titles))
	}

	// verify ordered list of tuples (title, id)
	for i, elem := range titles {
		expectedTitle, expectedId := title(i), id(i)
		actualTitle, actualId := elem[0], elem[1]
		if actualTitle != expectedTitle {
			t.Fatalf("expected=%v actual=%v", expectedTitle, actualTitle)
		}
		if actualId != expectedId {
			t.Fatalf("expected=%v actual=%v", expectedId, actualId)
		}
	}
}

func TestGetTitlesOfNotesWithTags(t *testing.T) {
	num := 10
	notebook := newFullNoteBook(num)

	titles := notebook.Titles(tag(2))
	if len(titles) != 2 {
		t.Fatalf("expected=%d actual=%d %v", 2, len(titles), titles)
	}

	// verify ordered list of tuples (title, id)
	for i, elem := range titles {
		expectedTitle, expectedId := title(i+1), id(i+1)
		actualTitle, actualId := elem[0], elem[1]
		if actualTitle != expectedTitle {
			t.Fatalf("expected=%v actual=%v", expectedTitle, actualTitle)
		}
		if actualId != expectedId {
			t.Fatalf("expected=%v actual=%v", expectedId, actualId)
		}
	}
}

func TestGetTitlesOfNotesWithMissingTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2"))
	notebook.SetNote(*NewNote("title", "body #tag2 #tag3"))
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2 #tag3"))

	expected := map[string]int{"tag1": 1, "tag2": 2, "tag3": 2}
	actual := *notebook.RelatedTags("not a tag", "tag3", "also not a tag")
	if !maps.Equal(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestGetAllTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2"))
	notebook.SetNote(*NewNote("title", "body #tag2 #tag3"))
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2 #tag3"))

	expected := map[string]int{"tag1": 2, "tag2": 3, "tag3": 2}
	actual := *notebook.RelatedTags()
	if !maps.Equal(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestGetRelatedTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2"))
	notebook.SetNote(*NewNote("title", "body #tag2 #tag3"))
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2 #tag3"))

	expected := map[string]int{"tag1": 1, "tag2": 2, "tag3": 2}
	actual := *notebook.RelatedTags("tag2", "tag3")
	if !maps.Equal(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestGetRelatedTagsWithSomeMissing(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2"))
	notebook.SetNote(*NewNote("title", "body #tag2 #tag3"))
	notebook.SetNote(*NewNote("title", "body #tag1 #tag2 #tag3"))

	expected := map[string]int{"tag1": 1, "tag2": 2, "tag3": 2}
	actual := *notebook.RelatedTags("not a tag", "tag2", "also not a tag", "tag3", "again not a tag")
	if !maps.Equal(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestRenameTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(Note{Id: "uuid1", Title: "title1", Body: "body1 #tag1 #tag2"})
	notebook.SetNote(Note{Id: "uuid2", Title: "title2", Body: "body2 #tag2 #tag3 #tag4"})
	notebook.SetNote(Note{Id: "uuid3", Title: "title3", Body: "body3 #tag1 #tag4"})

	notebook.RenameTag("tag1", "tag4")
	expectedTags := map[string]int{"tag2": 2, "tag3": 1, "tag4": 3}
	actualTags := *notebook.RelatedTags()
	if !maps.Equal(expectedTags, actualTags) {
		t.Fatalf("expected=%v actual=%v", expectedTags, actualTags)
	}
	checkNoteTags(notebook, "uuid1", *sets.New("tag2", "tag4"), t)
	checkNoteTags(notebook, "uuid2", *sets.New("tag2", "tag3", "tag4"), t)
	checkNoteTags(notebook, "uuid3", *sets.New("tag4"), t)
}

func TestDeleteTags(t *testing.T) {
	notebook := NewNoteBook()
	notebook.SetNote(Note{Id: "uuid1", Title: "title1", Body: "body1 #tag1 #tag2"})
	notebook.SetNote(Note{Id: "uuid2", Title: "title2", Body: "body2 #tag2 #tag3 #tag4"})
	notebook.SetNote(Note{Id: "uuid3", Title: "title3", Body: "body3 #tag4"})

	notebook.DeleteTag("tag4")
	expectedTags := map[string]int{"tag1": 1, "tag2": 2, "tag3": 1}
	actualTags := *notebook.RelatedTags()
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
		notebook.SetNote(Note{Id: id(i), Title: title(i), Body: "body #" + tag(i) + " #" + tag(i+1)})
	}
	return notebook
}

func id(num int) string {
	return number("id", num)
}

func title(num int) string {
	return number("title", num)
}

func tag(num int) string {
	return number("tag", num)
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

	actual := note.Tags()
	if !sets.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}
