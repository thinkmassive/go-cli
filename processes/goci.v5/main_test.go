package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "testing"
)

func TestRun(t *testing.T) {
  _, err := exec.LookPath("git")
  if err != nil {
    t.Skip("Git not installed. Skipping test.")
  }

  var testCases = []struct {
    name     string
    proj     string
    out      string
    errMsg   string
    setupGit bool
  }{
    {name: "success", proj: "./testdata/tool/",
      out:      "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCESS\nGit Push: SUCESS\n",
      errMsg:   "",
      setupGit: true},
    {name: "fail", proj: "./testdata/toolErr",
      out:      "",
      errMsg:   "'go build' failed",
      setupGit: false},
    {name: "failFormat", proj: "./testdata/toolFmtErr",
      out:      "",
      errMsg:   "'go fmt' failed",
      setupGit: false},
  }
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      if tc.setupGit {
        cleanup := setupGit(t, tc.proj)
        defer cleanup()
      }

      var out bytes.Buffer
      err := run(tc.proj, &out)

      if tc.errMsg != "" {
        if err == nil {
          t.Errorf("Expected error: %q. Got 'nil' instead.", tc.errMsg)
          return
        }

        if !strings.Contains(err.Error(), tc.errMsg) {
          t.Errorf("Expected error: %q. Got %q.", tc.errMsg, err)
        }
        return
      }

      if err != nil {
        t.Errorf("Unexpected error: %q", err)
      }

      if out.String() != tc.out {
        t.Errorf("Expected output: %q. Got %q", tc.out, out.String())
      }
    })
  }
}

func setupGit(t *testing.T, proj string) func() {
  t.Helper()

  gitExec, err := exec.LookPath("git")
  if err != nil {
    t.Fatal(err)
  }

  wd, err := os.Getwd()
  if err != nil {
    t.Fatal(err)
  }

  tempDir, err := ioutil.TempDir("", "gocitest")
  if err != nil {
    t.Fatal(err)
  }

  projPath := filepath.Join(wd, proj)
  remoteURI := fmt.Sprintf("file://%s", tempDir)

  var gitCmdList = []struct {
    args []string
    dir  string
    env  []string
  }{
    {[]string{"init", "--bare"}, tempDir, nil},
    {[]string{"init"}, projPath, nil},
    {[]string{"remote", "add", "origin", remoteURI}, projPath, nil},
    {[]string{"add", "."}, projPath, nil},
    {[]string{"commit", "-m", "test"}, projPath,
      []string{
        "GIT_COMMITTER_NAME=test",
        "GIT_COMMITTER_EMAIL=test@example.com",
        "GIT_AUTHOR_NAME=test",
        "GIT_AUTHOR_EMAIL=test@example.com",
      }},
  }

  for _, g := range gitCmdList {
    gitCmd := exec.Command(gitExec, g.args...)
    gitCmd.Dir = g.dir

    if g.env != nil {
      gitCmd.Env = append(os.Environ(), g.env...)
    }

    if err := gitCmd.Run(); err != nil {
      t.Fatal(err)
    }
  }

  return func() {
    os.RemoveAll(tempDir)
    os.RemoveAll(filepath.Join(projPath, ".git"))
  }
}
