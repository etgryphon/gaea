# Gaea: 
### _Give birth to amazing Google App Engine (GAE) projects for Go_
_*Version: ALPHA*_

## The Problem

If you have ever tried to create Golang applications for GAE, you can find it a little unwieldy to use.  I wanted to create an easier faster way to get up and running with GAE.

## The Solution

This is why I created Gaea.  The idea is that it is the one stop tool shop for getting your project off the ground.  The idea is to hide some of the complexity that is in the structure of the Google App Engine format and make it more like a native Golang application development experience

## API

`gaea init <name>`

This function will create completely workable project that will run and can be deployed if you would like.  If you don't put a `<name>`, it will ask you for a name.

`gaea run <path>`

This function runs the `dev_appserver.py` on the `<path>`
  
## License: (BSD New...)
```plaintext
Copyright (c) 2012, Evin T. Grano
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of Evin T. Grano nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```