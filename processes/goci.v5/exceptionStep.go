package main

import (
  "bytes"
  "fmt"
  "os/exec"
)

type exceptionStep struct {
  *step
}

func newExceptionStep(name, exe, message, proj string, args []string) *exceptionStep {
  s := &exceptionStep{}

  s.step = newStep(name, exe, message, proj, args)

  return s
}

func (s *exceptionStep) execute() (string, error) {
  cmd := exec.Command(s.exe, s.args...)

  var out bytes.Buffer
  cmd.Stdout = &out

  cmd.Dir = s.proj

  if err := cmd.Run(); err != nil {
    return "", err
  }

  if out.Len() > 0 {
    return "", fmt.Errorf("'%s' failed: %s", s.name, out.String())
  }

  return s.message, nil
}
