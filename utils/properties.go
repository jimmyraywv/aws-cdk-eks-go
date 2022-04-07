package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

type Properties map[string]string

func (p Properties) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return string(out)
}

func ReadProperties(f string) (Properties, error) {
	p := Properties{}

	if len(f) == 0 {
		return nil, errors.New("file path empty")
	}

	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				p[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return p, nil
}
