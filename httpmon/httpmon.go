package httpmon

import (
	"fmt"
	"net/http"
	"time"
)

func getMillisecondsSince(start time.Time) float64 {
	duration := time.Since(start)
	return float64(duration) / float64(time.Millisecond)
}

type Call struct {
	URL        string
	Success    bool
	Error      error
	CreatedAt  time.Time
	DurationMS float64
}

func IsURLUp(url string) *Call {
	call := &Call{
		URL:       url,
		Success:   false,
		CreatedAt: time.Now(),
	}
	start := time.Now()
	resp, err := http.Head(url)
	call.DurationMS = getMillisecondsSince(start)
	if err != nil {
		call.Error = err
		return call
	}

	call.Success = resp.StatusCode == http.StatusOK
	if !call.Success {
		call.Error = fmt.Errorf("call failed: %q", resp.Status)
	}
	return call
}

func Schedule(trigger func(), delay time.Duration) chan bool {
	ticker := time.NewTicker(delay)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				trigger()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}
