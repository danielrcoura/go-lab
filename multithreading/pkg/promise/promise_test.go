package promise

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Teste para a função Soma
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
