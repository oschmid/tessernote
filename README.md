#Tessernote
I'm making a note taking web app to teach myself web programming and to play around with different algorithms.

Tessernote is:
- minimal: notes are just text, organization is done with tags, and the interface is equally simple
- organized on-the-fly: tagging is done by writing hashtags in notes

###How to install
1. install go 1.0.3 ([install guide](http://golang.org/doc/install))
2. download [go-appengine-sdk 1.7.3](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
3. [copy appengine sdk and goprotobuf into go](http://stackoverflow.com/questions/11286534/test-cases-for-go-and-appengine)
4. setup a go workspace with bin/ pkg/ and src/ directories
5. run 'go get github.com/oschmid/tessernote'
6. create a symlink in src/ to app.yaml so you can run it locally using dev_appserver.py

###Fetchnotes
It turns out my idea of organizing notes by hashtag isn't as original as I thought. So if you want a note taking app
that works this way right now give [Fetchnotes](http://www.fetchnotes.com/) a try.
