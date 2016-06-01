# DojoBuilder
DojoBuilder is a tool to automate the build of [dojo] (https://dojotoolkit.org/) based projects in Golang.

# Installation
Make sur you have a working Go environment.

To install DojoBuilder, run:
```
$ go get github.com/tbaud0n/dojoBuilder
```
# Getting started
To run DojoBuilder, simply call dojoBuilder.Run(c *Config, names[]string, reset bool).

This method needs : 
- A dojoBuilder.Config which contains all the needed information to run the build or install files for non-built mode (see code for more details).
- An optional array of build name to execute (for buildMode). If nil, all the build configs will be executed.
- A boolean indicating if the destination folder has to be emptied. (The destination folder has to be emptied when switching between non-built and build mode)


# Example
An example is provided in the example folder.

To test the example, simply run the initExample.sh which will download go in it's client folder.
Then run:
```
$ go build main.go
```
To execute the example in non-built mode, simply run:
```
$ ./main
```
To execute the example in build mode, simply run:
```
$ ./main --buildMode
```

See the result going to [http://127.0.0.1:8080](http://127.0.0.1:8080)
