package sampler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResponseTime(t *testing.T) {
	const responseTimeJSON = `{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				{
					"metric": { "destination_service_name": "productservice" },
					"value": [ 1678886400, "150.2" ]
				},
				{
					"metric": { "destination_service_name": "shippingservice" },
					"value": [ 1678886400, "75.0" ]
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, responseTimeJSON)
	}))

	defer server.Close()

	client, err := api.NewClient(api.Config{
		Address: server.URL,
	})

	require.NoError(t, err, "failed to create the client")

	promAPI := v1.NewAPI(client)
	provider := NewPrometheusProvider(promAPI)

	actual, err := provider.getResponseTime(context.Background())
	expected := []responseTime{
		{serviceName: "productservice", value: 150.2},
		{serviceName: "shippingservice", value: 75.0},
	}

	assert.NoError(t, err)
	assert.Len(t, actual, 2)
	assert.Equal(t, expected, actual)
}

func TestGetThroughput(t *testing.T) {
	const throughputJSON = `{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				{
					"metric": { "destination_service_name": "productservice" },
					"value": [ 1678886400, "150.2" ]
				},
				{
					"metric": { "destination_service_name": "shippingservice" },
					"value": [ 1678886400, "75.0" ]
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, throughputJSON)
	}))

	defer server.Close()

	client, err := api.NewClient(api.Config{
		Address: server.URL,
	})

	require.NoError(t, err, "failed to create the client")

	promAPI := v1.NewAPI(client)
	provider := NewPrometheusProvider(promAPI)

	actual, err := provider.getThroughput(context.Background())
	expected := []throughput{
		{serviceName: "productservice", value: 150.2},
		{serviceName: "shippingservice", value: 75.0},
	}

	assert.NoError(t, err)
	assert.Len(t, actual, 2)
	assert.Equal(t, expected, actual)
}
