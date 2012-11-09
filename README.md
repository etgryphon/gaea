# Gaea: 
### _Give birth to amazing Google App Engine (GAE) projects for Go_

## Status
This is very much a dev-preview and is in __ALPHA__.  I haven't tested this extensively, just want to see if it is helpful to other people.  I am still currently working on this. It assumes a Linux file system...will abstract to Windows soon...

## The Problem

If you have ever tried to create Golang applications for GAE, you can find it a little unwieldy to use.  I wanted to create an easier faster way to get up and running with GAE.

## The Solution

This is why I created Gaea.  The idea is that it is the one stop tool shop for getting your project off the ground.  The idea is to hide some of the complexity that is in the structure of the Google App Engine format and make it more like a native Golang application development experience

## API

`gaea init <name>`

This function will create completely workable project that will run and can be deployed if you would like.  If you don't put a `<name>`, it will ask you for a name.

`gaea run <path>`

This function runs the `dev_appserver.py` on the `<path>`
  
`gaea get <package_name>`

This function works just like `go get`, but it creates a local package for use in GAE projects.  This is one of the most frustrating parts of writing GAE projects with Golang.  This is supposed to make it a little more transparent to the user and make it much less manual intensive.

`gaea help`

A help message for the user...

```plaintext
#########################################################################
           .oooooo.          .o.       oooooooooooo        .o.
          d8P'  'Y8b        .888.      '888'     '8       .888.
         888               .8"888.      888              .8"888.
         888              .8\' '888.     888oooo8       .8' '888.
         888     ooooo   .88ooo8888.    888     "      .88ooo8888.
         '88.    .88'   .8'     '888.   888       o   .8'     '888.
          'Y8bood8P'   o88o     o8888o o888ooooood8  o88o     o8888o'
#########################################################################


Gaea is a tool for managing and creating GAE code written in Golang.

Usage:

	gaea command [arguments]

The commands are:

    init        create an project from scratch
    run         run the GAE dev_appserver.py
    help        printout this help message
    get         fetch a Golang Package for use locally

```

## Installation

`go get github.com/etgryphon/gaea`

This should install it into your `$GOPATH` and now you should be able to run it from the command line...

## License: FreeBSD
```plaintext
Copyright (c) 2012, Evin T. Grano
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met: 

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer. 
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution. 

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies, 
either expressed or implied, of the FreeBSD Project.
```