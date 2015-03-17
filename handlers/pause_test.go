package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/mothership/handlers"
)

var paused bool

func (c mockClient) Pause(pause bool) error {
	called = true
	paused = pause
	return nil
}

func (c mockFailingClient) Pause(pause bool) error {
	called = true
	paused = pause
	return errors.New("Pause() failed")
}

var _ = Describe("PauseHandler", func() {
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
			handler = handlers.PauseHandler(client)
			req, _ := http.NewRequest("GET", "/pause", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Pause(true) and responds OK", func() {
			Expect(called).To(BeTrue())
			Expect(paused).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
		})
	})

	Context("When the MPD command fails", func() {
		var client *mockFailingClient

		BeforeEach(func() {
			client = &mockFailingClient{}
			handler = handlers.PauseHandler(client)
			req, _ := http.NewRequest("GET", "/pause", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Pause(true) and responds 500", func() {
			Expect(called).To(BeTrue())
			Expect(paused).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal(""))
		})
	})
})
