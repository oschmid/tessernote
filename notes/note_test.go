/**
 * Created with IntelliJ IDEA.
 * User: oschmid
 * Date: 30/09/12
 * Time: 12:47 PM
 * To change this template use File | Settings | File Templates.
 */
package notes

import (
	"testing"
)

func TestTagsAsString(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	note := Note{Title: "title", Body: []byte("body"), Tags: tags}
	stringTags, err := note.TagsAsString()
	if *stringTags != tags[0]+tagSeparator+tags[1] {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
