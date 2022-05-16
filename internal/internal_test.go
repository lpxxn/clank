package internal

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
