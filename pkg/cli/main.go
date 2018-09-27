package main

import (
	"fmt"
	"os"

	"github.com/alexesDev/cas/pkg/cas"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	scale, err := cas.Connect(os.Getenv("CAS_ADDR"))

	if err != nil {
		panic(err)
	}

	// erasePLU(conn, 0, 0)

	for i := 10; i >= 0; i -= 1 {
		name, _ := charmap.Windows1251.NewEncoder().String(fmt.Sprintf("button %d", i))
		var data cas.PLUData
		copy(data.PLUName1[:], name)
		data.DepartmentNumber = 1
		data.PLUNumber = uint32(i)
		data.PLUType = 1
		scale.DownloadPLU(0, data)
	}

	dec := charmap.Windows1251.NewDecoder()

	for i := 1; i < 10; i += 1 {
		data := scale.UploadPLU(0, uint32(i))
		out, _ := dec.Bytes(data.PLUName1[0:])
		fmt.Printf("%d %s %d %+v \n", i, string(out), data.PLUNumber, data)
	}
}
