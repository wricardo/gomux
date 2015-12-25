gomux
=====
[![Build Status](https://travis-ci.org/wricardo/gomux.svg?branch=master)](https://travis-ci.org/wricardo/gomux) [![Coverage Status](https://coveralls.io/repos/wricardo/gomux/badge.svg?branch=master&service=github)](https://coveralls.io/github/wricardo/gomux?branch=master) [![GoDoc](https://godoc.org/github.com/wricardo/gomux?status.png)](https://godoc.org/github.com/wricardo/gomux)


Go wrapper to create tmux sessions, windows and panes.

### Example
example.go:
```go
package main
import (
	"os"
	"github.com/wricardo/gomux"
)
func main() {
	session_name := "SESSION_NAME"

	s := gomux.NewSession(session_name, os.Stdout)

	//WINDOW 1
	w1 := s.AddWindow("Monitoring")

	w1p0 := w1.Pane(0)
	w1p0.Exec("htop")

	w1p1 := w1.Pane(0).Split()
	w1p1.Exec("tail -f /var/log/syslog")

	//WINDOW 2
	w2 := s.AddWindow("Vim")
	w2p0 := w2.Pane(0)

	w2p0.Exec("echo \"this is to vim\" | vim -")

	w2p1 := w2p0.Vsplit()
	w2p1.Exec("cd /tmp/")
	w2p1.Exec("ls -la")

	w2p0.ResizeRight(30)
	w1.Select()
}
```
To print the tmux commands:
```
go run example.go 
```
```
tmux new-session -d -s "SESSION_NAME" -n tmp
tmux rename-window -t "SESSION_NAME:0" "Monitoring"
tmux send-keys -t "SESSION_NAME:0.0" "htop" C-m
tmux split-window -v -t "SESSION_NAME:0.0"
tmux send-keys -t "SESSION_NAME:0.1" "tail -f /var/log/syslog" C-m
tmux new-window -t "SESSION_NAME:1" -n "Vim"
tmux rename-window -t "SESSION_NAME:1" "Vim"
tmux send-keys -t "SESSION_NAME:1.0" "echo \"this is to vim\" | vim -" C-m
tmux split-window -h -t "SESSION_NAME:1.0"
tmux send-keys -t "SESSION_NAME:1.1" "cd /tmp/" C-m
tmux send-keys -t "SESSION_NAME:1.1" "ls -la" C-m
tmux resize-pane -t "SESSION_NAME:1.0" -R 30
tmux select-window -t "SESSION_NAME:0"
```

To create and attach to the tmux session:
```
go run example.go | bash
tmux attach -t SESSION_NAME
```
![example screenshot](https://raw.githubusercontent.com/wricardo/gomux/master/examples/screenshot_example.png)
