/*
This file is part of Notes.

Notes is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Notes is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
 */
package notes

import (
	"testing"
)

func TestTagsAsString(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{Title: "title", Body: []byte("body"), Tags: tags}
	stringTags, err := note.TagsAsString()
	if *stringTags != tags[0]+tagSeparator+tags[1] {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
