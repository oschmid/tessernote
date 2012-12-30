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
	"appengine/memcache"
	"appengine/user"
	"encoding/gob"
	"errors"
	"github.com/oschmid/cachestore"
	"time"
)

var Debug = false // If true, print debug info

func init() {
	gob.Register(*new(Notebook))
	gob.Register(*new(Note))
	gob.Register(*new(Tag))
	gob.Register(*new(time.Time))
	gob.Register(*new(datastore.Key))
}

// TODO use Key sets (and implement PropertyLoadSaver) not arrays
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

// Key returns a datastore.Key for Notebook
func (notebook Notebook) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Notebook", notebook.ID, 0, nil)
}

// Tags returns all tags used to sort this Notebook's notes
func (notebook *Notebook) Tags(c appengine.Context) ([]Tag, error) {
	if len(notebook.tags) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.tags = make([]Tag, len(notebook.TagKeys))
		err := cachestore.GetMulti(c, notebook.TagKeys, notebook.tags)
		if err != nil {
			c.Errorf("getting notebook tags: %s", err)
			return notebook.tags, err
		}
	}
	return notebook.tags, nil
}

// Note returns a note by its ID.
func (notebook *Notebook) Note(id string, c appengine.Context) (note Note, err error) {
	key, err := datastore.DecodeKey(id)
	if err == nil {
		if containsKey(notebook.NoteKeys, key) {
			err = cachestore.Get(c, key, &note)
		} else {
			err = errors.New("notebook does not contain note with ID: " + id)
		}
	}
	return note, err
}

// Notes returns all notes in this Notebook including untagged notes
func (notebook *Notebook) Notes(c appengine.Context) ([]Note, error) {
	if len(notebook.notes) == 0 && len(notebook.NoteKeys) > 0 {
		notebook.notes = make([]Note, len(notebook.NoteKeys))
		err := cachestore.GetMulti(c, notebook.NoteKeys, notebook.notes)
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

// UntaggedNotes returns all untagged notes in this Notebook
func (notebook *Notebook) UntaggedNotes(c appengine.Context) ([]Note, error) {
	if len(notebook.untaggedNotes) == 0 && len(notebook.UntaggedNoteKeys) > 0 {
		notebook.untaggedNotes = make([]Note, len(notebook.UntaggedNoteKeys))
		err := cachestore.GetMulti(c, notebook.UntaggedNoteKeys, notebook.untaggedNotes)
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

// TagsFrom returns tags in this Notebook by name. Returns an error if a tag is missing  
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
			if Debug {
				c.Debugf("user missing tag: %s", name)
			}
			return tags, errors.New("tessernote: missing tag (" + name + ")")
		}
	}
	return tags, nil
}

// TagsOf returns the Tags of a Note in this Notebook
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

// RelatedTags returns all Tags in this Notebook that refer to the Notes referred to by a subset of Tags.
// 
// For example: if Tags A and B refer to Note C and only Tag A is given as input, the output will be A and B.
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

// Put updates a Note or creates it if it doesn't already exist and sorts out all Tag relationships.
func (notebook *Notebook) Put(note Note, c appengine.Context) (Note, error) {
	if note.ID != "" && containsKey(notebook.NoteKeys, note.Key(c)) {
		return notebook.updateNote(note, c)
	}
	return notebook.addNote(note, c)
}

// addNote adds a Note to this Notebook, updating existing Tags to point to it if they're mentioned
// and adding any new Tags
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
		if Debug {
			tc.Debugf("updating note (with tags): %#v", note)
		}
		key, err = cachestore.Put(tc, key, &note)
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
	if err != nil {
		memcache.Flush(c)
	}
	return note, err
}

// addNoteWithoutTags adds a Note to the datastore so that a unique Key is created for it. That Key is
// then used for adding/updating Tags.
func (notebook Notebook) addNoteWithoutTags(note *Note, c appengine.Context) (*datastore.Key, error) {
	note.Created = time.Now()
	note.LastModified = note.Created
	note.NotebookKeys = []*datastore.Key{notebook.Key(c)}
	key := notebook.newNoteKey(note, c)
	if Debug {
		c.Debugf("adding note (without tags): %#v", note)
	}
	key, err := cachestore.Put(c, key, note)
	if err != nil {
		c.Errorf("adding note (without tags): %s", err)
		return nil, err
	}
	note.ID = key.Encode()
	return key, nil
}

// newNoteKey returns a unique Key for a new Note. New Notes can be created with Keys generated outside of 
// Tessernote. However if this Key is not unique (e.g. if it refers to a Note in another Notebook) then a new
// Key is generated.
func (notebook Notebook) newNoteKey(note *Note, c appengine.Context) *datastore.Key {
	if note.ID != "" {
		key, err := datastore.DecodeKey(note.ID)
		if err == nil {
			err = cachestore.Get(c, key, new(interface{}))
			if err == datastore.ErrNoSuchEntity {
				return key
			}
		}
		// note exists (probably in another notebook), use a different key
		note.ID = ""
	}
	return datastore.NewIncompleteKey(c, "Note", notebook.Key(c))
}

// updateTags calculates the Tag difference between oldNote and note. It adds missing Tags to the datastore,
// removes note from Tags that refered to oldNote but not note, and cleans up Tags that no longer refer to any Notes.
func (notebook *Notebook) updateTags(key *datastore.Key, oldNote, note *Note, c appengine.Context) error {
	tagKeys, tags, count, deleted, err := notebook.updateTagKeys(oldNote, note, c)
	if err != nil {
		return err
	}
	if len(tagKeys) > 0 {
		if Debug {
			c.Debugf("adding/updating tags: %#v", tags)
		}
		tagKeys, err := cachestore.PutMulti(c, tagKeys, tags)
		if err != nil {
			c.Errorf("adding/updating tags: %s", err)
			return err
		}
		// update note tags
		note.TagKeys = tagKeys[:count]
	}
	if len(deleted) > 0 {
		if Debug {
			c.Debugf("deleting empty tags: %#v", deleted)
		}
		err = cachestore.DeleteMulti(c, deleted)
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

// updateTagKeys updates tags in memory to reflect the changes of turning oldNote into note and returns
// the objects to commit to the datastore to make these changes permanent.
func (notebook *Notebook) updateTagKeys(oldNote, note *Note, c appengine.Context) (keys []*datastore.Key, tags []Tag, count int, deleted []*datastore.Key, err error) {
	// get note tags
	keys, tags, names, err := notebook.parseTagsOf(*note, c)
	count = len(keys)

	// get remove tags
	removedFromTagKeys, removedFromTags, deleted, err := notebook.removedTags(oldNote, names, c)
	keys = append(keys, removedFromTagKeys...)
	tags = append(tags, removedFromTags...)
	return keys, tags, count, deleted, err
}

// parseTagsOf parses the hashtags of note.Body, and returns the associated Tags. Missing Tags are also created with
// incomplete Keys.
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
			tags = append(tags, *NewTag(name, note, *notebook, c))
		}
	}
	return keys, tags, names, nil
}

// removedTags returns the Tags in oldNote that are not named in names and the Tags that can be cleaned up because they no longer
// refer to any Notes.
func (notebook *Notebook) removedTags(oldNote *Note, names []string, c appengine.Context) (removedFromKeys []*datastore.Key, removedFromTags []Tag, deleteKeys []*datastore.Key, err error) {
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

// removeTagKeys removes tag Keys from this Notebook. Tags not in this Notebook are ignored.
func (notebook *Notebook) removeTagKeys(tagKeys []*datastore.Key) {
	notebook.tags = *new([]Tag)
	for _, key := range tagKeys {
		notebook.TagKeys = removeKey(notebook.TagKeys, key)
	}
}

// addTagKeys adds missing tag Keys to this Notebook. Keys of tags that already exist are ignored.
func (notebook *Notebook) addTagKeys(tagKeys []*datastore.Key) {
	notebook.tags = *new([]Tag)
	for _, key := range tagKeys {
		notebook.TagKeys = addKey(notebook.TagKeys, key)
	}
}

// PutAll adds or updates notes and sorts out all Tag relationships.
func (notebook *Notebook) PutAll(notes []Note, c appengine.Context) ([]Note, error) {
	// TODO
	return notes, errors.New("not yet implemented")
}

// save updates this Notebook in the datastore
func (notebook *Notebook) save(c appengine.Context) error {
	if Debug {
		c.Debugf("updating notebook: %#v", *notebook)
	}
	_, err := cachestore.Put(c, notebook.Key(c), notebook)
	if err != nil {
		c.Errorf("updating notebook: %s", err)
	}
	return err
}

// updateNote updates a Note in this Notebook, updating existing Tags to either start or stop pointing to it,
// cleaning up Tags that no longer point to any Note, and adding any new Tags.
func (notebook *Notebook) updateNote(note Note, c appengine.Context) (Note, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		// get old note
		var oldNote Note
		key := note.Key(tc)
		err := cachestore.Get(tc, key, &oldNote)
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
		if Debug {
			tc.Debugf("updating note: %#v", note)
		}
		key, err = cachestore.Put(tc, key, &note)
		if err != nil {
			tc.Errorf("updating note: %s", err)
			return err
		}

		// update notebook
		return notebook.save(tc)
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		memcache.Flush(c)
	}
	return note, err
}

// Delete deletes a Note from this Notebook, removes it from any Tags that refer to it and deletes any Tags
// that no longer refer to any Notes
func (notebook *Notebook) Delete(id string, c appengine.Context) (bool, error) {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		note := Note{ID: id}
		noteKey := note.Key(c)
		err := cachestore.Get(tc, noteKey, &note)
		if err != nil {
			c.Errorf("getting note: %s", err)
			return err
		}

		// remove note
		if Debug {
			tc.Debugf("deleting note: %#v", note)
		}
		err = cachestore.Delete(tc, noteKey)
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
	if err != nil {
		memcache.Flush(c)
	}
	return err == nil, err
}

// DeleteAll deletes all Notes and Tags from this Notebook.
func (notebook *Notebook) DeleteAll(c appengine.Context) error {
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		err := cachestore.DeleteMulti(tc, notebook.NoteKeys)
		if err != nil {
			return err
		}
		return cachestore.DeleteMulti(tc, notebook.TagKeys)
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		memcache.Flush(c)
	}
	return nil
}

// GetNotebook returns a user's unique Notebook
func GetNotebook(c appengine.Context) (*Notebook, error) {
	notebook := new(Notebook)
	u := user.Current(c)
	if u == nil {
		return notebook, errors.New("user is null")
	}
	notebook.ID = u.ID
	key := notebook.Key(c)
	err := cachestore.Get(c, key, notebook)
	if err != nil {
		if err != datastore.ErrNoSuchEntity {
			c.Warningf(err.Error())
		}
		// create new user
		if Debug {
			c.Debugf("adding new notebook for: %s", u.Email)
		}
		notebook.Name = u.Email
		key, err = cachestore.Put(c, key, notebook)
	}
	return notebook, err
}
