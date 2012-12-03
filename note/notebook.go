/*
This file is part of Tessernote.

Tessernote is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Tessernote is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Tessernote.  If not, see <http://www.gnu.org/licenses/>.
*/

package note

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
	tags     []Tag            // cache
	NoteKeys []*datastore.Key
	notes    []Note // cache
}

func (notebook Notebook) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Notebook", notebook.ID, 0, nil)
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
			log.Println("getMulti:notes", err)
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
		log.Println("tags:", err)
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
		log.Println("tags:", err)
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
		var err error
		note, err = notebook.addMissingTags(note, tc)
		if err != nil {
			log.Println("addMissingTags:", err)
			return err
		}

		note, err = notebook.createNote(note, tc)
		if err != nil {
			log.Println("createNote:", err)
			return err
		}

		err = note.addKeyToTags(tc)
		if err != nil {
			log.Println("addKeyToTags:", err)
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

// creates missing tags, adds their keys to notebook and note
// returns the note
func (notebook *Notebook) addMissingTags(note Note, c appengine.Context) (Note, error) {
	names := ParseTagNames(note.Body)
	if len(names) == 0 {
		return note, nil
	}

	notebookTags, err := notebook.Tags(c)
	if err != nil {
		log.Println("tags:", err)
		return note, err
	}

	// find missing tags, named tags
	notebookKey := notebook.Key(c)
	missingTagKeys := *new([]*datastore.Key)
	missingTags := *new([]Tag)
	for _, name := range names {
		i := indexOfTag(notebookTags, name)
		if i >= 0 {
			note.TagKeys = append(note.TagKeys, notebook.TagKeys[i])
			note.tags = append(note.tags, notebookTags[i])
		} else {
			tag := Tag{Name: name, NotebookKeys: []*datastore.Key{notebookKey}}
			missingTags = append(missingTags, tag)
			missingTagKeys = append(missingTagKeys, datastore.NewIncompleteKey(c, "Tag", nil))
		}
	}

	// create missing tags
	if len(missingTags) > 0 {
		missingTagKeys, err = datastore.PutMulti(c, missingTagKeys, missingTags)
		if err != nil {
			log.Println("putMulti:missingTags", err)
			return note, err
		}
	}

	// add missing tags to notebook
	notebook.TagKeys = append(notebook.TagKeys, missingTagKeys...)
	notebook.tags = append(notebook.tags, missingTags...)

	// add tags to note
	note.TagKeys = append(note.TagKeys, missingTagKeys...)
	note.tags = append(note.tags, missingTags...)
	return note, nil
}

// assumes "note" has keys for all its tags,
// creates new note, adds its key to notebook and
// returns new note
func (notebook *Notebook) createNote(note Note, c appengine.Context) (Note, error) {
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
	notebook.notes = append(notebook.notes, note)
	return note, nil
}

func (notebook *Notebook) updateNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		oldNote, err := GetNote(note.ID, tc)
		if err != nil {
			log.Println("getNote:", err)
			return err
		}

		note, err = notebook.removeNoteFromOldTags(*oldNote, note, tc)
		if err != nil {
			log.Println("removeNoteFromStaleTags:", err)
			return err
		}

		note, err = notebook.addMissingTags(note, tc)
		if err != nil {
			log.Println("addMissingTags:", err)
			return err
		}

		note, err := oldNote.Update(note, tc)
		if err != nil {
			log.Println("setBody:", err)
			return err
		}

		err = note.addKeyToTags(tc)
		if err != nil {
			log.Println("addKeyToTags:", err)
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

// removes note's key from tags it no longer has and
// removes empty tags from notebook
func (notebook *Notebook) removeNoteFromOldTags(oldNote, note Note, c appengine.Context) (Note, error) {
	oldTags, err := oldNote.Tags(c)
	if err != nil {
		log.Println("tags:", err)
		return note, err
	}

	// remove from old tags
	names := ParseTagNames(note.Body)
	tagsToUpdate := *new([]Tag)
	keysToUpdate := *new([]*datastore.Key)
	keysToRemove := *new([]*datastore.Key)
	for i := range oldTags {
		if len(oldTags[i].NoteKeys) == 1 {
			keysToRemove = append(keysToRemove, oldNote.TagKeys[i])
		} else if !containsString(names, oldTags[i].Name) {
			key, err := datastore.DecodeKey(note.ID)
			if err != nil {
				log.Println("decodeKey:note", err)
				return note, err
			}
			oldTags[i].NoteKeys = removeKey(oldTags[i].NoteKeys, key)
			tagsToUpdate = append(tagsToUpdate, oldTags[i])
			keysToUpdate = append(keysToUpdate, oldNote.TagKeys[i])
			note.TagKeys = append(note.TagKeys, oldNote.TagKeys[i])
			note.tags = append(note.tags, oldNote.tags[i])
		}
	}

	// update tags
	if len(keysToUpdate) > 0 {
		keysToUpdate, err = datastore.PutMulti(c, keysToUpdate, tagsToUpdate)
		if err != nil {
			log.Println("putMulti:update", err)
			return note, err
		}
	}

	// remove empty tags
	if len(keysToRemove) > 0 {
		err = datastore.DeleteMulti(c, keysToRemove)
		if err != nil {
			log.Println("deleteMulti:removeKeys", err)
			return note, err
		}
	}

	// update notebook tags
	cachedTags := len(notebook.tags) > 0
	tagKeys := *new([]*datastore.Key)
	tags := *new([]Tag)
	for i := range notebook.TagKeys {
		if !containsKey(keysToRemove, notebook.TagKeys[i]) {
			if containsKey(keysToUpdate, notebook.TagKeys[i]) && cachedTags {
				index := indexOfKey(keysToUpdate, notebook.TagKeys[i])
				notebook.tags[i] = tagsToUpdate[index]
			} else {
				tagKeys = append(tagKeys, notebook.TagKeys[i])
				if cachedTags {
					tags = append(tags, notebook.tags[i])
				}
			}
		}
	}
	notebook.TagKeys = tagKeys
	if cachedTags {
		notebook.tags = tags
	}
	return note, nil
}

func GetNotebook(c appengine.Context) (*Notebook, error) {
	u := user.Current(c)
	notebook := &Notebook{ID: u.ID}
	key := notebook.Key(c)
	err := datastore.Get(c, key, notebook)
	if err != nil {
		// create new user
		log.Println("new user", u.Email)
		notebook.Name = u.Email
		key, err = datastore.Put(c, key, notebook)
	}
	return notebook, err
}
