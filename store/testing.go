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

package store

import (
	"appengine"
	"appengine/datastore"
)

// TODO fake datastore interaction for testing

func Get(c appengine.Context, key *datastore.Key, dst interface{}) error {
	return datastore.Get(c, key, dst)
}

func GetMulti(c appengine.Context, keys []*datastore.Key, dst interface{}) error {
	return datastore.GetMulti(c, keys, dst)
}

func Put(c appengine.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	return datastore.Put(c, key, src)
}

func PutMulti(c appengine.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	return datastore.PutMulti(c, keys, src)
}

func Delete(c appengine.Context, key *datastore.Key) error {
	return datastore.Delete(c, key)
}

func DeleteMulti(c appengine.Context, keys []*datastore.Key) error {
	return datastore.DeleteMulti(c, keys)
}
