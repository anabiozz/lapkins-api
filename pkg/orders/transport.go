package orders

import (
	"context"
	"encoding/json"
	"net/http"
)

// encodeError writes a Service error to the given http.ResponseWriter.
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	defaultErr := &Error{
		ErrorCode: http.StatusInternalServerError,
	}
	if err, ok := err.(*Error); ok {
		requestID, _ := getRequestID(ctx)
		defaultErr.ErrorCode = err.ErrorCode
		defaultErr.RequestID = requestID
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultErr.ErrorCode)
	json.NewEncoder(w).Encode(defaultErr)
}

func decodeAddOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddOrderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addOrderResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}
