/*
This file is part of Tessernote.

Tessernote is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Tessernote is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Tessernote.  If not, see <http://www.gnu.org/licenses/>.
*/

package filepath

import "os"

// Merge 2 parts of a filepath so that any overlap is ignored
// e.g. /a/b/c + b/c/d = /a/b/c/d
func Merge(start, end string) string {
	overlap := stringOverlapLength(start, end)
	if overlap == 0 {
		return start + string(os.PathSeparator) + end
	}
	return start + end[overlap:]
}

// Calculates how much the end of s1 overlaps the beginning of s2
// e.g. stringOverlapLength("abcdef", "defghl") = 3
func stringOverlapLength(s1, s2 string) int {
	//Trim s1 so it isn't longer than s2
	if len(s1) > len(s2) {
		s1 = s1[len(s1)-len(s2):]
	}
	T := computeBackTrackTable(s2) //O(n)
	m := 0
	i := 0
	for m+i < len(s1) {
		if s2[i] == s1[m+i] {
			i += 1
		} else {
			m += i - T[i]
			if i > 0 {
				i = T[i]
			}
		}
	}
	return i
}

func computeBackTrackTable(s string) []int {
	T := make([]int, len(s))
	cnd := 0
	T[0] = -1
	T[1] = 0
	pos := 2
	for pos < len(s) {
		if s[pos-1] == s[cnd] {
			T[pos] = cnd + 1
			pos += 1
			cnd += 1
		} else if cnd > 0 {
			cnd = T[cnd]
		} else {
			T[pos] = 0
			pos += 1
		}
	}
	return T
}
