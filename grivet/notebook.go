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
	tags     []Tag             // cache
	notes    []Note            // cache
	context  appengine.Context `datastore:",noindex"`
}

func (notebook *Notebook) Tags() ([]Tag, error) {
	if len(notebook.tags) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.tags = make([]Tag, len(notebook.TagKeys))
		err := datastore.GetMulti(notebook.context, notebook.TagKeys, notebook.tags)
		if err != nil {
			return notebook.tags, err
		}
		for i := 0; i < len(notebook.tags); i++ {
			tag := &notebook.tags[i]
			tag.context = notebook.context
		}
	}
	return notebook.tags, nil
}

func (notebook *Notebook) Notes() ([]Note, error) {
	if len(notebook.notes) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.notes = make([]Note, len(notebook.NoteKeys))
		err := datastore.GetMulti(notebook.context, notebook.NoteKeys, notebook.notes)
		if err != nil {
			return notebook.notes, err
		}
		for i := 0; i < len(notebook.notes); i++ {
			note := &notebook.notes[i]
			note.ID = notebook.NoteKeys[i].Encode()
			note.context = notebook.context
		}
	}
	return notebook.notes, nil
}

// returns a subset of a user's tags by name
// missing tags result in errors
func (notebook *Notebook) TagsFrom(names []string) ([]Tag, error) {
	tags := *new([]Tag)
	if len(names) == 0 {
		return tags, nil
	}
	sort.Strings(names)
	allTags, err := notebook.Tags()
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

func (notebook *Notebook) RelatedTags(tags []Tag) ([]Tag, error) {
	relatedNoteKeys := make(map[string]datastore.Key)
	for _, tag := range tags {
		for _, key := range tag.NoteKeys {
			relatedNoteKeys[key.Encode()] = *key
		}
	}
	tags = *new([]Tag)
	allTags, err := notebook.Tags()
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
func (notebook *Notebook) Note(id string) (Note, error) {
	// TODO binary search
	notes, err := notebook.Notes()
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

func (notebook *Notebook) NewNote(body string) (*Note, error) {
	key := datastore.NewIncompleteKey(notebook.context, "Note", nil)
	note := &Note{
		Body:         body,
		Created:      time.Now(),
		LastModified: time.Now(),
		NotebookKeys: []*datastore.Key{notebook.Key()}}
	key, err := datastore.Put(notebook.context, key, note)
	if err != nil {
		return note, err
	}

	note.ID = key.Encode()
	note.context = notebook.context

	notebook.NoteKeys = append(notebook.NoteKeys, key)
	return note, notebook.Save()
}

func (notebook Notebook) Key() *datastore.Key {
	return datastore.NewKey(notebook.context, "Notebook", notebook.ID, 0, nil)
}

func (notebook Notebook) Save() error {
	_, err := datastore.Put(notebook.context, notebook.Key(), &notebook)
	return err
}

func GetNotebook(c appengine.Context) (*Notebook, error) {
	u := user.Current(c)
	notebook := &Notebook{ID: u.ID, context: c}
	key := notebook.Key()
	err := datastore.Get(c, key, notebook)
	if err != nil {
		// store new user
		log.Println("new user", u.Email)
		notebook.Name = u.Email
		key, err = datastore.Put(c, key, notebook)
	}
	notebook.context = c
	return notebook, err
}
