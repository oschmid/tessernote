var baseURL = "http://localhost:8080"
var getTagsURL = "/tags/get"
var getTitlesURL = "/titles"

function updateTags(tags) {
    $.post(baseURL+getTagsURL, tags, function(data) {
        var value = ""
        for (var tag in data) {
            if (data.hasOwnProperty(tag)) {
                value += (tag + " " + data[tag] + "</br>")
            }
        }
        $("#tags").html(value)
    }, "json")
}

function updateTitles(tags) {
    $.post(baseURL+getTitlesURL, tags, function(data) {
        var value = ""
        for (var title in data) {
            value += (title + "</br>")
        }
        $("#titles").html(value)
    }, "json")
}

$(document).ready(function() {
    updateTags("null")
    updateTitles("null")
})