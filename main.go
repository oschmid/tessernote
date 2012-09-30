/*
This file is part of Notes.

Notes is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Notes is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"notes"
	"net/http"
)

func main() {
	http.HandleFunc("/view/", notes.MakeHandler(notes.ViewHandler))
	http.HandleFunc("/edit/", notes.MakeHandler(notes.EditHandler))
	http.HandleFunc("/save/", notes.MakeHandler(notes.SaveHandler))
	http.HandleFunc("/", notes.RootHandler)
	http.ListenAndServe(":8080", nil)
}
