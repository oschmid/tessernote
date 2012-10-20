var baseURL = "http://localhost:8080";
var getTagsURL = "/tags/get";

$(document).ready(function() {
    $.post(baseURL+getTagsURL, "null", function(data) {
        alert(data);
    });
});