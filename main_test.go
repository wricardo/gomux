package gomux

import "testing"

type FakeWriter struct {
	data []byte
}

func (f *FakeWriter) Write(p []byte) (int, error) {
	f.data = append(f.data, p...)
	return len(p), nil
}

func TestNewSessionSimple(t *testing.T) {
	w := FakeWriter{
		data: make([]byte, 0),
	}
	s := NewSession("mysession", &w)

	if s.Name != "mysession" {
		t.Fatal("problem1")
	}
	if len(s.windows) != 0 {
		t.Fatal("problem1")
	}

	if string(w.data) != "tmux kill-session -t \"mysession\"\ntmux new-session -d -s \"mysession\" -n tmp\n" {
		t.Fatal("problem")
	}
}

func TestNewSessionDirectory(t *testing.T) {
	w := FakeWriter{
		data: make([]byte, 0),
	}
	NewSessionParams(NewSessionAttr{
		Name:      "mysession",
		Directory: "/tmp/a",
	}, &w)

	if string(w.data) != "tmux kill-session -t \"mysession\"\ntmux new-session -d -s \"mysession\" -n tmp -c /tmp/a\n" {
		t.Fatal("problem")
	}
}
