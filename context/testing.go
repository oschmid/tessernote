// +build !appengine

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
	"github.com/oschmid/appenginetesting"
	"net/http"
	"html/template"
)

type Context appengine.Context

var cntx appengine.Context
var Templates = template.Must(template.ParseFiles("templates/main.html"))

// New returns a new testing context.
func NewContext(r *http.Request) appengine.Context {
	if cntx == nil {
		var err error
		cntx, err = appenginetesting.NewContext(nil)
		if err != nil {
			panic(err)
		}
	}
	return cntx
}

// Close closes a testing context registered when New() is called.
func Close() {
	if cntx != nil {
		cntx.(*appenginetesting.Context).Close()
		cntx = nil
	}
}
