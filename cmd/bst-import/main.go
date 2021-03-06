package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dxw/beanstalk-tools"
	"github.com/kr/beanstalk"
)

const (
	defaultPriority = (2 ^ 32) - 1
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: bst-import hostname:11300 tubename")
	}

	host := os.Args[1]
	tubeName := os.Args[2]

	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	tube := &beanstalk.Tube{conn, tubeName}

	decoder := json.NewDecoder(os.Stdin)

	for decoder.More() {
		var item common.Item
		if err := decoder.Decode(&item); err != nil {
			log.Fatal(err)
		}

		_, err = tube.Put([]byte(item.Content), defaultPriority, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}
