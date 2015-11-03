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
	NewSessionAttr(SessionAttr{
		Name:      "mysession",
		Directory: "/tmp/a",
	}, &w)

	if string(w.data) != "tmux kill-session -t \"mysession\"\ntmux new-session -d -s \"mysession\" -n tmp -c /tmp/a\n" {
		t.Fatal("problem")
	}
}

func TestNewWindow(t *testing.T) {
	w := FakeWriter{
		data: make([]byte, 0),
	}
	sess := NewSession("mySession", &w)
	sess.AddWindow("myWindow")
	sess.AddWindow("myWindow2")

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux new-window -t "mySession:1" -n "myWindow2"
tmux rename-window -t "mySession:1" "myWindow2"
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestNewWindowDirectory(t *testing.T) {

	attr := WindowAttr{
		Directory: "/tmp/a",
		Name:      "myWindow2",
	}
	w := FakeWriter{
		data: make([]byte, 0),
	}
	sess := NewSession("mySession", &w)
	sess.AddWindow("myWindow")
	sess.AddWindowAttr(attr)

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux new-window -t "mySession:1" -n "myWindow2" -c /tmp/a
tmux rename-window -t "mySession:1" "myWindow2"
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestVsplit(t *testing.T) {
	w := FakeWriter{
		data: make([]byte, 0),
	}
	sess := NewSession("mySession", &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)
	p.Vsplit()

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -h -t "mySession:0.0"
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestSplit(t *testing.T) {
	w := FakeWriter{
		data: make([]byte, 0),
	}
	sess := NewSession("mySession", &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)
	p.Split()

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -v -t "mySession:0.0"
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestSplitWAttr(t *testing.T) {

	attr := SplitAttr{
		Directory: "/tmp/c",
	}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSession("mySession", &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)

	p.SplitWAttr(attr)

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -v -t "mySession:0.0" -c /tmp/c
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestVsplitWAttr(t *testing.T) {

	attr := SplitAttr{
		Directory: "/tmp/c",
	}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSession("mySession", &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)

	p.VsplitWAttr(attr)

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -h -t "mySession:0.0" -c /tmp/c
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestVsplitWAttrBubbleWindow(t *testing.T) {
	attr := SplitAttr{}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSession("mySession", &w)
	window := sess.AddWindowAttr(WindowAttr{
		Name:      "myWindow",
		Directory: "/tmp/window",
	})
	p := window.Pane(0)

	p.VsplitWAttr(attr)

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -h -t "mySession:0.0" -c /tmp/window
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestSplitWAttrBubbleWindow(t *testing.T) {

	attr := SplitAttr{}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSession("mySession", &w)
	window := sess.AddWindowAttr(WindowAttr{
		Name:      "myWindow",
		Directory: "/tmp/window",
	})
	p := window.Pane(0)

	p.SplitWAttr(attr)

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -v -t "mySession:0.0" -c /tmp/window
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestVsplitWAttrBubbleSession(t *testing.T) {
	attr := SessionAttr{
		Name:      "mySession",
		Directory: "/tmp/session",
	}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSessionAttr(attr, &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)

	p.VsplitWAttr(SplitAttr{})

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp -c /tmp/session
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -h -t "mySession:0.0" -c /tmp/session
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}

func TestSplitWAttrBubbleSession(t *testing.T) {
	attr := SessionAttr{
		Name:      "mySession",
		Directory: "/tmp/session",
	}
	w := FakeWriter{
		data: make([]byte, 0),
	}

	sess := NewSessionAttr(attr, &w)
	window := sess.AddWindow("myWindow")
	p := window.Pane(0)

	p.SplitWAttr(SplitAttr{})

	expected := `tmux kill-session -t "mySession"
tmux new-session -d -s "mySession" -n tmp -c /tmp/session
tmux rename-window -t "mySession:0" "myWindow"
tmux split-window -v -t "mySession:0.0" -c /tmp/session
`

	if string(w.data) != expected {
		t.Fatal("problem")
	}
}
