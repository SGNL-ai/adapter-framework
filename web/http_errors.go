// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"net/http"
	"strconv"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

// HTTPError returns a detailed error if the given HTTP response status code
// indicates that the HTTP request failed, and nil otherwise.
func HTTPError(statusCode int, retryAfterHeader string) (adapterErr *api_adapter_v1.Error) {
	if statusCode >= 200 && statusCode < 300 { // Success.
		return nil
	}

	adapterErr = new(api_adapter_v1.Error)

	if retryAfterHeader != "" {
		// Cf. https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.3,
		// Retry-After can be either an HTTP-date or a number of seconds.
		seconds, err := strconv.ParseInt(retryAfterHeader, 10, 64)
		if err == nil {
			adapterErr.RetryAfter = &durationpb.Duration{
				Seconds: seconds,
				Nanos:   0,
			}
		} else {
			afterTime, err := time.Parse(time.RFC1123, retryAfterHeader)
			if err == nil {
				duration := afterTime.UTC().Sub(time.Now().UTC())
				adapterErr.RetryAfter = &durationpb.Duration{
					Seconds: int64(duration / time.Second),
					Nanos:   0,
				}
			}
		}
	}

	if statusCode >= 300 && statusCode < 400 {
		// In the case of 3xx, this is an internal error since the adapter
		// should be responsible for handling redirects if the datasource
		// supports them.
		adapterErr.Message = api_adapter_v1.ErrorMsgUnhandledStatusCode
		adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL
	} else if statusCode >= 400 && statusCode < 500 {
		switch statusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			adapterErr.Message = api_adapter_v1.ErrorMsgFailedToAuthenticate
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED
		case http.StatusTooManyRequests:
			adapterErr.Message = api_adapter_v1.ErrorMsgDatasourceTooManyRequests
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_TOO_MANY_REQUESTS
		default:
			// In the case of other 4xx responses, indicate the adapter
			// constructed an invalid request.
			adapterErr.Message = api_adapter_v1.ErrorMsgDatasourceRejectedRequest
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL
		}
	} else if statusCode >= 500 && statusCode < 600 {
		switch statusCode {
		case http.StatusInternalServerError:
			adapterErr.Message = api_adapter_v1.ErrorMsgDatasourceInternalError
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			adapterErr.Message = api_adapter_v1.ErrorMsgDatasourceTemporarilyUnavailable
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_TEMPORARILY_UNAVAILABLE
		default:
			// In the case of other 5xx responses, indicate that the datasource
			// is permanently unavailable.
			adapterErr.Message = api_adapter_v1.ErrorMsgDatasourcePermanentlyUnavailable
			adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_PERMANENTLY_UNAVAILABLE
		}
	} else {
		adapterErr.Message = api_adapter_v1.ErrorMsgUnexpectedErrorCode
		adapterErr.Code = api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL
	}

	return
}
