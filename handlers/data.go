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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"string/collections/slice"
)

const TAGS = "/data/tags"

/*
const TITLES = "/data/titles"
const NOTE = "/note/"
const SAVE = "/save/"
*/

// Returns tags (optional tags in POST/JSON)
func TagsHandler(w http.ResponseWriter, r *http.Request) {
	// convert JSON -> map[string]bool
	var tags map[string]bool
	body, _ := readBody(r) // TODO handle errors
	_ = json.Unmarshal(body, &tags)

	tags = notebook.Tags(*slice.FromSet(tags)...)

	// convert map[string]bool -> JSON
	response, _ := json.Marshal(tags)
	_, _ = w.Write(response)
}

// TODO get titles (optional tags in POST/JSON)

// TODO get note from UUID in GET

// TODO set note in POST

// helper functions

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return []byte{}, err
	}
	return []byte(buf.String()), nil
}
