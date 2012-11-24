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
	"appengine"
	"appengine/user"
	"bytes"
	"encoding/json"
	"grivet"
	"io"
	"log"
	"net/http"
)

const (
	SaveNoteURL = "/note/save"
)

func isDataURL(url string) bool {
	return url == SaveNoteURL
}

func serveData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	notebook, err := grivet.GetNotebook(c)
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
	case SaveNoteURL:
		saveNote(w, body, notebook)
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

// Reads a JSON formatted Note in from POST and writes it to the datastore.
// Returns the new or updated Note in JSON format.
func saveNote(w http.ResponseWriter, body []byte, notebook *grivet.Notebook) {
	var note grivet.Note
	err := json.Unmarshal(body, &note)
	if err != nil {
		log.Println("save:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if note.ID == "" {
		// save new note
		n, err := notebook.NewNote(note.Body)
		if err != nil {
			log.Println("newnote:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		note = *n
	} else {
		// update existing note
		n, err := notebook.Note(note.ID)
		if err != nil {
			log.Println("note:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = n.SetBody(note.Body)
		if err != nil {
			log.Println("setbody:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		note = n
	}

	bytes, err := json.Marshal(note)
	if err != nil {
		log.Println("marshal:", err, note)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// TODO remove old api
/*

// Returns a map of tags -> note count in JSON format.
// Request can optionally specify a list of tags in JSON format in POST.
func GetTags(w http.ResponseWriter, body []byte) {
	// convert JSON -> []string
	var tagsIn []string
	err := json.Unmarshal(body, &tagsIn)
	if err != nil {
		log.Println("GetTags", err.Error(), string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert map[string]bool -> JSON
	tagsOut := *notebook.RelatedTags(tagsIn...)
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
func RenameTags(w http.ResponseWriter, body []byte) {
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

func DeleteTags(w http.ResponseWriter, body []byte) {
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
func GetTitles(w http.ResponseWriter, body []byte) {
	// convert JSON -> []string
	var tags []string
	err := json.Unmarshal(body, &tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert [][]string -> JSON
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
func GetNote(w http.ResponseWriter, id string) {
	note, contained := notebook.Notes[id]
	if !contained {
		http.Error(w, "note with id="+id+" exists in notebook", http.StatusBadRequest)
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
// New Notes (with empty Ids) will be given a UUID.
// Response consists of the Note's UUID
func SaveNote(w http.ResponseWriter, body []byte) {
	// convert JSON -> Note
	var note notes.Note
	err := json.Unmarshal(body, &note)
	if err != nil {
		log.Println("SaveNote cannot unmarshal:", err, "\n", string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add UUID to new Notes
	if note.Id == "" {
		note.Id = uuid.New()
	}
	notebook.SetNote(note)

	// write response
	w.Write([]byte(note.Id))
}

func DeleteNote(w http.ResponseWriter, id string) {
	notebook.Delete(id)
	w.WriteHeader(http.StatusOK)
}
*/
