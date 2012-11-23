package api

import (
	"appengine"
	"appengine/user"
	"grivet"
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

		notebook, err := grivet.GetNotebook(c)
		if err != nil {
			log.Println("notebook error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page := &Page{Tags: notebook.Tags()}
		tags, err := parseSelectedTags(w, r, notebook);
		if err != nil {
			log.Println("page error:", err)
			return
		}
		page.SelectedTags = tags
		page.RelatedTags = notebook.RelatedTags(tags)
		if len(tags) == 0 {
			notes, err := notebook.Notes()
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
			log.Println("template error:", err)
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
func parseSelectedTags(w http.ResponseWriter, r *http.Request, notebook *grivet.Notebook) ([]grivet.Tag, error) {
	var names []string
	if r.URL.Path != "/" {
		names = strings.Split(r.URL.Path[1:], tagSeparator)
	}
	tags, err := notebook.TagsFrom(names)
	if err != nil {
		names = namesFrom(tags)
		tagString := strings.Join(names, tagSeparator)
		http.Redirect(w, r, "/"+tagString, http.StatusFound)
	}
	return tags, err
}

func namesFrom(tags []grivet.Tag) []string {
	names := *new([]string)
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}
