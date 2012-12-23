var notesURL = '/notes/'

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

function showDelete() {
    $(this).children("div.delete:first").show();
}

function deleteNote(e) {
    if (!e) var e = window.event;
    e.cancelBuble = true;
    if (e.stopPropagation) e.stopPropagation();

    textarea = $(this).next('textarea')
    note = new Object();
    note.ID = textarea.attr('noteid')
    $.ajax({url: notesURL + note.ID, type:'DELETE', success: function(response) {
        // TODO remove note without reloading page
        location.reload();
    }});
}

function hideDelete() {
    $(this).children("div.delete:first").hide();
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
        if (note.ID) {
            $.ajax({url:notesURL+note.ID, type:"PUT", data:JSON.stringify(note)})
        } else {
            $.post(notesURL, JSON.stringify(note), function(note) {
                // TODO insert note without reloading page
                location.reload();
            });
        }
    }
    $(this).hide();
}

$(document).ready(function() {
    $('div.tag').click(filterByTag);
    $('textarea.resize').autosize();
    $('div.note').click(startEdit);
    $('div.note').not('#new').mouseenter(showDelete).mouseleave(hideDelete);
    $('div.delete').click(deleteNote);
    $('input.save').click(saveNote);
})