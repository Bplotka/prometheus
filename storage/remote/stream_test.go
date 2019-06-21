package remote

import (
	"bytes"
	"github.com/prometheus/prometheus/util/testutil"
	"io"
	"testing"
)

func TestStreamReaderCanReadWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewStreamWriter(b)
	r := NewStreamReader(b)

	msgs := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
		[]byte("test3"),
		[]byte("test4"),
		[]byte{}, // This is ignored by writer.
		[]byte("test5-after-empty"),
	}

	for _, msg := range msgs {
		n, err := w.Write(msg)
		testutil.Ok(t, err)
		testutil.Equals(t, len(msg), n)
	}

	i := 0
	for ; i < 4; i++ {
		msg, err := r.Next()
		testutil.Ok(t, err)
		testutil.Assert(t, i < len(msgs), "more messages then expected")
		testutil.Equals(t, msgs[i], msg)
	}

	// Empty byte slice is skipped.
	i++

	msg, err := r.Next()
	testutil.Ok(t, err)
	testutil.Assert(t, i < len(msgs), "more messages then expected")
	testutil.Equals(t, msgs[i], msg)

	_, err = r.Next()
	testutil.NotOk(t, err, "expected io.EOF")
	testutil.Equals(t, io.EOF, err)
}
