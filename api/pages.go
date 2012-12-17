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
	"github.com/oschmid/tessernote"
	"github.com/oschmid/tessernote/filepath"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	tagPattern   = "[a-zA-Z0-9]" // TODO use twitter regex
	tagSeparator = ","
	tagsPattern  = "(" + tagPattern + "+\\" + tagSeparator + ")*" + tagPattern + "+"
	untaggedURL  = "/untagged/"
	validPageURL = regexp.MustCompile("^(/|(" + untaggedURL + ")|(/" + tagsPattern + "))$")
	templates    = getTemplates()
)

func init() {
	http.HandleFunc("/", serve)
}

func serve(w http.ResponseWriter, r *http.Request) {
	if validDataURL.MatchString(r.URL.Path) {
		serveData(w, r)
	} else if validPageURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		if !loggedIn(w, r, c) {
			return
		}

		notebook, err := tessernote.GetNotebook(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page := new(tessernote.Page)
		page.Tags, err = notebook.Tags(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page.SelectedTags, err = parseSelectedTags(w, r, notebook, c)
		if err != nil {
			return
		}

		page.RelatedTags, err = notebook.RelatedTags(page.SelectedTags, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page.UntaggedNotes = len(notebook.UntaggedNoteKeys) > 0
		if r.URL.Path == untaggedURL {
			page.Notes, err = notebook.UntaggedNotes(c)
		} else if len(page.SelectedTags) == 0 {
			page.Notes, err = notebook.Notes(c)
		} else {
			page.Notes, err = tessernote.RelatedNotes(page.SelectedTags, c)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "main.html", page)
		if err != nil {
			c.Errorf("executing template: %s", err)
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
			c.Errorf("logging in: %s", err)
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
func parseSelectedTags(w http.ResponseWriter, r *http.Request, notebook *tessernote.Notebook, c appengine.Context) ([]tessernote.Tag, error) {
	var names []string
	if r.URL.Path != "/" && r.URL.Path != untaggedURL {
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

func namesFrom(tags []tessernote.Tag) []string {
	names := *new([]string)
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}

func getTemplates() *template.Template {
	pwd, _ := os.Getwd()
	main := strings.Join([]string{"github.com", "oschmid", "tessernote", "api", "templates", "main.html"}, string(os.PathSeparator))
	main = filepath.Merge(pwd, main)
	return template.Must(template.ParseFiles(main))
}
