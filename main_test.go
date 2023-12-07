package main

import "testing"

func Test_idFromURL(t *testing.T) {
	res, err := idFromURL("https://www.instagram.com/reel/CFQ4YJYn7ZS/")
	if err != nil {
		t.Error(err)
	}
	if res != "CFQ4YJYn7ZS" {
		t.Error("url != CFQ4YJYn7ZS")
	}
}
