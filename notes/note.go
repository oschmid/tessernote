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
	"strings"
	"unicode"
)

// TODO merge Title and Body
type Note struct {
	Id    string
	Title string
	Body  string // use SetBody() to update value
	tags  map[string]bool
}

func NewNote(title string, body string) *Note {
	note := &Note{Title: title, Body: body}
	note.Id = uuid.New()
	return note
}

func (n *Note) SetBody(body string) {
	n.Body = body
	n.tags = nil
}

func (n *Note) Tags() map[string]bool {
	if n.tags == nil {
		n.tags = parseTags(n.Body)
	}
	return n.tags
}

func parseTags(body string) map[string]bool {
	tags := make(map[string]bool)
	matches := Hashtag.FindAllString(body, len(body))
	for _, tag := range matches {
		tags[strings.TrimFunc(tag, isHashtagDecoration)] = true
	}
	return tags
}

func isHashtagDecoration(r rune) bool {
	return r == '#' || unicode.IsSpace(r)
}

func (n Note) Equal(note Note) bool {
	return n.Id == note.Id && n.Title == note.Title && n.Body == note.Body
}
