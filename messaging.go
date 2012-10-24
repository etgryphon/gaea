package main

/*
	Strings for messaging...
*/
const HELP_MESSAGE string = `
Gaea is a tool for managing and creating GAE code written in Golang.

Usage:

	gaea command [arguments]

The commands are:

    init        create an project from scratch
    run         run the GAE dev_appserver.py
    help        printout this help message
    get         fetch a Golang Package for use locally
`

const LOGO = `
#########################################################################
           .oooooo.          .o.       oooooooooooo        .o.
          d8P'  'Y8b        .888.      '888'     '8       .888.
         888               .8"888.      888              .8"888.
         888              .8\' '888.     888oooo8       .8' '888.
         888     ooooo   .88ooo8888.    888     "      .88ooo8888.
         '88.    .88'   .8'     '888.   888       o   .8'     '888.
          'Y8bood8P'   o88o     o8888o o888ooooood8  o88o     o8888o'
#########################################################################
`

const YML_TEMPLATE = `application: {{.}}
version: 1
runtime: go
api_version: go1

handlers:
- url: /stylesheets
  static_dir: public/stylesheets

- url: /(.*\.(gif|png|jpg))
  static_files: public/images/\1
  upload: public/images/(.*\.(gif|png|jpg))

- url: /(.*\.js)
  static_files: public/js/\1
  upload: public/js/(.*\.js)
  
- url: /.*
  script: _go_app
`

const BASE_GAE_APP_TEMPLATE = `package	{{.}}

import (
  "net/http"
  "fmt"
)

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, World!")
}
`
