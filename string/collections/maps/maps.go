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

func New(elements ...string) *map[string]int {
	s := make(map[string]int)
	for _, elem := range elements {
		s[elem] = 1
	}
	return &s
}

func Union(a map[string]int, b map[string]int) *map[string]int {
	c := make(map[string]int)
	for elem, count := range a {
		c[elem] = count
	}
	for elem, count := range b {
		prevCount, contained := c[elem]
		if contained {
			c[elem] = count + prevCount
		} else {
			c[elem] = count
		}
	}
	return &c
}

func Equal(a map[string]int, b map[string]int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for elemA, countA := range a {
		countB, contained := b[elemA]
		if !contained || countA != countB {
			return false
		}
	}
	return true
}
