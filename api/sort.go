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

package api

import (
	"errors"
	"github.com/oschmid/tessernote"
	"regexp"
)

var sortOrderURL = regexp.MustCompile("\\?sort=" + tessernote.Orders.String() + "$")

// parseSortOrder will parse a URL and return one of AlphaAscending, AlphaDescending, LastModified, FirstModified,
// LastCreated, or FirstCreated. If no sort order is specified, the empty string is returned.
func parseSortOrder(url string) string {
	if sortOrderURL.MatchString(url) {
		return url[len(url)-2:]
	}
	return ""
}
