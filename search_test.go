package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		result := worker(context.Background(), test.urls, test.str)
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
			m:         Store{Map: map[string]int{"https://pingocean.com/": 10}},
			url:       "https://github.com/",
			count:     10,
			wantError: nil,
		},
		"Wait Error Url is already exists": {
			m:         Store{Map: map[string]int{"https://pingocean.com/": 0}},
			url:       "https://pingocean.com/",
			count:     0,
			wantError: errors.New("https://pingocean.com/ url is already exists"),
		},
	}
	for _, test := range tests {
		err := test.m.setMap(test.url, test.count)
		assert.Equal(t, err, test.wantError)
	}
}
