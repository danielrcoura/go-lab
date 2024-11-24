package promise

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	result, _ := All(dummyFunc("f1", 100), dummyFunc("f2", 500), dummyFunc("f3", 1000))
	assert.Equal(t, []any{"f1", "f2", "f3"}, result)
}

func TestAllError(t *testing.T) {
	errorFn := func() (any, error) {
		return "", errors.New("error")
	}
	_, err := All(errorFn, dummyFunc("f2", 500), dummyFunc("f3", 1000))
	assert.Error(t, err, "Expected an error but got nil")
}

func TestAllSettled(t *testing.T) {
	f1 := func() string {
		time.Sleep(100 * time.Millisecond)
		return "f1"
	}
	f2 := func() string {
		time.Sleep(500 * time.Millisecond)
		return "f2"
	}
	f3 := func() string {
		time.Sleep(1000 * time.Millisecond)
		return "f3"
	}
	result := AllSettled(f1, f2, f3)
	assert.Equal(t, []string{"f1", "f2", "f3"}, result)
}

func dummyFunc(name string, milliseconds int64) func() (any, error) {
	return func() (any, error) {
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		return name, nil
	}
}
