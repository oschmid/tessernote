package api

import (
	"appengine"
	"appengine/user"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var (
	hexPattern  = "[a-fA-F0-9]"
	tagPattern  = "[a-zA-Z0-9]" // TODO use twitter regex
	tagsPattern = "((" + tagPattern + "+\\,)*" + tagPattern + "+)"
	uuidPattern = "(" + hexPattern + "{8}(-" + hexPattern + "{4}){3}-" + hexPattern + "{12})"
	validURL    = regexp.MustCompile("(^/$|" +
		"^/" + uuidPattern + "$|" +
		"^/" + uuidPattern + "/edit$|" +
		"^/" + tagsPattern + "/$|" +
		"^/" + tagsPattern + "/" + uuidPattern + "$|" +
		"^/" + tagsPattern + "/" + uuidPattern + "/edit$)")
	uuidRegexp = regexp.MustCompile(uuidPattern)
	tagsRegexp = regexp.MustCompile(tagsPattern)
	templates  = template.Must(template.ParseFiles("templates/main.html"))
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if validURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			// TODO store temp notebook based on IP address
		}
		matchedTags := strings.Split(tagsRegexp.FindString(r.URL.Path), ",")
		matchedId := uuidRegexp.FindString(r.URL.Path)
		page := Page{Tags: matchedTags, NoteBody: matchedId}
		err := templates.ExecuteTemplate(w, "main.html", page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}
