var baseURL = "http://localhost:8080";
var getTagsURL = "/tags/get";

function updateTags(tags) {
    $.post(baseURL+getTagsURL, tags, function(data) {
        var value = ""
        for (var key in eval(data)) {
            if (data.hasOwnProperty(key)) {
                value += (key + " " + data[key] + "</br>")
            }
        }
        $("#tags").html(value)
    })
}

$(document).ready(updateTags("null"))