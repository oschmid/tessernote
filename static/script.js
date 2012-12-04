var saveNoteURL = '/note/save'

function filterByTag(e) {
    if ($(this).text() == 'All Notes') {
        location.pathname = '/'
    } else if ($(this).text() == 'Untagged Notes') {
        location.pathname = '/untagged/'
    } else if (e.shiftKey && location.pathname != '/') {
        location.pathname += ',' + $(this).text()
    } else {
        location.pathname = '/' + $(this).text()
    }
}

function startEdit() {
    $(this).children('textarea:first').focus();
    $(this).children('input:first').show();
}

function saveNote() {
    textarea = $(this).prev('textarea')
    note = new Object();
    note.ID = textarea.attr('noteid')
    note.Body = textarea.attr('value')
    if (note.Body != '') {
        $.post(saveNoteURL, JSON.stringify(note), function(data) {
            // TODO insert note without reloading page
            location.reload();
        });
    }
    $(this).hide();
}

$(document).ready(function() {
    $('textarea.resize').autosize();
    $('div.note').click(startEdit);
    $('input.save').click(saveNote);
    $('div.tag').click(filterByTag);
})