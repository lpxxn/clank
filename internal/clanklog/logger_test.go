package clanklog

import (
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//os.Setenv("CLANK_LOG_LEVEL", "info")
	//os.Setenv("CLANK_LOG_FORMATTER", "json")
	NewLogger()
	os.Exit(m.Run())
}

func TestLogger1(t *testing.T) {
	Info("a", "b")
	Infof("a %s", "b")
	Error("err1", "err2")
	Errorf("hello err: %+v", errors.New("hahahha"))
	Fatalf("wtf: %+v", errors.New("surprise"))
}
