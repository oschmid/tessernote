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

// returns a user's note by ID
func (notebook *Notebook) Note(id string, c appengine.Context) (Note, error) {
	// TODO binary search
	notes, err := notebook.Notes(c)
	if err != nil {
		return *new(Note), err
	}

	for _, note := range notes {
		if note.ID == id {
			return note, nil
		}
	}
	return *new(Note), errors.New("note does not exist")
}

func (notebook *Notebook) SetNote(note Note, c appengine.Context) (Note, error) {
	if note.ID == "" {
		return notebook.addNewNote(note, c)
	}
	return notebook.updateNote(note, c)
}

func (notebook *Notebook) addNewNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		notebookKey := datastore.NewKey(tc, "Notebook", notebook.ID, 0, nil)

		// add note
		key := datastore.NewIncompleteKey(tc, "Note", nil)
		note.Created = time.Now()
		note.LastModified = note.Created
		note.NotebookKeys = []*datastore.Key{notebookKey}
		key, err := datastore.Put(tc, key, &note)
		if err != nil {
			return err
		}
		note.ID = key.Encode()

		// TODO add/update tags
		//names := note.ParseTagNames()

		// update notebook
		notebook.NoteKeys = append(notebook.NoteKeys, key)
		_, err = datastore.Put(tc, notebookKey, &notebook)
		return err
	}, nil)
	return note, err
}

func (notebook *Notebook) updateNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		// update note
		key, err := datastore.DecodeKey(note.ID)
		if err != nil {
			return err
		}

		existing := new(Note)
		err = datastore.Get(tc, key, existing)
		if err != nil {
			return err
		}

		// TODO remove from previous tags

		existing.SetBody(note.Body)
		note = *existing
		_, err = datastore.Put(tc, key, &note)
		if err != nil {
			return err
		}

		// TODO add/update/remove tags
		// TODO update notebook
		return nil
	}, nil)
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
