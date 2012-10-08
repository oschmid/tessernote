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

func New(elements ...string) *map[string]bool {
	s := make(map[string]bool)
	for _, elem := range elements {
		s[elem] = true
	}
	return &s
}

// Returns the intersection of sets "a" and "b" i.e. all the elements that are in both "a" and "b"
func Intersection(s map[string]bool, t map[string]bool) *map[string]bool {
	u := make(map[string]bool)
	for elem, sContains := range s {
		if sContains && t[elem] {
			u[elem] = true
		}
	}
	return &u
}

// Returns a-b i.e. all the elements of a that are not in b
func Difference(a map[string]bool, b map[string]bool) *map[string]bool {
	c := make(map[string]bool)
	for elem, _ := range a {
		if !b[elem] {
			c[elem] = true
		}
	}
	return &c
}

// Returns the union of sets "a" and "b" i.e. all the elements that are in "a", "b" or both "a" and "b"
func Union(a map[string]bool, b map[string]bool) *map[string]bool {
	c := make(map[string]bool)
	for elem, _ := range a {
		c[elem]=true
	}
	for elem, _ := range b {
		c[elem]=true
	}
	return &c
}

func Equals(a map[string]bool, b map[string]bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for elem, _ := range a {
		if !b[elem] {
			return false
		}
	}

	return true
}
