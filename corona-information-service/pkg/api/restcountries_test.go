package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCountryNameByAlphaCode(t *testing.T) {
	type args struct {
		alpha3 string
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic-request-with-real-alpha3",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{\"name\":{\"common\":\"Norway\",\"official\":\"Kingdom of Norway\",\"nativeName\":{\"nno\":{\"official\":\"Kongeriket Noreg\",\"common\":\"Noreg\"},\"nob\":{\"official\":\"Kongeriket Norge\",\"common\":\"Norge\"},\"smi\":{\"official\":\"Norgga gonagasriika\",\"common\":\"Norgga\"}}}}"))
			})),
			args:    args{alpha3: "NOR"},
			want:    "Norway",
			wantErr: false,
		},
		{
			name: "basic-request-with-fake-alpha3",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{\"status\":404,\"message\":\"Not Found\"}"))
			})),
			args:    args{alpha3: "test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCountryNameByAlphaCode(tt.args.alpha3)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCountryNameByAlphaCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCountryNameByAlphaCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
