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
	"errors"
	"regexp"
	"sort"
	"strings"
)

const (
	AlphaAscending  = "aa"
	AlphaDescending = "ad"
	LastModified    = "lm"
	FirstModified   = "fm"
	LastCreated     = "lc"
	FirstCreated    = "fc"
	DefaultOrder    = AlphaAscending
	PastOrderWeight = 0.9
	OrderSetBonus   = 50
	OrderGetBonus   = 10
)

// SortOrder matches all legal methods of ordering notes
var SortOrder = regexp.MustCompile("(" + AlphaAscending + "|" + AlphaDescending + "|" +
	LastModified + "|" + FirstModified + "|" + LastCreated + "|" + FirstCreated + ")")

// NoteOrder tracks the preferred sort order of the notes of a set of tags.
// Preference is based on frecency and whether is set explicitly or accepted
// implicitly.
type NoteOrder struct {
	Weight map[string]*orderWeight
	Last   string
}

// Get returns the preferred order for the notes of a set of tags.
func (order *NoteOrder) Get(tag []Tag) string {
	name := groupName(tag)
	if _, ok := order.Weight[name]; !ok {
		order.Weight[name] = newOrderWeight()
	}
	order.Last = order.Weight[name].Get(order.Last)
	return order.Last
}

// groupName creates a unique name for a group of tags.
func groupName(tag []Tag) string {
	name := Name(tag)
	sort.Strings(name)
	return strings.Join(name, "#")
}

// Set updates the order preferences for the notes of a set of tags.
func (order *NoteOrder) Set(tag []Tag, o string) error {
	if !SortOrder.MatchString(o) {
		return errors.New("invalid sort order: " + o)
	}
	name := groupName(tag)
	if _, ok := order.Weight[name]; !ok {
		order.Weight[name] = newOrderWeight()
	}
	order.Weight[name].Set(o)
	return nil
}

// Cleanup removes order weights referring to any Tag t in tag.
func (order *NoteOrder) Cleanup(tag []Tag) {
	name := Name(tag)
	for on, _ := range order.Weight {
		for _, n := range name {
			if strings.Contains(on, n) {
				delete(order.Weight, on)
				break
			}
		} 
	}
}

// NewNoteOrder returns a NoteOrder with default order.
func NewNoteOrder() NoteOrder {
	return NoteOrder{
		Weight: make(map[string]*orderWeight),
		Last:   DefaultOrder,
	}
}

// orderWeight tracks the weights of different note orders.
type orderWeight map[string]float32

// Get reweighs all order weights and returns an order based on frecency of past orders.
func (weight orderWeight) Get(lastOrder string) string {
	weight.reweigh()
	weight[lastOrder] += OrderGetBonus
	return weight.Max()
}

// Max returns the order with the greatest weight. 
func (weight orderWeight) Max() string {
	var max float32 = 0.0
	order := ""
	for o, w := range weight {
		if w > max {
			max = w
			order = o
		}
	}
	return order
}

// reweigh updates the weights of previous orders.
func (weight orderWeight) reweigh() {
	for o, w := range weight {
		weight[o] = w * PastOrderWeight
	}
}

// Set reweighs all order weights and sets the current order.
func (weight orderWeight) Set(order string) {
	weight.reweigh()
	weight[order] += OrderSetBonus
}

// newOrderWeight returns an orderWeight with weight given to the default order.
func newOrderWeight() *orderWeight {
	order := make(orderWeight)
	order[DefaultOrder] = OrderGetBonus
	return &order
}
