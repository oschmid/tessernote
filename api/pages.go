package api

import (
	"appengine"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var (
	hex       = "[a-fA-F0-9]"
	tag       = "[a-zA-Z0-9]" // TODO use twitter regex
	tags      = "((" + tag + "+\\,)*" + tag + "+)"
	uniqueId  = "(" + hex + "{8}(-" + hex + "{4}){3}-" + hex + "{12})"
	validURL  = regexp.MustCompile("(^/$|" + // grivetnotes.appspot.com/
		"^/" + uniqueId + "$|" + // grivetnotes.appspot.com/<note UUID>
		"^/" + uniqueId + "/edit$|" + // grivetnotes.appspot.com/<note UUID>/edit
		"^/" + tags + "/$|" + // grivetnotes.appspot.com/<tags>/
		"^/" + tags + "/" + uniqueId + "$|" + // grivetnotes.appspot.com/<tags>/<note UUID>
		"^/" + tags + "/" + uniqueId + "/edit$)") //grivetnotes.appspot.com/<tags>/<note UUID>/edit
	idPattern = regexp.MustCompile(uniqueId)
	tagsPattern = regexp.MustCompile(tags)
	templates = template.Must(template.ParseFiles("templates/main.html"))
)

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	if validURL.MatchString(r.URL.Path) {
		appengine.NewContext(r) // TODO use
		matchedTags := strings.Split(tagsPattern.FindString(r.URL.Path), ",")
		matchedId := idPattern.FindString(r.URL.Path)
		page := Page{Tags:matchedTags,NoteBody:matchedId}
		err := templates.ExecuteTemplate(w, "main.html", page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}
