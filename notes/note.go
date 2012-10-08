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
package notes

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

const TAG_SEPARATOR = ", "

type Note struct {
	Id    string
	Title string
	Body  string
	Tags  map[string]bool
}

func NewNote(title string, body string, tags map[string]bool) *Note {
	note := &Note{Title: title, Body: body, Tags: tags}
	note.Id = uuid.New()
	return note
}

// Returns a string representation of this Note's tags
func (note Note) TagString() string {
	if len(note.Tags) == 0 {
		return ""
	}

	// sort tags in alphabetical order
	tags := []string{}
	for tag, _ := range note.Tags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	tagString := *new(string)
	for _, tag := range tags {
		tagString += tag + TAG_SEPARATOR
	}

	return tagString[:len(tagString)-len(TAG_SEPARATOR)]
}

// TODO save NoteBooks instead of individual Notes

func SaveNote(note Note) error {
	fileName := "data/" + note.Title + ".txt"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("saveNote error: %v\n", err)
		return err
	}

	defer file.Close()
	gob.NewEncoder(file).Encode(note)
	return nil
}

func LoadNote(title string) (*Note, error) {
	fileName := "data/" + title + ".txt"
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	var note Note
	err = gob.NewDecoder(file).Decode(&note)
	return &note, err
}
