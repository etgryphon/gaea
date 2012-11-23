all: local

local:
	@@go build -o ../bin/gaea ./*.go
	@@echo 'build local version'

public: 
	@@go build -o $GOPATH/bin/gaea ./*.go
	@@echo 'build public version'