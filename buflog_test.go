package buflog_test

import (
	log "github.com/bitxel/buflog"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log.Id("1").Info("test log")
	log.Id("2").Info("test log2")
	log.Id("1").Info("test log2")
	time.Sleep(time.Second * 2)
}
