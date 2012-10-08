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
package set

import "testing"

func TestIntersection(t *testing.T) {
	a := New("one", "two", "three", "four")
	b := New("zero", "two", "four", "five")
	expected := New("two", "four")
	actual := Intersection(a, b)
	if !Equals(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}

	actual = Intersection(b, a)
	if !Equals(actual, expected) {
		t.Fatalf("expected=%v actual=%v", expected, actual)
	}
}

func TestDifference(t *testing.T) {
	a := New("one", "two", "three")
	b := New("two", "four", "one", "five")
	c := Difference(a, b)
	d := Difference(b, a)

	expectedC := New("three")
	if !Equals(c, expectedC) {
		t.Fatalf("expected=%v actual=%v", expectedC, c)
	}

	expectedD := New("four", "five")
	if !Equals(d, expectedD) {
		t.Fatalf("expected=%v actual=%v", expectedD, d)
	}
}

func TestUnion(t *testing.T) {
	a := New("one", "two", "three")
	b := New("two", "four", "five", "one")
	c := Union(a, b)
	d := Union(b, a)

	expected := New("one", "two", "three", "four", "five")
	if !Equals(c, expected) {
		t.Fatalf("expected=%v actual=%v", expected, c)
	}
	if !Equals(d, expected) {
		t.Fatalf("expected=%v actual=%v", expected, d)
	}
}

func TestEquals(t *testing.T) {
	a := New("one", "two", "three")
	b := New("one", "two", "three")
	c := New("one", "two", "four")
	d := New("one", "two", "three", "four")
	e := New("two", "three", "one")

	if !Equals(a, b) {
		t.FailNow()
	}
	if Equals(a, c) {
		t.FailNow()
	}
	if Equals(a, d) {
		t.FailNow()
	}
	if !Equals(a, e) {
		t.FailNow()
	}
}
