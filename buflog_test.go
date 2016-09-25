package buflog_test

import (
	log "github.com/bitxel/buflog"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log.ID("1").Info("test log")
	log.ID("2").Info("test log2")
	log.ID("1").Info("test log2")
	time.Sleep(time.Second * 2)
}

func TestPrefix(t *testing.T) {
	prefix := "this is prefix"
	log.SetPrefix(prefix)
	if log.Prefix() != prefix {
		t.Fatal("set prefix failed")
	}
	log.Print("test prefix")
	time.Sleep(time.Second * 2)
}
