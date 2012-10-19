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
package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func MakePostHandler(url string, fn func(http.ResponseWriter, []byte)) (string, func(http.ResponseWriter, *http.Request)) {
	return url, func(w http.ResponseWriter, r *http.Request) {
		log.Println(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), r.RemoteAddr) // TODO remove, slow
		post, err := readPost(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fn(w, post)
	}
}

func readPost(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func MakeGetHandler(url string, fn func(http.ResponseWriter, string)) (string, func(http.ResponseWriter, *http.Request)) {
	return url, func(w http.ResponseWriter, r *http.Request) {
		log.Println(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), r.RemoteAddr) // TODO remove, slow
		get := r.URL.Path[len(url):]
		if !titleValidator.MatchString(get) {
			http.NotFound(w, r)
			return
		}
		fn(w, get)
	}
}
