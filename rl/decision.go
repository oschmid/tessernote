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

package rl

type Decision struct {
	Weight     []float64   // actions and their weights
	Discount   float64     // value between 0-1 inclusive for how much previous actions should be discounted
	SetBonus   float64     // bonus when decision is explicit (i.e. user chosen)
	GetBonus   float64     // bonus when decision is implicit (i.e. system chosen)
	Dependence []*Decision `datastore:"-"` // Decisions on which this depends, len(d.Weight) must be the same
}

func (d *Decision) GetAction(last int) int {
	d.reweigh(last, d.GetBonus)
	return d.max()
}

func (d *Decision) reweigh(action int, bonus float64) {
	for a, w := range d.Weight {
		d.Weight[a] = w * d.Discount
	}
	d.Weight[action] += bonus
	for _, p := range d.Dependence {
		p.reweigh(action, bonus)
	}
}

func (d Decision) max() int {
	var max float64 = -1.0
	action := -1
	for a, w := range d.Weight {
		dw := d.dependenceWeight(a)
		if w+dw > max {
			max = w + dw
			action = a
		}
	}
	return action
}

func (d Decision) dependenceWeight(a int) float64 {
	var w float64 = 0.0
	for _, p := range d.Dependence {
		w += p.Weight[a]
	}
	return w
}

func (d *Decision) SetAction(a int) {
	d.reweigh(a, d.SetBonus)
}
