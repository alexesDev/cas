package main

import (
	"io/ioutil"
	"os"

	"github.com/alexesDev/cas/pkg/cas"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	output, err := cas.ProcessJSON(data)

	if err != nil {
		panic(err)
	}

	os.Stdout.Write(output)
}
