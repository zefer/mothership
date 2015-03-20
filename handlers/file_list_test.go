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

type mockFLClient struct{}

var listedURI string = ""

func (c mockFLClient) ListInfo(uri string) ([]mpd.Attrs, error) {
	listedURI = uri
	files := []mpd.Attrs{
		{
			"file":          "random/Gorilla.mp3",
			"last-modified": "2014-06-29T13:11:49Z",
		},
		{
			"directory":     "random/Autechre - Incunabula",
			"last-modified": "2014-04-29T13:11:49Z",
		},
		{
			"playlist":      "random/ZeferRadio.m3u",
			"last-modified": "2014-05-29T13:11:49Z",
		},
	}
	return files, nil
}

var _ = Describe("FileListHandler", func() {
	var handler http.Handler
	var w *httptest.ResponseRecorder
	var client *mockFLClient

	BeforeEach(func() {
		w = httptest.NewRecorder()
		client = &mockFLClient{}
		handler = handlers.FileListHandler(client)
	})

	Context("with disallowed HTTP methods", func() {
		It("responds with 405 method not allowed", func() {
			for _, method := range []string{"POST", "PUT", "PATCH", "DELETE"} {
				req, _ := http.NewRequest(method, "/files", nil)
				handler.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(w.Body.String()).To(Equal(""))
			}
		})
	})

	Context("when missing the required uri param", func() {
		It("responds 400 bad request", func() {
			req, _ := http.NewRequest("GET", "/files", nil)
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("with an empty 'uri' field", func() {
		It("responds 400 bad request", func() {
			req, _ := http.NewRequest("GET", "/files?uri=", nil)
			handler.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	It("queries MPD for listings at the given URI", func() {
		listedURI = ""
		req, _ := http.NewRequest("GET", "/files?uri=space/spaceman", nil)
		handler.ServeHTTP(w, req)
		Expect(listedURI).To(Equal("space/spaceman"))
	})

	Describe("the response", func() {
		BeforeEach(func() {
			req, _ := http.NewRequest("GET", "/files?uri=random", nil)
			handler.ServeHTTP(w, req)
		})

		It("responds with a 200 OK", func() {
			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("responds with the JSON content-type", func() {
			Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
		})

		It("returns a JSON array of the items", func() {
			var files []map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&files); err != nil {
				Fail(fmt.Sprintf("Could not parse JSON %v", err))
			}
			Expect(len(files)).To(Equal(3))
			for _, file := range files {
				Expect(len(file)).To(Equal(4))
				for _, field := range []string{"path", "type", "base", "lastModified"} {
					Expect(file[field]).NotTo(BeNil())
				}
			}
		})
	})

	Describe("sorting", func() {
		var files []map[string]interface{}

		It("defaults to sort by name asc (ascending)", func() {
			req, _ := http.NewRequest(
				"GET", "/files?uri=random", nil,
			)
			handler.ServeHTTP(w, req)
			json.NewDecoder(w.Body).Decode(&files)
			Expect(files[0]["base"]).To(Equal("Autechre - Incunabula"))
			Expect(files[1]["base"]).To(Equal("Gorilla.mp3"))
			Expect(files[2]["base"]).To(Equal("ZeferRadio.m3u"))
		})

		It("sorts by name asc (ascending)", func() {
			req, _ := http.NewRequest(
				"GET", "/files?uri=random&sort=name&direction=asc", nil,
			)
			handler.ServeHTTP(w, req)
			json.NewDecoder(w.Body).Decode(&files)
			Expect(files[0]["base"]).To(Equal("Autechre - Incunabula"))
			Expect(files[1]["base"]).To(Equal("Gorilla.mp3"))
			Expect(files[2]["base"]).To(Equal("ZeferRadio.m3u"))
		})

		It("sorts by name desc (descending)", func() {
			req, _ := http.NewRequest(
				"GET", "/files?uri=random&sort=name&direction=desc", nil,
			)
			handler.ServeHTTP(w, req)
			json.NewDecoder(w.Body).Decode(&files)
			Expect(files[0]["base"]).To(Equal("ZeferRadio.m3u"))
			Expect(files[1]["base"]).To(Equal("Gorilla.mp3"))
			Expect(files[2]["base"]).To(Equal("Autechre - Incunabula"))
		})

		It("sorts by date asc (ascending)", func() {
			req, _ := http.NewRequest(
				"GET", "/files?uri=random&sort=date&direction=asc", nil,
			)
			handler.ServeHTTP(w, req)
			json.NewDecoder(w.Body).Decode(&files)
			Expect(files[0]["base"]).To(Equal("Autechre - Incunabula"))
			Expect(files[1]["base"]).To(Equal("ZeferRadio.m3u"))
			Expect(files[2]["base"]).To(Equal("Gorilla.mp3"))
		})

		It("sorts by date desc (descending)", func() {
			req, _ := http.NewRequest(
				"GET", "/files?uri=random&sort=date&direction=desc", nil,
			)
			handler.ServeHTTP(w, req)
			json.NewDecoder(w.Body).Decode(&files)
			Expect(files[0]["base"]).To(Equal("Gorilla.mp3"))
			Expect(files[1]["base"]).To(Equal("ZeferRadio.m3u"))
			Expect(files[2]["base"]).To(Equal("Autechre - Incunabula"))
		})
	})
})
