package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultReturn_b7d4f8063b(t *testing.T) {
	type args struct {
		resp interface{}
		err  error
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "Success case",
			args: args{
				resp: map[string]string{"message": "success"},
				err:  nil,
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"success"}`,
		},
		{
			name: "Error case",
			args: args{
				resp: nil,
				err:  errors.New("internal error"),
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defaultReturn(w, r, tt.args.resp, tt.args.err)
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				body := rr.Body.String()
				body = strings.TrimSuffix(body, "\n")
				if body != tt.wantBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						body, tt.wantBody)
				}
			}
		})
	}
}
