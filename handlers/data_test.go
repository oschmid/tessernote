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
	"strings"
	"testing"
)

func setUp() {
	body = nil
	notebook = *notes.NewNoteBook()
	notebook.Set(*notes.NewNote("title", "body", *sets.New("tag1", "tag2", "tag3")))
	notebook.Set(*notes.NewNote("title", "body", *sets.New("tag1", "tag3", "tag4")))
	notebook.Set(*notes.NewNote("title", "body", *sets.New("tag5")))
}

func TestNoTagsHandler(t *testing.T) {
	setUp()

	// build request
	contents, err := json.Marshal(new([]string))
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	w := new(responseWriter)
	r, err := http.NewRequest("POST", "http://www.grivet.com"+TAGS, strings.NewReader(string(contents)))
	if err != nil {
		t.Fatal(err)
	}
	TagsHandler(w, r)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1, "tag5": 1}
	var actual map[string]int
	err = json.Unmarshal(body, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestTagsInPostHandler(t *testing.T) {
	setUp()

	// build request
	contents, err := json.Marshal([]string{"tag1", "tag3"})
	if err != nil {
		t.Fatal(err)
	}

	// handle request
	w := new(responseWriter)
	r, err := http.NewRequest("POST", "http://www.grivet.com"+TAGS, strings.NewReader(string(contents)))
	if err != nil {
		t.Fatal(err)
	}
	TagsHandler(w, r)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1}
	var actual map[string]int
	err = json.Unmarshal(body, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

// http.ResponseWriter for testing
type responseWriter struct{}

var body []byte

// Not Implemented
func (w responseWriter) Header() http.Header {
	return nil
}

// Contents are written to "Body"
func (w responseWriter) Write(b []byte) (int, error) {
	body = b
	return len(body), nil
	//w.Body = make([]byte, len(b))
	//return copy(w.Body, b), nil
}

// Not Implemented
func (w responseWriter) WriteHeader(int) {

}
