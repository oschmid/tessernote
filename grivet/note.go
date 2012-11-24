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
package grivet

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Note struct {
	ID           string `datastore:"-"` // datastore.Key.Encode()
	Body         string
	Created      time.Time
	LastModified time.Time
	NotebookKeys []*datastore.Key
	context      appengine.Context `datastore:"-"`
}

func (n Note) Notebooks() []Notebook {
	var users []Notebook
	datastore.GetMulti(n.context, n.NotebookKeys, users)
	for _, u := range users {
		u.context = n.context
	}
	return users
}

// TODO get tags by parsing body for hashtags

func (n *Note) SetBody(body string) error {
	n.Body = body
	n.LastModified = time.Now()
	k, err := datastore.DecodeKey(n.ID)
	if err != nil {
		return err
	}
	_, err = datastore.Put(n.context, k, n)
	return err
}
