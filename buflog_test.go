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

func TestLogCompatible(t *testing.T) {
	log.Print("Print one")
	log.Printf("Printf test %d\n", 123)
	log.Println("Println test %d")
	log.SetFlags(15)
	log.Print("test")
	log.ID("t").Info("bbb")
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
