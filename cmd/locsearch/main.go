package main

import (
	"github.com/surtexx/locsearch/pkg/http/rest"
)

func main() {
	r := rest.Handler()
	r.Run()
}
