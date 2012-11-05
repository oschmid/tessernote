/*
This file is part of Grivet.

Grivet is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Grivet is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Grivet.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"api"
	"flag"
	"log"
	"net/http"
)

func main() {
	p := flag.Bool("p", false, "if true, populates notebook with test data")
	flag.Parse()
	file := "notebook"
	if len(flag.Args()) > 0 {
		file = flag.Arg(0)
	}
	if *p {
		api.PopulateNotebook(file)
	} else {
		api.LoadNotebook(file)
	}

	http.HandleFunc(api.MakePostHandler(api.GetTagsUrl, api.GetTags))
	http.HandleFunc(api.MakePostHandler(api.RenameTagsUrl, api.RenameTags))
	http.HandleFunc(api.MakePostHandler(api.DeleteTagsUrl, api.DeleteTags))
	http.HandleFunc(api.MakePostHandler(api.GetTitlesUrl, api.GetTitles))
	http.HandleFunc(api.MakeGetHandler(api.GetNoteUrl, api.GetNote))
	http.HandleFunc(api.MakePostHandler(api.SaveNoteUrl, api.SaveNote))
	http.HandleFunc(api.MakeGetHandler(api.DeleteNoteUrl, api.DeleteNote))
	http.Handle("/notes/", http.StripPrefix("/notes/", http.FileServer(http.Dir("web"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
