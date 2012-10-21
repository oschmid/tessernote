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

type StringSliceSlice [][]string

func (p StringSliceSlice) Len() int {
	return len(p)
}

func (p StringSliceSlice) Less(i, j int) bool {
	return p[i][0] < p[j][0]
}

func (p StringSliceSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
