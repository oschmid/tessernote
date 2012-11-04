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
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"notes"
	"os"
	"regexp"
)

var notebook = *notes.NewNoteBook()
var notebookFileName = "../data/notebook"
var uuidValidator = regexp.MustCompile("[a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}")

func MakePostHandler(url string, fn func(http.ResponseWriter, []byte)) (string, func(http.ResponseWriter, *http.Request)) {
	return url, func(w http.ResponseWriter, r *http.Request) {
		post, err := readPost(r)
		if err != nil {
			log.Println("error reading post.", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fn(w, post)
		SaveNotebook()
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

func MakeGetHandler(url string, fn func(http.ResponseWriter, string)) (string, func(http.ResponseWriter, *http.Request)) {
	return url, func(w http.ResponseWriter, r *http.Request) {
		get := r.URL.Path[len(url):]
		if !uuidValidator.MatchString(get) {
			log.Println("id not uuid", get)
			http.NotFound(w, r)
			return
		}
		fn(w, get)
	}
}

func PopulateNotebook() {
	_, err := os.Create(notebookFileName)
	if err != nil {
		log.Fatal(err)
	}

	notebook.SetNote(*notes.NewNote("title1", "body1\n #tag1 #tag2 #tag3"))
	notebook.SetNote(*notes.NewNote("title2", "body2\n #tag1 #tag3 #tag4"))
	notebook.SetNote(*notes.NewNote("title3", "body3\n #tag5"))
}

func LoadNotebook() {
	file, err := os.Open(notebookFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	err = gob.NewDecoder(file).Decode(&notebook)
	if err != nil {
		log.Println(err)
	}
}

func SaveNotebook() {
	file, err := os.OpenFile(notebookFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()
	err = gob.NewEncoder(file).Encode(notebook)
	if err != nil {
		log.Println(err)
	}
}
