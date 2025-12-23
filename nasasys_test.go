package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	testApodJSON = `
{
    "date": "2000-01-01",
    "explanation": "Welcome to the millennial year at the threshold of millennium three.  During millennium two, humanity continually redefined its concept of \"Universe\": first as spheres centered on the Earth, in mid-millennium as the Solar System, a few centuries ago as the Galaxy, and within the last century as the matter emanating from the Big Bang.  During millennium three humanity may hope to discover alien life, to understand the geometry and composition of our present concept of Universe, and even to travel through this Universe.  Whatever our accomplishments, humanity will surely find adventure and discovery in the space above and beyond, and possibly define the surrounding Universe in ways and colors we cannot yet imagine by the threshold of millennium four.",
    "hdurl": "https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor.gif",
    "media_type": "image",
    "service_version": "v1",
    "title": "The Millennium that Defines Universe",
    "url": "https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor_big.gif"
}
`
	mesosticApod = `
                 welcome To
                         Humanity continually r
             first as sphEr
                      in M
             a few centurIes ago as the ga
          and within the Last century as the matter emanating from the big bang
                during miL
                   to undErsta
                        aNd eve
whatever our accomplishmeNts
                    humanIty w
and possibly define the sUrro`
)

func TestGetAPOD(t *testing.T) {
	mockWWW := makeMockWebServBody(0*time.Millisecond, testApodJSON)
	mockFS := MockFS{}

	got, err := GetAPOD(mockWWW.URL, mockFS)
	assertError(t, err, nil)
	assertStringContains(t, got, mesosticApod)
	t.Log(got)
}

// TestSingleFetch should handle single URLs
func TestSingleFetch(t *testing.T) {
	mockWWW := makeMockWebServBody(0*time.Millisecond, "craquemattic")
	urlWWW := mockWWW.URL

	t.Run("Fetches a single URL", func(t *testing.T) {
		want := "craquemattic"
		_, get, err := SingleFetch(urlWWW)

		got := string(get)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})

	t.Run("Returns Status 200", func(t *testing.T) {
		got, _, _ := SingleFetch(urlWWW)
		assertStatus(t, got, 200)
	})

	// Close this mock server to run additional tests
	mockWWW.Close()

	t.Run("Returns Error after Server Close", func(t *testing.T) {
		_, _, got := SingleFetch(urlWWW)
		assertGotError(t, got)
		fmt.Println(got)
	})

	t.Run("Returns Error after Host Unreachable", func(t *testing.T) {
		_, _, err := SingleFetch("http://badhost:4420")
		assertGotError(t, err)
		assertStringContains(t, err.Error(), "no such host")
	})

	t.Run("Returns 500 Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", 500)
		}))
		defer server.Close()

		statusCode, _, err := SingleFetch(server.URL)
		assertStatus(t, statusCode, 500)
		assertError(t, err, nil)
	})

	/* No tests for the extremely robust io.ReadAll, it is very difficult to break. */
}

func TestSingleFetchWithClient_Timeout(t *testing.T) {
	c := &mockTimeoutClient{}
	_, _, err := SingleFetchWithClient("http://test", c)
	assertGotError(t, err)
}

/// Helpers

// MockFS for dependency injection on FileSystem to test config file
type MockFS struct {
	OpenError bool
	StatError bool
	FileSize  int64
}

func (fs MockFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile("/dev/null", data, perm)
}

func (fs MockFS) Open(name string) (*os.File, error) {
	if fs.OpenError {
		return nil, errors.New("mock: could not open file")
	}

	if runtime.GOOS == "windows" {
		return os.Open("NUL")
	}
	return os.Open("/dev/null")
}

func (fs MockFS) Stat(name string) (os.FileInfo, error) {
	if fs.StatError {
		return nil, errors.New("mock: could not stat file")
	}

	return &MockFileInfo{size: fs.FileSize}, nil
}

type MockFileInfo struct {
	size int64
}

func (fi MockFileInfo) Size() int64        { return fi.size }
func (fi MockFileInfo) Name() string       { return "mock-file" }
func (fi MockFileInfo) Mode() os.FileMode  { return 0644 }
func (fi MockFileInfo) ModTime() time.Time { return time.Now() }
func (fi MockFileInfo) IsDir() bool        { return false }
func (fi MockFileInfo) Sys() interface{}   { return nil }

// Mock web server configurations
type mockTimeoutClient struct{}

func (mc *mockTimeoutClient) Get(u string) (*http.Response, error) {
	response := &url.Error{
		Op:  "Get",
		URL: u,
		Err: context.DeadlineExceeded,
	}

	return nil, response
}

type FailingReader struct {
	data      []byte
	position  int
	failAfter int
}

func (fr *FailingReader) Read(p []byte) (n int, err error) {
	if fr.position >= fr.failAfter {
		return 0, fmt.Errorf("simulated I/O error")
	}

	remaining := len(fr.data) - fr.position
	if remaining == 0 {
		return 0, io.EOF
	}

	toCopy := len(p)
	if toCopy > remaining {
		toCopy = remaining
	}

	copy(p, fr.data[fr.position:fr.position+toCopy])
	fr.position += toCopy
	return toCopy, nil
}

// Mock responder for external API calls with configurable body content
func makeMockWebServBody(delay time.Duration, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testAnswer := []byte(body)
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write(testAnswer)
		if err != nil {
			log.Fatalf("ERROR: Could not write to output.")
		}
	}))
}

// Assertions
func assertError(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertGotError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Errorf("Expected an error but got %q", got)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct value, got %d, want %d", got, want)
	}
}

func assertInt64(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct value, got %d, want %d", got, want)
	}
}

func assertStringContains(t *testing.T, full, want string) {
	t.Helper()
	if !strings.Contains(full, want) {
		t.Errorf("Did not find %q, expected string contains %q", want, full)
	}
}
