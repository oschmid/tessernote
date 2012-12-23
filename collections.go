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

package tessernote

import (
	"appengine/datastore"
)

// addKey appends add to keys if add doesn't already exist in keys
func addKey(keys []*datastore.Key, add *datastore.Key) []*datastore.Key {
	if !containsKey(keys, add) {
		keys = append(keys, add)
	}
	return keys
}

// removeKey removes remove from keys
func removeKey(keys []*datastore.Key, remove *datastore.Key) []*datastore.Key {
	i := indexOfKey(keys, remove)
	if i >= 0 {
		copy(keys[i:], keys[i + 1:])
		keys[len(keys) - 1] = nil
		return keys[:len(keys) - 1]
	}
	return keys
}

func containsKey(keys []*datastore.Key, key *datastore.Key) bool {
	return indexOfKey(keys, key) >= 0
}

func indexOfKey(keys []*datastore.Key, key *datastore.Key) int {
	for i := range keys {
		if keys[i].Encode() == key.Encode() {
			return i
		}
	}
	return -1
}

func containsString(strings []string, s string) bool {
	for _, elem := range strings {
		if elem == s {
			return true
		}
	}
	return false
}

func indexOfTag(tags []Tag, name string) int {
	for i, tag := range tags {
		if tag.Name == name {
			return i
		}
	}
	return -1
}

func unionKeys(a, b []*datastore.Key) []*datastore.Key {
	c := *new([]*datastore.Key)
	for _, elem := range a {
		j := indexOfKey(b, elem)
		if j >= 0 {
			c = append(c, b[j])
		}
	}
	return c
}
