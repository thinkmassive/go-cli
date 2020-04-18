package main

import (
  "context"
  "fmt"
  "os/exec"
  "time"
)

type timeoutStep struct {
  *step
  timeout time.Duration
}

func newTimeoutStep(name, exe, message, proj string,
  args []string, timeout time.Duration) *timeoutStep {
  s := &timeoutStep{}

  s.step = newStep(name, exe, message, proj, args)

  s.timeout = timeout
  if s.timeout == 0 {
    s.timeout = 30 * time.Second
  }

  return s
}

func (s *timeoutStep) execute() (string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
  defer cancel()

  cmd := exec.CommandContext(ctx, s.exe, s.args...)
  cmd.Dir = s.proj

  if err := cmd.Run(); err != nil {
    if ctx.Err() == context.DeadlineExceeded {
      return "", fmt.Errorf("'%s' failed: time out", s.name)
    }

    return "", fmt.Errorf("'%s' failed: %s", s.name, err)
  }

  return s.message, nil
}
