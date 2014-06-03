package tmux

import (
	"fmt"
	"os/exec"
)

type Script string

func (this *Script) append(s string) {
	*this = Script(string(*this) + s + "\n")
}

var tmux_script Script

type Pane struct {
	number   int
	commands []string
	window   *Window
}

func NewPane(number int, window *Window) *Pane {
	p := new(Pane)
	p.number = number
	p.commands = make([]string, 0)
	p.window = window
	return p
}
func (this *Pane) Exec(command string) {
	Exec("tmux", "send-keys", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number), command, "C-m")
}

func (this *Pane) Vsplit() *Pane {
	Exec("tmux", "split-window", "-h", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number))
	return this.window.AddPane()
}

func (this *Pane) Split() *Pane {
	Exec("tmux", "split-window", "-v", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number))
	return this.window.AddPane()
}

func (this *Pane) ResizeRight(num int) {
	Exec("tmux", "resize-pane", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number), "-R", fmt.Sprint(num))
}

type Window struct {
	number           int
	name             string
	session          *Session
	panes            []*Pane
	split_commands   []string
	next_pane_number int
}

func newWindow(number int, name string, session *Session) *Window {
	w := new(Window)
	w.name = name
	w.number = number
	w.session = session
	w.next_pane_number = 0
	w.panes = make([]*Pane, 0)
	w.split_commands = make([]string, 0)
	w.AddPane()
	return w
}
func NewWindow(number int, name string, session *Session) *Window {
	w := newWindow(number, name, session)
	Exec("tmux", "new-window", "-t", w.session.name+":"+fmt.Sprint(w.number), "-n", w.name)
	Exec("tmux", "rename-window", "-t", w.session.name+":"+fmt.Sprint(w.number), w.name)
	return w
}

func (this *Window) AddPane() *Pane {
	pane := NewPane(this.next_pane_number, this)
	this.panes = append(this.panes, pane)
	this.next_pane_number = this.next_pane_number + 1
	return pane
}

func (this *Window) Pane(number int) *Pane {
	return this.panes[number]
}

func (this *Window) Exec(command string) {
	this.Pane(0).Exec(command)
}

func (this *Window) Select() {
	Exec("tmux", "select-window", "-t", this.session.name+":"+fmt.Sprint(this.number))
}

type Session struct {
	name               string
	next_window_number int
	windows            []*Window
}

func NewSession(name string) *Session {
	s := new(Session)
	s.name = name
	s.windows = make([]*Window, 0)
	Exec("tmux", "kill-session", "-t", s.name)
	Exec("tmux", "new-session", "-d", "-s", s.name, "-n tmp")
	return s
}

func (this *Session) AddWindow(name string) *Window {
	w := NewWindow(this.next_window_number, name, this)
	this.windows = append(this.windows, w)
	this.next_window_number = this.next_window_number + 1
	return w
}

func Exec(args ...string) {
	exec.Command(args[0], args[1:]...).Run()
}
