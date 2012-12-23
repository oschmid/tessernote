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
	"html/template"
)

type Page struct {
	Tags          []Tag
	Notes         []Note
	UntaggedNotes bool
	relatedTag    map[string]bool
	selectedTag   map[string]bool
}

func (p *Page) SetRelatedTags(tags []Tag) {
	p.relatedTag = make(map[string]bool)
	for _, tag := range tags {
		p.relatedTag[tag.Name] = true
	}
}

func (p *Page) SetSelectedTags(tags []Tag) {
	p.selectedTag = make(map[string]bool)
	for _, tag := range tags {
		p.selectedTag[tag.Name] = true
	}
}

// format tags as html
func (p Page) HtmlTags() template.HTML {
	html := p.createTagDiv("All Notes")
	for _, tag := range p.Tags {
		html += p.createTagDiv(tag.Name)
	}
	if p.UntaggedNotes {
		html += p.createTagDiv("Untagged Notes")
	}
	return template.HTML(html)
}

func (p Page) createTagDiv(name string) string {
	div := "<div class='tag"
	if p.relatedTag[name] {
		div += " related"
	}
	if p.selectedTag[name] {
		div += " selected"
	}
	div += "'>" + name + "</div>"
	return div
}

// format notes as html
func (p Page) HtmlNotes() template.HTML {
	html := ""
	for _, note := range p.Notes {
		html += "<div class='note'><div class='delete'>x</div><textarea noteid=\"" + note.ID + "\" class='resize'>" + note.Body + "</textarea><input type='button' class='save' value='Save'></div>"
	}
	return template.HTML(html)
}
