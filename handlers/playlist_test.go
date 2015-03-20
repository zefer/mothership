package handlers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/gompd/mpd"
	"github.com/zefer/mothership/handlers"
)

type mockPlClient struct{}

func (c mockPlClient) Status() (mpd.Attrs, error) {
	return map[string]string{}, nil
}

func (c mockPlClient) PlaylistInfo(start, end int) ([]mpd.Attrs, error) {
	return make([]mpd.Attrs, 0), nil
}

func (c mockPlClient) Clear() error {
	return nil
}

func (c mockPlClient) PlaylistLoad(name string, start, end int) error {
	return nil
}

func (c mockPlClient) Add(uri string) error {
	return nil
}

func (c mockPlClient) Play(pos int) error {
	return nil
}

var _ = Describe("PlayListHandler", func() {
	var handler http.Handler
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		called = false
		w = httptest.NewRecorder()
	})

	Context("with disallowed HTTP methods", func() {
		var client *mockPlClient

		BeforeEach(func() {
			client = &mockPlClient{}
			handler = handlers.PlayListHandler(client)
		})

		It("responds with 405 method not allowed", func() {
			for _, method := range []string{"PUT", "PATCH", "DELETE"} {
				req, _ := http.NewRequest(method, "/playlist", nil)
				handler.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(w.Body.String()).To(Equal(""))
			}
		})
	})
})
