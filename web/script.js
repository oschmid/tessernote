var baseURL = "http://localhost:8080"
var getTagsURL = "/tags/get"
var getTitlesURL = "/titles"
var getNoteURL = "/note/get/"

var currentNoteId = ""

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
    tag.checked = "true"
    tag.onclick = onTagClick
    return tag
}

function onTagClick(event) {
    // TODO updateTags
    // TODO updateTitles
    // TODO deselect note if no longer in titles
    alert("tag clicked")
}

function updateTitles(tags) {
    $.post(baseURL+getTitlesURL, tags, function(data) {
        for (var i = 0; i < data.length; i++) {
            $("#titles").append(createTitle(data[i]), "<br>")
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
    // TODO set currentNoteId
    // TODO updateNote(id)
    alert("title clicked")
}

function updateNote(id) {
    $.getJSON(baseURL+getTitlesURL+id, function(data) {
        // TODO display note contents
    })
}

function onNewNoteClick(event) {
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
    $("[name='newNote']").click(onNewNoteClick)
})