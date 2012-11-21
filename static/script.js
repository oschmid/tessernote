$(document).ready(function(){
    $("textarea.resize").autosize();
    $("div.note").click(function(){
        $(this).children('textarea:first').focus();
    })
    // TODO on textarea blur save new note
})