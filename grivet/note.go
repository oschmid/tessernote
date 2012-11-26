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
	"log"
	"time"
)

type Note struct {
	ID           string `datastore:"-"` // datastore.Key.Encode()
	Body         string
	Created      time.Time
	LastModified time.Time
	TagKeys      []*datastore.Key
	tags         []Tag // cache
	NotebookKeys []*datastore.Key
}

func (note *Note) Tags(c appengine.Context) ([]Tag, error) {
	if len(note.tags) == 0 && len(note.TagKeys) > 0 {
		note.tags = make([]Tag, len(note.TagKeys))
		err := datastore.GetMulti(c, note.TagKeys, note.tags)
		if err != nil {
			log.Println("getMulti:tags", err)
			return note.tags, err
		}
	}
	return note.tags, nil
}

func (note Note) Notebooks(c appengine.Context) ([]Notebook, error) {
	notebooks := make([]Notebook, len(note.NotebookKeys))
	if len(notebooks) > 0 {
		err := datastore.GetMulti(c, note.NotebookKeys, notebooks)
		if err != nil {
			log.Println("getMulti:notebooks", err)
			return notebooks, err
		}
	}
	return notebooks, nil
}

// updates "note.Body" and "note.LastModified" and
// returns itself
func (note *Note) SetBody(body string, c appengine.Context) (Note, error) {
	key, err := datastore.DecodeKey(note.ID)
	if err != nil {
		log.Println("decodeKey:note", err)
		return *note, err
	}

	note.Body = body
	note.LastModified = time.Now()
	_, err = datastore.Put(c, key, note)
	if err != nil {
		log.Println("put:note", err)
	}
	return *note, err
}

// adds note's key to those of note's tags that are missing it
func (note Note) addKeyToTags(c appengine.Context) error {
	noteKey, err := datastore.DecodeKey(note.ID)
	if err != nil {
		log.Println("decodeKey:note", err)
		return err
	}

	tags, err := note.Tags(c)
	if err != nil {
		log.Println("tags:", err)
		return err
	}

	put := false
	for i := range tags {
		if !containsKey(tags[i].NoteKeys, noteKey) {
			tags[i].NoteKeys = append(tags[i].NoteKeys, noteKey)
			put = true
		}
	}

	if put {
		_, err = datastore.PutMulti(c, note.TagKeys, tags)
		if err != nil {
			log.Println("putMulti:tags", err)
		}
	}
	return err
}

func containsKey(keys []*datastore.Key, key *datastore.Key) bool {
	for i := range keys {
		if keys[i].Encode() == key.Encode() {
			return true
		}
	}
	return false
}

func GetNote(id string, c appengine.Context) (*Note, error) {
	note := new(Note)
	key, err := datastore.DecodeKey(id)
	if err != nil {
		log.Println("decodeKey:", err)
		return note, err
	}

	err = datastore.Get(c, key, note)
	if err != nil {
		log.Println("get:", err)
	}
	note.ID = id
	return note, err
}
