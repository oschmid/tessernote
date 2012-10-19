$(document).ready(function() {
    $.post("localhost:8080/tags/get", "null", function(data) {
        alert(data);
    })
});