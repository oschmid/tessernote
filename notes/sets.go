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

func set(elements ...string) map[string]bool {
	s := make(map[string]bool)
	for _, elem := range elements {
		s[elem] = true
	}
	return s
}

func intersection(s []string, t []string) []string {
	u := make([]string, len(s))
	for _, sElem := range s {
		for _, tElem := range t {
			if sElem == tElem {
				u = append(u, sElem)
			}
		}
	}
	return u
}

/* TODO uses sets not lists
func intersection2(s map[string]bool, t map[string]bool) map[string]bool {
	u := make(map[string]bool)
	for elem, sContains := range s {
		if sContains && t[elem] {
			u[elem]=true
		}
	}
	return u
}*/