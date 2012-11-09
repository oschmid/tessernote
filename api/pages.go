package api

import (
	"fmt"
	"net/http"
	"regexp"
)

var (
	hex                     = "[a-fA-F0-9]"
	tag                     = "[a-zA-Z0-9]" // TODO use twitter regex
	tags                    = "(" + tag + "+\\,)*" + tag + "+"
	uniqueId                = hex + "{8}(-" + hex + "{4}){3}-" + hex + "{12}"
	displayNoteWithAllTags  = regexp.MustCompile("^/" + uniqueId + "$")                   // grivet.ca/<note UUID>
	editNoteWithAllTags     = regexp.MustCompile("^/" + uniqueId + "/edit$")              // grivet.ca/<note UUID>?edit
	displayTags             = regexp.MustCompile("^/" + tags + "/$")                      // grivet.ca/<tags>/
	displayNoteWithSomeTags = regexp.MustCompile("^/" + tags + "/" + uniqueId + "$")      // grivet.ca/<tags>/<note UUID
	editNoteWithSomeTags    = regexp.MustCompile("^/" + tags + "/" + uniqueId + "/edit$") // grivet.ca/<tags>/<note UUID
)

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	if editNoteWithSomeTags.MatchString(r.URL.Path) {
		fmt.Fprintf(w, "edit note with some tags")
	} else if displayNoteWithSomeTags.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display note with some tags")
	} else if editNoteWithAllTags.MatchString(r.URL.Path) {
		fmt.Fprint(w, "edit note with all tags")
	} else if displayNoteWithAllTags.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display note with all tags")
	} else if displayTags.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display tags")
	} else {
		fmt.Fprint(w, r.URL.Path)
	}
}
