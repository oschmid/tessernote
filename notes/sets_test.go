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
	a := set("one", "two", "three", "four")
	b := set("zero", "two", "four", "five")
	expected := set("two", "four")
	actual := intersection(a, b)
	if !equals(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}

	actual = intersection(b, a)
	if !equals(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestDifference(t *testing.T) {
	a := set("one", "two", "three")
	b := set("two", "four", "one", "five")
	c := difference(a, b)
	d := difference(b, a)

	expectedC := set("three")
	if !equals(c, expectedC) {
		t.Fatalf("expected=%v actual=%v", expectedC, c)
	}

	expectedD := set("four", "five")
	if !equals(d, expectedD) {
		t.Fatalf("expected=%v actual=%v", expectedD, d)
	}
}

func TestEquals(t *testing.T) {
	a := set("one", "two", "three")
	b := set("one", "two", "three")
	c := set("one", "two", "four")
	d := set("one", "two", "three", "four")
	e := set("two", "three", "one")

	if !equals(a, b) {
		t.FailNow()
	}
	if equals(a, c) {
		t.FailNow()
	}
	if equals(a, d) {
		t.FailNow()
	}
	if !equals(a, e) {
		t.FailNow()
	}
}
