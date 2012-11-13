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
	"sort"
)

type User struct {
	ID       string // user.User.ID
	Name     string
	TagKeys  []*datastore.Key // sorted by Tag.Name
	NoteKeys []*datastore.Key
	context  appengine.Context `datastore:",noindex"`
}

func (u User) Tags() []Tag {
	var tags []Tag
	datastore.GetMulti(u.context, u.TagKeys, tags)
	for _, t := range tags {
		t.context = u.context
	}
	return tags
}

func (u User) Notes() []Note {
	var notes []Note
	datastore.GetMulti(u.context, u.NoteKeys, notes)
	for _, n := range notes {
		n.context = u.context
	}
	return notes
}

// returns a subset of a user's tags by name
// missing tags result in errors
func (u User) TagsFrom(names []string) ([]Tag, error) {
	tags := *new([]Tag)
	if len(names) == 0 {
		return tags, nil
	}
	sort.Strings(names)
	allTags := u.Tags()
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

// returns a user's note by ID, missing note will be skipped
func (u User) Note(id string) (Note, error) {
	// TODO binary search
	var note Note
	for _, key := range u.NoteKeys {
		if key.Encode() == id {
			datastore.Get(u.context, key, note)
			return note, nil
		}
	}
	return note, errors.New("note does not exist")
}

func (u User) RelatedTags(tags []Tag) []Tag {
	relatedNoteKeys := make(map[string]datastore.Key)
	for _, tag := range tags {
		for _, key := range tag.NoteKeys {
			relatedNoteKeys[key.Encode()] = *key
		}
	}
	tags = *new([]Tag)
	for _, tag := range u.Tags() {
		for _, key := range tag.NoteKeys {
			if _, contained := relatedNoteKeys[key.Encode()]; contained {
				tags = append(tags, tag)
				break
			}
		}
	}
	return tags
}

func CurrentUser(c appengine.Context) *User {
	g := new(User)
	u := user.Current(c)
	k := datastore.NewKey(c, "User", u.ID, 0, nil)
	if err := datastore.Get(c, k, u); err != nil {
		// store new user
		g.Name = u.Email
		k, err = datastore.Put(c, k, g)
	}
	g.context = c
	return g
}

func PutUser(u User) {
	k := datastore.NewKey(u.context, "User", u.ID, 0, nil)
	if _, err := datastore.Put(u.context, k, &u); err != nil {
		// TODO handle failed put
	}
}
