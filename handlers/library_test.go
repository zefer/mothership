package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/mothership/handlers"
)

type mockLibClient struct{}

var updatedURI string = ""

func (c mockLibClient) Update(uri string) (int, error) {
	updatedURI = uri
	return 1, nil
}

var _ = Describe("LibraryUpdateHandler", func() {
	var handler http.Handler
	var w *httptest.ResponseRecorder
	var client *mockLibClient
	var validParams map[string]interface{}

	BeforeEach(func() {
		w = httptest.NewRecorder()
		client = &mockLibClient{}
		handler = handlers.LibraryUpdateHandler(client)
		validParams = map[string]interface{}{"uri": "/bananas"}
	})

	Context("with valid params", func() {
		It("responds 202 accepted", func() {
			json, _ := json.Marshal(validParams)
			req, _ := http.NewRequest("POST", "/library/updated", bytes.NewBuffer(json))
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusAccepted))
			Expect(w.Body.String()).To(Equal(""))
		})

		It("tells MPD to update the library with the given URI", func() {
			updatedURI = ""
			json, _ := json.Marshal(validParams)
			req, _ := http.NewRequest("POST", "/library/updated", bytes.NewBuffer(json))
			handler.ServeHTTP(w, req)
			Expect(updatedURI).To(Equal("/bananas"))
		})
	})

	Context("with disallowed HTTP methods", func() {
		It("responds with 405 method not allowed", func() {
			for _, method := range []string{"GET", "PATCH", "DELETE"} {
				req, _ := http.NewRequest(method, "/library/updated", nil)
				handler.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(w.Body.String()).To(Equal(""))
			}
		})
	})

	Context("with un-parseable JSON", func() {
		It("responds 400 bad request", func() {
			var json = []byte(`{not-json`)
			req, _ := http.NewRequest("POST", "/library/updated", bytes.NewBuffer(json))
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when missing the required uri field", func() {
		It("responds 400 bad request", func() {
			delete(validParams, "uri")
			json, _ := json.Marshal(validParams)
			req, _ := http.NewRequest("POST", "/library/updated", bytes.NewBuffer(json))
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	// Without URI, we don't know what to update.
	Context("with an empty 'uri' field", func() {
		It("responds 400 bad request", func() {
			validParams["uri"] = ""
			json, _ := json.Marshal(validParams)
			req, _ := http.NewRequest("POST", "/library/updated", bytes.NewBuffer(json))
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
