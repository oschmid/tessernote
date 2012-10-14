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
package handlers

import (
	"encoding/json"
	"net/http"
	"notes"
	"string/collections/maps"
	"string/collections/sets"
	"string/collections/slices"
	"strings"
	"testing"
)

func setUp() {
	responseBody = nil
	w = *new(responseWriter)
	notebook = *notes.NewNoteBook()
	notebook.Set(notes.Note{"uuid1", "title1", "body1", *sets.New("tag1", "tag2", "tag3")})
	notebook.Set(notes.Note{"uuid2", "title2", "body2", *sets.New("tag1", "tag3", "tag4")})
	notebook.Set(notes.Note{"uuid3", "title3", "body3", *sets.New("tag5")})
}

func TestTagsHandlerNoPost(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal(*new([]string))
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTagsGet, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	GetTagsHandler(w, r, body)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1, "tag5": 1}
	var actual map[string]int
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestTagsHandler(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal([]string{"tag1", "tag3"})
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTagsGet, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	GetTagsHandler(w, r, body)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1}
	var actual map[string]int
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestRenameTagsHandler(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal(map[string]string{"tag1": "tag6", "tag3": "tag4"})
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTagsRename, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	RenameTagsHandler(w, r, body)

	// verify tags
	expectedTags := map[string]int{"tag2": 1, "tag4": 2, "tag5": 1, "tag6": 2}
	actualTags := *notebook.Tags()
	if !maps.Equal(expectedTags, actualTags) {
		t.Fatalf("expected=%v actual=%v", expectedTags, actualTags)
	}

	// verify notes
	checkNoteTags("uuid1", *sets.New("tag2", "tag4", "tag6"), t)
	checkNoteTags("uuid2", *sets.New("tag4", "tag6"), t)
	checkNoteTags("uuid3", *sets.New("tag5"), t)
}

func TestDeleteTagsHandler(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal([]string{"tag3", "tag5"})
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTagsDelete, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	DeleteTagsHandler(w, r, body)

	// verify tags
	expectedTags := map[string]int{"tag1": 2, "tag2": 1, "tag4": 1}
	actualTags := *notebook.Tags()
	if !maps.Equal(expectedTags, actualTags) {
		t.Fatalf("expected=%v actual=%v", expectedTags, actualTags)
	}

	// verify notes
	checkNoteTags("uuid1", *sets.New("tag1", "tag2"), t)
	checkNoteTags("uuid2", *sets.New("tag1", "tag4"), t)
	checkNoteTags("uuid3", *sets.New(), t)
}

func TestTitlesHandlerNoPost(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal(*new([]string))
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTitles, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	TitlesHandler(w, r, body)

	// verify response
	expected := []string{"title1", "title2", "title3"}
	var actual []string
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestTitlesHandler(t *testing.T) {
	setUp()

	// build request
	body, err := json.Marshal([]string{"tag3"})
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	r, err := http.NewRequest("POST", "http://www.grivet.com"+UrlTitles, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	TitlesHandler(w, r, body)

	// verify response
	expected := []string{"title1", "title2"}
	var actual []string
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestGetNoteHandler(t *testing.T) {
	t.Fatal("not yet implemented")
}

func TestSaveNoteHandler(t *testing.T) {
	t.Fatal("not yet implemented")
}

// http.ResponseWriter for testing
type responseWriter struct{}

var responseBody []byte
var w responseWriter

// Not Implemented
func (w responseWriter) Header() http.Header {
	return nil
}

// Contents are written to "Body"
func (w responseWriter) Write(b []byte) (int, error) {
	responseBody = b
	return len(responseBody), nil
}

// Not Implemented
func (w responseWriter) WriteHeader(int) {

}

func checkNoteTags(uuid string, expected map[string]bool, t *testing.T) {
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
