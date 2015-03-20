package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/gompd/mpd"
	"github.com/zefer/mothership/handlers"
)

type mockPlClient struct{}

var mockStatus map[string]string = map[string]string{}

func (c mockPlClient) Status() (mpd.Attrs, error) {
	return mockStatus, nil
}

var requestedRange [2]int

func (c mockPlClient) PlaylistInfo(start, end int) ([]mpd.Attrs, error) {
	requestedRange = [2]int{start, end}
	pls := []mpd.Attrs{
		{
			"file":          "Led Zeppelin - Houses Of The Holy/08 - Led Zeppelin - The Ocean.mp3",
			"Artist":        "Led Zeppelin",
			"Title":         "The Ocean",
			"Album":         "Houses of the Holy",
			"Last-Modified": "2010-12-09T21:32:02Z",
			"Pos":           "0",
		},
		{
			"file":          "Johnny Cash – Unchained/Johnny Cash – Sea Of Heartbreak.mp3",
			"Last-Modified": "2011-10-09T11:45:11Z",
			"Pos":           "1",
		},
	}
	return pls, nil
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

	Context("with a GET request (list the current playlist)", func() {
		var client *mockPlClient

		BeforeEach(func() {
			client = &mockPlClient{}
			handler = handlers.PlayListHandler(client)
		})

		Describe("the MPD query", func() {
			Context("when there are less than 500 items on the playlist", func() {
				BeforeEach(func() {
					mockStatus = map[string]string{"playlistlength": "12"}
					req, _ := http.NewRequest("GET", "/playlist", nil)
					handler.ServeHTTP(w, req)
				})
				It("requests the full playlist from MPD", func() {
					Expect(requestedRange[0]).To(Equal(-1))
					Expect(requestedRange[1]).To(Equal(-1))
				})
			})

			Context("when there are more than 500 items on the playlist", func() {
				BeforeEach(func() {
					mockStatus = map[string]string{"playlistlength": "501", "song": "123"}
					req, _ := http.NewRequest("GET", "/playlist", nil)
					handler.ServeHTTP(w, req)
				})
				It("requests a slice of the playlist from MPD. Current pos -1 to +500", func() {
					Expect(requestedRange[0]).To(Equal(122))
					Expect(requestedRange[1]).To(Equal(623))
				})
			})
		})

		Describe("the response", func() {
			It("responds with 200 OK", func() {
				req, _ := http.NewRequest("GET", "/playlist", nil)
				handler.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("responds with the JSON content-type", func() {
				req, _ := http.NewRequest("GET", "/playlist", nil)
				handler.ServeHTTP(w, req)
				Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
			})

			It("responds with a JSON array of playlist items", func() {
				req, _ := http.NewRequest("GET", "/playlist", nil)
				handler.ServeHTTP(w, req)
				var pls []map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&pls); err != nil {
					Fail(fmt.Sprintf("Could not parse JSON %v", err))
				}
				// Item 1 has artist & track parts, so we expect "artist - track".
				Expect(len(pls[0])).To(Equal(2))
				Expect(pls[0]["pos"]).To(BeEquivalentTo(1))
				Expect(pls[0]["name"]).To(Equal("Led Zeppelin - The Ocean"))
				// Item 2 doesn't have artist & track parts, so we expect "file.mp3".
				Expect(len(pls[1])).To(Equal(2))
				Expect(pls[1]["pos"]).To(BeEquivalentTo(2))
				Expect(pls[1]["name"]).To(Equal("Johnny Cash – Sea Of Heartbreak.mp3"))
			})
		})
	})
})
