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
package maps

import "testing"

func TestUnion(t *testing.T) {
	a := *New("one", "two", "three")
	b := *New("two", "four", "five", "one")
	c := *Union(a, b)
	d := *Union(b, a)

	expected := map[string]int{"one": 2, "two": 2, "three": 1, "four": 1, "five": 1}
	if !Equals(c, expected) {
		t.Fatalf("expected=%v actual=%v", expected, c)
	}
	if !Equals(d, expected) {
		t.Fatalf("expected=%v actual=%v", expected, d)
	}
}

func TestEquals(t *testing.T) {
	a := *New("one", "two", "three")
	b := *New("one", "two", "three")
	c := *New("one", "two", "four")
	d := *New("one", "two", "three", "four")
	e := *New("two", "three", "one")
	f := map[string]int{"one": 2, "two": 1, "three": 3}

	if !Equals(a, b) {
		t.Fatal()
	}
	if Equals(a, c) {
		t.Fatal()
	}
	if Equals(a, d) {
		t.Fatal()
	}
	if !Equals(a, e) {
		t.Fatal()
	}
	if Equals(a, f) {
		t.Fatal()
	}
}
