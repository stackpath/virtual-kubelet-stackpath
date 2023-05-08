package provider

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	workloads "github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/workloads"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestBuildAnInvalidClientIDError(t *testing.T) {
	err := NewStackPathError(&url.Error{
		Err: &oauth2.RetrieveError{
			Response: &http.Response{
				StatusCode: 404,
			},
		},
	})

	switch err.(type) {
	case *InvalidClientIDError:
		if err.Error() != "invalid or unknown StackPath client ID" {
			t.Errorf("An invalid client ID error has the unexpected error message \"%s\"", err.Error())
		}
	default:
		t.Errorf("NewStackPathError() built an incorrect error type. Expected InvalidClientIDError but got %t", err)
	}
}

func TestBuildAnInvalidClientSecretError(t *testing.T) {
	err := NewStackPathError(&url.Error{
		Err: &oauth2.RetrieveError{
			Response: &http.Response{
				StatusCode: 401,
			},
		},
	})

	switch err.(type) {
	case *InvalidClientSecretError:
		if err.Error() != "invalid StackPath client secret" {
			t.Errorf("An invalid client ID error has the unexpected error message \"%s\"", err.Error())
		}
	default:
		t.Errorf("NewStackPathError() built an incorrect error type. Expected InvalidClientSecretError but got %t", err)
	}
}

func TestBuildANonStackPathError(t *testing.T) {
	type TestError struct {
		error
	}

	err := NewStackPathError(&TestError{errors.New("foo")})

	switch err.(type) {
	case *TestError:
	default:
		t.Errorf("NewStackPathError built an incorrect error type. Expected TestError but got %t", err)
	}
}

func TestHTTPStatusFromCode(t *testing.T) {
	testCases := []struct {
		code       int32
		httpStatus int
	}{
		{
			code:       0,
			httpStatus: http.StatusOK,
		},
		{
			code:       1,
			httpStatus: http.StatusRequestTimeout,
		},
		{
			code:       2,
			httpStatus: http.StatusInternalServerError,
		},
		{
			code:       3,
			httpStatus: http.StatusBadRequest,
		},
		{
			code:       4,
			httpStatus: http.StatusGatewayTimeout,
		},
		{
			code:       5,
			httpStatus: http.StatusNotFound,
		},
		{
			code:       6,
			httpStatus: http.StatusConflict,
		},
		{
			code:       7,
			httpStatus: http.StatusForbidden,
		},
		{
			code:       8,
			httpStatus: http.StatusTooManyRequests,
		},
		{
			code:       9,
			httpStatus: http.StatusBadRequest,
		},
		{
			code:       10,
			httpStatus: http.StatusConflict,
		},
		{
			code:       11,
			httpStatus: http.StatusBadRequest,
		},
		{
			code:       12,
			httpStatus: http.StatusNotImplemented,
		},
		{
			code:       13,
			httpStatus: http.StatusInternalServerError,
		},
		{
			code:       14,
			httpStatus: http.StatusServiceUnavailable,
		},
		{
			code:       15,
			httpStatus: http.StatusInternalServerError,
		},
		{
			code:       16,
			httpStatus: http.StatusUnauthorized,
		},
	}

	for _, c := range testCases {
		httpStatus := HTTPStatusFromCode(c.code)
		assert.Equal(t, c.httpStatus, httpStatus)
	}
}

func TestStackPathError(t *testing.T) {
	response := workloads.NewGetWorkloadsDefault(3)
	response.Payload = &workload_models.StackpathapiStatus{
		Code:    int32(HTTPStatusFromCode(3)),
		Message: "Payload message",
	}
	response.Payload.SetDetails([]workload_models.APIStatusDetail{
		&workload_models.StackpathRPCBadRequest{FieldViolations: []*workload_models.StackpathRPCBadRequestFieldViolation{
			{
				Description: "Description",
				Field:       "Field",
			}}},
		&workload_models.StackpathRPCRequestInfo{
			RequestID: "123",
		},
		&workload_models.StackpathRPCHelp{},
	})

	err := NewStackPathError(response)

	assert.IsType(t, &APIError{}, err)
	assert.Equal(t, err.(*APIError).NotFound(), false)
	assert.Equal(t, err.(*APIError).Cause().Error(), "a 500 error was returned from StackPath: \"Payload message\". The following fields have errors: Field: Description (request ID 123)")
}
