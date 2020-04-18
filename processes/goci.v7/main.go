package main

import (
  "flag"
  "fmt"
  "io"
  "os"

  "os/signal"
  "syscall"

  "time"
)

type executer interface {
  execute() (string, error)
}

func main() {
  proj := flag.String("p", "", "Project directory")
  flag.Parse()

  if err := run(*proj, os.Stdout); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func run(proj string, out io.Writer) error {
  if proj == "" {
    return fmt.Errorf("Project directory is required")
  }

  pipeline := make([]executer, 4)
  pipeline[0] = newStep(
    "go build",
    "go",
    "Go Build: SUCCESS",
    proj,
    []string{"build", ".", "errors"},
  )

  pipeline[1] = newStep(
    "go test",
    "go",
    "Go Test: SUCCESS",
    proj,
    []string{"test", "-v"},
  )

  pipeline[2] = newExceptionStep(
    "go fmt",
    "gofmt",
    "Gofmt: SUCESS",
    proj,
    []string{"-l", "."},
  )

  pipeline[3] = newTimeoutStep(
    "git push",
    "git",
    "Git Push: SUCESS",
    proj,
    []string{"push", "origin", "master"},
    10*time.Second,
  )

  sig := make(chan os.Signal, 1)
  errCh := make(chan error)
  done := make(chan struct{})

  signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    for _, s := range pipeline {
      msg, err := s.execute()
      if err != nil {
        errCh <- err
      }

      _, err = fmt.Fprintln(out, msg)
      if err != nil {
        errCh <- err
      }
    }
    done <- struct{}{}
  }()

  for {
    select {
    case rec := <-sig:
      signal.Stop(sig)
      return fmt.Errorf("Received signal: %s. Exiting", rec)
    case err := <-errCh:
      return err
    case <-done:
      return nil
    }
  }
}
