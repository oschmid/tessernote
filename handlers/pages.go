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
	"collections/set"
	"html/template"
	"net/http"
	"notes"
	"regexp"
	"strings"
)

const PATH_LENGTH = len("/view/")

var Templates *template.Template
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	note, err := notes.LoadNote(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", note)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	note, err := notes.LoadNote(title)
	if err != nil {
		note = notes.NewNote(title, "", make(map[string]bool))
	}
	renderTemplate(w, "edit", note)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	tags := *set.New(strings.Split(r.FormValue("tags"), notes.TAG_SEPARATOR)...)
	note := *notes.NewNote(title, body, tags)
	err := notes.SaveNote(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// TODO list view of notes?
	http.Redirect(w, r, "/view/FrontNote", http.StatusFound)
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[PATH_LENGTH:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *notes.Note) {
	err := Templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}