package tmux

import(
"os/exec"
"fmt"
)

type Script string
func (this *Script) append(s string){
	*this = Script(string(*this) + s + "\n")
}

var tmux_script Script 

type Pane struct {
	number int
	commands []string
	window *Window
}
func NewPane(number int, window *Window) *Pane{
	p:= new(Pane)
	p.number = number
	p.commands = make([]string,0)
	p.window = window
	return p
}
func (this *Pane) Run(command string){
	this.commands = append(this.commands, command)
}

func (this *Pane) Vsplit() *Pane{
	this.window.split_commands = append(this.window.split_commands, "tmux split-window -h -t "+this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number))
	exec.Command("tmux", "split-window", "-h", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number)).Run()
	return this.window.AddPane()
}

func (this *Pane) Split() *Pane{
	this.window.split_commands = append(this.window.split_commands, "tmux split-window -v -t "+this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number))
	return this.window.AddPane()
}

func (this *Pane) start() {
	for _,command := range this.commands {
		exec.Command("tmux", "send-keys", "-t", this.window.session.name+":"+fmt.Sprint(this.window.number)+"."+fmt.Sprint(this.number), command, "C-m").Run()
	}
}


type Window struct{
	number int
	name string
	session *Session
	panes []*Pane
	split_commands []string
	next_pane_number int
}
func NewWindow(number int, name string, session *Session) *Window{
	w := new(Window)
	w.name = name
	w.number = number
	w.session = session
	w.next_pane_number = 0
	w.panes = make([]*Pane,0)
	w.split_commands = make([]string,0)
	w.AddPane()
	return w
}


func (this *Window)AddPane() *Pane{
	pane := NewPane(this.next_pane_number, this)
	this.panes = append(this.panes, pane)
	this.next_pane_number = this.next_pane_number + 1
	return pane
}

func (this *Window)Pane(number int) *Pane{
	return this.panes[number]
}

func (this *Window) start() {
	if this.number > 0 {
		exec.Command("tmux", "new-window", "-t", this.session.name+":"+fmt.Sprint(this.number), "-n", this.name).Run()
	}

	for _, pane := range this.panes{
		pane.start()
	}
}
func (this *Window) Run(command string){
	this.Pane(0).Run(command)
}

type Session struct{
	name string
	next_window_number int
	windows []*Window
}

func NewSession(name string) *Session{
	s := new(Session)
	s.name = name
	s.windows = make([]*Window, 0)
	return s
}

func (this *Session) AddWindow(name string) *Window{
	w := NewWindow(this.next_window_number, name, this)
	this.windows = append(this.windows, w)
	this.next_window_number = this.next_window_number + 1
	return w
}

func (this *Session) Start(){
	exec.Command("tmux", "kill-session", "-t", this.name).Run()
	exec.Command("tmux", "new-session", "-d", "-s",  this.name , "-n",  this.windows[0].name).Run()
	for _, window := range this.windows {
		window.start()
	}
	exec.Command("tmux", "select-window", "-t", this.name+":0").Run()
}

