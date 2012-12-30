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

// serveData handles requests to Tessernote's RESTful data API
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
		default:
			http.NotFound(w, r)
		}
	} else {
		switch r.Method {
		case "GET":
			GetNote(w, r, c, notebook)
		case "PUT":
			ReplaceNote(w, r, c, notebook)
		case "DELETE":
			DeleteNote(w, r, c, notebook)
		default:
			http.NotFound(w, r)
		}
	}
}

// GetAllNotes writes a JSON formatted list of all Note IDs in the authorized User's Notebook to w.
func GetAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	notes, err := notebook.Notes(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(notes)
	if err != nil {
		c.Errorf("marshaling notes (%d): %s", len(notes), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// ReplaceAllNotes replaces the Notes of the authorized User's Notebook with a new set of Notes. It takes as
// input a JSON formatted list of Notes and writes the added notes if succeeded or an error message otherwise to w.
// Notes may be written back with different IDs than those submitted, see ReplaceNote().
func ReplaceAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	notes, err := readNotes(w, r)
	if err != nil {
		return
	}
	err = notebook.DeleteAll(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	notes, err = notebook.PutAll(notes, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(notes)
	if err != nil {
		c.Errorf("marshaling notes (%d): %s", len(notes), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// readNotes decodes a JSON formatted []Note from the request body.
func readNotes(w http.ResponseWriter, r *http.Request) (notes []tessernote.Note, err error) {
	body, err := readRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return notes, err
	}
	err = json.Unmarshal(body, &notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return notes, err
	}
	return notes, nil
}

// CreateNote creates a new Note in the authorized User's Notebook. It takes as input a JSON formatted Note 
// and writes the new Note (with its automatically assigned unique ID) in JSON format to w.
func CreateNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	note, err := readNote(w, r)
	if err != nil {
		return
	}
	note, err = notebook.Put(note, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(note)
	if err != nil {
		c.Errorf("marshaling note (%#v): %s", note, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// readNote decodes a JSON formatted Note from the request body.
func readNote(w http.ResponseWriter, r *http.Request) (note tessernote.Note, err error) {
	body, err := readRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return note, err
	}
	err = json.Unmarshal(body, &note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return note, err
	}
	return note, nil
}

// readRequestBody reads the bytes from r's body
func readRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

// DeleteAllNotes deletes all Notes from the authorized User's Notebook. It writes true if Notes were deleted,
// and false if the Notebook was empty to w.
func DeleteAllNotes(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	empty := len(notebook.NoteKeys) == 0
	err := notebook.DeleteAll(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(!empty)
	if err != nil {
		c.Errorf("marshaling delete all response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// GetNote retrieves a note from the authorized User's Notebook by ID. The Note is written in JSON format to w.
func GetNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	id := r.URL.Path[len(NotesURL):]
	note, err := notebook.Note(id, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reply, err := json.Marshal(note)
	if err != nil {
		c.Errorf("marshaling note (%#v): %s", note, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// ReplaceNote replaces a Note in the authorized User's Notebook by its ID. If the Note doesn't exist it is created.
// If the Note's ID has already been assigned (e.g. in another Notebook) a new one is generated for this Note.
// The Note is written in JSON format to w.
func ReplaceNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	id := r.URL.Path[len(NotesURL):]
	note, err := readNote(w, r)
	if err != nil {
		return
	}
	if id != note.ID {
		http.Error(w, "mismatched note.ID and URL", http.StatusBadRequest)
		return
	}
	note, err = notebook.Put(note, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(note)
	if err != nil {
		c.Errorf("marshaling note (%#v): %s", note, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// DeleteNote deletes a Note by the ID in the URL. Uses w to write true if the Note was deleted, false
// if it never existed.
func DeleteNote(w http.ResponseWriter, r *http.Request, c appengine.Context, notebook *tessernote.Notebook) {
	id := r.URL.Path[len(NotesURL):]
	deleted, err := notebook.Delete(id, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := json.Marshal(deleted)
	if err != nil {
		c.Errorf("marshaling delete response (%t): %s", deleted, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
