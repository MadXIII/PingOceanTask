package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		urls       []string
		str        string
		wantResult map[string]int
	}{
		"Success": {
			urls:       []string{"https://pingocean.com/", "https://github.com/"},
			str:        "pingocean",
			wantResult: map[string]int{"https://pingocean.com/": 284, "https://github.com/": 0},
		},
	}

	for _, test := range tests {
		result := goSender(context.Background(), test.urls, test.str)

		if !reflect.DeepEqual(result, test.wantResult) {
			t.Errorf("Wait %v, but got %v", result, test.wantResult)
		}
	}
}

func TestSetMap(t *testing.T) {
	tests := map[string]struct {
		m         Store
		url       string
		count     int
		wantError error
	}{
		"Success": {
			m:         Store{Map: map[string]int{"https://pingocean.com/": 284}},
			url:       "https://github.com/",
			count:     0,
			wantError: nil,
		},
		"Wait Error Url is already exists": {
			m:         Store{Map: map[string]int{"https://pingocean.com/": 284}},
			url:       "https://pingocean.com/",
			count:     284,
			wantError: errors.New("url is already exists"),
		},
	}

	for _, test := range tests {
		if err := test.m.setMap(test.url, test.count); err != test.wantError {
			fmt.Println(test.wantError)
			fmt.Println(err)
			t.Errorf("Wait %v, but got %v", test.wantError, err)
		}
	}
}
