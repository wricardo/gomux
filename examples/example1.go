package main

import (
	"fmt"

	"github.com/wricardo/gomux"
)

func main() {
	session_name := "SESSION_NAME"

	s := tmux.NewSession(session_name)

	//WINDOW 1
	w1 := s.AddWindow("LOGS")

	w1p0 := w1.Pane(0)
	w1p0.Exec("tail -f /var/log/authd.log")

	w1p1 := w1.Pane(0).Split()
	w1p1.Exec("tail -f /var/log/system.log")

	//WINDOW 2
	w2 := s.AddWindow("Vim")
	w2p0 := w2.Pane(0)

	w2p0.Exec("echo \"this is to vim\" | vim -")

	w2p1 := w2p0.Vsplit()
	w2p1.Exec("cd /tmp/")
	w2p1.Exec("ls -la")

	w2p0.ResizeRight(30)

	fmt.Println("Tmux Session \"", session_name, "\" created\n")
	fmt.Println("Now you can run:")
	fmt.Println("tmux attach -t ", session_name)

}
