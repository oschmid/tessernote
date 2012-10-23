var baseURL = "http://localhost:8080"
var getTagsURL = "/tags/get"
var getTitlesURL = "/titles"
var getNoteURL = "/note/get/"

var currentTags = []
var currentNoteId

function updateTags(tags) {
    $.post(baseURL+getTagsURL, tags, function(data) {
        for (var tag in data) {
            if (data.hasOwnProperty(tag)) {
                $("#tags").append(createTag(tag), tag+" ("+data[tag]+")<br>")
            }
        }
    }, "json")
}

function createTag(name) {
    var tag = document.createElement("input")
    tag.type = "checkbox"
    tag.name = "tag"
    tag.value = name
    tag.checked = false
    tag.onclick = onTagClick
    return tag
}

function onTagClick(event) {
    var selectedTags = []
    $('[name="tag"]:checked').each(function(index, element) {
        selectedTags[index] = element.value
    })
    // TODO colour tags differently depending on if they will further narrow the selection
    updateTitles(JSON.stringify(selectedTags))
}

function updateTitles(tags) {
    $.post(baseURL+getTitlesURL, tags, function(data) {
        $("#titles").empty()
        for (var i = 0; i < data.length; i++) {
            $("#titles").append(createTitle(data[i]), "<br>")
        }

        if (!noteInNotes()) {
            updateNote()
        }
    }, "json")
}

function createTitle(info) {
    var title = document.createElement("input")
    title.type = "button"
    title.name = "title"
    title.value = info[0]
    title.noteId = info[1]
    title.onclick = onTitleClick
    return title
}

function onTitleClick(event) {
    updateNote(event.target.noteId)
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
        $("#noteTitle").empty()
        $("#noteBody").empty()
    } else {
        $.getJSON(baseURL+getNoteURL+id, function(data) {
            $("#noteTitle").html(data.Title)
            $("#noteBody").html(data.Body)
            // TODO parse note and replace hashtags with clickable links
        })
    }
}

function onNewNoteClick() {
    // TODO create new note
    // TODO makeNoteEditable()
    alert("new note clicked")
}

function makeNoteEditable() {
    // TODO convert note display into text area
}

function makeNoteNonEditable() {
    // TODO convert note text area into display
}

$(document).ready(function() {
    updateTags("null")
    updateTitles("null")
    $("#notePanel").click(makeNoteEditable)
    $("#newNote").click(onNewNoteClick)
})