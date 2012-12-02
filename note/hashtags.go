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

package note

import (
	"regexp"
	"strings"
	"unicode"
)

// hashtag regexp pattern from https://github.com/twitter/twitter-text-java/blob/master/src/com/twitter/Regex.java
const (
	latinAccentsChars = "\u00c0-\u00d6\u00d8-\u00f6\u00f8-\u00ff" + // Latin-1
		"\u0100-\u024f" + // Latin Extended A and B
		"\u0253\u0254\u0256\u0257\u0259\u025b\u0263\u0268\u026f\u0272\u0289\u028b" + // IPA Extensions
		"\u02bb" + // Hawaiian
		"\u0300-\u036f" + // Combining diacritics
		"\u1e00-\u1eff" // Latin Extended Additional (mostly for Vietnamese)
	hashtagAlphaChars = "a-z" + latinAccentsChars +
		"\u0400-\u04ff\u0500-\u0527" + // Cyrillic
		"\u2de0-\u2dff\ua640-\ua69f" + // Cyrillic Extended A/B
		"\u0591-\u05bf\u05c1-\u05c2\u05c4-\u05c5\u05c7" +
		"\u05d0-\u05ea\u05f0-\u05f4" + // Hebrew
		"\ufb1d-\ufb28\ufb2a-\ufb36\ufb38-\ufb3c\ufb3e\ufb40-\ufb41" +
		"\ufb43-\ufb44\ufb46-\ufb4f" + // Hebrew Pres. Forms
		"\u0610-\u061a\u0620-\u065f\u066e-\u06d3\u06d5-\u06dc" +
		"\u06de-\u06e8\u06ea-\u06ef\u06fa-\u06fc\u06ff" + // Arabic
		"\u0750-\u077f\u08a0\u08a2-\u08ac\u08e4-\u08fe" + // Arabic Supplement and Extended A
		"\ufb50-\ufbb1\ufbd3-\ufd3d\ufd50-\ufd8f\ufd92-\ufdc7\ufdf0-\ufdfb" + // Pres. Forms A
		"\ufe70-\ufe74\ufe76-\ufefc" + // Pres. Forms B
		"\u200c" + // Zero-Width Non-Joiner
		"\u0e01-\u0e3a\u0e40-\u0e4e" + // Thai
		"\u1100-\u11ff\u3130-\u3185\uA960-\uA97F\uAC00-\uD7AF\uD7B0-\uD7FF" + // Hangul (Korean)
		"\\p{Hiragana}\\p{Katakana}" + // Japanese Hiragana and Katakana
		// TODO "\\p{Unified_Ideograph}" + // Japanese Kanji / Chinese Han (could also be CJK, CJK_Unified_Ideograph, Unified_Ideograph)
		"\u3003\u3005\u303b" + // Kanji/Han iteration marks
		"\uff21-\uff3a\uff41-\uff5a" + // full width Alphabet
		"\uff66-\uff9f" + // half width Katakana
		"\uffa1-\uffdc" // half width Hangul (Korean)
	hashtagAlphaNumericChars = "0-9\uff10-\uff19_" + hashtagAlphaChars
	hashtagAlphaNumeric      = "[" + hashtagAlphaNumericChars + "]"
	hashtagAlpha             = "[" + hashtagAlphaChars + "]"
	hashtag                  = "(^|[^&" + hashtagAlphaNumericChars + "])(#|\uFF03)(" + hashtagAlphaNumeric + "*" + hashtagAlpha + hashtagAlphaNumeric + "*)"
)

var Hashtag = regexp.MustCompile(hashtag)

func ParseTagNames(text string) []string {
	var names []string
	matches := Hashtag.FindAllString(text, len(text))
	for _, match := range matches {
		name := strings.TrimFunc(match, isHashtagDecoration)
		names = append(names, name)
	}
	return names
}

func isHashtagDecoration(r rune) bool {
	return r == '#' || r == '\uFF03' || unicode.IsSpace(r)
}
