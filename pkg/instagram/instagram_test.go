package instagram

import "testing"

func Test_idFromURL(t *testing.T) {
	res, err := idFromURL("https://www.instagram.com/reel/C0Y6YmMIWSs/?igshid=ZDE1MWVjZGVmZQ==")
	if err != nil {
		t.Error(err)
	}
	if res != "C0Y6YmMIWSs" {
		t.Error("url != C0Y6YmMIWSs")
	}
}
