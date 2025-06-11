package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain_FindsElementByID(t *testing.T) {
	htmlContent := `
	<html>
		<body>
			<div id="foo">Hello</div>
			<span id="bar">World</span>
		</body>
	</html>`
	tmpfile, err := ioutil.TempFile("", "test*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(htmlContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cmd := exec.Command(os.Args[0], tmpfile.Name(), "foo")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	// Use a helper process to avoid recursion
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		main()
		os.Exit(0)
	}
	err = cmd.Run()
	if err != nil {
		t.Fatalf("cmd.Run() failed: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Found element with ID 'foo': <div id=\"foo\">") {
		t.Errorf("expected found message, got: %q", got)
	}
}

func TestMain_ElementNotFound(t *testing.T) {
	htmlContent := `<html><body><div id="foo">Hello</div></body></html>`
	tmpfile, err := ioutil.TempFile("", "test*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(htmlContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cmd := exec.Command(os.Args[0], tmpfile.Name(), "notfound")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		main()
		os.Exit(0)
	}
	err = cmd.Run()
	if err != nil {
		t.Fatalf("cmd.Run() failed: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Element with ID 'notfound' not found.") {
		t.Errorf("expected not found message, got: %q", got)
	}
}

func TestMain_MissingArguments(t *testing.T) {
	cmd := exec.Command(os.Args[0])
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		main()
		os.Exit(0)
	}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for missing arguments, got none")
	}
	got := out.String()
	if !strings.Contains(got, "Usage:") {
		t.Errorf("expected usage message, got: %q", got)
	}
}
