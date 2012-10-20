package main

/*
	Strings for messaging...
*/
var helpMessage string = 
`
This is a test
With lots of Lines
To print out...
`

var logo string = 
`
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

var ymlTemplate string = 
`
application: %s
version: 1
runtime: go
api_version: go1

handlers:
- url: /.*
  script: _go_app
 `