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
