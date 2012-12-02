var saveNoteURL = '/note/save'

function filterByTag(e) {
    if (e.shiftKey && location.pathname != '/') {
        location.pathname += ',' + $(this).text()
    } else {
        location.pathname = "/" + $(this).text()
    }
}

function focusTextArea() {
    $(this).children('textarea:first').focus();
}

function saveNote() {
    note = new Object();
    note.ID = $(this).attr('noteid')
    note.Body = $(this).attr('value')
    if (note.Body != "") {
        $.post(saveNoteURL, JSON.stringify(note), function(data) {
            // TODO insert note without reloading page
            location.reload();
        });
    }
}

$(document).ready(function() {
    $('textarea.resize').autosize();
    $('div.note').click(focusTextArea)
    $('#notes textarea').blur(saveNote)
    $('div.tag').click(filterByTag)
})