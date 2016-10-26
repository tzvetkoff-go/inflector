# The Inflector transforms words in various ways

[![GoDoc](https://godoc.org/github.com/go2c/inflector?status.svg)](http://godoc.org/github.com/go2c/inflector)
[![Build Status](https://travis-ci.org/go2c/inflector.svg?branch=master)](https://travis-ci.org/go2c/inflector)

A Go port of the [Rails](http://rubyonrails) Inflector.

## Examples

``` go
package main

import (
	"github.com/go2c/inflector"
)

func main() {
	println(inflector.Pluralize("person"))      // "people"
	println(inflector.Singularize("men"))       // "man"
}
```

## License

The code is subject to the [MIT license](https://opensource.org/licenses/MIT).
