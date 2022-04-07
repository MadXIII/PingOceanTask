package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	url := flag.String("url", "https://pingocean.com/", "wrong flag")
	str := flag.String("str", "pingocean", "wrong flag")
	flag.Parse()
	count, err := searchAndCount(url, str)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(count)
	}
}

func searchAndCount(url, str *string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, NewRequest: %w", err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, clientDo: %w", err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, ReadAll: %w", err)
	}

	return strings.Count(string(bytes), *str), nil
}
