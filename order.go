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
	"github.com/oschmid/tessernote/rl"
	"regexp"
	"sort"
	"strings"
)

const (
	AlphaAscending      = "aa"
	AlphaDescending     = "ad"
	LastModified        = "lm"
	FirstModified       = "fm"
	LastCreated         = "lc"
	FirstCreated        = "fc"
	alphaAscendingIndex = iota
	alphaDescendingIndex
	lastModifiedIndex
	firstModifiedIndex
	lastCreatedIndex
	firstCreatedIndex
	DefaultOrderIndex    = alphaAscendingIndex
	DefaultOrderDiscount = 0.9
	DefaultOrderSetBonus = 30.0
	DefaultOrderGetBonus = 10.0
)

// Orders matches all note order types.
var Orders = regexp.MustCompile("(" + AlphaAscending + "|" + AlphaDescending + "|" +
	LastModified + "|" + FirstModified + "|" + LastCreated + "|" + FirstCreated + ")")

// Order tracks the preferred order for notes base on selected tags
// and which device is being used to look at them.
type Order struct {
	Tags    map[string]*rl.Decision
	Device  map[string]*rl.Decision
	Default *rl.Decision
	Last    int
}

// NewOrder creates a 
func NewOrder() Order {
	return Order{
		Tags:   make(map[string]*rl.Decision),
		Device: make(map[string]*rl.Decision),
		Default: &rl.Decision{
			Weight:     make([]float64, 6),
			Discount:   DefaultOrderDiscount,
			SetBonus:   DefaultOrderSetBonus,
			GetBonus:   DefaultOrderGetBonus,
			Dependence: []*rl.Decision{},
		},
		Last: DefaultOrderIndex,
	}
}

// Get returns the preferred order for tag on device.
func (order *Order) Get(tag []Tag, device string) int {
	name := groupName(tag)
	decision := order.getTagsDecision(name, device)
	return decision.GetAction(order.Last)
}

// groupName creates a unique name for a group of tags.
func groupName(tag []Tag) string {
	name := Name(tag)
	sort.Strings(name)
	return strings.Join(name, "#")
}

// getTagsDecision lazily creates a tags decision that depends on a
// device decision.
func (order *Order) getTagsDecision(tags, device string) *rl.Decision {
	if _, ok := order.Tags[tags]; !ok {
		order.Tags[tags] = &rl.Decision{
			Weight:     make([]float64, 6),
			Discount:   DefaultOrderDiscount,
			SetBonus:   DefaultOrderSetBonus,
			GetBonus:   DefaultOrderGetBonus,
			Dependence: []*rl.Decision{order.getDeviceDecision(device)},
		}
	}
	return order.Tags[tags]
}

// getDeviceDecision lazily creates a device decision that depends
// on the default decision.
func (order *Order) getDeviceDecision(device string) *rl.Decision {
	if _, ok := order.Device[device]; !ok {
		order.Device[device] = &rl.Decision{
			Weight:     make([]float64, 6),
			Discount:   DefaultOrderDiscount,
			SetBonus:   DefaultOrderSetBonus,
			GetBonus:   DefaultOrderGetBonus,
			Dependence: []*rl.Decision{order.Default},
		}
	}
	return order.Device[device]
}

// Cleanup removes tags decision for groups containing any Tag in tag.
func (order *Order) Cleanup(tag []Tag) {
	name := Name(tag)
	for group, _ := range order.Tags {
		for _, n := range name {
			if strings.Contains(group, n) {
				delete(order.Tags, group)
				break
			}
		}
	}
}
