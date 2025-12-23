package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	apiJSON = `{"text": "the quick brown; fox jumps over; the lazy dog", "spinestring": "cra"}`
)

func TestSetupMux(t *testing.T) {
	sp := &ServePoems{}
	mux := sp.SetupMux()

	tests := []struct {
		name     string
		target   string
		method   string
		wantCode int
		expect   string
		jsonbody string
	}{
		{name: "Healthz endpoint answers",
			target:   "/healthz",
			method:   "GET",
			wantCode: http.StatusOK,
			expect:   "ok",
			jsonbody: "",
		},
		{name: "Homepage frontend answers",
			target:   "/",
			method:   "GET",
			wantCode: http.StatusOK,
			expect:   "",
			jsonbody: "",
		},
		{name: "Errors on incorrect JSON body",
			target:   "/app",
			method:   "POST",
			wantCode: http.StatusBadRequest,
			expect:   "empty",
			jsonbody: testApodJSON, // does not have DataAPI fields
		},
		{name: "Errors on unreadable request body",
			target:   "/app",
			method:   "POST",
			wantCode: http.StatusBadRequest,
			expect:   "",
			jsonbody: strings.Repeat("a", 2*1024*1024),
		},
		{name: "Errors on unmarshall-able JSON body",
			target:   "/app",
			method:   "POST",
			wantCode: http.StatusInternalServerError,
			expect:   "",
			jsonbody: "not_json",
		},
		{name: "Retrieves mesostic",
			target:   "/app",
			method:   "POST",
			wantCode: http.StatusOK,
			expect:   "quiCk",
			jsonbody: apiJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.jsonbody))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			defer r.Body.Close()

			assertStatus(t, w.Code, tt.wantCode)
			assertStringContains(t, w.Body.String(), tt.expect)
		})
	}
}
