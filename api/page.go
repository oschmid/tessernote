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
	"grivet"
	"html/template"
	"strings"
)

type Page struct {
	Tags         []grivet.Tag
	RelatedTags  []grivet.Tag
	SelectedTags []grivet.Tag
	Notes        []grivet.Note
	Note         grivet.Note
	Edit         bool
}

// format tags as html
func (p Page) HtmlTags() template.HTML {
	html := ""
	for _, tag := range p.Tags {
		html += tag.Name + "<br>"
	}
	return template.HTML(html)
}

// format titles as html
func (p Page) HtmlTitles() template.HTML {
	html := ""
	for _, note := range p.Notes {
		html += note.Title + "<br>"
	}
	return template.HTML(html)
}

// format note
func (p Page) HtmlNote() template.HTML {
	html := ""
	if p.Edit {
		html = "<form action='" + p.url(saveSuffix) + "' method='POST'><input type='submit' value='Save' class='button'><div class='textwrap'><textarea id='noteTextArea' name='note'>" + p.Note.Title + "\n" + p.Note.Body + "</textarea></div></form>"
	} else {
		if p.Note.ID == nil {
			html = "<form action='" + p.url(newSuffix) + "' method='POST'><input type='submit' value='New' class='button'></form>"
		} else {
			html = "<form action='" + p.url(editSuffix) + "' method='POST'><input type='submit' value='Edit' class='button'></form><b>" + p.Note.Title + "</b><br>" + p.Note.Body
		}
	}
	return template.HTML(html)
}

func (p Page) url(suffix string) string {
	url := ""
	if len(p.SelectedTags) > 0 {
		url += "/" + p.tagString()
	}
	if p.Note.ID != nil {
		url += "/" + p.Note.ID.Encode()
	}
	return url + suffix
}

func (p Page) tagString() string {
	names := grivet.TagNames(p.SelectedTags)
	return strings.Join(names, tagSeparator)
}
