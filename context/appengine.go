// +build appengine

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

package context

import (
	"appengine"
	"net/http"
	"html/template"
)

var Templates = template.Must(template.ParseFiles("github.com/oschmid/tessernote/api/templates/main.html"))

type Context appengine.Context

// NewContext returns a new context for an in-flight HTTP request.
func NewContext(r *http.Request) appengine.Context {
	return appengine.NewContext(r)
}

func Close() {
	// do nothing
}
