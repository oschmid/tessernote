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

package tessernote

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Note struct {
	ID           string `datastore:"-"` // datastore.Key.Encode()
	Body         string
	Created      time.Time
	LastModified time.Time
	TagKeys      []*datastore.Key
	NotebookKeys []*datastore.Key
	tags         []Tag // cache
}

func (note *Note) Tags(c appengine.Context) ([]Tag, error) {
	if len(note.tags) == 0 && len(note.TagKeys) > 0 {
		note.tags = make([]Tag, len(note.TagKeys))
		err := datastore.GetMulti(c, note.TagKeys, note.tags)
		if err != nil {
			c.Errorf("getting note tags:", err)
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
			c.Errorf("getting note notebooks:", err)
			return notebooks, err
		}
	}
	return notebooks, nil
}

// updates "note.Body" and "note.LastModified" and
// returns itself
func (note *Note) Update(new Note, c appengine.Context) (Note, error) {
	key, err := datastore.DecodeKey(note.ID)
	if err != nil {
		c.Errorf("decoding note key:", err)
		return *note, err
	}

	note.Body = new.Body
	note.LastModified = time.Now()
	note.TagKeys = new.TagKeys
	note.tags = new.tags
	note.NotebookKeys = note.NotebookKeys
	_, err = datastore.Put(c, key, note)
	if err != nil {
		c.Errorf("updating note:", err)
	}
	return *note, err
}

// adds note's key to those of note's tags that are missing it
func (note Note) addKeyToTags(c appengine.Context) error {
	noteKey, err := datastore.DecodeKey(note.ID)
	if err != nil {
		c.Errorf("decoding note key:", err)
		return err
	}

	tags, err := note.Tags(c)
	if err != nil {
		return err
	}

	tagKeysToUpdate := *new([]*datastore.Key)
	tagsToUpdate := *new([]Tag)
	for i := range tags {
		if !containsKey(tags[i].NoteKeys, noteKey) {
			tags[i].NoteKeys = append(tags[i].NoteKeys, noteKey)
			tagKeysToUpdate = append(tagKeysToUpdate, note.TagKeys[i])
			tagsToUpdate = append(tagsToUpdate, tags[i])
		}
	}

	if len(tagsToUpdate) > 0 {
		_, err = datastore.PutMulti(c, tagKeysToUpdate, tagsToUpdate)
		if err != nil {
			c.Errorf("updating note tags:", err)
		}
	}
	return err
}

func GetNote(id string, c appengine.Context) (*Note, error) {
	note := new(Note)
	key, err := datastore.DecodeKey(id)
	if err != nil {
		c.Errorf("decoding note key:", err)
		return note, err
	}

	err = datastore.Get(c, key, note)
	if err != nil {
		c.Errorf("getting note:", err)
	}
	note.ID = id
	return note, err
}
