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

// TODO serve only: / and /<tags>

func serve(w http.ResponseWriter, r *http.Request) {
	if validURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}

		g, err := grivet.GetUser(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page := &Page{Tags: g.Tags()}

		var names []string
		if r.URL.Path != "/" {
			names = strings.Split(r.URL.Path[1:], tagSeparator)
		}
		tags, err := g.TagsFrom(names)
		if err != nil {
			log.Println("length:", len(names))
			log.Println(err)
			names = namesFrom(tags)
			tagString := strings.Join(names, tagSeparator)
			http.Redirect(w, r, "/"+tagString, http.StatusFound)
			return
		}
		page.RelatedTags = g.RelatedTags(tags)
		page.SelectedTags = tags
		// TODO get related notes
		err = templates.ExecuteTemplate(w, "main.html", page)
		if err != nil {
			log.Println("template error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
		return
	}
}

func namesFrom(tags []grivet.Tag) []string {
	names := *new([]string)
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}
