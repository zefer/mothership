package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/mothership/handlers"
)

var randomed bool

func (c mockClient) Random(random bool) error {
	called = true
	randomed = random
	return nil
}

func (c mockFailingClient) Random(random bool) error {
	called = true
	randomed = random
	return errors.New("Random() failed")
}

var _ = Describe("RandomOnHandler", func() {
	var handler http.Handler
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		called = false
		w = httptest.NewRecorder()
	})

	Describe("RandomOn", func() {
		Context("When the MPD command succeeds", func() {
			var client *mockClient

			BeforeEach(func() {
				client = &mockClient{}
				handler = handlers.RandomOnHandler(client)
				req, _ := http.NewRequest("GET", "/randomOn", nil)
				handler.ServeHTTP(w, req)
			})

			It("calls Random(true) and responds OK", func() {
				Expect(called).To(BeTrue())
				Expect(randomed).To(BeTrue())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(Equal(""))
			})
		})

		Context("When the MPD command fails", func() {
			var client *mockFailingClient

			BeforeEach(func() {
				client = &mockFailingClient{}
				handler = handlers.RandomOnHandler(client)
				req, _ := http.NewRequest("GET", "/randomOn", nil)
				handler.ServeHTTP(w, req)
			})

			It("calls Random(true) and responds 500", func() {
				Expect(called).To(BeTrue())
				Expect(randomed).To(BeTrue())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(""))
			})
		})
	})

	Describe("RandomOff", func() {
		Context("When the MPD command succeeds", func() {
			var client *mockClient

			BeforeEach(func() {
				client = &mockClient{}
				handler = handlers.RandomOffHandler(client)
				req, _ := http.NewRequest("GET", "/randomOff", nil)
				handler.ServeHTTP(w, req)
			})

			It("calls Random(false) and responds OK", func() {
				Expect(called).To(BeTrue())
				Expect(randomed).To(BeFalse())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(Equal(""))
			})
		})

		Context("When the MPD command fails", func() {
			var client *mockFailingClient

			BeforeEach(func() {
				client = &mockFailingClient{}
				handler = handlers.RandomOffHandler(client)
				req, _ := http.NewRequest("GET", "/randomOff", nil)
				handler.ServeHTTP(w, req)
			})

			It("calls Random(true) and responds 500", func() {
				Expect(called).To(BeTrue())
				Expect(randomed).To(BeFalse())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(""))
			})
		})
	})
})
