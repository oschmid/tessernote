var saveNoteURL = '/note/save'

function focusTextArea() {
    $(this).children('textarea:first').focus();
}

function saveNote() {
    note = new Object();
    note.ID = $(this).attr('noteid')
    note.Body = $(this).attr('value')
    if (note.Body != "") {
        $.post(saveNoteURL, JSON.stringify(note), function(data) {
            location.reload();
        });
    }
}

$(document).ready(function() {
    $('textarea.resize').autosize();
    $('div.note').click(focusTextArea)
    $('#notes textarea').blur(saveNote)
})