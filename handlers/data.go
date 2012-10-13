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
)

const GET_TAGS = "/tags/get"

//const DELETE_TAGS = "/tags/delete"
//const RENAME_TAGS = "/tags/rename"
const TITLES = "/titles"

//const GET_NOTE = "/note/get"
//const SAVE_NOTE = "/note/save"

// Returns a map of tags -> note count in JSON format.
// Request can optionally specify a list of tags in JSON format in POST.
func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	// read request
	body, err := readBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert JSON -> []string
	var tagsIn []string
	err = json.Unmarshal(body, &tagsIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert map[string]bool -> JSON
	tagsOut := *notebook.Tags(tagsIn...)
	response, err := json.Marshal(tagsOut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Returns a slice of titles in JSON format.
// Request can optionally specify a list of tags in JSON format in POST.
func TitlesHandler(w http.ResponseWriter, r *http.Request) {
	// read request
	body, err := readBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert JSON -> []string
	var tags []string
	err = json.Unmarshal(body, &tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert []string -> JSON
	titles := notebook.Titles(tags...)
	response, err := json.Marshal(titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
