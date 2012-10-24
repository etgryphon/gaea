# Gaea: 
### _Give birth to amazing Google App Engine (GAE) projects for Go_
## Version: *ALPHA*

## The Problem

If you have ever tried to create Golang applications for GAE, you can find it a little unwieldy to use.  I wanted to create an easier faster way to get up and running with GAE.

## The Solution

This is why I created Gaea.  The idea is that it is the one stop tool shop for getting your project off the ground.  The idea is to hide some of the complexity that is in the structure of the Google App Engine format and make it more like a native Golang application development experience

## API

`gaea init <name>`

This function will create completely workable project that will run and can be deployed if you would like.  If you don't put a `<name>`, it will ask you for a name.

`gaea run <path>`

This function runs the `dev_appserver.py` on the `<path>`