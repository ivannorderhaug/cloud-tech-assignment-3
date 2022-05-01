package cases

import (
	"bytes"
	"corona-information-service/internal/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCaseHandler(t *testing.T) {
	client := http.Client{}
	endpoint := httptest.NewServer(http.HandlerFunc(CaseHandler(&client)))

	tests := []struct {
		name       string
		server     *httptest.Server
		args       string
		want       *model.Case
		statusCode int
	}{
		{
			name: "basic-request-with-existing-country",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//This is the nested structure that gets returned when doing a graphql request
				_, _ = w.Write([]byte("{\n  \"data\": {\n    \"country\": {\n      \"name\": \"Norway\",\n      \"mostRecent\": {\n        \"date\": \"2022-04-05\",\n        \"confirmed\": 1411550,\n        \"recovered\": 0,\n        \"deaths\": 2518,\n        \"growthRate\": 0.0010630821154695824\n      }\n    }\n  }\n}"))
			})),
			args: "Norway",
			want: &model.Case{
				Country:        "Norway",
				Date:           "2022-04-05",
				ConfirmedCases: 1411550,
				Recovered:      0,
				Deaths:         2518,
				GrowthRate:     0.0010630821154695824,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "basic-request-with-non-existing-country ",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("{\n  \"errors\": [\n    {\n      \"message\": \"Couldn't find data from country test\",\n      \"locations\": [\n        {\n          \"line\": 2,\n          \"column\": 3\n        }\n      ],\n      \"path\": [\n        \"country\"\n      ],\n      \"extensions\": {\n        \"code\": \"INTERNAL_SERVER_ERROR\"\n      }\n    }\n  ],\n  \"data\": {\n    \"country\": null\n  }\n}"))
			})),
			args:       "test",
			want:       &model.Case{},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			model.CASES_URL = tt.server.URL

			request, _ := http.NewRequest(http.MethodGet, endpoint.URL+model.CASE_PATH+tt.args, nil)
			res, _ := client.Do(request)

			decoder := json.NewDecoder(res.Body)

			var c model.Case

			_ = decoder.Decode(&c)

			if !reflect.DeepEqual(tt.want, &c) {
				t.Errorf("FAILED: expected %v, got %v\n", tt.want, &c)
			}
			if !reflect.DeepEqual(tt.statusCode, res.StatusCode) {
				t.Errorf("FAILED: expected %v, got %v\n", tt.statusCode, res.StatusCode)
			}

		})
	}
}
func Test_getCase(t *testing.T) {
	type args struct {
		res *http.Response
	}
	tests := []struct {
		name  string
		args  args
		want  *model.Case
		want1 error
		want2 int
	}{
		{
			name: "existing-country",
			args: args{
				res: &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString("{\n  \"data\": {\n    \"country\": {\n      \"name\": \"Norway\",\n      \"mostRecent\": {\n        \"date\": \"2022-04-05\",\n        \"confirmed\": 1411550,\n        \"recovered\": 0,\n        \"deaths\": 2518,\n        \"growthRate\": 0.0010630821154695824\n      }\n    }\n  }\n}")),
					StatusCode: http.StatusOK,
				},
			},
			want: &model.Case{
				Country:        "Norway",
				Date:           "2022-04-05",
				ConfirmedCases: 1411550,
				Recovered:      0,
				Deaths:         2518,
				GrowthRate:     0.0010630821154695824,
			},
			want1: nil,
			want2: 0,
		},
		{
			name: "non-existing-country",
			args: args{
				res: &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString("{\n  \"errors\": [\n    {\n      \"message\": \"Couldn't find data from country test\",\n      \"locations\": [\n        {\n          \"line\": 2,\n          \"column\": 3\n        }\n      ],\n      \"path\": [\n        \"country\"\n      ],\n      \"extensions\": {\n        \"code\": \"INTERNAL_SERVER_ERROR\"\n      }\n    }\n  ],\n  \"data\": {\n    \"country\": null\n  }\n}")),
					StatusCode: http.StatusOK,
				},
			},
			want:  &model.Case{},
			want1: errors.New("could not find a country with that name. the external api is very sensitive, please see README for more information"),
			want2: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getCase(tt.args.res)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCase() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getCase() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getCase() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
