package main

import (
  "fmt"
  "os/exec"
)

type step struct {
  name    string
  exe     string
  args    []string
  message string
  proj    string
}

func newStep(name, exe, message, proj string, args []string) *step {
  s := &step{
    name:    name,
    exe:     exe,
    message: message,
    args:    args,
    proj:    proj,
  }

  return s
}

func (s *step) execute() (string, error) {
  cmd := exec.Command(s.exe, s.args...)
  cmd.Dir = s.proj

  if err := cmd.Run(); err != nil {
    return "", fmt.Errorf("'%s' failed: %s", s.name, err)
  }

  return s.message, nil
}
