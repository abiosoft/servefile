package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
testfiles
├── dir
│   ├── dirfile0
└── file1.txt
*/

var requests []string = []string{
	"",
	"/",
	"dir",
	"dir/",
	"dir/dirfile0",
	"file1.txt",
	"invalid",
	"dir/invalid",
}

var responses = func(path string) []byte {
	switch path {
	case "":
		fallthrough
	case "/":
		return []byte(`<pre>
<a href="dir/">dir/</a>
<a href="file1.txt">file1.txt</a>
</pre>
`)
	case "dir":
		fallthrough
	case "dir/":
		return []byte(`<pre>
<a href="dirfile0">dirfile0</a>
</pre>
`)
	case "invalid":
		fallthrough
	case "dir/invalid":
		return []byte(`404 page not found`)
	default:
		return []byte{}
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestDirHandling(t *testing.T) {
	dir := "testfiles"
	handler := fileHandler(dir)

	for _, r := range requests {
		req, err := http.NewRequest("GET", "/"+r, nil)
		check(t, err)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		expected := responses(r)
		if len(expected) == 0 {
			expected, err = ioutil.ReadFile(dir + "/" + r)
			check(t, err)
		}
		expected = bytes.TrimSpace(expected)
		value := bytes.TrimSpace(recorder.Body.Bytes())

		if !bytes.Equal(expected, value) {
			t.Fatalf("Invalid response for %v. Expected: %v, Found: %v", r, string(expected), string(value))
		}
	}
}

func TestFileHandling(t *testing.T) {
	file := "testfiles/dir/dirfile0"
	handler := fileHandler(file)

	expected, err := ioutil.ReadFile(file)
	check(t, err)
	expected = bytes.TrimSpace(expected)

	for _, r := range requests {
		req, err := http.NewRequest("GET", "/"+r, nil)
		check(t, err)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		value := bytes.TrimSpace(recorder.Body.Bytes())

		if !bytes.Equal(expected, value) {
			t.Fatalf("Invalid response for %v. Expected: %v, Found: %v", r, string(expected), string(value))
		}
	}
}
