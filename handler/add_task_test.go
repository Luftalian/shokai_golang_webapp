package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luftalian/shokai_golang_webapp/entity"
	"github.com/Luftalian/shokai_golang_webapp/store"
	"github.com/Luftalian/shokai_golang_webapp/testutil"

	"github.com/go-playground/validator/v10"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want want
	}{
		"ok": {
			reqFile: "testdata/ok_req.json.golden",
			want: want{
				status: http.StatusOK,
				rspFile: "testdata/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/bad_req.json.golden",
			want: want{
				status: http.StatusBadRequest,
				rspFile: "testdata/bad_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T){
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}

}