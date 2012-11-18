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
	tagPattern                  = "[a-zA-Z0-9]" // TODO use twitter regex
	tagSeparator                = ","
	tagsPattern                 = "((" + tagPattern + "+\\" + tagSeparator + ")*" + tagPattern + "+)"
	idPattern                   = "[a-zA-Z0-9-_]+"
	newSuffix                   = "/new"
	editSuffix                  = "/edit"
	saveSuffix                  = "/save"
	newNoteURL                  = regexp.MustCompile("^" + newSuffix)
	displayNoteURL              = regexp.MustCompile("^/" + idPattern + "$")
	editNoteURL                 = regexp.MustCompile("^/" + idPattern + editSuffix + "$")
	saveNoteURL                 = regexp.MustCompile("^/" + idPattern + saveSuffix + "$")
	selectTagsURL               = regexp.MustCompile("^/" + tagsPattern + "/$")
	selectTagsAndNewNoteURL     = regexp.MustCompile("^/" + tagsPattern + newSuffix + "$")
	selectTagsAndDisplayNoteURL = regexp.MustCompile("^/" + tagsPattern + "/" + idPattern + "$")
	selectTagsAndEditNoteURL    = regexp.MustCompile("^/" + tagsPattern + "/" + idPattern + editSuffix + "$")
	selectTagsAndSaveNoteURL    = regexp.MustCompile("^/" + tagsPattern + "/" + idPattern + saveSuffix + "$")
	templates                   = template.Must(template.ParseFiles("templates/main.html"))
)

func init() {
	http.HandleFunc("/", serve)
}

func serve(w http.ResponseWriter, r *http.Request) {
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

	g := grivet.GetUser(c)
	page := new(Page)
	page.Tags = g.Tags()

	if r.URL.Path == "/" {
		page.Notes = g.Notes()
	} else if r.URL.Path == saveSuffix {
		text := r.FormValue("note")
		note, err := g.NewNote(text)
		if err != nil {
			// TODO handle error
		}
		http.Redirect(w, r, "/"+note.Id.Encode(), http.StatusFound)
	} else if newNoteURL.MatchString(r.URL.Path) {
		page.Notes = g.Notes()
		page.Edit = true
	} else if displayNoteURL.MatchString(r.URL.Path) {
		page.Notes = g.Notes()
		urlSplit := strings.Split(r.URL.Path, "/")
		noteID := urlSplit[1]
		note, err := g.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		page.Note = note
	} else if editNoteURL.MatchString(r.URL.Path) {
		page.Notes = g.Notes()
		urlSplit := strings.Split(r.URL.Path, "/")
		noteID := urlSplit[1]
		note, err := g.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		page.Note = note
		page.Edit = true
	} else if saveNoteURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		noteID := urlSplit[1]
		// TODO save note
		http.Redirect(w, r, "/"+noteID, http.StatusFound)
	} else if selectTagsURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		tags, err := g.TagsFrom(strings.Split(tagString, tagSeparator))
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
		page.RelatedTags = g.RelatedTags(tags)
		page.SelectedTags = tags
		// TODO get related notes
	} else if selectTagsAndNewNoteURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		tags, err := g.TagsFrom(strings.Split(tagString, tagSeparator))
		if err != nil {
			if len(tags) > 0 {
				names := grivet.TagNames(tags)
				tagString := strings.Join(names, tagSeparator)
				http.Redirect(w, r, "/"+tagString+newSuffix, http.StatusFound)
			} else {
				http.Redirect(w, r, newSuffix, http.StatusFound)
			}
			return
		}
		page.RelatedTags = g.RelatedTags(tags)
		page.SelectedTags = tags
		// TODO get related notes
	} else if selectTagsAndDisplayNoteURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		noteID := urlSplit[2]
		tags, err := g.TagsFrom(strings.Split(tagString, tagSeparator))
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
		page.RelatedTags = g.RelatedTags(tags)
		page.SelectedTags = tags
		// TODO get related notes
		note, err := g.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/"+tagString+"/", http.StatusFound)
		}
		page.Note = note
	} else if selectTagsAndEditNoteURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		noteID := urlSplit[2]
		tags, err := g.TagsFrom(strings.Split(tagString, tagSeparator))
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
		page.RelatedTags = g.RelatedTags(tags)
		page.SelectedTags = tags
		// TODO get related notes
		note, err := g.Note(noteID)
		if err != nil {
			http.Redirect(w, r, "/"+tagString+"/", http.StatusFound)
			return
		}
		page.Note = note
		page.Edit = true
	} else if selectTagsAndSaveNoteURL.MatchString(r.URL.Path) {
		urlSplit := strings.Split(r.URL.Path, "/")
		tagString := urlSplit[1]
		noteID := urlSplit[2]
		tags, err := g.TagsFrom(strings.Split(tagString, tagSeparator))
		if err != nil {
			if len(tags) > 0 {
				names := grivet.TagNames(tags)
				tagString = strings.Join(names, tagSeparator)
				http.Redirect(w, r, "/"+tagString+"/"+noteID+saveSuffix, http.StatusFound)
			} else {
				http.Redirect(w, r, "/"+noteID+saveSuffix, http.StatusFound)
			}
			return
		}
		// TODO save note
		http.Redirect(w, r, "/"+tagString+"/"+noteID, http.StatusFound)
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
