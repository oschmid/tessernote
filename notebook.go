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
	"appengine/user"
	"errors"
	"github.com/oschmid/tessernote/store"
	"time"
)

// TODO use sets (with datastore loaders) not arrays
type Notebook struct {
	ID               string // user.User.ID
	Name             string
	TagKeys          []*datastore.Key // sorted by Tag.Name
	NoteKeys         []*datastore.Key
	UntaggedNoteKeys []*datastore.Key
	tags             []Tag  // cache
	notes            []Note // cache
	untaggedNotes    []Note //cache
}

func (notebook Notebook) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Notebook", notebook.ID, 0, nil)
}

func (notebook *Notebook) Tags(c appengine.Context) ([]Tag, error) {
	if len(notebook.tags) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.tags = make([]Tag, len(notebook.TagKeys))
		err := store.GetMulti(c, notebook.TagKeys, notebook.tags)
		if err != nil {
			c.Errorf("getting notebook tags: %s", err)
			return notebook.tags, err
		}
	}
	return notebook.tags, nil
}

func (notebook *Notebook) Notes(c appengine.Context) ([]Note, error) {
	if len(notebook.notes) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.notes = make([]Note, len(notebook.NoteKeys))
		err := store.GetMulti(c, notebook.NoteKeys, notebook.notes)
		if err != nil {
			c.Errorf("getting notebook notes: %s", err)
			return notebook.notes, err
		}
		for i := 0; i < len(notebook.notes); i++ {
			note := &notebook.notes[i]
			note.ID = notebook.NoteKeys[i].Encode()
		}
	}
	return notebook.notes, nil
}

func (notebook *Notebook) UntaggedNotes(c appengine.Context) ([]Note, error) {
	if len(notebook.untaggedNotes) == 0 && len(notebook.UntaggedNoteKeys) > 0 {
		notebook.untaggedNotes = make([]Note, len(notebook.UntaggedNoteKeys))
		err := store.GetMulti(c, notebook.UntaggedNoteKeys, notebook.untaggedNotes)
		if err != nil {
			c.Errorf("getting notebook untagged notes: %s", err)
			return notebook.untaggedNotes, err
		}
		for i := 0; i < len(notebook.untaggedNotes); i++ {
			note := &notebook.untaggedNotes[i]
			note.ID = notebook.UntaggedNoteKeys[i].Encode()
		}
	}
	return notebook.untaggedNotes, nil
}

// returns a subset of a user's tags by name
// missing tags result in errors
func (notebook *Notebook) TagsFrom(names []string, c appengine.Context) (tags []Tag, err error) {
	allTags, err := notebook.Tags(c)
	if err != nil {
		return tags, err
	}
	for _, name := range names {
		i := indexOfTag(allTags, name)
		if i >= 0 {
			tags = append(tags, allTags[i])
		} else {
			c.Debugf("user missing tag: %s", name)
			return tags, errors.New("user missing tag: " + name)
		}
	}
	return tags, nil
}

func (notebook *Notebook) TagsOf(note Note, c appengine.Context) (tags []Tag, err error) {
	allTags, err := notebook.Tags(c)
	if err != nil {
		return tags, err
	}
	for _, key := range note.TagKeys {
		i := indexOfKey(notebook.TagKeys, key)
		if i >= 0 {
			tags = append(tags, allTags[i])
		} else {
			c.Errorf("notebook missing tag: %s", key)
			return tags, errors.New("notebook missing tag: " + key.Encode())
		}
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

func (notebook *Notebook) Put(note Note, c appengine.Context) (Note, error) {
	if note.ID != "" && containsKey(notebook.NoteKeys, note.Key(c)) {
		return notebook.updateNote(note, c)
	}
	return notebook.addNote(note, c)
}

func (notebook *Notebook) addNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		// add note (without tags) TODO add existing tags
		key, err := notebook.addNoteWithoutTags(&note, tc)
		if err != nil {
			return err
		}

		// add/update tags
		err = notebook.updateTags(key, new(Note), &note, tc)
		if err != nil {
			return err
		}

		// update note (with tags) TODO skip if no new tags
		tc.Debugf("updating note (with tags): %+v", note)
		key, err = store.Put(tc, key, &note)
		if err != nil {
			tc.Errorf("updating note (with tags): %s", err)
			return err
		}

		// update notebook
		notebook.NoteKeys = append(notebook.NoteKeys, key)
		if len(note.TagKeys) > 0 {
			notebook.addTagKeys(note.TagKeys)
		} else {
			notebook.UntaggedNoteKeys = append(notebook.UntaggedNoteKeys, key)
		}
		return notebook.save(tc)
	}, &datastore.TransactionOptions{XG: true})
	return note, err
}

func (notebook Notebook) addNoteWithoutTags(note *Note, c appengine.Context) (*datastore.Key, error) {
	note.Created = time.Now()
	note.LastModified = note.Created
	note.NotebookKeys = []*datastore.Key{notebook.Key(c)}
	key := notebook.newNoteKey(note, c)
	c.Debugf("adding note (without tags): %+v", note)
	key, err := store.Put(c, key, note)
	if err != nil {
		c.Errorf("adding note (without tags): %s", err)
		return nil, err
	}
	note.ID = key.Encode()
	return key, nil
}

func (notebook Notebook) newNoteKey(note *Note, c appengine.Context) *datastore.Key {
	if note.ID != "" {
		key, err := datastore.DecodeKey(note.ID)
		if err == nil {
			err = store.Get(c, key, new(interface{}))
			if err == datastore.ErrNoSuchEntity {
				return key
			}
		}
		// note exists (probably in another notebook), use a different key
		note.ID = ""
	}
	return datastore.NewIncompleteKey(c, "Note", notebook.Key(c))
}

// Gets the tag difference between 'oldNote' and 'note'.
// Adds missing tags to the datastore, removes note from unused tags, adds note to new existing tags, removes empty tags.
// Updates 'notebook' and 'note' tag keys.
func (notebook *Notebook) updateTags(key *datastore.Key, oldNote, note *Note, c appengine.Context) error {
	tagKeys, tags, count, deleted, err := notebook.updateTagKeys(oldNote, note, c)
	if err != nil {
		return err
	}
	if len(tagKeys) > 0 {
		c.Debugf("adding/updating tags: %+v", tags)
		tagKeys, err := store.PutMulti(c, tagKeys, tags)
		if err != nil {
			c.Errorf("adding/updating tags: %s", err)
			return err
		}
		// update note tags
		note.TagKeys = tagKeys[:count]
	}
	if len(deleted) > 0 {
		c.Debugf("deleting empty tags: %+v", deleted)
		err = store.DeleteMulti(c, deleted)
		if err != nil {
			c.Errorf("deleting empty tags: %s", err)
		}
	}
	// update notebook tags
	if len(oldNote.TagKeys) > 0 && len(note.TagKeys) > 0 {
		notebook.removeTagKeys(deleted)
		notebook.addTagKeys(note.TagKeys)
	} else if len(oldNote.TagKeys) == 0 && len(note.TagKeys) > 0 {
		notebook.UntaggedNoteKeys = removeKey(notebook.UntaggedNoteKeys, key)
		notebook.addTagKeys(note.TagKeys)
	} else if len(oldNote.TagKeys) > 0 && len(note.TagKeys) == 0 {
		notebook.removeTagKeys(deleted)
		if note.ID != "" {
			notebook.UntaggedNoteKeys = append(notebook.UntaggedNoteKeys, key)
		}
	} else if note.ID == "" {
		notebook.UntaggedNoteKeys = removeKey(notebook.UntaggedNoteKeys, key)
	}
	return err
}

// Updates the tag keys of 'notebook' and 'note' to reflect the changes of turning 'oldNote' into 'note.
// Returns the datastore changes that make this change permanent.
func (notebook *Notebook) updateTagKeys(oldNote, note *Note, c appengine.Context) (keys []*datastore.Key, tags []Tag, count int, deleted []*datastore.Key, err error) {
	// get note tags
	keys, tags, names, err := notebook.parseTagsOf(*note, c)
	count = len(keys)

	// get remove tags
	removedFromTagKeys, removedFromTags, deleted, err := notebook.getRemovedTags(oldNote, names, c)
	keys = append(keys, removedFromTagKeys...)
	tags = append(tags, removedFromTags...)
	return keys, tags, count, deleted, err
}

// Parses the hashtags of a note.Body. Finds the associated Tags. Creates missing tags.
// Returns tags, their keys and names, creating new incomplete keys for those that don't yet exist.
func (notebook *Notebook) parseTagsOf(note Note, c appengine.Context) (keys []*datastore.Key, tags []Tag, names []string, err error) {
	names = ParseTagNames(note.Body)
	allTags, err := notebook.Tags(c)
	if err != nil {
		return keys, tags, names, err
	}
	notebookKey := notebook.Key(c)
	for _, name := range names {
		i := indexOfTag(allTags, name)
		if i >= 0 {
			allTags[i].NoteKeys = addKey(allTags[i].NoteKeys, note.Key(c))
			keys = append(keys, notebook.TagKeys[i])
			tags = append(tags, allTags[i])
		} else {
			keys = append(keys, datastore.NewIncompleteKey(c, "Tag", notebookKey))
			tags = append(tags, *NewTag(name, *notebook, note, c))
		}
	}
	return keys, tags, names, nil
}

// Gets the tags that note was removed from and the tags that can be deleted because they no longer refer to any notes.
func (notebook *Notebook) getRemovedTags(oldNote *Note, names []string, c appengine.Context) (removedFromKeys []*datastore.Key, removedFromTags []Tag, deleteKeys []*datastore.Key, err error) {
	oldTags, err := notebook.TagsOf(*oldNote, c)
	if err != nil {
		return removedFromKeys, removedFromTags, deleteKeys, err
	}
	// remove from old tags
	for i := range oldTags {
		if !containsString(names, oldTags[i].Name) {
			if len(oldTags[i].NoteKeys) == 1 {
				deleteKeys = append(deleteKeys, oldNote.TagKeys[i])
			} else {
				oldTags[i].NoteKeys = removeKey(oldTags[i].NoteKeys, oldNote.Key(c))
				removedFromKeys = append(removedFromKeys, oldNote.TagKeys[i])
				removedFromTags = append(removedFromTags, oldTags[i])
			}
		}
	}
	return removedFromKeys, removedFromTags, deleteKeys, nil
}

// Removes tag keys from notebook. Ignores tag keys not in notebook.
func (notebook *Notebook) removeTagKeys(tagKeys []*datastore.Key) {
	notebook.tags = *new([]Tag)
	for _, key := range tagKeys {
		notebook.TagKeys = removeKey(notebook.TagKeys, key)
	}
}

// Adds missing tag keys to notebook. Ignores tag keys that already exist.
func (notebook *Notebook) addTagKeys(tagKeys []*datastore.Key) {
	notebook.tags = *new([]Tag)
	for _, key := range tagKeys {
		notebook.TagKeys = addKey(notebook.TagKeys, key)
	}
}

func (notebook *Notebook) save(c appengine.Context) error {
	c.Debugf("updating notebook: %+v", *notebook)
	_, err := store.Put(c, notebook.Key(c), notebook)
	if err != nil {
		c.Errorf("updating notebook: %s", err)
	}
	return err
}

func (notebook *Notebook) updateNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		// get old note
		var oldNote Note
		key := note.Key(tc)
		err := store.Get(tc, key, &oldNote)
		if err != nil {
			tc.Errorf("getting old note: %s", err)
			return err
		}
		oldNote.ID = note.ID

		// add/update/delete tags
		err = notebook.updateTags(key, &oldNote, &note, tc)
		if err != nil {
			return err
		}

		// update note
		note.Created = oldNote.Created
		note.LastModified = time.Now()
		note.NotebookKeys = oldNote.NotebookKeys
		tc.Debugf("updating note: %+v", note)
		key, err = store.Put(tc, key, &note)
		if err != nil {
			tc.Errorf("updating note: %s", err)
			return err
		}

		// update notebook
		return notebook.save(tc)
	}, &datastore.TransactionOptions{XG: true})
	return note, err
}

func (notebook *Notebook) Delete(id string, c appengine.Context) (bool, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		note := Note{ID: id}
		noteKey := note.Key(c)
		err := store.Get(tc, noteKey, &note)
		if err != nil {
			c.Errorf("getting note: %s", err)
			return err
		}

		// remove note
		tc.Debugf("deleting note: %+v", note)
		err = store.Delete(tc, noteKey)
		if err != nil {
			c.Errorf("deleting note: %s", err)
			return err
		}

		// remove note from tags
		err = notebook.updateTags(noteKey, &note, new(Note), tc)
		if err != nil {
			return err
		}

		// remove note from notebook
		notebook.NoteKeys = removeKey(notebook.NoteKeys, noteKey)
		return notebook.save(tc)
	}, &datastore.TransactionOptions{XG: true})
	return err == nil, err
}

func GetNotebook(c appengine.Context) (*Notebook, error) {
	notebook := new(Notebook)
	u := user.Current(c)
	if u == nil {
		return notebook, errors.New("user is null")
	}
	notebook.ID = u.ID
	key := notebook.Key(c)
	err := store.Get(c, key, notebook)
	if err != nil {
		// create new user
		c.Infof("adding new notebook for:", u.Email)
		notebook.Name = u.Email
		key, err = store.Put(c, key, notebook)
	}
	return notebook, err
}
