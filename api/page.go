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

package api

import (
	"html/template"
	"note"
)

type Page struct {
	Tags         []note.Tag
	RelatedTags  []note.Tag
	SelectedTags []note.Tag
	Notes        []note.Note
	Edit         bool
}

// format tags as html
func (p Page) HtmlTags() template.HTML {
	html := "<div class='tag'>All Notes</div>"
	for _, tag := range p.Tags {
		html += "<div class='tag'>" + tag.Name + "</div>"
	}
	// TODO untagged notes
	return template.HTML(html)
}

// format notes as html
func (p Page) HtmlNotes() template.HTML {
	html := ""
	for _, note := range p.Notes {
		html += "<div class='note'><textarea noteid=\"" + note.ID + "\" class='resize'>" + note.Body + "</textarea></div>"
	}
	return template.HTML(html)
}
