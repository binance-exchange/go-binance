package binance

import (
	"testing"
	"context"
)

func TestErrorHandler(t *testing.T) {
	as := NewAPIService("", "", nil, nil, context.Background()).(*apiService)
	err := as.handleError([]byte(`{"code":-1105,"msg":"Parameter 'side' was was empty."}`))
	tErr, ok := err.(*Error)
	if !ok {
		t.Errorf("invalid type of error returned: %T", tErr)
	}
	if tErr.Code != -1105 {
		t.Errorf("invalid error code extracted")
	}
	if tErr.Message != "Parameter 'side' was was empty." {
		t.Errorf("invalid error message extracted")
	}
}
