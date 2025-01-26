package gomux

import (
	"fmt"
	"io"
	"strings"
)

// Pane represents a tmux pane inside a Window.
type Pane struct {
	Number   int
	commands []string // Not used in the current snippet, but might be used for queued commands later.
	window   *Window
}

// NewPane creates and returns a pointer to a new Pane.
func NewPane(number int, window *Window) *Pane {
	return &Pane{
		Number:   number,
		commands: make([]string, 0),
		window:   window,
	}
}

// SplitAttr holds optional attributes when splitting a pane, such as the working directory.
type SplitAttr struct {
	Directory string
}

// Exec sends a shell command to this Pane, ending with an Enter key press.
func (p *Pane) Exec(command string) {
	escapedCmd := strings.Replace(command, "\"", "\\\"", -1)
	fmt.Fprintf(p.window.session.writer,
		"tmux send-keys -t \"%s\" \"%s\" C-m\n",
		p.getTargetName(), escapedCmd,
	)
}

// Vsplit splits the current pane vertically (left-right split).
func (p *Pane) Vsplit() *Pane {
	// splitWindow is presumably a type that implements String() for tmux command generation.
	fmt.Fprint(p.window.session.writer, splitWindow{h: true, t: p.getTargetName()})
	return p.window.AddPane(p.Number + 1)
}

// VsplitWAttr is a vertical split with an optional working directory.
func (p *Pane) VsplitWAttr(attr SplitAttr) *Pane {
	c := p.resolveDirectory(attr)
	fmt.Fprint(p.window.session.writer, splitWindow{h: true, t: p.getTargetName(), c: c})
	return p.window.AddPane(p.Number + 1)
}

// Split splits the current pane horizontally (top-bottom split).
func (p *Pane) Split() *Pane {
	fmt.Fprint(p.window.session.writer, splitWindow{v: true, t: p.getTargetName()})
	return p.window.AddPane(p.Number + 1)
}

// SplitWAttr is a horizontal split with an optional working directory.
func (p *Pane) SplitWAttr(attr SplitAttr) *Pane {
	c := p.resolveDirectory(attr)
	fmt.Fprint(p.window.session.writer, splitWindow{v: true, t: p.getTargetName(), c: c})
	return p.window.AddPane(p.Number + 1)
}

// ResizeRight grows the pane to the right by the given number of cells.
func (p *Pane) ResizeRight(num int) {
	p.resize("R", num)
}

// ResizeLeft shrinks the pane to the left by the given number of cells.
func (p *Pane) ResizeLeft(num int) {
	p.resize("L", num)
}

// ResizeUp grows the pane upwards by the given number of cells.
func (p *Pane) ResizeUp(num int) {
	p.resize("U", num)
}

// ResizeDown grows the pane downwards by the given number of cells.
// BUG FIX: changed "U" to "D" in the command flag.
func (p *Pane) ResizeDown(num int) {
	p.resize("D", num)
}

// getTargetName returns a string identifying this pane in tmux notation, e.g. "mySession:1.0".
func (p *Pane) getTargetName() string {
	return fmt.Sprintf("%s:%d.%d", p.window.session.Name, p.window.Number, p.Number)
}

// resolveDirectory determines which directory should be used when splitting a pane.
func (p *Pane) resolveDirectory(attr SplitAttr) string {
	switch {
	case attr.Directory != "":
		return attr.Directory
	case p.window.Directory != "":
		return p.window.Directory
	case p.window.session.Directory != "":
		return p.window.session.Directory
	default:
		return ""
	}
}

// resize runs the tmux resize-pane command with the given direction flag and size.
func (p *Pane) resize(prefix string, num int) {
	fmt.Fprintf(p.window.session.writer,
		"tmux resize-pane -t \"%s\" -%s %d\n",
		p.getTargetName(), prefix, num,
	)
}

// Window represents a tmux window, which can contain multiple panes.
//
// Typically you create a Window through a Session's AddWindow method.
type Window struct {
	Number    int
	Name      string
	Directory string

	session        *Session
	panes          []*Pane
	split_commands []string // Unused in snippet, but you might use it for later expansions.
}

// WindowAttr holds optional attributes for creating a Window.
type WindowAttr struct {
	Name      string
	Directory string
}

// createWindow is an internal helper to set up a new Window in a given Session.
func createWindow(number int, attr WindowAttr, session *Session) *Window {
	w := &Window{
		Name:      attr.Name,
		Directory: attr.Directory,
		Number:    number,
		session:   session,
		panes:     make([]*Pane, 0),
		// Possibly used for storing commands to run after creation:
		split_commands: make([]string, 0),
	}

	// Create the window in tmux, if number != 0.
	// Typically the first window is created with the session, so you might skip.
	if number != 0 {
		fmt.Fprint(session.writer, newWindow{t: w.t(), n: w.Name, c: w.Directory})
	}

	// Rename the window to the desired name. (By default, tmux might assign something else.)
	fmt.Fprint(session.writer, renameWindow{t: w.t(), n: w.Name})

	// By default, add an initial pane (pane #0).
	w.AddPane(0)

	return w
}

// t is an internal helper that builds the tmux target string for this window, e.g. -t "mySession:1".
func (w *Window) t() string {
	return fmt.Sprintf("-t \"%s:%d\"", w.session.Name, w.Number)
}

// AddPane creates a new Pane in this Window. The `withNumber` is
// the pane index (tmux's ID). Typically it's sequentially assigned.
func (w *Window) AddPane(withNumber int) *Pane {
	pane := NewPane(withNumber, w)
	w.panes = append(w.panes, pane)
	return pane
}

// Pane returns the pane by its index in the window's internal slice (panes array).
func (w *Window) Pane(number int) *Pane {
	return w.panes[number]
}

// Exec sends a command to the first pane of this window.
func (w *Window) Exec(command string) {
	w.Pane(0).Exec(command)
}

// Select makes this window the active window in the session.
func (w *Window) Select() {
	fmt.Fprint(w.session.writer, selectWindow{t: fmt.Sprintf("%s:%d", w.session.Name, w.Number)})
}

// Session represents a tmux session. Use NewSession or NewSessionAttr to create.
type Session struct {
	Name      string
	Directory string

	windows            []*Window
	next_window_number int
	// writer is where tmux commands are sent.
	writer io.Writer
}

// SessionAttr holds optional attributes for creating a new Session.
type SessionAttr struct {
	Name      string
	Directory string
}

// NewSession kills any existing session with the same name, then creates a new one.
func NewSession(name string, writer io.Writer) *Session {
	return NewSessionAttr(SessionAttr{Name: name}, writer)
}

// NewSessionAttr kills any existing session with the same name, then creates a new one with attributes.
func NewSessionAttr(attr SessionAttr, writer io.Writer) *Session {
	s := &Session{
		Name:      attr.Name,
		Directory: attr.Directory,
		windows:   make([]*Window, 0),
		writer:    writer,
	}
	fmt.Fprint(writer, newSession{d: true, s: attr.Name, c: attr.Directory, n: "tmp"})
	return s
}

// KillSession sends a command to kill an existing tmux session by name.
func KillSession(name string, writer io.Writer) {
	fmt.Fprint(writer, killSession{t: name})
}

// AddWindow creates a new Window in the session with the specified name.
func (s *Session) AddWindow(name string) *Window {
	return s.AddWindowAttr(WindowAttr{Name: name})
}

// AddWindowAttr creates a new Window in the session with additional attributes (e.g., directory).
func (s *Session) AddWindowAttr(attr WindowAttr) *Window {
	w := createWindow(s.next_window_number, attr, s)
	s.windows = append(s.windows, w)
	s.next_window_number++
	return w
}
