package policy

import (
	"corona-information-service/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestPolicyHandler(t *testing.T) {
	client := http.Client{}
	endpoint := httptest.NewServer(http.HandlerFunc(PolicyHandler(&client)))

	tests := []struct {
		name       string
		server     *httptest.Server
		args       string
		want       *model.Policy
		statusCode int
	}{
		{
			name: "basic-request-with-real-alpha3-no-date",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{\"policyActions\":[{\"policy_type_code\":\"NONE\",\"policy_type_display\":\"No data.  Data may be inferred for last 7 days.\",\"flag_value_display_field\":\"\",\"policy_value_display_field\":\"No data.  Data may be inferred for last 7 days.\",\"policyvalue\":0,\"flagged\":null,\"notes\":null}],\"stringencyData\":{\"msg\":\"Data unavailable\"}}"))
			})),
			args: "NOR",
			want: &model.Policy{
				CountryCode: "NOR",
				Scope:       time.Now().Format("2006-01-02"), // Will use today's date
				Stringency:  -1,
				Policies:    0,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "basic-request-with-real-alpha3-with-date",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{\"policyActions\":[{\"policy_type_code\":\"C1\",\"policy_type_display\":\"School closing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C2\",\"policy_type_display\":\"Workplace closing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C3\",\"policy_type_display\":\"Cancel public events\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C4\",\"policy_type_display\":\"Restrictions on gatherings\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No restrictions\"},{\"policy_type_code\":\"C5\",\"policy_type_display\":\"Close public transport\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C6\",\"policy_type_display\":\"Stay at home requirements\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C7\",\"policy_type_display\":\"Restrictions on internal movement\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C8\",\"policy_type_display\":\"International travel controls\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"E1\",\"policy_type_display\":\"Income support\",\"policyvalue\":\"2\",\"policyvalue_actual\":2,\"flagged\":false,\"is_general\":false,\"notes\":null,\"flag_value_display_field\":\"All workers\",\"policy_value_display_field\":\">50% lost income\"},{\"policy_type_code\":\"E2\",\"policy_type_display\":\"Debt/contract relief\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"None\"},{\"policy_type_code\":\"H1\",\"policy_type_display\":\"Public information campaigns\",\"policyvalue\":\"2\",\"policyvalue_actual\":2,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Coordinated public campaign\"},{\"policy_type_code\":\"H2\",\"policy_type_display\":\"Testing policy\",\"policyvalue\":\"3\",\"policyvalue_actual\":3,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"Generally available\"},{\"policy_type_code\":\"H3\",\"policy_type_display\":\"Contact tracing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"No contact tracing\"},{\"policy_type_code\":\"H5\",\"policy_type_display\":\"Investment in vaccines\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"USD Value\"},{\"policy_type_code\":\"H6\",\"policy_type_display\":\"Facial Coverings\",\"policyvalue\":\"1\",\"policyvalue_actual\":1,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Recommended\"},{\"policy_type_code\":\"H7\",\"policy_type_display\":\"Vaccination policy\",\"policyvalue\":\"5\",\"policyvalue_actual\":5,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"Government funded\",\"policy_value_display_field\":\"Universal\"},{\"policy_type_code\":\"H8\",\"policy_type_display\":\"Protection of elderly people\",\"policyvalue\":\"1\",\"policyvalue_actual\":1,\"flagged\":false,\"is_general\":false,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Recommended protections\"},{\"policy_type_code\":\"V1\",\"policy_type_display\":\"Vaccine Prioritisation\",\"policyvalue\":\"-3\",\"policyvalue_actual\":-3,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V2\",\"policy_type_display\":\"Vaccine Availability\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V3\",\"policy_type_display\":\"Vaccine Financial Support\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V4\",\"policy_type_display\":\"Mandatory Vaccination\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null}],\"stringencyData\":{\"date_value\":\"2022-03-07\",\"country_code\":\"NOR\",\"confirmed\":1315692,\"deaths\":1664,\"stringency_actual\":11.11,\"stringency\":11.11}}"))
			})),
			args: "NOR?scope=2022-03-07",
			want: &model.Policy{
				CountryCode: "NOR",
				Scope:       "2022-03-07",
				Stringency:  11.11,
				Policies:    21,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "basic-request-with-invalid-alpha3",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{\"policyActions\":[{\"policy_type_code\":\"C1\",\"policy_type_display\":\"School closing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C2\",\"policy_type_display\":\"Workplace closing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C3\",\"policy_type_display\":\"Cancel public events\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C4\",\"policy_type_display\":\"Restrictions on gatherings\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No restrictions\"},{\"policy_type_code\":\"C5\",\"policy_type_display\":\"Close public transport\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C6\",\"policy_type_display\":\"Stay at home requirements\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No measures\"},{\"policy_type_code\":\"C7\",\"policy_type_display\":\"Restrictions on internal movement\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"C8\",\"policy_type_display\":\"International travel controls\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"No Measures\"},{\"policy_type_code\":\"E1\",\"policy_type_display\":\"Income support\",\"policyvalue\":\"2\",\"policyvalue_actual\":2,\"flagged\":false,\"is_general\":false,\"notes\":null,\"flag_value_display_field\":\"All workers\",\"policy_value_display_field\":\">50% lost income\"},{\"policy_type_code\":\"E2\",\"policy_type_display\":\"Debt/contract relief\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"None\"},{\"policy_type_code\":\"H1\",\"policy_type_display\":\"Public information campaigns\",\"policyvalue\":\"2\",\"policyvalue_actual\":2,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Coordinated public campaign\"},{\"policy_type_code\":\"H2\",\"policy_type_display\":\"Testing policy\",\"policyvalue\":\"3\",\"policyvalue_actual\":3,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"Generally available\"},{\"policy_type_code\":\"H3\",\"policy_type_display\":\"Contact tracing\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"No contact tracing\"},{\"policy_type_code\":\"H5\",\"policy_type_display\":\"Investment in vaccines\",\"policyvalue\":\"0\",\"policyvalue_actual\":0,\"flagged\":null,\"is_general\":null,\"notes\":null,\"policy_value_display_field\":\"USD Value\"},{\"policy_type_code\":\"H6\",\"policy_type_display\":\"Facial Coverings\",\"policyvalue\":\"1\",\"policyvalue_actual\":1,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Recommended\"},{\"policy_type_code\":\"H7\",\"policy_type_display\":\"Vaccination policy\",\"policyvalue\":\"5\",\"policyvalue_actual\":5,\"flagged\":true,\"is_general\":true,\"notes\":null,\"flag_value_display_field\":\"Government funded\",\"policy_value_display_field\":\"Universal\"},{\"policy_type_code\":\"H8\",\"policy_type_display\":\"Protection of elderly people\",\"policyvalue\":\"1\",\"policyvalue_actual\":1,\"flagged\":false,\"is_general\":false,\"notes\":null,\"flag_value_display_field\":\"General\",\"policy_value_display_field\":\"Recommended protections\"},{\"policy_type_code\":\"V1\",\"policy_type_display\":\"Vaccine Prioritisation\",\"policyvalue\":\"-3\",\"policyvalue_actual\":-3,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V2\",\"policy_type_display\":\"Vaccine Availability\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V3\",\"policy_type_display\":\"Vaccine Financial Support\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null},{\"policy_type_code\":\"V4\",\"policy_type_display\":\"Mandatory Vaccination\",\"policyvalue\":\"25,126\",\"policyvalue_actual\":25126,\"flagged\":null,\"is_general\":null,\"notes\":null}],\"stringencyData\":{\"date_value\":\"2022-03-07\",\"country_code\":\"NOR\",\"confirmed\":1315692,\"deaths\":1664,\"stringency_actual\":11.11,\"stringency\":11.11}}"))
			})),
			args:       "chipmunk",
			want:       &model.Policy{},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			model.STRINGENCY_URL = tt.server.URL + "/%s/%s" //If not added, tests break

			request, _ := http.NewRequest(http.MethodGet, endpoint.URL+model.POLICY_PATH+tt.args, nil)
			res, _ := client.Do(request)

			decoder := json.NewDecoder(res.Body)

			var p model.Policy

			_ = decoder.Decode(&p)

			if !reflect.DeepEqual(tt.want, &p) {
				t.Errorf("FAILED: expected %v, got %v\n", tt.want, &p)
			}
			if !reflect.DeepEqual(tt.statusCode, res.StatusCode) {
				t.Errorf("FAILED: expected %v, got %v\n", tt.statusCode, res.StatusCode)
			}

		})
	}
}
