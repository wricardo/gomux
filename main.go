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

type SplitAttr struct {
	Directory string
}

func (this *Pane) Exec(command string) {
	fmt.Fprintf(this.window.session.writer, "tmux send-keys -t \"%s\" \"%s\" %s\n", this.getTargetName(), strings.Replace(command, "\"", "\\\"", -1), "C-m")
}

func (this *Pane) Vsplit() *Pane {
	fmt.Fprint(this.window.session.writer, splitWindow{h: true, t: this.getTargetName()})
	return this.window.AddPane()
}

func (this *Pane) VsplitWAttr(attr SplitAttr) *Pane {
	var c string
	if attr.Directory != "" {
		c = attr.Directory
	} else if this.window.Directory != "" {
		c = this.window.Directory
	} else if this.window.session.Directory != "" {
		c = this.window.session.Directory
	}

	fmt.Fprint(this.window.session.writer, splitWindow{h: true, t: this.getTargetName(), c: c})
	return this.window.AddPane()
}

func (this *Pane) Split() *Pane {
	fmt.Fprint(this.window.session.writer, splitWindow{v: true, t: this.getTargetName()})
	return this.window.AddPane()
}

func (this *Pane) SplitWAttr(attr SplitAttr) *Pane {
	var c string
	if attr.Directory != "" {
		c = attr.Directory
	} else if this.window.Directory != "" {
		c = this.window.Directory
	} else if this.window.session.Directory != "" {
		c = this.window.session.Directory
	}

	fmt.Fprint(this.window.session.writer, splitWindow{v: true, t: this.getTargetName(), c: c})
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
	Directory        string
	session          *Session
	panes            []*Pane
	split_commands   []string
	next_pane_number int
}

type WindowAttr struct {
	Name      string
	Directory string
}

func createWindow(number int, attr WindowAttr, session *Session) *Window {
	w := new(Window)
	w.Name = attr.Name
	w.Directory = attr.Directory
	w.Number = number
	w.session = session
	w.next_pane_number = 0
	w.panes = make([]*Pane, 0)
	w.split_commands = make([]string, 0)
	w.AddPane()

	if number != 0 {
		fmt.Fprint(session.writer, newWindow{t: w.t(), n: w.Name, c: attr.Directory})
	}

	fmt.Fprint(session.writer, renameWindow{t: w.t(), n: w.Name})
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
	fmt.Fprint(this.session.writer, selectWindow{t: this.session.Name + ":" + fmt.Sprint(this.Number)})
}

// Session represents a tmux session.
//
// Use the method NewSession to create a Session instance.
type Session struct {
	Name               string
	Directory          string
	windows            []*Window
	directory          string
	next_window_number int
	writer             io.Writer
}

// Creates a new Tmux Session. It will kill any existing session with the provided name.
func NewSession(name string, writer io.Writer) *Session {
	p := SessionAttr{
		Name: name,
	}
	return NewSessionAttr(p, writer)
}

type SessionAttr struct {
	Name      string
	Directory string
}

// Creates a new Tmux Session based on NewSessionAttr. It will kill any existing session with the provided name.
func NewSessionAttr(p SessionAttr, writer io.Writer) *Session {
	s := new(Session)
	s.writer = writer
	s.Name = p.Name
	s.Directory = p.Directory
	s.windows = make([]*Window, 0)

	fmt.Fprint(writer, killSession{t: s.Name})
	fmt.Fprint(writer, newSession{d: true, s: p.Name, c: p.Directory, n: "tmp"})
	return s
}

// Creates window with provided name for this session
func (this *Session) AddWindow(name string) *Window {

	attr := WindowAttr{
		Name: name,
	}

	return this.AddWindowAttr(attr)
}

// Creates window with provided name for this session
func (this *Session) AddWindowAttr(attr WindowAttr) *Window {
	w := createWindow(this.next_window_number, attr, this)
	this.windows = append(this.windows, w)
	this.next_window_number = this.next_window_number + 1
	return w
}
