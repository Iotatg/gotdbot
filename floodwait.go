package gotdbot

import (
	"errors"
	"time"
)

func AsTDLibError(err error) *Error {
	if err == nil {
		return nil
	}
	var ptr *Error
	if errors.As(err, &ptr) && ptr != nil {
		return ptr
	}
	var val Error
	if errors.As(err, &val) {
		return &val
	}
	return nil
}

func IsFloodWait(err error) bool {
	tdErr := AsTDLibError(err)
	return tdErr != nil && tdErr.GetRetryAfter() > 0
}

func FloodWaitSeconds(err error) int {
	tdErr := AsTDLibError(err)
	if tdErr == nil {
		return 0
	}
	return tdErr.GetRetryAfter()
}

func SleepFloodWait(err error) bool {
	seconds := FloodWaitSeconds(err)
	if seconds <= 0 {
		return false
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	return true
}

func RetryAfterFloodWait(err error, fn func() error) error {
	if !SleepFloodWait(err) {
		return err
	}
	return fn()
}
