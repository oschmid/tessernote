package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	hex                    = "[a-fA-F0-9]"
	tag                    = "[a-zA-Z0-9]" // TODO use twitter regex
	tags                   = "((" + tag + "+\\,)*" + tag + "+)"
	uniqueId               = "(" + hex + "{8}(-" + hex + "{4}){3}-" + hex + "{12})"
	displayNoteURL         = regexp.MustCompile("^/" + uniqueId + "$")                   // grivet.ca/<note UUID>
	editNoteURL            = regexp.MustCompile("^/" + uniqueId + "/edit$")              // grivet.ca/<note UUID>/edit
	displayTagsURL         = regexp.MustCompile("^/" + tags + "/$")                      // grivet.ca/<tags>/
	displayNoteWithTagsURL = regexp.MustCompile("^/" + tags + "/" + uniqueId + "$")      // grivet.ca/<tags>/<note UUID>
	editNoteWithTagsURL    = regexp.MustCompile("^/" + tags + "/" + uniqueId + "/edit$") // grivet.ca/<tags>/<note UUID>/edit
	tagsPattern            = regexp.MustCompile(tags)
	uniqueIdPattern        = regexp.MustCompile(uniqueId)
)

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	matchedTags := strings.Split(tagsPattern.FindString(r.URL.Path), ",")
	matchedId := uniqueIdPattern.FindString(r.URL.Path)

	if r.URL.Path == "/" {
		displayHome(w, r)
	} else if displayTagsURL.MatchString(r.URL.Path) {
		displayTags(w, r, matchedTags)
	} else if displayNoteWithTagsURL.MatchString(r.URL.Path) {
		displayNoteWithTags(w, r, matchedTags, matchedId)
	} else if editNoteWithTagsURL.MatchString(r.URL.Path) {
		editNoteWithTags(w, r, matchedTags, matchedId)
	} else if displayNoteURL.MatchString(r.URL.Path) {
		displayNote(w, r, matchedId)
	} else if editNoteURL.MatchString(r.URL.Path) {
		editNote(w, r, matchedId)
	} else {
		http.NotFound(w, r)
	}
}

func editNoteWithTags(w http.ResponseWriter, r *http.Request, tags []string, id string) {
	fmt.Fprintf(w, "edit\ntags:%v\nid:%v", tags, id)
}

func displayNoteWithTags(w http.ResponseWriter, r *http.Request, tags []string, id string) {
	fmt.Fprintf(w, "display\ntags:%v\nid:%v", tags, id)
}

func editNote(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Fprintf(w, "edit\nid:%v", id)
}

func displayNote(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Fprintf(w, "display\nid:%v", id)
}

func displayTags(w http.ResponseWriter, r *http.Request, tags []string) {
	fmt.Fprintf(w, "display\ntags:%v", tags)
}

func displayHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home")
}
