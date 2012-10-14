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
)

const (
	UrlGetTags    = "/tags/get"
	UrlDeleteTags = "/tags/delete"
	UrlRenameTags = "/tags/rename"
	UrlGetTitles  = "/titles"
	UrlGetNote    = "/note/get/"
	UrlSaveNote   = "/note/save"
)

// TODO get untagged notes?

// Returns a map of tags -> note count in JSON format.
// Request can optionally specify a list of tags in JSON format in POST.
func GetTagsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	// convert JSON -> []string
	var tagsIn []string
	err := json.Unmarshal(body, &tagsIn)
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
		return
	}
}

// Renames tags.
// Request specifies oldTags -> newTags in JSON format in POST.
// OldTags that don't exist are skipped. NewTags that already exist will create a union.
func RenameTagsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	// convert JSON -> map[string]string
	var tags map[string]string
	err := json.Unmarshal(body, &tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// rename tags
	for old, new := range tags {
		notebook.RenameTag(old, new)
	}

	// write response
	w.WriteHeader(http.StatusOK)
}

func DeleteTagsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	// convert JSON -> []string
	var tags []string
	err := json.Unmarshal(body, &tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// delete tags
	for _, tag := range tags {
		notebook.DeleteTag(tag)
	}

	// write response
	w.WriteHeader(http.StatusOK)
}

// Returns a slice of titles in JSON format.
// Request can optionally specify a list of tags in JSON format in POST.
func TitlesHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	// convert JSON -> []string
	var tags []string
	err := json.Unmarshal(body, &tags)
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

// Returns Note in JSON format.
// Request must specify Note.Id in URL
func GetNoteHandler(w http.ResponseWriter, r *http.Request, id string) {
	note, err := notebook.Note(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert Note -> JSON
	response, err := json.Marshal(note)
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

// Saves Note
// Request must specify Note in JSON format in POST
func SaveNoteHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	// convert JSON -> Note
	var note notes.Note
	err := json.Unmarshal(body, &note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// set note
	notebook.Set(note)
	w.WriteHeader(http.StatusOK)
}
