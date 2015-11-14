# Gobot

Gobot is a dead simple run and forget build daemon for your golang projects.
Simply run gobot from your project's topmost directory and gobot will track changes in your **.go* files.

Each change will trigger a rebuild and call golint (if available) on the changed files.

Gobot is a 100% pure go solution.

![screenshot](https://github.com/fredjeck/gobot/raw/master/img/screenshot.png)

## Installation
Assuming your GOPATH/bin is in your PATH, simply get gobot : 
```
go get github.com/fredjeck/gobot
``` 

```
go install
``` 
Then go in your project's directory and run ```gobot``` 
