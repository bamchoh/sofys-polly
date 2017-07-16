package main

import "os/exec"

type SofTalk struct {
	ExecDir string
	InExec  chan int
}

func NewSofTalk(dir string) *SofTalk {
	return &SofTalk{
		ExecDir: dir,
		InExec:  make(chan int, 1),
	}
}

func (s *SofTalk) ReadAloud(text string) error {
	s.InExec <- 0
	cmd := exec.Command(s.ExecDir, "/X:1", "/W:"+text)
	err := cmd.Run()
	<-s.InExec
	return err
}
