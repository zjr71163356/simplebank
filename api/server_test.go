package api

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := fmt.Errorf("test error")
	fmt.Printf("%v", err.Error())
}
