package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alexesDev/cas/pkg/cas"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output, err := cas.ProcessJSON(data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
