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
	"collections/set"
	"testing"
)

func TestTagString(t *testing.T) {
	tag1, tag2 := "tag1", "tag2"
	tags := *set.New(tag1, tag2)
	note := NewNote("title", "body", tags)
	tagString := note.TagString()
	expected := tag1 + TAG_SEPARATOR + tag2
	if tagString != expected {
		t.Fatalf("expected=%v actual=%v", expected, tagString)
	}
}
