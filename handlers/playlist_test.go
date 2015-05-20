package handlers_test

import (
	"bytes"
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
		{
			"file":          "http://somestream",
			"Name":          "HTTP stream from pls",
			"Last-Modified": "2011-10-09T11:45:11Z",
			"Pos":           "2",
		},
	}
	return pls, nil
}

var clearCalled bool = false

func (c mockPlClient) Clear() error {
	clearCalled = true
	return nil
}

var loadedURI string = ""

func (c mockPlClient) PlaylistLoad(uri string, start, end int) error {
	loadedURI = uri
	return nil
}

var addedURI string = ""

func (c mockPlClient) Add(uri string) error {
	addedURI = uri
	return nil
}

var playCalled bool = false
var playedPos int = 0

func (c mockPlClient) Play(pos int) error {
	playCalled = true
	playedPos = pos
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
				Context("when the current playlist position isn't the first song", func() {
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

				Context("when the current playlist position is the first song", func() {
					BeforeEach(func() {
						mockStatus = map[string]string{"playlistlength": "501", "song": "0"}
						req, _ := http.NewRequest("GET", "/playlist", nil)
						handler.ServeHTTP(w, req)
					})
					// Checking we don't query with a negative start index.
					It("requests a slice of the playlist from MPD. 0 to 500", func() {
						Expect(requestedRange[0]).To(Equal(0))
						Expect(requestedRange[1]).To(Equal(500))
					})
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
				Expect(len(pls)).To(Equal(3))
				// Item 1 has artist & track parts, so we expect "artist - track".
				Expect(len(pls[0])).To(Equal(2))
				Expect(pls[0]["pos"]).To(BeEquivalentTo(1))
				Expect(pls[0]["name"]).To(Equal("Led Zeppelin - The Ocean"))
				// Item 2 doesn't have artist & track parts, so we expect "file.mp3".
				Expect(len(pls[1])).To(Equal(2))
				Expect(pls[1]["pos"]).To(BeEquivalentTo(2))
				Expect(pls[1]["name"]).To(Equal("Johnny Cash – Sea Of Heartbreak.mp3"))
				// Item 3 has a 'name' field, such as from a loaded pls playlist.
				Expect(len(pls[2])).To(Equal(2))
				Expect(pls[2]["pos"]).To(BeEquivalentTo(3))
				Expect(pls[2]["name"]).To(Equal("HTTP stream from pls"))
			})
		})
	})

	Context("with a POST request (update the current playlist)", func() {
		var validParams map[string]interface{}

		BeforeEach(func() {
			validParams = map[string]interface{}{
				"uri": "gorilla.mp3", "type": "file", "replace": true, "play": true,
			}
		})

		Describe("POST data validation", func() {
			Context("with valid params", func() {
				It("responds 204 no content", func() {
					json, _ := json.Marshal(validParams)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})
			})

			Context("with un-parseable JSON", func() {
				It("responds 400 bad request", func() {
					var json = []byte(`{not-json`)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("with missing required fields", func() {
				It("responds 400 bad request", func() {
					for _, f := range []string{"uri", "type", "replace", "play"} {
						// d = map[string]string{"uri": "", "type": "", "replace": "", "play": ""}
						params := make(map[string]interface{})
						for k, v := range validParams {
							params[k] = v
						}
						delete(params, f)
						json, _ := json.Marshal(params)
						req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
						handler.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					}
				})
			})

			// Without URI, we don't know what to add to the playlist.
			Context("with an empty 'uri' field", func() {
				It("responds 400 bad request", func() {
					validParams["uri"] = ""
					json, _ := json.Marshal(validParams)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})

		Context("with replace=true", func() {
			It("clears the playlist", func() {
				clearCalled = false
				validParams["replace"] = true
				json, _ := json.Marshal(validParams)
				req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
				handler.ServeHTTP(w, req)
				Expect(clearCalled).To(BeTrue())
			})
		})

		Context("with replace=false", func() {
			It("does not clear the playlist", func() {
				clearCalled = false
				validParams["replace"] = false
				json, _ := json.Marshal(validParams)
				req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
				handler.ServeHTTP(w, req)
				Expect(clearCalled).To(BeFalse())
			})
		})

		Context("when type='playlist'", func() {
			It("loads the given URI", func() {
				loadedURI = ""
				addedURI = ""
				validParams["type"] = "playlist"
				validParams["uri"] = "http://gorillas"
				json, _ := json.Marshal(validParams)
				req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
				handler.ServeHTTP(w, req)
				Expect(loadedURI).To(Equal("http://gorillas"))
				Expect(addedURI).To(Equal(""))
			})
		})

		Context("when type='directory' or type='file'", func() {
			It("adds the given URI", func() {
				for _, t := range []string{"directory", "file"} {
					loadedURI = ""
					addedURI = ""
					validParams["type"] = t
					validParams["uri"] = "http://gorillas"
					json, _ := json.Marshal(validParams)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(addedURI).To(Equal("http://gorillas"))
					Expect(loadedURI).To(Equal(""))
				}
			})
		})

		Context("when play=true", func() {
			BeforeEach(func() {
				validParams["play"] = true
				mockStatus = map[string]string{"playlistlength": "66"}
				playedPos = 123
				playCalled = false
			})

			Context("and replace=true", func() {
				It("it tells MPD to play from position 0", func() {
					validParams["replace"] = true
					json, _ := json.Marshal(validParams)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(playCalled).To(BeTrue())
					Expect(playedPos).To(Equal(0))
				})
			})

			Context("and replace=false", func() {
				It("it tells MPD to play from the start of the new added items", func() {
					validParams["replace"] = false
					json, _ := json.Marshal(validParams)
					req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
					handler.ServeHTTP(w, req)
					Expect(playCalled).To(BeTrue())
					Expect(playedPos).To(Equal(66))
				})
			})
		})

		Context("when play=false", func() {
			It("it does not tell MPD to play", func() {
				playCalled = false
				validParams["play"] = false
				json, _ := json.Marshal(validParams)
				req, _ := http.NewRequest("POST", "/playlist", bytes.NewBuffer(json))
				handler.ServeHTTP(w, req)
				Expect(playCalled).To(BeFalse())
			})
		})
	})
})
