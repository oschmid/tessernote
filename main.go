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
package main

import (
	"bytes"
	"handlers"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func main() {
	// data handlers
	addPostHandler(handlers.UrlGetTags, handlers.GetTagsHandler)
	addPostHandler(handlers.UrlRenameTags, handlers.RenameTagsHandler)
	addPostHandler(handlers.UrlDeleteTags, handlers.DeleteTagsHandler)
	addPostHandler(handlers.UrlGetTitles, handlers.TitlesHandler)
	addGetHandler(handlers.UrlGetNote, handlers.GetNoteHandler)
	addPostHandler(handlers.UrlSaveNote, handlers.SaveNoteHandler)

	// page handlers
	handlers.Templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
	addGetHandler("/view/", handlers.ViewHandler)
	addGetHandler("/edit/", handlers.EditHandler)
	addGetHandler("/save/", handlers.SaveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addGetHandler(url string, fn func(http.ResponseWriter, *http.Request, string)) {
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		get := r.URL.Path[len(url):]
		if !titleValidator.MatchString(get) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, get)
	}
	http.HandleFunc(url, getHandler)
}

func addPostHandler(url string, fn func(http.ResponseWriter, *http.Request, []byte)) {
	postHandler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		post := []byte(buf.String())
		fn(w, r, post)
	}
	http.HandleFunc(url, postHandler)
}
