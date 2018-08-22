package configuration

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNewService(t *testing.T) {
	config := Config{
		Host:  "host",
		Token: "token",
	}

	testService := service{
		config: config,
	}

	service := NewService(config)

	assert.Equal(t, testService, service)
}

func Test_service_Consents(t *testing.T) {
	service := service{
		config: Config{
			Host: "https://api.globaldatacompany.com/configuration/v1",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.globaldatacompany.com/configuration/v1/consents/Identity%20Verification/AU",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`[
    "Australia Driver Licence",
    "Australia Passport",
    "Birth Registry",
    "Visa Verification",
    "DVS Driver License Search",
    "DVS Medicare Search",
    "DVS Passport Search",
    "DVS Visa Search",
    "DVS ImmiCard Search",
    "DVS Citizenship Certificate Search",
    "DVS Certificate of Registration by Descent Search",
    "Credit Agency"
]`,
			), nil
		},
	)

	consents, err := service.Consents("AU")
	if assert.NoError(t, err) {
		assert.Equal(t, Consents{
			"Australia Driver Licence",
			"Australia Passport",
			"Birth Registry",
			"Visa Verification",
			"DVS Driver License Search",
			"DVS Medicare Search",
			"DVS Passport Search",
			"DVS Visa Search",
			"DVS ImmiCard Search",
			"DVS Citizenship Certificate Search",
			"DVS Certificate of Registration by Descent Search",
			"Credit Agency",
		}, consents)
	}
}

func Test_service_Consents_Error(t *testing.T) {
	service := service{
		config: Config{
			Host: "https://api.globaldatacompany.com/configuration/v1",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.globaldatacompany.com/configuration/v1/consents/Identity%20Verification/AU",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusBadRequest,
				`{
    "Message": "Country code does not exist"
}`,
			), nil
		},
	)

	consents, err := service.Consents("AU")
	assert.Nil(t, consents)
	assert.Error(t, err)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.globaldatacompany.com/configuration/v1/consents/Identity%20Verification/AU",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusBadRequest,
				`{}`,
			), nil
		},
	)

	consents, err = service.Consents("AU")
	assert.Nil(t, consents)
	assert.Error(t, err)
	assert.Equal(t, "Unknown error", err.Error())

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.globaldatacompany.com/configuration/v1/consents/Identity%20Verification/AU",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusBadRequest,
				`{`,
			), nil
		},
	)

	consents, err = service.Consents("AU")
	assert.Nil(t, consents)
	assert.Error(t, err)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.globaldatacompany.com/configuration/v1/consents/Identity%20Verification/AU",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	consents, err = service.Consents("AU")
	assert.Nil(t, consents)
	assert.Error(t, err)

	consents, err = service.Consents("")
	assert.Nil(t, consents)
	assert.Error(t, err)
}
