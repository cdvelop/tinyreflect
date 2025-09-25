package tinyreflect_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestStructCacheRace triggers concurrent cache population to reproduce the
// known data race in TinyReflect's struct cache. Run with `go test -race` to
// observe the detector reporting the race.
func TestStructCacheRace(t *testing.T) {
	type raceStruct struct {
		Name string
		Age  int
		Flag bool
	}

	const (
		goroutines = 16
		iterations = 200
	)

	for iter := 0; iter < iterations; iter++ {
		tr := tinyreflect.New()

		start := make(chan struct{})
		var wg sync.WaitGroup
		errChan := make(chan error, goroutines)

		for g := 0; g < goroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for i := 0; i < 32; i++ {
					if typ := tr.TypeOf(raceStruct{}); typ == nil {
						errChan <- fmt.Errorf("TypeOf returned nil")
						return
					}
				}
			}()
		}

		close(start)
		wg.Wait()
		close(errChan)

		// Check for errors after all goroutines complete
		for err := range errChan {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
