package icloud

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type ICloudHandler struct {
	cache  CacheInterface
	config Config
	getWebStream   func(base_url string) (*StreamResponse, error)
	getWebAssetUrls func(base_url string, photoGuids []string) (*AssetUrlsResponse, error)
}

type Config struct {
	Token string
}

func NewICloudHandler(config Config) *ICloudHandler {
	handler := &ICloudHandler{
		cache:  NewCache(),
		config: config,
	}
	handler.getWebStream = handler.defaultGetWebStream
	handler.getWebAssetUrls = handler.defaultGetWebAssetUrls

	return handler
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

func (h *ICloudHandler) defaultGetWebStream(base_url string) (*StreamResponse, error) {
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

func (h *ICloudHandler) defaultGetWebAssetUrls(base_url string, photoGuids []string) (*AssetUrlsResponse, error) {
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

	partition := getPartitionFromToken(h.config.Token)
	base_url := "https://p" + partition + "-sharedstreams.icloud.com/" + h.config.Token + "/sharedstreams"

	cacheKeyStream := "stream_" + h.config.Token
	cacheKeyAssets := "assets_" + h.config.Token

	var stream *StreamResponse
	var err error

	if cachedStream, found := h.cache.Get(cacheKeyStream); found {
		log.Print("Cache hit for stream")
		stream = cachedStream.(*StreamResponse)
	} else {
		stream, err = h.getWebStream(base_url)
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
		assetsUrl, err = h.getWebAssetUrls(base_url, photoGuids)
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
