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
	"note"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	tagPattern   = "[a-zA-Z0-9]" // TODO use twitter regex
	tagSeparator = ","
	validURL     = regexp.MustCompile("^/((" + tagPattern + "+\\" + tagSeparator + ")*" + tagPattern + "+)*$")
	templates    = template.Must(template.ParseFiles("templates/main.html"))
)

func init() {
	http.HandleFunc("/", serve)
}

func serve(w http.ResponseWriter, r *http.Request) {
	if isDataURL(r.URL.Path) {
		serveData(w, r)
	} else if validURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		if !loggedIn(w, r, c) {
			return
		}

		notebook, err := note.GetNotebook(c)
		if err != nil {
			log.Println("getNotebook:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tags, err := notebook.Tags(c)
		if err != nil {
			log.Println("tags:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page := &Page{Tags: tags}
		tags, err = parseSelectedTags(w, r, notebook, c)
		if err != nil {
			log.Println("page:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page.SelectedTags = tags
		page.RelatedTags, err = notebook.RelatedTags(tags, c)
		if err != nil {
			log.Println("relatedTags:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if len(tags) == 0 {
			notes, err := notebook.Notes(c)
			if err != nil {
				log.Println("notes:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			page.Notes = notes
		} else {
			// TODO get related notes
		}

		err = templates.ExecuteTemplate(w, "main.html", page)
		if err != nil {
			log.Println("executeTemplate:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

// checks if user is logged in, redirects user to login if they aren't
func loggedIn(w http.ResponseWriter, r *http.Request, c appengine.Context) bool {
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return false
	}
	return true
}

// parses url for selected tags, redirects if it refers to missing tags
func parseSelectedTags(w http.ResponseWriter, r *http.Request, notebook *note.Notebook, c appengine.Context) ([]note.Tag, error) {
	var names []string
	if r.URL.Path != "/" {
		names = strings.Split(r.URL.Path[1:], tagSeparator)
	}
	tags, err := notebook.TagsFrom(names, c)
	if err != nil {
		names = namesFrom(tags)
		tagString := strings.Join(names, tagSeparator)
		http.Redirect(w, r, "/"+tagString, http.StatusFound)
	}
	return tags, err
}

func namesFrom(tags []note.Tag) []string {
	names := *new([]string)
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}
