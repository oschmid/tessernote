/*
This file is part of Grivet.

Grivet is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Grivet is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Grivet.  If not, see <http://www.gnu.org/licenses/>.
*/
package grivet

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"errors"
	"log"
	"sort"
	"time"
)

type Notebook struct {
	ID       string // user.User.ID
	Name     string
	TagKeys  []*datastore.Key // sorted by Tag.Name
	NoteKeys []*datastore.Key
	tags     []Tag  // cache
	notes    []Note // cache
}

func (notebook *Notebook) Tags(c appengine.Context) ([]Tag, error) {
	if len(notebook.tags) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.tags = make([]Tag, len(notebook.TagKeys))
		err := datastore.GetMulti(c, notebook.TagKeys, notebook.tags)
		if err != nil {
			log.Println("tags:", err)
			return notebook.tags, err
		}
	}
	return notebook.tags, nil
}

func (notebook *Notebook) Notes(c appengine.Context) ([]Note, error) {
	if len(notebook.notes) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.notes = make([]Note, len(notebook.NoteKeys))
		err := datastore.GetMulti(c, notebook.NoteKeys, notebook.notes)
		if err != nil {
			return notebook.notes, err
		}
		for i := 0; i < len(notebook.notes); i++ {
			note := &notebook.notes[i]
			note.ID = notebook.NoteKeys[i].Encode()
		}
	}
	return notebook.notes, nil
}

// returns a subset of a user's tags by name
// missing tags result in errors
func (notebook *Notebook) TagsFrom(names []string, c appengine.Context) ([]Tag, error) {
	tags := *new([]Tag)
	if len(names) == 0 {
		return tags, nil
	}
	sort.Strings(names)
	allTags, err := notebook.Tags(c)
	if err != nil {
		return tags, err
	}
	tagsIndex := 0
	namesIndex := 0
	for tagsIndex < len(allTags) && namesIndex < len(names) {
		tag := allTags[tagsIndex]
		name := names[namesIndex]
		if tag.Name == name {
			tags = append(tags, tag)
			namesIndex++
		} else if tag.Name > name {
			namesIndex++
		}
		tagsIndex++
	}
	if len(tags) != len(names) {
		return tags, errors.New("user missing tag(s)")
	}
	return tags, nil
}

func (notebook *Notebook) RelatedTags(tags []Tag, c appengine.Context) ([]Tag, error) {
	relatedNoteKeys := make(map[string]datastore.Key)
	for _, tag := range tags {
		for _, key := range tag.NoteKeys {
			relatedNoteKeys[key.Encode()] = *key
		}
	}
	tags = *new([]Tag)
	allTags, err := notebook.Tags(c)
	if err != nil {
		return tags, err
	}
	for _, tag := range allTags {
		for _, key := range tag.NoteKeys {
			if _, contained := relatedNoteKeys[key.Encode()]; contained {
				tags = append(tags, tag)
				break
			}
		}
	}
	return tags, nil
}

func (notebook *Notebook) SetNote(note Note, c appengine.Context) (Note, error) {
	if note.ID == "" {
		return notebook.addNote(note, c)
	}
	return notebook.updateNote(note, c)
}

func (notebook *Notebook) addNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		names := note.ParseTagNames()
		var err error
		note.TagKeys, err = notebook.addMissingTags(names, tc)
		if err != nil {
			return err
		}

		note, err = notebook.addNewNote(note, tc)
		if err != nil {
			return err
		}

		err = note.addKeyToTags(tc)
		if err != nil {
			return err
		}

		key := notebook.Key(tc)
		_, err = datastore.Put(tc, key, notebook)
		if err != nil {
			log.Println("put:notebook", err)
		}
		return err
	}, &datastore.TransactionOptions{XG: true})
	return note, err
}

func (notebook Notebook) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Notebook", notebook.ID, 0, nil)
}

// creates missing tags, adds their keys to notebook
// returns keys for all tags in "names"
func (notebook *Notebook) addMissingTags(names []string, c appengine.Context) ([]*datastore.Key, error) {
	namedKeys := *new([]*datastore.Key)
	if len(names) == 0 {
		return namedKeys, nil
	}

	notebookTags, err := notebook.Tags(c)
	if err != nil {
		return namedKeys, err
	}

	// find missing tags, named tags
	notebookKey := notebook.Key(c)
	missingTagKeys := *new([]*datastore.Key)
	missingTags := *new([]Tag)
	for _, name := range names {
		missing := false
		for i, tag := range notebookTags {
			if tag.Name == name {
				missing = false
				namedKeys = append(namedKeys, notebook.TagKeys[i])
			}
		}
		if missing {
			tag := Tag{Name: name, NotebookKeys: []*datastore.Key{notebookKey}}
			missingTags = append(missingTags, tag)
			missingTagKeys = append(missingTagKeys, datastore.NewIncompleteKey(c, "Tag", nil))
		}
	}

	// create missing tags
	missingTagKeys, err = datastore.PutMulti(c, missingTagKeys, missingTags)
	if err != nil {
		log.Println("putMulti:missingTags", err)
		return namedKeys, err
	}

	// add missing tag keys to notebook
	notebook.TagKeys = append(notebook.TagKeys, missingTagKeys...)

	// return keys for all tags in "names"
	return namedKeys, nil
}

// assumes "note" has keys for all its tags
// creates new note, adds its key to notebook
// returns new note's key
func (notebook *Notebook) addNewNote(note Note, c appengine.Context) (Note, error) {
	noteKey := datastore.NewIncompleteKey(c, "Note", nil)
	note.Created = time.Now()
	note.LastModified = note.Created
	note.NotebookKeys = []*datastore.Key{notebook.Key(c)}
	noteKey, err := datastore.Put(c, noteKey, &note)
	if err != nil {
		log.Println("put:note", err)
		return note, err
	}
	note.ID = noteKey.Encode()
	notebook.NoteKeys = append(notebook.NoteKeys, noteKey)
	return note, nil
}

func (notebook *Notebook) updateNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		// update note
		key, err := datastore.DecodeKey(note.ID)
		if err != nil {
			log.Println("decodeKey:", err)
			return err
		}

		existing := new(Note)
		err = datastore.Get(tc, key, existing)
		if err != nil {
			log.Println("get:", err)
			return err
		}

		// TODO remove from previous tags

		existing.SetBody(note.Body)
		note = *existing
		_, err = datastore.Put(tc, key, &note)
		if err != nil {
			log.Println("put:", err)
			return err
		}

		// TODO add/update tags
		// TODO update notebook
		return nil
	}, &datastore.TransactionOptions{XG: true})
	return note, err
}

func GetNotebook(c appengine.Context) (*Notebook, error) {
	u := user.Current(c)
	notebook := &Notebook{ID: u.ID}
	key := datastore.NewKey(c, "Notebook", notebook.ID, 0, nil)
	err := datastore.Get(c, key, notebook)
	if err != nil {
		// store new user
		log.Println("new user", u.Email)
		notebook.Name = u.Email
		key, err = datastore.Put(c, key, notebook)
	}
	return notebook, err
}
