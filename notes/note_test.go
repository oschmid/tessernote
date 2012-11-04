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
	"string/collections/sets"
	"testing"
)

func TestParseTags(t *testing.T) {
	note := NewNote("title", "body #these #are not all #tags")
	expected := *sets.New("these", "are", "tags")
	actual := note.Tags()
	if !sets.Equal(expected, actual) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestSetBody(t *testing.T) {
	note := NewNote("title", "body #these #are not all #tags")
	expected := *sets.New("these", "are", "tags")
	old := note.Tags()
	if !sets.Equal(expected, old) {
		t.Fatalf("tags not initialized. expected=%v actual=%v", expected, old)
	}

	expectedBody := "body with #some but not all #new #tags"
	note.SetBody(expectedBody)
	actualBody := note.Body
	if expectedBody != actualBody {
		t.Fatalf("body not updated. expected=%v actual=%v", expectedBody, actualBody)
	}

	expected = *sets.New("some", "new", "tags")
	actual := note.Tags()
	if !sets.Equal(expected, actual) {
		t.Fatalf("tags not updated. expected=%v actual=%v", expected, actual)
	}
}
