package debounce

import "time"

func New(wait time.Duration, fn func()) func() {
	var timer *time.Timer

	return func() {
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, fn)
	}
}
