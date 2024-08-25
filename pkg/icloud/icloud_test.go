package icloud

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(key string, value interface{}, duration time.Duration) {
	m.Called(key, value, duration)
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func mockGetWebStream(base_url string) (*StreamResponse, error) {
	return &StreamResponse{
		Photos: []Photo{
			{
				PhotoGuid:      "guid1",
				DateCreated:    "2023-01-01",
				MediaAssetType: "image",
				Derivatives: map[string]Derivative{
					"thumb": {Checksum: "checksum1"},
				},
			},
		},
	}, nil
}

func mockGetWebAssetUrls(base_url string, photoGuids []string) (*AssetUrlsResponse, error) {
	return &AssetUrlsResponse{
		Items: map[string]AssetItem{
			"checksum1": {UrlLocation: "location1", UrlPath: "/path1"},
		},
		Locations: map[string]Location{
			"location1": {Scheme: "https", Hosts: []string{"host1.com"}},
		},
	}, nil
}

func TestICloudHandler_Photos(t *testing.T) {
	var cache *MockCache
	var handler *ICloudHandler

	setup := func() {
		cache = new(MockCache)
		handler = NewICloudHandler(Config{Token: "testtoken"})
		handler.cache = cache
	}

	t.Run("Valid request with cache hit", func(t *testing.T) {
		setup()
		req := httptest.NewRequest(http.MethodGet, "/photos", nil)
		w := httptest.NewRecorder()

		streamResponse := &StreamResponse{
			Photos: []Photo{
				{
					PhotoGuid:      "guid1",
					DateCreated:    "2023-01-01",
					MediaAssetType: "image",
					Derivatives: map[string]Derivative{
						"thumb": {Checksum: "checksum1"},
					},
				},
			},
		}

		assetsUrlResponse := &AssetUrlsResponse{
			Items: map[string]AssetItem{
				"checksum1": {UrlLocation: "location1", UrlPath: "/path1"},
			},
			Locations: map[string]Location{
				"location1": {Scheme: "https", Hosts: []string{"host1.com"}},
			},
		}

		cache.On("Get", "stream_testtoken").Return(streamResponse, true)
		cache.On("Get", "assets_testtoken").Return(assetsUrlResponse, true)

		handler.Photos(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		expected := `[{"dateCreated":"2023-01-01","mediaAssetType":"image","derivatives":{"thumb":{"mediaUrl":"https://host1.com/path1","dateCreated":"2023-01-01"}}}]`
		assert.JSONEq(t, expected, string(body))

		cache.AssertExpectations(t)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		setup()
		req := httptest.NewRequest(http.MethodPost, "/photos", nil)
		w := httptest.NewRecorder()

		handler.Photos(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		assert.Equal(t, "Method is not supported.", strings.TrimSpace(string(body)))
	})

	t.Run("Failed to marshal response", func(t *testing.T) {
		setup()
		req := httptest.NewRequest(http.MethodGet, "/photos", nil)
		w := httptest.NewRecorder()

		streamResponse := &StreamResponse{
			Photos: []Photo{
				{
					PhotoGuid:      "guid1",
					DateCreated:    "2023-01-01",
					MediaAssetType: "image",
					Derivatives: map[string]Derivative{
						"thumb": {Checksum: "checksum1"},
					},
				},
			},
		}

		assetsUrlResponse := &AssetUrlsResponse{
			Items: map[string]AssetItem{
				"checksum1": {UrlLocation: "location1", UrlPath: "/path1"},
			},
			Locations: map[string]Location{
				"location1": {Scheme: "https", Hosts: []string{"host1.com"}},
			},
		}

		cache.On("Get", "stream_testtoken").Return(streamResponse, true)
		cache.On("Get", "assets_testtoken").Return(assetsUrlResponse, true)

		handler.Photos(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "Failed to marshal response.", strings.TrimSpace(string(body)))
	})

	t.Run("Failed to get web asset URLs", func(t *testing.T) {
		setup()
		req := httptest.NewRequest(http.MethodGet, "/photos", nil)
		w := httptest.NewRecorder()

		streamResponse := &StreamResponse{
			Photos: []Photo{
				{
					PhotoGuid:      "guid1",
					DateCreated:    "2023-01-01",
					MediaAssetType: "image",
					Derivatives: map[string]Derivative{
						"thumb": {Checksum: "checksum1"},
					},
				},
			},
		}

		cache.On("Get", "stream_testtoken").Return(streamResponse, true)
		cache.On("Get", "assets_testtoken").Return(nil, false)

		handler.getWebAssetUrls = func(base_url string, photoGuids []string) (*AssetUrlsResponse, error) {
			return nil, errors.New("failed to get web asset URLs")
		}

		handler.Photos(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "Failed to get web asset URLs.", strings.TrimSpace(string(body)))
	})
}
