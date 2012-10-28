var baseURL = 'http://localhost:8080'
var getTagsURL = baseURL+'/tags/get'
var getTitlesURL = baseURL+'/titles'
var getNoteURL = baseURL+'/note/get/'
var saveNoteURL = baseURL+'/note/save'

// hashtag regexp pattern from https://github.com/twitter/twitter-text-java/blob/master/src/com/twitter/Regex.java
var latinAccentsChars = "\\u00c0-\\u00d6\\u00d8-\\u00f6\\u00f8-\\u00ff" + // Latin-1
                        "\\u0100-\\u024f" + // Latin Extended A and B
                        "\\u0253\\u0254\\u0256\\u0257\\u0259\\u025b\\u0263\\u0268\\u026f\\u0272\\u0289\\u028b" + // IPA Extensions
                        "\\u02bb" + // Hawaiian
                        "\\u0300-\\u036f" + // Combining diacritics
                        "\\u1e00-\\u1eff"; // Latin Extended Additional (mostly for Vietnamese)
var hashTagAlphaChars = "a-z" + latinAccentsChars +
                        "\\u0400-\\u04ff\\u0500-\\u0527" +  // Cyrillic
                        "\\u2de0-\\u2dff\\ua640-\\ua69f" +  // Cyrillic Extended A/B
                        "\\u0591-\\u05bf\\u05c1-\\u05c2\\u05c4-\\u05c5\\u05c7" +
                        "\\u05d0-\\u05ea\\u05f0-\\u05f4" + // Hebrew
                        "\\ufb1d-\\ufb28\\ufb2a-\\ufb36\\ufb38-\\ufb3c\\ufb3e\\ufb40-\\ufb41" +
                        "\\ufb43-\\ufb44\\ufb46-\\ufb4f" + // Hebrew Pres. Forms
                        "\\u0610-\\u061a\\u0620-\\u065f\\u066e-\\u06d3\\u06d5-\\u06dc" +
                        "\\u06de-\\u06e8\\u06ea-\\u06ef\\u06fa-\\u06fc\\u06ff" + // Arabic
                        "\\u0750-\\u077f\\u08a0\\u08a2-\\u08ac\\u08e4-\\u08fe" + // Arabic Supplement and Extended A
                        "\\ufb50-\\ufbb1\\ufbd3-\\ufd3d\\ufd50-\\ufd8f\\ufd92-\\ufdc7\\ufdf0-\\ufdfb" + // Pres. Forms A
                        "\\ufe70-\\ufe74\\ufe76-\\ufefc" + // Pres. Forms B
                        "\\u200c" +                        // Zero-Width Non-Joiner
                        "\\u0e01-\\u0e3a\\u0e40-\\u0e4e" + // Thai
                        "\\u1100-\\u11ff\\u3130-\\u3185\\uA960-\\uA97F\\uAC00-\\uD7AF\\uD7B0-\\uD7FF" + // Hangul (Korean)
                        "\\p{InHiragana}\\p{InKatakana}" +  // Japanese Hiragana and Katakana
                        "\\p{InCJKUnifiedIdeographs}" +     // Japanese Kanji / Chinese Han
                        "\\u3003\\u3005\\u303b" +           // Kanji/Han iteration marks
                        "\\uff21-\\uff3a\\uff41-\\uff5a" +  // full width Alphabet
                        "\\uff66-\\uff9f" +                 // half width Katakana
                        "\\uffa1-\\uffdc";                  // half width Hangul (Korean)
var hashTagAlphaNumericChars = "0-9\\uff10-\\uff19_" + hashTagAlphaChars
var hashTagAlphaNumeric = "[" + hashTagAlphaNumericChars + "]"
var hashTagAlpha = "[" + hashTagAlphaChars + "]"
var validHashTag = "(^|[^&" + hashTagAlphaNumericChars + "])(#|\uFF03)(" + hashTagAlphaNumeric + "*" + hashTagAlpha + hashTagAlphaNumeric + "*)"

var currentTags = []
var currentNoteId // TODO make attribute of notePanel

function getTags(tags, replyHandler) {
    $.post(getTagsURL, tags, replyHandler, 'json')
}

// style relevant tags
function updateRelatedTags(relatedTags) {
    $('[name="tagCheckbox"]').each(function(index, tag) {
        if (relatedTags[tag.value]) {
            $(this).parent().addClass('relatedTag')
        } else {
            $(this).parent().removeClass('relatedTag')
        }
    })
}

// update listed tags
function updateTags(tags) {
    $('#tags').empty()
    for (var tag in tags) {
        if (tags.hasOwnProperty(tag)) {
            $('#tags').append('<div id="tag"><input type="checkbox" name="tagCheckbox" value="'+tag+'" onclick="onTagClick()">'+tag+' ('+tags[tag]+")<br></div>")
        }
    }
}

function getSelectedTags() {
    var selectedTags = []
    $('[name="tagCheckbox"]:checked').each(function(index, element) {
        selectedTags[index] = element.value
    })
    return JSON.stringify(selectedTags)
}

function onTagClick() {
    tags = getSelectedTags()
    getTags(tags, updateRelatedTags)
    updateTitles(tags)
}

function updateTitles(tags) {
    $.post(getTitlesURL, tags, function(data) {
        $('#titles').empty()
        for (var i = 0; i < data.length; i++) {
            $('#titles').append('<input type="button" name="title" value="'+data[i][0]+'" noteId="'+data[i][1]+'"><br>')
        }

        $('input[name="title"]').click(function() {
            updateNote($(this).attr('noteId'))
        })

        if (!noteInNotes()) {
            updateNote()
        }
    }, 'json')
}

// returns true if note being displayed is in the list of titles displayed
function noteInNotes() {
    var contained = false
    $('[name="title"]').each(function(index, element) {
        if (element.noteId == currentNoteId) {
            contained = true
            return
        }
    })
    return contained
}

function updateNote(id) {
    currentNoteId = id
    if (!id) {
        $('#noteTitle').empty()
        $('#noteBody').empty()
    } else {
        $.getJSON(getNoteURL+id, function(note) {
            $('#noteTitle').html(note.Title)
            $('#noteBody').html(linkHashTags(note.Body))
        })
    }
}

// replaces hash tags with links to all the notes of that tag
function linkHashTags(body) {
    // TODO
    return body
}

function unlinkHashTags(body) {
    // TODO
    return body
}

function onNewNoteClick() {
    tags = getSelectedTags()
    $.post(saveNoteURL, '{"Title":"Untitled","Body":"","Tags":'+tags+'}', function(id) {
        alert(id)
    })
    updateTitles(tags)
    startEditing()
}

function startEditing() {
    // remove click handler
    $("#notePanel").off('click')

    // change to textarea
    title = $('#noteTitle').text()
    body = unlinkHashTags($('#noteBody').text())
    $('#noteEditor').append('<textarea id="noteTextArea">'+title+'\n'+body+'</textarea>')
    $('#noteTitle').empty()
    $('#noteBody').empty()

    // setup end of edit
    $('#noteTextArea').blur(stopEditing)
    $('#noteTextArea').focus()
}

function stopEditing() {
    text = $('#noteTextArea').val().split('\n', 2)
    title = text[0]
    body = text[1]
    tags = '' // TODO parse tags?
    $.post(saveNoteURL, '{"Id":"'+currentNoteId+'","Title":"'+title+'","Body":"'+body+'","Tags":{'+tags+'}}', function(id) {
        $('#noteTitle').html(title)
        $('#noteBody').html(body)
        $('#noteEditor').empty()
        $('#notePanel').click(startEditing)
    })
}

$(document).ready(function() {
    getTags('null', updateTags)
    getTags('null', updateRelatedTags)
    updateTitles('null')
    $('#notePanel').click(startEditing)
    $('#newNote').click(onNewNoteClick)
})