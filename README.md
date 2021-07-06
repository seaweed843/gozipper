# gozipper [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
Golang zipper

zip a folder or file 

## Adding dependency
```
import "github.com/seaweed843/gozipper"
```

## Example:

init module:
```
$ go mod init main
```

main.go
```
package main

import "github.com/seaweed843/gozipper"
import "fmt"

func main() {
	err := gozipper.ZipPath("./srcFolder")
  //expected result: ./srcFolder.zip
  
  err = gozipper.ZipPath("./srcFile.ext")
  //expected result: ./srcFile.ext.zip
  
  err = gozipper.ZipPath("./srcFolder", "./dstFolder", "dstFileName.zip")
  //expected result: ./dstFolder/dstFileName.zip
  
  if err != nil {
	  fmt.Println(err)
  }

}
```

run:
```
$ go mod tidy
$ go run main.go
```
