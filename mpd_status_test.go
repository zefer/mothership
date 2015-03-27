package main

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zefer/gompd/mpd"
)

func TestMpdStatus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MpdStatus")
}

type mockStatusClient struct{}

var mockStatus mpd.Attrs = make(mpd.Attrs)

func (c mockStatusClient) Status() (mpd.Attrs, error) {
	return mockStatus, nil
}

var mockCurrentSong mpd.Attrs = make(mpd.Attrs)

func (c mockStatusClient) CurrentSong() (mpd.Attrs, error) {
	return mockCurrentSong, nil
}

var _ = Describe("mpdStatusJSON", func() {
	var client *mockStatusClient

	BeforeEach(func() {
		client = &mockStatusClient{}
		mockStatus["volume"] = "11"
		mockStatus["gorillas"] = "bananas"
		mockCurrentSong["spaceman"] = "nope"
	})

	It("returns a JSON map combining MPD's `status` & `currentsong` data", func() {
		b, err := mpdStatusJSON(client)
		if err != nil {
			Fail(fmt.Sprintf("mpdStatusJSON failed %v", err))
		}
		var status map[string]string
		if err := json.Unmarshal(b, &status); err != nil {
			Fail(fmt.Sprintf("Could not parse JSON %v", err))
		}
		Expect(len(status)).To(Equal(3))
		Expect(status["volume"]).To(Equal("11"))
		Expect(status["spaceman"]).To(Equal("nope"))
		Expect(status["gorillas"]).To(Equal("bananas"))
	})

})
