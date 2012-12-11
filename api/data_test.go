/*
This file is part of Tessernote.

Tessernote is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Tessernote is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Tessernote.  If not, see <http://www.gnu.org/licenses/>.
*/

package api

import (
	"encoding/json"
	"github.com/oschmid/tessernote"
	"github.com/oschmid/tessernote/context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveNewNote(t *testing.T) {
	note := tessernote.Note{Body: "body"}
	bytes, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "https://tessernote.appspot.com"+SaveNoteURL, strings.NewReader(string(bytes)))
	if err != nil {
		t.Fatal(err)
	}
	serve(w, r)

	c := context.NewContext(r)
	defer context.Close()

	// TODO fake a user

	// check note was added
	notebook, err := tessernote.GetNotebook(c)
	if err != nil {
		t.Fatal(err)
	}
	notes, err := notebook.Notes(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 {
		t.Fatal(err)
	}
	if notes[0].Body != note.Body {
		t.Fatal("expected=%s actual=%s", notes[0].Body, note.Body)
	}

	// check response ID is the same
	response := []byte(w.Body.String())
	err = json.Unmarshal(response, note)
	if err != nil {
		t.Fatal(err, string(response))
	}
	if notes[0].ID != note.ID {
		t.Fatal("expected=%s actual=%s", notes[0].ID, note.ID)
	}
}
