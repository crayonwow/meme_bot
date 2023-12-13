package bot

import (
	"testing"
	"time"
)

func TestBot(t *testing.T) {
	t.Run("rate_limiter", func(t *testing.T) {
		t.Skip()
		b := NewBot(nil, "", 3*time.Second)
		for i := 0; i < 10; i++ {
			t.Log(b.isSilent(false))
			time.Sleep(time.Second)
		}
	})
}
