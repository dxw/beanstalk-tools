package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kr/beanstalk"
)

const (
	defaultPriority = (2 ^ 32) - 1
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: bst-import hostname:11300 tubename /path/to/export.txt")
	}

	host := os.Args[1]
	tubeName := os.Args[2]
	filePath := os.Args[3]

	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	tube := &beanstalk.Tube{conn, tubeName}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)

	for {
		line, err := fileReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if line == "---\n" {
			continue
		}

		line = strings.TrimRight(line, "\n")

		_, err = tube.Put([]byte(line), defaultPriority, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}
