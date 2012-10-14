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
	"html/template"
	"net/http"
	"notes"
	"string/collections/sets"
	"strings"
)

var Templates *template.Template
var notebook = *notes.NewNoteBook()

func ViewHandler(w http.ResponseWriter, r *http.Request, id string) {
	note, err := notebook.Note(id)
	if err != nil {
		http.Redirect(w, r, "/edit/"+id, http.StatusFound)
		return
	}
	renderTemplate(w, "view", note)
}

func EditHandler(w http.ResponseWriter, r *http.Request, id string) {
	note, err := notebook.Note(id)
	if err != nil {
		note = notes.NewNote("", "", *sets.New())
		note.Id = id
		notebook.Set(*note)
	}
	renderTemplate(w, "edit", note)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, id string) {
	titleAndBody := strings.SplitN(r.FormValue("title_body"), notes.TitleBodySeparator, 2)
	var title, body string
	if len(titleAndBody) == 0 {
		title, body = "Untitled", ""
	} else if len(titleAndBody) == 1 {
		title, body = titleAndBody[0], ""
	} else {
		title, body = titleAndBody[0], titleAndBody[1]
	}

	tags := *sets.New(strings.Split(r.FormValue("tags"), notes.TagSeparator)...)
	note := *notes.NewNote(title, body, tags)
	note.Id = id
	notebook.Set(note)
	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *notes.Note) {
	err := Templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
