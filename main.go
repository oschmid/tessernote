/**
 * Created with IntelliJ IDEA.
 * User: oschmid
 * Date: 21/09/12
 * Time: 5:22 PM
 */
package main

import (
	"notes"
	"net/http"
)

func main() {
	http.HandleFunc("/view/", notes.MakeHandler(notes.ViewHandler))
	http.HandleFunc("/edit/", notes.MakeHandler(notes.EditHandler))
	http.HandleFunc("/save/", notes.MakeHandler(notes.SaveHandler))
	http.HandleFunc("/", notes.RootHandler)
	http.ListenAndServe(":8080", nil)
}
