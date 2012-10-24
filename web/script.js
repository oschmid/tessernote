var baseURL = 'http://localhost:8080'
var getTagsURL = baseURL+'/tags/get'
var getTitlesURL = baseURL+'/titles'
var getNoteURL = baseURL+'/note/get/'
var saveNoteURL = baseURL+'/note/save'

var currentTags = []
var currentNoteId

function getTags(tags, replyHandler) {
    $.post(getTagsURL, tags, replyHandler, 'json')
}

// style narrowing tags
function updateNarrowingTags(narrowingTags) {
    $('[name="tagCheckbox"]').each(function(index, tag) {
        if (narrowingTags[tag.value]) {
            $(this).parent().addClass('narrowingTag')
        } else {
            $(this).parent().removeClass('narrowingTag')
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

function onTagClick() {
    var selectedTags = []
    $('[name="tagCheckbox"]:checked').each(function(index, element) {
        selectedTags[index] = element.value
    })
    tags = JSON.stringify(selectedTags)

    getTags(tags, updateNarrowingTags)
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
            // TODO parse note and replace hashtags with clickable links
            $('#noteBody').html(note.Body)
        })
    }
}

function onNewNoteClick() {
    // TODO add current tags
    $.post(saveNoteURL, '{"Title":"Untitled","Body":"","Tags":{}}', function(id) {
        alert(id)
    })
    // TODO makeNoteEditable()
    // TODO updateTitles()
}

function makeNoteEditable() {
    // TODO convert note display into text area
}

function makeNoteNonEditable() {
    // TODO convert note text area into display
}

$(document).ready(function() {
    getTags('null', updateTags)
    updateTitles('null')
    $("#notePanel").click(makeNoteEditable)
    $("#newNote").click(onNewNoteClick)
})