package gomux

import (
	"fmt"
	"io"
	"strings"
)

type Pane struct {
	Number   int
	commands []string
	window   *Window
}

func NewPane(number int, window *Window) *Pane {
	p := new(Pane)
	p.Number = number
	p.commands = make([]string, 0)
	p.window = window
	return p
}

func (this *Pane) Exec(command string) {
	fmt.Fprintf(this.window.session.writer, "tmux send-keys -t \"%s\" \"%s\" %s\n", this.getTargetName(), strings.Replace(command, "\"", "\\\"", -1), "C-m")
}

func (this *Pane) Vsplit() *Pane {
	fmt.Fprintf(this.window.session.writer, "tmux split-window -h -t \"%s\"\n", this.getTargetName())
	return this.window.AddPane()
}

func (this *Pane) Split() *Pane {
	fmt.Fprintf(this.window.session.writer, "tmux split-window -v -t \"%s\"\n", this.getTargetName())
	return this.window.AddPane()
}

func (this *Pane) ResizeRight(num int) {
	this.resize("R", num)
}

func (this *Pane) ResizeLeft(num int) {
	this.resize("L", num)
}

func (this *Pane) ResizeUp(num int) {
	this.resize("U", num)
}

func (this *Pane) ResizeDown(num int) {
	this.resize("U", num)
}

func (this *Pane) resize(prefix string, num int) {
	fmt.Fprintf(this.window.session.writer, "tmux resize-pane -t \"%s\" -%s\n", this.getTargetName(), prefix, fmt.Sprint(num))
}

func (this *Pane) getTargetName() string {
	return this.window.session.Name + ":" + fmt.Sprint(this.window.Number) + "." + fmt.Sprint(this.Number)
}

// Window Represent a tmux window. You usually should not create an instance of Window directly.
type Window struct {
	Number           int
	Name             string
	session          *Session
	panes            []*Pane
	split_commands   []string
	next_pane_number int
}

func newWindow(number int, name string, session *Session) *Window {
	w := new(Window)
	w.Name = name
	w.Number = number
	w.session = session
	w.next_pane_number = 0
	w.panes = make([]*Pane, 0)
	w.split_commands = make([]string, 0)
	w.AddPane()
	if number != 0 {
		fmt.Fprintf(session.writer, "tmux new-window %s -n \"%s\"\n", w.t(), w.Name)
	}
	fmt.Fprintf(session.writer, "tmux rename-window %s \"%s\"\n", w.t(), w.Name)
	return w
}

func (this *Window) t() string {
	return fmt.Sprintf("-t \"%s:%s\"", this.session.Name, fmt.Sprint(this.Number))
}

// Create a new Pane and add to this window
func (this *Window) AddPane() *Pane {
	pane := NewPane(this.next_pane_number, this)
	this.panes = append(this.panes, pane)
	this.next_pane_number = this.next_pane_number + 1
	return pane
}

// Find and return the Pane object by the number
func (this *Window) Pane(number int) *Pane {
	return this.panes[number]
}

// Executes a command on the first pane of this window
//
// // example
// // example
func (this *Window) Exec(command string) {
	this.Pane(0).Exec(command)
}

func (this *Window) Select() {
	fmt.Fprintf(this.session.writer, "tmux select-window -t \"%s:%s\"\n", this.session.Name, fmt.Sprint(this.Number))
}

// Session represents a tmux session.
//
// Use the method NewSession to create a Session instance.
type Session struct {
	Name               string
	windows            []*Window
	next_window_number int
	writer             io.Writer
}

// Creates a new Tmux Session. It kill any existing session with the provided name.
func NewSession(name string, writer io.Writer) *Session {
	s := new(Session)
	s.writer = writer
	s.Name = name
	s.windows = make([]*Window, 0)
	fmt.Fprintf(writer, "tmux kill-session -t \"%s\"\n", s.Name)
	fmt.Fprintf(writer, "tmux new-session -d -s \"%s\" -n tmp\n", s.Name)
	return s
}

// Creates window with provided name for this session
func (this *Session) AddWindow(name string) *Window {
	w := newWindow(this.next_window_number, name, this)
	this.windows = append(this.windows, w)
	this.next_window_number = this.next_window_number + 1
	return w
}
