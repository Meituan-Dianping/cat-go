package cat

import (
	"testing"
)

func TestHitSample(t *testing.T) {
	var (
		total  = 100000
		sample = 0.01
		count  = 0
	)

	for i := 0; i < total; i++ {
		if manager.hitSample(sample) {
			count++
		}
	}

	if count != int(float64(total)*sample) {
		t.Fail()
	}
}
