package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/notassigned/tweb/client"
)

func StartClient(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	app, err := client.CreateClient(string(content))
	t := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	app.Start()
	fmt.Println(t)
}
