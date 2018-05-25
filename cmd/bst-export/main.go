package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dxw/beanstalk-tools"
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

	encoder := json.NewEncoder(os.Stdout)

	for {
		id, data, err := tubeSet.Reserve(1 * time.Second)
		if e, ok := err.(beanstalk.ConnError); ok && e.Err == beanstalk.ErrTimeout {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		stats, err := conn.StatsJob(id)
		if err != nil {
			log.Fatal(err)
		}

		age, err := strconv.Atoi(stats["age"])
		if err != nil {
			log.Fatal(err)
		}

		// Subtract age (reported by beanstalkd) from the current
		// time (reported by the OS) to get the creation time
		creation := time.Now().Add(-time.Duration(age) * time.Second).UTC()

		err = encoder.Encode(&common.Item{
			Content:   string(data),
			CreatedAt: creation,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = conn.Delete(id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
