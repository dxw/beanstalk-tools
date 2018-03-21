package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kr/beanstalk"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: bst-export hostname:11300 tubename")
	}

	host := os.Args[1]
	tubeName := os.Args[2]

	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	tubeSet := beanstalk.NewTubeSet(conn, tubeName)

	for {
		id, data, err := tubeSet.Reserve(1 * time.Second)
		if e, ok := err.(beanstalk.ConnError); ok && e.Err == beanstalk.ErrTimeout {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("---")
		fmt.Println(string(data))

		err = conn.Delete(id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
