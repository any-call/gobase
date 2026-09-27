package mynet

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestCopyHalfCloseForwardsBothSides(t *testing.T) {
	src := &halfCloseReader{Buffer: bytes.NewBufferString("hello")}
	dst := &halfCloseWriter{}

	n, err := copyHalfClose(dst, src)
	if err != nil {
		t.Fatalf("copyHalfClose error: %v", err)
	}
	if n != 5 || dst.String() != "hello" {
		t.Fatalf("copied (%d, %q), want (5, hello)", n, dst.String())
	}
	if !src.readClosed || !dst.writeClosed {
		t.Fatal("half-close was not forwarded to both endpoints")
	}
}

func TestRelayWithIdleTimeoutStopsIdleConnection(t *testing.T) {
	left, leftPeer := net.Pipe()
	right, rightPeer := net.Pipe()
	defer leftPeer.Close()
	defer rightPeer.Close()

	started := time.Now()
	_, _, err := RelayWithIdleTimeout(left, right, 40*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("RelayWithIdleTimeout error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(started); elapsed < 30*time.Millisecond {
		t.Fatalf("relay stopped too early after %v", elapsed)
	}
}

func TestRelayWithIdleTimeoutKeepsActiveDownloadAlive(t *testing.T) {
	left, client := net.Pipe()
	right, server := net.Pipe()
	done := make(chan error, 1)
	go func() {
		_, _, err := RelayWithIdleTimeout(left, right, 50*time.Millisecond)
		done <- err
	}()

	chunk := []byte("data")
	buf := make([]byte, len(chunk))
	for i := 0; i < 5; i++ {
		writeDone := make(chan error, 1)
		go func() {
			_, err := server.Write(chunk)
			writeDone <- err
		}()

		if _, err := io.ReadFull(client, buf); err != nil {
			t.Fatalf("read active download: %v", err)
		}
		if err := <-writeDone; err != nil {
			t.Fatalf("write active download: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}

	select {
	case err := <-done:
		t.Fatalf("relay stopped while data was active: %v", err)
	default:
	}

	_ = client.Close()
	_ = server.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("relay did not stop after endpoints closed")
	}
}

func TestSmartCopyRejectsZeroLengthWrite(t *testing.T) {
	_, err := SmartCopy(zeroWriter{}, bytes.NewBufferString("data"))
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("SmartCopy error = %v, want io.ErrShortWrite", err)
	}
}

type halfCloseReader struct {
	*bytes.Buffer
	readClosed bool
}

func (r *halfCloseReader) CloseRead() error {
	r.readClosed = true
	return nil
}

type halfCloseWriter struct {
	bytes.Buffer
	writeClosed bool
}

func (w *halfCloseWriter) CloseWrite() error {
	w.writeClosed = true
	return nil
}

type zeroWriter struct{}

func (zeroWriter) Write([]byte) (int, error) { return 0, nil }
