package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	hexPattern             = "[a-fA-F0-9]"
	tag                    = "[a-zA-Z0-9]" // TODO use twitter regex
	tags                   = "((" + tag + "+\\,)*" + tag + "+)"
	uniqueId               = "(" + hexPattern + "{8}(-" + hexPattern + "{4}){3}-" + hexPattern + "{12})"
	displayNoteURL         = regexp.MustCompile("^/" + uniqueId + "$")                   // grivet.ca/<note UUID>
	editNoteURL            = regexp.MustCompile("^/" + uniqueId + "/edit$")              // grivet.ca/<note UUID>/edit
	displayTagsURL         = regexp.MustCompile("^/" + tags + "/$")                      // grivet.ca/<tags>/
	displayNoteWithTagsURL = regexp.MustCompile("^/" + tags + "/" + uniqueId + "$")      // grivet.ca/<tags>/<note UUID>
	editNoteWithTagsURL    = regexp.MustCompile("^/" + tags + "/" + uniqueId + "/edit$") // grivet.ca/<tags>/<note UUID>/edit
)

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	if editNoteWithTagsURL.MatchString(r.URL.Path) {
		matches := editNoteWithTagsURL.FindStringSubmatch(r.URL.Path)
		editNoteWithTags(w, r, strings.Split(matches[1], ","), matches[3])
	} else if displayNoteWithTagsURL.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display note with some tags")
	} else if editNoteURL.MatchString(r.URL.Path) {
		fmt.Fprint(w, "edit note with all tags")
	} else if displayNoteURL.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display note with all tags")
	} else if displayTagsURL.MatchString(r.URL.Path) {
		fmt.Fprint(w, "display tags")
	} else {
		fmt.Fprint(w, r.URL.Path)
	}
}

func editNoteWithTags(w http.ResponseWriter, r *http.Request, tags []string, id string) {
	fmt.Fprintf(w, "tags:%v\nid:%v", tags, id)
}
