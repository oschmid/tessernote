/**
 * Created with IntelliJ IDEA.
 * User: oschmid
 * Date: 30/09/12
 * Time: 12:45 PM
 * To change this template use File | Settings | File Templates.
 */
package notes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type Note struct {
	Title string
	Body  []byte
	Tags  []string
}

func (note Note) TagsAsString() (*string, error) {
	tags := new(string)
	buffer := bytes.NewBufferString(*tags)
	for _, tag := range note.Tags {
		_, err := buffer.WriteString(tag + tagSeparator)
		if err != nil {
			return nil, err
		}
	}
	return tags, nil // TODO remove last separator
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
