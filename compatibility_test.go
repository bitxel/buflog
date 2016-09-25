package buflog_test

import (
	"bytes"
	log "github.com/bitxel/buflog"
	"runtime"
	"strings"
	"testing"
)

func testPrint(t *testing.T, buf *bytes.Buffer, in, out string) {
	log.Print(in)
	s, err := buf.ReadString('\n')
	if err != nil || !strings.Contains(s, out) {
		t.Fatalf("print func error: %s contain expected: %s ret: %s", err, out, s)
	}
	t.Log(s)
}

func TestLogCompatible(t *testing.T) {
	log.SetBufTimeout(0)
	var buf = new(bytes.Buffer)
	log.SetOutput(buf)

	testPrint(t, buf, "Print func", "Print func")

	log.Printf("Printf func %d", 123)
	s, err := buf.ReadString('\n')
	if err != nil || !strings.Contains(s, "Printf func 123") {
		t.Fatalf("printf func error: %s ret: %s", err, s)
	}

	_, fullpath, _, _ := runtime.Caller(0)
	log.SetFlags(15)
	testPrint(t, buf, "test long filename", fullpath)

	filename := "compatibility_test"
	log.SetFlags(log.Lshortfile)
	testPrint(t, buf, "test sort filename", filename)

	log.SetPrefix("pf")
	testPrint(t, buf, "test prefix", "pf")
}
