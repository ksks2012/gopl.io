package reader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestStringReader_Read_Basic(t *testing.T) {
	r := NewReader("hello")
	buf := make([]byte, 5)
	n, err := r.Read(buf)
	if n != 5 {
		t.Errorf("expected n=5, got %d", n)
	}
	if err != nil && err != io.EOF {
		t.Errorf("expected err=nil or io.EOF, got %v", err)
	}
	if string(buf) != "hello" {
		t.Errorf("expected buf='hello', got '%s'", string(buf))
	}
}

func TestStringReader_Read_Partial(t *testing.T) {
	r := NewReader("hello")
	buf := make([]byte, 2)
	n, err := r.Read(buf)
	if n != 2 || err != nil {
		t.Errorf("expected n=2, err=nil, got n=%d, err=%v", n, err)
	}
	if string(buf) != "he" {
		t.Errorf("expected buf='he', got '%s'", string(buf))
	}

	n, err = r.Read(buf)
	if n != 2 || err != nil {
		t.Errorf("expected n=2, err=nil, got n=%d, err=%v", n, err)
	}
	if string(buf) != "ll" {
		t.Errorf("expected buf='ll', got '%s'", string(buf))
	}

	n, err = r.Read(buf)
	if n != 1 || err != io.EOF {
		t.Errorf("expected n=1, err=io.EOF, got n=%d, err=%v", n, err)
	}
	if string(buf[:1]) != "o" {
		t.Errorf("expected buf='o', got '%s'", string(buf[:1]))
	}
}

func TestStringReader_Read_EmptyString(t *testing.T) {
	r := NewReader("")
	buf := make([]byte, 10)
	n, err := r.Read(buf)
	if n != 0 {
		t.Errorf("expected n=0, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected err=io.EOF, got %v", err)
	}
}

func TestStringReader_Read_AfterEOF(t *testing.T) {
	r := NewReader("a")
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if n != 1 || err != io.EOF {
		t.Errorf("expected n=1, err=io.EOF, got n=%d, err=%v", n, err)
	}
	n, err = r.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF after EOF, got n=%d, err=%v", n, err)
	}
}

func redirectStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = oldStdout
	return <-outC
}

func TestForEachNode(t *testing.T) {
	htmlContent := `
<html>
<head>
    <title>Test</title>
</head>
<body>
    <h1>Hello</h1>
    <p>World</p>
</body>
</html>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("html.Parse failed: %v", err)
	}

	expectedOutput := strings.TrimSpace(`
<html>
  <head>
    <title>
    </title>
  </head>
  <body>
    <h1>
    </h1>
    <p>
    </p>
  </body>
</html>`)

	output := redirectStdout(func() {
		depth := 0
		startElement := func(n *html.Node) {
			if n.Type == html.ElementNode {
				fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
				depth++
			}
		}
		endElement := func(n *html.Node) {
			if n.Type == html.ElementNode {
				depth--
				fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
			}
		}
		forEachNode(doc, startElement, endElement)
	})

	if strings.TrimSpace(output) != expectedOutput {
		t.Errorf("TestForEachNode: unexpected output\nExpected:\n%s\nActual:\n%s", expectedOutput, strings.TrimSpace(output))
	}
}

func TestOutline(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
<html>
<head>
    <title>Example Page</title>
</head>
<body>
    <div id="main">
        <p>This is a test.</p>
        <ul>
            <li>Item 1</li>
            <li>Item 2</li>
        </ul>
    </div>
</body>
</html>`)
	}))
	defer ts.Close()

	expectedOutput := strings.TrimSpace(`
<html>
  <head>
    <title>
    </title>
  </head>
  <body>
    <div>
      <p>
      </p>
      <ul>
        <li>
        </li>
        <li>
        </li>
      </ul>
    </div>
  </body>
</html>
`)

	var outlineErr error
	output := redirectStdout(func() {
		resp, err := http.Get(ts.URL)
		if err != nil {
			outlineErr = fmt.Errorf("http.Get failed: %v", err)
			return
		}
		defer resp.Body.Close()
		outlineErr = outline(resp.Body)
	})
	if outlineErr != nil {
		t.Fatalf("outline failed: %v", outlineErr)
	}

	if strings.TrimSpace(output) != expectedOutput {
		t.Errorf("TestOutline: unexpected output\nExpected:\n%s\nActual:\n%s", expectedOutput, strings.TrimSpace(output))
	}

	t.Run("Invalid URL", func(t *testing.T) {
		r := strings.NewReader("not html")
		err := outline(r)
		if err != nil {
			t.Errorf("expected no error for invalid HTML, got %v", err)
		}
	})
}

func TestLimitedReader_Read_Basic(t *testing.T) {
	src := strings.NewReader("abcdef")
	lr := LimitReader(src, 4)
	buf := make([]byte, 10)
	n, err := lr.Read(buf)
	if n != 4 {
		t.Errorf("expected n=4, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected err=io.EOF, got %v", err)
	}
	if string(buf[:n]) != "abcd" {
		t.Errorf("expected buf='abcd', got '%s'", string(buf[:n]))
	}

	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF after limit reached, got n=%d, err=%v", n, err)
	}
}

func TestLimitedReader_Read_Partial(t *testing.T) {
	src := strings.NewReader("abcdef")
	lr := LimitReader(src, 3)
	buf := make([]byte, 2)
	n, err := lr.Read(buf)
	if n != 2 || err != nil {
		t.Errorf("expected n=2, err=nil, got n=%d, err=%v", n, err)
	}
	if string(buf[:n]) != "ab" {
		t.Errorf("expected buf='ab', got '%s'", string(buf[:n]))
	}

	n, err = lr.Read(buf)
	if n != 1 || err != io.EOF {
		t.Errorf("expected n=1, err=io.EOF, got n=%d, err=%v", n, err)
	}
	if string(buf[:n]) != "c" {
		t.Errorf("expected buf='c', got '%s'", string(buf[:n]))
	}

	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF after EOF, got n=%d, err=%v", n, err)
	}
	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF on subsequent reads after EOF, got n=%d, err=%v", n, err)
	}
}

func TestLimitedReader_Read_ZeroLimit(t *testing.T) {
	src := strings.NewReader("abcdef")
	lr := LimitReader(src, 0)
	buf := make([]byte, 5)
	n, err := lr.Read(buf)
	if n != 0 {
		t.Errorf("expected n=0, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected err=io.EOF, got %v", err)
	}
	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF on subsequent reads after EOF, got n=%d, err=%v", n, err)
	}
}

func TestLimitedReader_Read_ExactLimit(t *testing.T) {
	src := strings.NewReader("abcdef")
	lr := LimitReader(src, 6)
	buf := make([]byte, 6)
	n, err := lr.Read(buf)
	if n != 6 {
		t.Errorf("expected n=6, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected err=io.EOF, got %v", err)
	}
	if string(buf[:n]) != "abcdef" {
		t.Errorf("expected buf='abcdef', got '%s'", string(buf[:n]))
	}
	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF on subsequent reads after EOF, got n=%d, err=%v", n, err)
	}
}

func TestLimitedReader_Read_MultipleSmallReads(t *testing.T) {
	src := strings.NewReader("abcdef")
	lr := LimitReader(src, 4)
	buf := make([]byte, 2)
	n, err := lr.Read(buf)
	if n != 2 || err != nil {
		t.Errorf("expected n=2, err=nil, got n=%d, err=%v", n, err)
	}
	if string(buf[:n]) != "ab" {
		t.Errorf("expected buf='ab', got '%s'", string(buf[:n]))
	}
	n, err = lr.Read(buf)
	if n != 2 || err != io.EOF {
		t.Errorf("expected n=2, err=io.EOF, got n=%d, err=%v", n, err)
	}
	if string(buf[:n]) != "cd" {
		t.Errorf("expected buf='cd', got '%s'", string(buf[:n]))
	}
	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Errorf("expected n=0, err=io.EOF after EOF, got n=%d, err=%v", n, err)
	}

	t.Run("UnderlyingReaderEOFBeforeLimit", func(t *testing.T) {
		srcShort := strings.NewReader("xy")
		lrShort := LimitReader(srcShort, 5)

		n, err := lrShort.Read(buf)
		if n != 2 || err != nil {
			t.Errorf("expected n=2, err=nil, got n=%d, err=%v", n, err)
		}
		if string(buf[:n]) != "xy" {
			t.Errorf("expected buf='xy', got '%s'", string(buf[:n]))
		}

		n, err = lrShort.Read(buf)
		if n != 0 || err != io.EOF {
			t.Errorf("expected n=0, err=io.EOF, got n=%d, err=%v", n, err)
		}
	})
}
