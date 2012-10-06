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
	"testing"
)

func TestIntersection(t *testing.T) {
	a := []string{"one", "two", "three", "four"}
	b := []string{"zero", "two", "four", "five"}
	expected := []string{"two", "four"}
	actual := intersection(a, b)
	if equals(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestEquals(t *testing.T) {
	a := []string{"one", "two", "three"}
	b := []string{"one", "two", "three"}
	c := []string{"one", "two", "four"}
	d := []string{"one", "two", "three", "four"}
	if !equals(a, b) {
		t.FailNow()
	}
	if equals(a, c) {
		t.FailNow()
	}
	if equals(a, d) {
		t.FailNow()
	}
}
