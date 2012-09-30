/**
 * Created with IntelliJ IDEA.
 * User: oschmid
 * Date: 28/09/12
 * Time: 10:01 PM
 */
package notes

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

const lenPath = len("/view/")
const tagSeparator = ", "

var templates = template.Must(template.ParseFiles("src/tmpl/edit.html", "src/tmpl/view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	note, err := loadNote(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", note)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	note, err := loadNote(title)
	if err != nil {
		note = &Note{Title: title}
	}
	renderTemplate(w, "edit", note)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	tags := strings.Split(r.FormValue("tags"), tagSeparator)
	note := Note{Title: title, Body: []byte(body), Tags: tags}
	err := saveNote(note)
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
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *Note) {
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
