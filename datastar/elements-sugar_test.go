package datastar

import (
	"testing"
)

func TestActionSSE(t *testing.T) {
	type testCase struct {
		Expected string
		Actual   string
	}

	tests := []testCase{
		{Expected: "@get('/sse')", Actual: SSEGet("/sse")},
		{Expected: "@post('/sse')", Actual: SSEPost("/sse")},
		{Expected: "@put('/sse')", Actual: SSEPut("/sse")},
		{Expected: "@patch('/sse')", Actual: SSEPatch("/sse")},
		{Expected: "@delete('/sse')", Actual: SSEDelete("/sse")},
		{
			Expected: "@get('/sse', { contentType: 'json' })",
			Actual:   SSEGet("/sse", WithContentType(ContentTypeJSON)),
		},
		{
			Expected: "@post('/sse', { " +
				"contentType: 'form', " +
				"filterSignals: { include: '^foo' } " +
				"})",
			Actual: SSEPost("/sse",
				WithContentType(ContentTypeForm),
				WithFilterSignals("{ include: '^foo' }"),
			),
		},
		{
			Expected: "@put('/sse', { " +
				"contentType: 'form', " +
				"filterSignals: { exclude: '^foo' }, " +
				"selector: '#form-id', " +
				"headers: 'Authorization: Bearer token', " +
				"openWhenHidden: true, " +
				"payload: { bar: 'bar' }, " +
				"requestCancellation: 'cleanup', " +
				"retry: 'always', " +
				"retryInterval: 500, " +
				"retryMaxCount: 5, " +
				"retryMaxWait: 10000, " +
				"retryScaler: 4 " +
				"})",
			Actual: SSEPut("/sse",
				WithContentType(ContentTypeForm),
				WithFilterSignals("{ exclude: '^foo' }"),
				WithFormSelector("#form-id"),
				WithHeaders("Authorization: Bearer token"),
				WithOpenWhenHidden(true),
				WithPayload("{ bar: 'bar' }"),
				WithRequestCancellation("cleanup"),
				WithRetry(RetryAlways),
				WithRetryInterval(500),
				WithRetryMaxCount(5),
				WithRetryMaxWait(10000),
				WithRetryScaler(4),
			),
		},
	}

	for _, tc := range tests {
		if tc.Expected != tc.Actual {
			t.Errorf("\n\nExpected: %s,\n\nActual: %s", tc.Expected, tc.Actual)
		}
	}
}
