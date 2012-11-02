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
	"sort"
	"string/collections/sets"
	"strings"
	"unicode"
)

const TagSeparator = ", "

// TODO parse hashtags from body
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

func (note Note) TitleBodyString() string {
	return note.Title + TitleBodySeparator + note.Body
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
		tagString += tag + TagSeparator
	}

	return tagString[:len(tagString)-len(TagSeparator)]
}

// parses hashtags from note.Body
func (n Note) parseTags() []string {
	tags := Hashtag.FindAllString(n.Body, len(n.Body))
	for i, tag := range tags {
		tags[i] = strings.TrimFunc(tag, func(r rune) bool {
				if r == '#' {
					return true
				}
				return unicode.IsSpace(r)
			})
	}
	return tags
}

func (n Note) Equal(note Note) bool {
	return n.Id == note.Id && n.Title == note.Title && n.Body == note.Body && sets.Equal(n.Tags, note.Tags)
}
