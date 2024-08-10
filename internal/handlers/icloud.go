package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type ICloudHandler struct {
	cache *Cache
}

func NewICloudHandler() *ICloudHandler {
	return &ICloudHandler{
		cache: NewCache(),
	}
}

type StreamResponse struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	PhotoGuid      string                `json:"photoGuid"`
	DateCreated    string                `json:"dateCreated"`
	MediaAssetType string                `json:"mediaAssetType"`
	Derivatives    map[string]Derivative `json:"derivatives"`
}

type Derivative struct {
	Checksum string `json:"checksum"`
}

type AssetUrlsResponse struct {
	Items     map[string]AssetItem `json:"items"`
	Locations map[string]Location  `json:"locations"`
}

type AssetItem struct {
	UrlLocation string `json:"url_location"`
	UrlPath     string `json:"url_path"`
}

type Location struct {
	Scheme string   `json:"scheme"`
	Hosts  []string `json:"hosts"`
}

func getWebStream(base_url string) (*StreamResponse, error) {
	url := base_url + "/webstream"
	body := map[string]interface{}{
		"streamCtag": nil,
	}
	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var streamResponse StreamResponse
	err = json.Unmarshal(bodyBytes, &streamResponse)
	if err != nil {
		return nil, err
	}

	return &streamResponse, nil
}

func getWebAssetUrls(base_url string, photoGuids []string) (*AssetUrlsResponse, error) {
	url := base_url + "/webasseturls"
	body := map[string]interface{}{
		"photoGuids": photoGuids,
	}
	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var assetUrlsResponse AssetUrlsResponse
	err = json.Unmarshal(bodyBytes, &assetUrlsResponse)
	if err != nil {
		return nil, err
	}

	return &assetUrlsResponse, nil
}

func (h *ICloudHandler) Photos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	token := "B1A55Z2WMR2vuY"
	partition := "72"
	base_url := "https://p" + partition + "-sharedstreams.icloud.com/" + token + "/sharedstreams"

	cacheKeyStream := "stream_" + token
	cacheKeyAssets := "assets_" + token

	var stream *StreamResponse
	var err error

	if cachedStream, found := h.cache.Get(cacheKeyStream); found {
		log.Print("Cache hit for stream")
		stream = cachedStream.(*StreamResponse)
	} else {
		stream, err = getWebStream(base_url)
		if err != nil {
			http.Error(w, "Failed to get web stream.", http.StatusInternalServerError)
			return
		}
		h.cache.Set(cacheKeyStream, stream, 10*time.Minute)
	}

	photoGuids := []string{}
	for _, photo := range stream.Photos {
		photoGuids = append(photoGuids, photo.PhotoGuid)
	}

	var assetsUrl *AssetUrlsResponse
	if cachedAssets, found := h.cache.Get(cacheKeyAssets); found {
		log.Print("Cache hit for assets")
		assetsUrl = cachedAssets.(*AssetUrlsResponse)
	} else {
		assetsUrl, err = getWebAssetUrls(base_url, photoGuids)
		if err != nil {
			http.Error(w, "Failed to get web asset URLs.", http.StatusInternalServerError)
			return
		}
		h.cache.Set(cacheKeyAssets, assetsUrl, 10*time.Minute)
	}

	photos := []map[string]interface{}{}
	for _, photo := range stream.Photos {
		derivatives := map[string]interface{}{}
		for derivativeKey, derivative := range photo.Derivatives {
			asset := assetsUrl.Items[derivative.Checksum]
			location := assetsUrl.Locations[asset.UrlLocation]
			host := location.Scheme + "://" + location.Hosts[0]
			urlPath := asset.UrlPath
			derivatives[derivativeKey] = map[string]interface{}{
				"mediaUrl":    host + urlPath,
				"dateCreated": photo.DateCreated,
			}
		}
		photos = append(photos, map[string]interface{}{
			"dateCreated":    photo.DateCreated,
			"mediaAssetType": photo.MediaAssetType,
			"derivatives":    derivatives,
		})
	}

	response, err := json.Marshal(photos)
	if err != nil {
		http.Error(w, "Failed to marshal response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration int64
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(duration).UnixNano(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}
