package requests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shipments/servers"
)

// NewRequest creates a new http request
func NewRequest(server *servers.Server, method, path string, payload []byte) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	var (
		req *http.Request
		err error
	)
	if len(payload) > 0 {
		fmt.Printf("\n\nPayload: %s\n\n", payload)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		return nil, err
	}

	server.Router.ServeHTTP(w, req)

	return w, nil
}
