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
	"net/http"
	"note"
	"strings"
	"testing"
)

func TestSaveNewNote(t *testing.T) {
	note := note.Note{Body: "body"}
	bytes, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}
	_, err = http.NewRequest("POST", "https://tessernote.appspot.com"+SaveNoteURL, strings.NewReader(string(bytes)))
	if err != nil {
		t.Fatal(err)
	}
	// TODO make save request
}
