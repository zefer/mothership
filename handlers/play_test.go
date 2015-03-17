package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/mothership/handlers"
)

var playPos int

func (c mockClient) Play(pos int) error {
	called = true
	playPos = pos
	return nil
}

func (c mockFailingClient) Play(pos int) error {
	called = true
	playPos = pos
	return errors.New("Play() failed")
}

var _ = Describe("PlayHandler", func() {
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
			handler = handlers.PlayHandler(client)
			req, _ := http.NewRequest("GET", "/play", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Play(-1) and responds OK", func() {
			Expect(called).To(BeTrue())
			Expect(playPos).To(Equal(-1))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
		})
	})

	Context("When the MPD command fails", func() {
		var client *mockFailingClient

		BeforeEach(func() {
			client = &mockFailingClient{}
			handler = handlers.PlayHandler(client)
			req, _ := http.NewRequest("GET", "/play", nil)
			handler.ServeHTTP(w, req)
		})

		It("calls Play(-1) and responds 500", func() {
			Expect(called).To(BeTrue())
			Expect(playPos).To(Equal(-1))
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal(""))
		})
	})
})
