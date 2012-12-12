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
	"appengine"
	"appengine/user"
	"bytes"
	"encoding/json"
	"github.com/oschmid/tessernote"
	"io"
	"net/http"
	"regexp"
)

const (
	NotesURL = "/notes/"
)

var (
	idPattern    = "[0-9a-zA-Z-_]"
	validDataURL = regexp.MustCompile("^" + NotesURL + idPattern + "*$")
)

func serveData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	notebook, err := tessernote.GetNotebook(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path == NotesURL {
		switch r.Method {
		case "GET":
			GetAllNotes(w, r, c, notebook)
		case "PUT":
			ReplaceAllNotes(w, r, c, notebook)
		case "POST":
			CreateNote(w, r, c, notebook)
		case "DELETE":
			DeleteAllNotes(w, r, c, notebook)
		}
	} else {
		switch r.Method {
		case "GET":
			GetNote(w, r, c, notebook)
		case "PUT":
			ReplaceNote(w, r, c, notebook)
		case "DELETE":
			DeleteNote(w, r, c, notebook)
		}
	}
}

// Uses "w" to write a JSON formatted list of all Note IDs in "notebook"
func GetAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	// TODO return a list note IDs
	http.Error(w, "not yet implemented", http.StatusInternalServerError)
}

// Replaces the contents of "notebook" with the JSON formatted list of notes in "r"
// Uses "w" to write "true" if succeeded or an error message otherwise
func ReplaceAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	// TODO replace all notes with new notes
	http.Error(w, "not yet implemented", http.StatusInternalServerError)
}

// Creates a new note in "notebook" with the contents of the JSON formatted note in "r"
// Uses "w" to write the new note (with its automatically assigned ID) in JSON format
func CreateNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	body, err := readRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var note tessernote.Note
	err = json.Unmarshal(body, &note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note, err = notebook.Put(note, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply, err := json.Marshal(note)
	if err != nil {
		c.Errorf("marshaling note:", err, "\n", note)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// Deletes all notes in "notebook"
// Uses "w" to write "true" if notes were deleted, "false" if notebook was empty
func DeleteAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	// TODO delete all notes
	http.Error(w, "not yet implemented", http.StatusInternalServerError)
}

// Retrieves a note from "notebook" by its ID
// Uses "w" to write the note in JSON format
func GetNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	// TODO retrieve a note by ID
	http.Error(w, "not yet implemented", http.StatusInternalServerError)
}

// Replaces a note from "notebook" by its ID, creates it if it doesn't exist
// Uses "w" to write the note in JSON format
func ReplaceNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	// TODO replace note by ID, create if it doesn't exist
	http.Error(w, "not yet implemented", http.StatusInternalServerError)
}

// Reads a Note.ID from URL and deletes it from the Notebook
// Uses "w" to write "true" if note was deleted, "false" if it never existed
func DeleteNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	id := r.URL.Path[len(NotesURL):]
	deleted, err := notebook.Delete(id, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply, err := json.Marshal(deleted)
	if err != nil {
		c.Errorf("marshaling delete response:", err, "\n", deleted)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

func readRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}
