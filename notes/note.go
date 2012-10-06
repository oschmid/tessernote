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
	"encoding/gob"
	"fmt"
	"os"
)

type Note struct {
	Title string
	Body  string
	Tags  []string // TODO convert to set
}

func (note Note) TagsAsString() string {
	if len(note.Tags) == 0 {
		return ""
	}

	tags := note.Tags[0]
	for _, tag := range note.Tags[1:] {
		tags += tagSeparator + tag
	}

	return tags
}

func saveNote(note Note) error {
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

func loadNote(title string) (*Note, error) {
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
