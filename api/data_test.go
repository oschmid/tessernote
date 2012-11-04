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
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notes"
	"string/collections/maps"
	"string/collections/sets"
	"string/collections/slices"
	"strings"
	"testing"
)

const (
	uuid1 = "00000000-0000-0000-0000-000000000000"
	uuid2 = "00000000-0000-0000-0000-000000000001"
	uuid3 = "00000000-0000-0000-0000-000000000002"
)

func setUp() *httptest.ResponseRecorder {
	notebook = *notes.NewNoteBook()
	notebook.SetNote(notes.Note{Id: uuid1, Title: "title1", Body: "body1 #tag1 #tag2 #tag3"})
	notebook.SetNote(notes.Note{Id: uuid2, Title: "title2", Body: "body2 #tag1 #tag3 #tag4"})
	notebook.SetNote(notes.Note{Id: uuid3, Title: "title3", Body: "body3 #tag5"})
	return httptest.NewRecorder()
}

func TestGetAllTags(t *testing.T) {
	recorder := setUp()
	body := marshal(*new([]string), t) // "null"
	request := makePostRequest(GetTagsUrl, body, t)
	_, handle := MakePostHandler(GetTagsUrl, GetTags)
	handle(recorder, request)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1, "tag5": 1}
	var actual map[string]int
	unmarshal(recorder, &actual, t)
	compareMaps(expected, actual, t)
}

func TestGetSomeTags(t *testing.T) {
	recorder := setUp()
	body := marshal([]string{"tag1", "tag3"}, t)
	request := makePostRequest(GetTagsUrl, body, t)
	_, handle := MakePostHandler(GetTagsUrl, GetTags)
	handle(recorder, request)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1}
	var actual map[string]int
	unmarshal(recorder, &actual, t)
	compareMaps(expected, actual, t)
}

func TestGetTagsWithSomeMissing(t *testing.T) {
	recorder := setUp()
	body := marshal([]string{"bad1", "tag1", "bad2", "tag3", "bad3"}, t)
	request := makePostRequest(GetTagsUrl, body, t)
	_, handle := MakePostHandler(GetTagsUrl, GetTags)
	handle(recorder, request)

	// verify response
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1}
	var actual map[string]int
	unmarshal(recorder, &actual, t)
	compareMaps(expected, actual, t)
}

func TestRenameTags(t *testing.T) {
	recorder := setUp()
	body := marshal(map[string]string{"tag1": "tag6", "tag3": "tag4"}, t)
	request := makePostRequest(RenameTagsUrl, body, t)
	_, handle := MakePostHandler(RenameTagsUrl, RenameTags)
	handle(recorder, request)

	// verify tags
	expected := map[string]int{"tag2": 1, "tag4": 2, "tag5": 1, "tag6": 2}
	actual := *notebook.RelatedTags()
	compareMaps(expected, actual, t)

	// verify notes
	checkNoteTags(uuid1, *sets.New("tag2", "tag4", "tag6"), t)
	checkNoteTags(uuid2, *sets.New("tag4", "tag6"), t)
	checkNoteTags(uuid3, *sets.New("tag5"), t)
}

func TestDeleteTags(t *testing.T) {
	recorder := setUp()
	body := marshal([]string{"tag3", "tag5"}, t)
	request := makePostRequest(DeleteTagsUrl, body, t)
	_, handle := MakePostHandler(DeleteTagsUrl, DeleteTags)
	handle(recorder, request)

	// verify tags
	expected := map[string]int{"tag1": 2, "tag2": 1, "tag4": 1}
	actual := *notebook.RelatedTags()
	compareMaps(expected, actual, t)

	// verify notes
	checkNoteTags(uuid1, *sets.New("tag1", "tag2"), t)
	checkNoteTags(uuid2, *sets.New("tag1", "tag4"), t)
	checkNoteTags(uuid3, *sets.New(), t)
}

func TestGetAllTitles(t *testing.T) {
	recorder := setUp()
	body := marshal(*new([][]string), t)
	request := makePostRequest(GetTitlesUrl, body, t)
	_, handle := MakePostHandler(GetTitlesUrl, GetTitles)
	handle(recorder, request)

	// verify response
	expected := [][]string{[]string{"title1", uuid1}, []string{"title2", uuid2}, []string{"title3", uuid3}}
	var actual [][]string
	unmarshal(recorder, &actual, t)
	for i, _ := range expected {
		compareSlices(expected[i], actual[i], t)
	}
}

func TestGetSomeTitles(t *testing.T) {
	recorder := setUp()
	body := marshal([]string{"tag3"}, t)
	request := makePostRequest(GetTitlesUrl, body, t)
	_, handle := MakePostHandler(GetTitlesUrl, GetTitles)
	handle(recorder, request)

	// verify response
	expected := [][]string{[]string{"title1", uuid1}, []string{"title2", uuid2}}
	var actual [][]string
	unmarshal(recorder, &actual, t)
	for i, _ := range expected {
		compareSlices(expected[i], actual[i], t)
	}
}

func TestGetTitlesWithSomeMissingTags(t *testing.T) {
	recorder := setUp()
	body := marshal([]string{"bad1", "tag3", "bad2"}, t)
	request := makePostRequest(GetTitlesUrl, body, t)
	_, handle := MakePostHandler(GetTitlesUrl, GetTitles)
	handle(recorder, request)

	// verify response
	expected := [][]string{[]string{"title1", uuid1}, []string{"title2", uuid2}}
	var actual [][]string
	unmarshal(recorder, &actual, t)
	for i, _ := range expected {
		compareSlices(expected[i], actual[i], t)
	}
}

func TestGetNote(t *testing.T) {
	recorder := setUp()
	request := makeGetRequest(GetNoteUrl, uuid1, t)
	_, handle := MakeGetHandler(GetNoteUrl, GetNote)
	handle(recorder, request)

	// verify
	expected := notes.Note{Id: uuid1, Title: "title1", Body: "body1 #tag1 #tag2 #tag3"}
	var actual notes.Note
	unmarshal(recorder, &actual, t)
	compareNote(expected, actual, t)
}

func TestNewNote(t *testing.T) {
	recorder := setUp()
	note := notes.Note{Title: "untitled"}
	if note.Id != "" {
		t.Fatal(note)
	}
	body := marshal(note, t)
	request := makePostRequest(SaveNoteUrl, body, t)
	_, handle := MakePostHandler(SaveNoteUrl, SaveNote)
	handle(recorder, request)

	// verify response
	if !uuidValidator.MatchString(recorder.Body.String()) {
		t.Fatal("new note was not assigned a UUID")
	}
	note.Id = recorder.Body.String()

	// verify tags
	expectedTags := map[string]int{"tag1": 2, "tag2": 1, "tag3": 2, "tag4": 1, "tag5": 1}
	actualTags := *notebook.RelatedTags()
	compareMaps(expectedTags, actualTags, t)

	// verify note
	compareNoteInNoteBook(note, note.Id, t)
}

func TestEditNote(t *testing.T) {
	recorder := setUp()
	note := getNote(uuid2, t)
	note.Body = "new body"
	check := getNote(uuid2, t)
	if check.Body == note.Body {
		t.Fatal("note body passed by reference")
	}
	body := marshal(note, t)
	request := makePostRequest(SaveNoteUrl, body, t)
	_, handle := MakePostHandler(SaveNoteUrl, SaveNote)
	handle(recorder, request)

	// verify response
	if recorder.Body.String() != note.Id {
		t.Fatalf("expected=%v actual=%v", note.Id, recorder.Body.String())
	}

	// verify note
	compareNoteInNoteBook(*note, uuid2, t)
}

// helper functions

func makePostRequest(url string, body []byte, t *testing.T) *http.Request {
	request, err := http.NewRequest("POST", "http://www.grivet.ca"+url, strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func makeGetRequest(url string, value string, t *testing.T) *http.Request {
	request, err := http.NewRequest("GET", "http://www.grivet.ca"+url+value, nil)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func marshal(v interface{}, t *testing.T) []byte {
	body, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func unmarshal(recorder *httptest.ResponseRecorder, destination interface{}, t *testing.T) {
	body := []byte(recorder.Body.String())
	err := json.Unmarshal(body, destination)
	if err != nil {
		t.Fatal(err, ":", string(body))
	}
}

func getNote(id string, t *testing.T) *notes.Note {
	note, err := notebook.Note(id)
	if err != nil {
		t.Fatal(err)
	}
	return note
}

func compareMaps(expected map[string]int, actual map[string]int, t *testing.T) {
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func compareSlices(expected []string, actual []string, t *testing.T) {
	if !slices.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func compareNote(expected notes.Note, actual notes.Note, t *testing.T) {
	if !expected.Equal(actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func compareNoteInNoteBook(note notes.Note, id string, t *testing.T) {
	actual := getNote(id, t)
	compareNote(note, *actual, t)
}

func checkNoteTags(id string, expected map[string]bool, t *testing.T) {
	note := getNote(id, t)
	actual := note.Tags()
	if !sets.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}
