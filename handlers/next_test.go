package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/mothership/handlers"
)

func (c mockClient) Next() error {
	called = true
	return nil
}

func (c mockFailingClient) Next() error {
	called = true
	return errors.New("Next() failed")
}

var _ = Describe("NextHandler", func() {
	var handler http.Handler
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		called = false
		w = httptest.NewRecorder()
	})

	Context("When the MPD command succeeds", func() {
		var client *mockClient

		BeforeEach(func() {
			client = &mockClient{}
			handler = handlers.NextHandler(client)
			req, _ := http.NewRequest("GET", "/next", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Next() and responds OK", func() {
			Expect(called).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusNoContent))
			Expect(w.Body.String()).To(Equal(""))
		})
	})

	Context("When the MPD command fails", func() {
		var client *mockFailingClient

		BeforeEach(func() {
			client = &mockFailingClient{}
			handler = handlers.NextHandler(client)
			req, _ := http.NewRequest("GET", "/next", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Next() and responds 500", func() {
			Expect(called).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal(""))
		})
	})
})
