var baseURL = 'http://localhost:8080'
var getTagsURL = baseURL+'/tags/get'
var getTitlesURL = baseURL+'/titles'
var getNoteURL = baseURL+'/note/get/'
var saveNoteURL = baseURL+'/note/save'
var deleteNoteURL = baseURL+'/note/delete/'

var currentNoteId = ""
var isEditing = false
var addClicked = false
var deleteClicked = false

function getTags(tags, replyHandler) {
    $.post(getTagsURL, JSON.stringify(tags), replyHandler, 'json')
}

// style relevant tags
function updateRelatedTags(relatedTags) {
    count = 0
    $tagCheckboxes = $('input[name="tagCheckbox"]')
    $tagCheckboxes.each(function(index, tag) {
        if (relatedTags[tag.value]) {
            $(this).parent().addClass('relatedTag')
            count++
        } else {
            $(this).parent().removeClass('relatedTag')
        }
    })
    // no related tags means all tags are related
    if (count == 0) {
        $tagCheckboxes.each(function(index, tag) {
            $(this).parent().addClass('relatedTag')
        })
    }
}

// update listed tags
function updateListedTags(tags, allRelated) {
    $('#tags').empty()
    tmp = ''
    for (var tag in tags) {
        if (tags.hasOwnProperty(tag)) {
            tmp += '<div id="tag"'
            if (allRelated) {
                tmp += ' class="relatedTag"'
            }
            tmp += '><input type="checkbox" name="tagCheckbox" value="'+tag+'">'+tag+' ('+tags[tag]+')<br></div>'
        }
    }
    $('#tags').append(tmp)
}

// gets selected tags as JSON array
function getSelectedTags() {
    return $('input[name="tagCheckbox"]:checked').map(function() {
        return $(this).attr('value')
    }).toArray()
}

function setSelectedTags(tags) {
    $('input[name="tagCheckbox"]').each(function(index, element) {
        if ($.inArray(element.value, selectedTags) >= 0) {
            element.checked = true
        }
    })
}

// TODO trigger when clicking label
function onTagClick() {
    // clicking an unrelated tag should clear all other selections
    if (!$(this).parent().hasClass('relatedTag')) {
        deselectUnrelated($(this).attr('value'))
    }
    tags = getSelectedTags()
    getTags(tags, updateRelatedTags)
    updateTitles(tags)
}

function deselectUnrelated(selected) {
    $('input[name="tagCheckbox"]:checked').each(function(index, tag) {
        if (tag.value != selected) {
            tag.checked = false
        }
    })
}

// updates list of titles
function updateTitles(tags) {
    $.post(getTitlesURL, JSON.stringify(tags), function(data) {
        includesCurrentNote = false
        $('#titles').empty()
        tmp = ''
        for (var i = 0; i < data.length; i++) {
            tmp += '<input type="button" name="title" value="'+data[i][0]+'" noteId="'+data[i][1]+'"><br>'
            if (data[i][1]==currentNoteId) {
                includesCurrentNote = true
            }
        }
        $('#titles').append(tmp)
        if (!includesCurrentNote) {
            hideNote()
        }
    }, 'json')
}

function onTitleClick() {
    id = $(this).attr('noteId')
    if (id != currentNoteId) {
        showNote(id)
    }
}

function showNote(id) {
    currentNoteId = id
    $.getJSON(getNoteURL+id, function(note) {
        $('#deleteNote').attr('value', 'Delete').show()
        $('#noteTitle').html(note.Title)
        $('#noteBody').html(format(note.Body))
    })
}

function hideNote() {
    currentNoteId = ""
    $('#deleteNote').hide()
    $('#noteTitle').empty()
    $('#noteBody').empty()
    $('#noteEditor').empty()
}

// replaces hash tags with links to all the notes of that tag
// adds line breaks
function format(body) {
    // TODO implement
    return body
}

function unformat(body) {
    // TODO implement
    return body
}

// pull out an set of hash tags from a body of text
function parseTags(body) {
    // TODO implement
    return {}
}

function onNewNoteClick() {
    addClicked = true
    if (isEditing) {
        text = $('#noteTextArea').val().split('\n', 2)
        title = text[0]
        body = text[1]
        saveNote(title, body, parseTags(body), readyNewNote)
    } else {
        readyNewNote()
    }
}

function readyNewNote() {
    currentNoteId = ""
    $('#deleteNote').attr('value', 'Cancel')
    $('#noteTitle').html('Untitled')
    $('#noteBody').empty()
    $('#noteEditor').empty()
    addSelectedHashTagsToNote()
    startEditing()
}

function saveNote(title, body, tags, reply) {
    tags = JSON.stringify(tags)
    $.post(saveNoteURL, '{"Id":"'+currentNoteId+'","Title":"'+title+'","Body":"'+body+'","Tags":'+tags+'}', reply)
}

function addSelectedHashTagsToNote() {
    tags = getSelectedTags()
    if (tags.length > 0) {
        hashTags = ''
        for (i = 0; i < tags.length; i++) {
            hashTags += '#' + tags[i]
        }
        $('#noteBody').html(hashTags)
    }
}

function onDeleteNoteClick() {
    deleteClicked = true
    if (currentNoteId != "") {
        $.get(deleteNoteURL+currentNoteId, updateTagsAndTitles)
    }
    hideNote()
}

function updateTagsAndTitles() {
    selectedTags = getSelectedTags()
    getTags(null, function(allTags) {
        updateListedTags(allTags)
        setSelectedTags(selectedTags)
        getTags(selectedTags, updateRelatedTags)
        updateTitles(selectedTags)
    })
}

function startEditing() {
    isEditing = true
    $("#notePanel").off('click')

    // change to textarea
    title = $('#noteTitle').text()
    body = unformat($('#noteBody').text())
    $('#noteEditor').append('<textarea id="noteTextArea">'+title+'\n'+body+'</textarea>')
    $('#noteTitle').empty().off('click')
    $('#noteBody').empty().off('click')

    // setup end of edit
    $('#noteTextArea').blur(function() {
        addClicked = false
        deleteClicked = false
        window.setTimeout(stopEditing, 100)
    }).focus()
}

function stopEditing() {
    isEditing = false
    if (!deleteClicked && !addClicked) {
        text = $('#noteTextArea').val().split('\n', 2)
        title = text[0]
        body = text[1]
        saveNote(title, body, parseTags(body), function(id) {
            currentNoteId = id
            $('#deleteNote').attr('value', 'Delete')
            $('#noteTitle').html(title).click(startEditing)
            $('#noteBody').html(body).click(startEditing)
            $('#noteEditor').empty()
            updateTagsAndTitles() // TODO update without server call (already have necessary information)
        })
    }
}

$(document).ready(function() {
    $('#tags').on('click', 'input', onTagClick)
    getTags(null, function(tags) { updateListedTags(tags, true); })
    $('#titles').on('click', 'input', onTitleClick)
    updateTitles(null)
    $('#noteTitle').click(startEditing)
    $('#noteBody').click(startEditing)
    $('#newNote').click(onNewNoteClick)
    $('#deleteNote').click(onDeleteNoteClick)
})