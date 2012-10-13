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
package slices

import "testing"

func TestEqual(t *testing.T) {
	original := []string{"one", "two", "three"}
	same := []string{"one", "two", "three"}
	if !Equal(original, same) {
		t.Fatal()
	}

	disorder := []string{"two", "three", "one"}
	if !Equal(original, disorder) {
		t.Fatal()
	}

	different := []string{"one", "two", "a"}
	if Equal(original, different) {
		t.Fatal()
	}

	plus := []string{"one", "two", "three", "four"}
	if Equal(original, plus) {
		t.Fatal()
	}

	minus := []string{"one", "two"}
	if Equal(original, minus) {
		t.Fatal()
	}

}

func TestContains(t *testing.T) {
	slice := []string{"one", "two", "three"}
	if !Contains(slice, "one") {
		t.Fatal()
	}
	if Contains(slice, "four") {
		t.Fatal()
	}
}
