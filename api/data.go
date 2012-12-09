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
	"api/context"
	"appengine"
	"appengine/user"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"note"
)

const (
	DeleteNoteURL = "/note/delete"
	SaveNoteURL   = "/note/save"
)

func isDataURL(url string) bool {
	return url == DeleteNoteURL || url == SaveNoteURL
}

func serveData(w http.ResponseWriter, r *http.Request) {
	c := context.NewContext(r)
	defer context.Close()
	u := user.Current(c)
	if u == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	notebook, err := note.GetNotebook(c)
	if err != nil {
		log.Println("getnotebook:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := readPost(r)
	if err != nil {
		log.Println("readPost:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case DeleteNoteURL:
		deleteNote(w, body, notebook, c)
	case SaveNoteURL:
		saveNote(w, body, notebook, c)
	}
}

func readPost(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

// Reads a JSON formatted Note in from POST and deletes the note with that ID from the datastore.
// Returns true if note was deleted, false otherwise
func deleteNote(w http.ResponseWriter, body []byte, notebook *note.Notebook, c appengine.Context) {
	var note note.Note
	err := json.Unmarshal(body, &note)
	if err != nil {
		log.Println("delete:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleted, err := notebook.Delete(note, c)
	if err != nil {
		log.Println("delete:note", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply, err := json.Marshal(deleted)
	if err != nil {
		log.Println("marshal:", err, deleted)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// Reads a JSON formatted Note in from POST and writes it to the datastore.
// Returns the new or updated Note in JSON format.
func saveNote(w http.ResponseWriter, body []byte, notebook *note.Notebook, c appengine.Context) {
	var note note.Note
	err := json.Unmarshal(body, &note)
	if err != nil {
		log.Println("save:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note, err = notebook.Put(note, c)
	if err != nil {
		log.Println("put:note", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply, err := json.Marshal(note)
	if err != nil {
		log.Println("marshal:", err, note)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
