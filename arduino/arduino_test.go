package arduino

import (
	"bytes"
	"testing"
	"time"
)

type FakeReaderWriterCloser struct {
	read  *bytes.Buffer
	write *bytes.Buffer
}

func (f *FakeReaderWriterCloser) Read(p []byte) (n int, err error) {
	data := f.write.String()
	switch data {
	case "get_data\n":
		f.read.Write([]byte("OK;1510986597;24.75;35.21;true\n"))
	}
	return f.read.Read(p)
}
func (f *FakeReaderWriterCloser) Write(p []byte) (n int, err error) {
	return f.write.Write(p)
}
func (f *FakeReaderWriterCloser) Close() error {
	return nil
}
func TestState(t *testing.T) {
	fakePort := &FakeReaderWriterCloser{
		read:  &bytes.Buffer{},
		write: &bytes.Buffer{},
	}
	client := Client{port: fakePort}
	state, err := client.State()
	if err != nil {
		t.Errorf("got error during running func: %v", err)
	}
	if state.Humidity != 35.21 {
		t.Errorf("wrong humidity: want 35.21, got: %f", state.Humidity)
	}
	if state.Temperature != 24.75 {
		t.Errorf("wrong temperature: want 24.75, got: %f", state.Temperature)
	}
	if state.IsManualHandling != true {
		t.Errorf("wrong isManualHandling: want true, got: %t", state.IsManualHandling)
	}

	if state.DateTime != time.Unix(1510986597, 0) {
		t.Errorf("wrong DateTime: want %s, got: %s", time.Unix(1510986597, 0), state.DateTime)
	}

}
