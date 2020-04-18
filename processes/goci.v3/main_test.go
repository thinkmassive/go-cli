package main

import (
  "bytes"
  "strings"
  "testing"
)

func TestRun(t *testing.T) {
  var testCases = []struct {
    name   string
    proj   string
    out    string
    errMsg string
  }{
    {name: "success", proj: "./testdata/tool/",
      out:    "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCESS\n",
      errMsg: ""},
    {name: "fail", proj: "./testdata/toolErr",
      out:    "",
      errMsg: "'go build' failed"},
    {name: "failFormat", proj: "./testdata/toolFmtErr",
      out:    "",
      errMsg: "'go fmt' failed"},
  }
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
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
