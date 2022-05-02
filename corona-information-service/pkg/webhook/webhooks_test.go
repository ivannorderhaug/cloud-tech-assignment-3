package webhook

import (
	"bytes"
	"corona-information-service/internal/model"
	"corona-information-service/pkg/db"
	"corona-information-service/tools/hash"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	COLLECTION = "testCollection"

	countries := []string{"Norway", "Denmark", "Sweden", "Germany", "Finland"}

	err := db.InitializeFirestore("../../service-account.json")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= 4; i++ {
		wh := model.Webhook{
			ID:          countries[i][:3],
			Url:         "http://localhost:8081",
			Country:     countries[i],
			Calls:       3,
			ActualCalls: 0,
		}
		_ = db.AddToFirestore(hash.Hash("testCollection"), hash.Hash(wh.ID), wh)
	}

	m.Run()
}

func TestGetAllWebhooks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "get-all-webhooks",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllWebhooks()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWebhooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetWebhook(t *testing.T) {
	type args struct {
		webhookId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "get-existing-webhook",
			args:    args{webhookId: "Nor"},
			wantErr: false,
		},
		{
			name:    "get-non-existing-webhook",
			args:    args{webhookId: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetWebhook(tt.args.webhookId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRegisterWebhook(t *testing.T) {
	//Make sure doc does not exist beforehand
	webhooks, _ := GetAllWebhooks()
	for _, wh := range webhooks {
		if wh.Country == "Iceland" {
			db.DeleteSingleDocumentFromFirestore(hash.Hash(COLLECTION), hash.Hash(wh.ID))
		}
	}

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "register-webhook-correctly",
			args: args{r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(
					`{
							"country": "Iceland",
							"url":"http://localhost:8081",
							"calls": 4
						}`)),
			},
			},
			want:    http.StatusCreated,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := RegisterWebhook(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	type args struct {
		webhookId string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "delete-existing-webhook",
			args: args{webhookId: "Nor"},
			want: true,
		},
		{
			name: "delete-non-existing-webhook",
			args: args{webhookId: "test"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteWebhook(tt.args.webhookId); got != tt.want {
				t.Errorf("DeleteWebhook() = %v, want %v", got, tt.want)
			}
		})
	}
}
