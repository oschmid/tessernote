package api

import (
	"appengine"
	"grivet"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	tagPattern                  = "[a-zA-Z0-9]" // TODO use twitter regex
	tagSeparator                = ","
	tagsPattern                 = "((" + tagPattern + "+\\" + tagSeparator + ")*" + tagPattern + "+)"
	idPattern                   = "[a-zA-Z0-9-_]+"
	editSuffix                  = "/edit"
	displayNoteURL              = regexp.MustCompile("^/" + idPattern + "$")
	editNoteURL                 = regexp.MustCompile("^/" + idPattern + editSuffix + "$")
	selectTagsURL               = regexp.MustCompile("^/" + tagsPattern + "/$")
	selectTagsAndDisplayNoteURL = regexp.MustCompile("^/" + tagsPattern + "/" + idPattern + "$")
	selectTagsAndEditNoteURL    = regexp.MustCompile("^/" + tagsPattern + "/" + idPattern + editSuffix + "$")
	templates                   = template.Must(template.ParseFiles("templates/main.html"))
)

func init() {
	http.HandleFunc("/", serve)
}

func serve(w http.ResponseWriter, r *http.Request) {
	page := new(Page)
	if r.URL.Path == "/" {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		page.Notes = u.Notes()
	} else if displayNoteURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		page.Notes = u.Notes()
		urlSplit := strings.Split(r.URL.Path, "/")
		noteID := urlSplit[1]
		note, err := u.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		page.Note = note
	} else if editNoteURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		page.Notes = u.Notes()
		urlSplit := strings.Split(r.URL.Path, "/")
		noteID := urlSplit[1]
		note, err := u.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		page.Note = note
		page.Edit = true
	} else if selectTagsURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		tags, err := u.TagsFrom(strings.Split(tagString, tagSeparator))
		if err != nil {
			if len(tags) > 0 {
				names := grivet.TagNames(tags)
				tagString := strings.Join(names, tagSeparator)
				http.Redirect(w, r, "/"+tagString+"/", http.StatusFound)
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
			}
			return
		}
		page.RelatedTags = tags
	} else if selectTagsAndDisplayNoteURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		noteID := urlSplit[2]
		tags, err := u.TagsFrom(strings.Split(tagString, tagSeparator))
		if err != nil {
			if len(tags) > 0 {
				names := grivet.TagNames(tags)
				tagString = strings.Join(names, tagSeparator)
				http.Redirect(w, r, "/"+tagString+"/"+noteID, http.StatusFound)
			} else {
				http.Redirect(w, r, "/"+noteID, http.StatusFound)
			}
			return
		}
		page.RelatedTags = tags
		note, err := u.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/"+tagString+"/", http.StatusFound)
		}
		page.Note = note
	} else if selectTagsAndEditNoteURL.MatchString(r.URL.Path) {
		c := appengine.NewContext(r)
		u := grivet.CurrentUser(c)
		page.Tags = u.Tags()
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		noteID := urlSplit[2]
		tags, err := u.TagsFrom(strings.Split(tagString, tagSeparator))
		if err != nil {
			if len(tags) > 0 {
				names := grivet.TagNames(tags)
				tagString = strings.Join(names, tagSeparator)
				http.Redirect(w, r, "/"+tagString+"/"+noteID+editSuffix, http.StatusFound)
			} else {
				http.Redirect(w, r, "/"+noteID+editSuffix, http.StatusFound)
			}
			return
		}
		page.RelatedTags = tags
		note, err := u.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/"+tagString+"/", http.StatusFound)
			return
		}
		page.Note = note
		page.Edit = true
	} else {
		http.NotFound(w, r)
		return
	}

	err := templates.ExecuteTemplate(w, "main.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error executing template", err, page)
	}
}
