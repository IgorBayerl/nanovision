// Package client talks to a team server. It implements history.Target, so the
// CLI saves runs and finds base runs the same way as with a local store.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/server"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

type Client struct {
	base  string
	token string
	http  *http.Client
	urls  map[int64]string
}

// New connects to the team server at baseURL. token is needed for uploads.
func New(baseURL, token string) *Client {
	return &Client{
		base:  strings.TrimSuffix(baseURL, "/"),
		token: token,
		http:  &http.Client{Timeout: 5 * time.Minute},
		urls:  make(map[int64]string),
	}
}

// Ping checks that the server answers, so the CLI can fail soft early.
func (c *Client) Ping(ctx context.Context) error {
	var info struct {
		Version string `json:"version"`
	}
	return c.getJSON(ctx, "/api/v1/info", &info)
}

func (c *Client) Missing(ctx context.Context, hashes []blob.Hash) ([]blob.Hash, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/blobs/missing", bytes.NewReader(server.EncodeHashes(hashes)), "application/octet-stream", true)
	if err != nil {
		return nil, err
	}
	return server.DecodeHashes(body)
}

func (c *Client) PutBlobs(ctx context.Context, blobs []store.Blob) error {
	var buf bytes.Buffer
	if err := server.EncodeBlobs(&buf, blobs); err != nil {
		return err
	}
	_, err := c.do(ctx, http.MethodPost, "/api/v1/blobs", &buf, "application/octet-stream", true)
	return err
}

func (c *Client) AddRun(ctx context.Context, run store.Run) (int64, error) {
	payload, err := json.Marshal(run)
	if err != nil {
		return 0, err
	}
	body, err := c.do(ctx, http.MethodPost, "/api/v1/runs", bytes.NewReader(payload), "application/json", true)
	if err != nil {
		return 0, err
	}
	var res struct {
		ID  int64  `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, err
	}
	c.urls[res.ID] = res.URL
	return res.ID, nil
}

// RunURL is the run page the server named when the run was registered.
func (c *Client) RunURL(id int64) string {
	if u, ok := c.urls[id]; ok && u != "" {
		return u
	}
	return fmt.Sprintf("%s/runs/%d", c.base, id)
}

func (c *Client) Candidates(ctx context.Context, q store.CandidateQuery) ([]store.Run, error) {
	v := url.Values{"project": {q.Project}, "profile": {q.Profile}, "revision": q.Revisions}
	if q.Stream != "" {
		v.Set("stream", q.Stream)
	}
	for _, k := range q.Kinds {
		v.Add("kind", string(k))
	}
	if q.Change > 0 {
		v.Set("change", strconv.FormatInt(q.Change, 10))
	}
	var runs []store.Run
	if err := c.getJSON(ctx, "/api/v1/base?"+v.Encode(), &runs); err != nil {
		return nil, err
	}
	return runs, nil
}

func (c *Client) Manifest(ctx context.Context, hash blob.Hash) (blob.Manifest, error) {
	data, err := c.do(ctx, http.MethodGet, "/api/v1/blobs/"+hash.String(), nil, "", false)
	if err != nil {
		return blob.Manifest{}, err
	}
	if blob.Sum(data) != hash {
		return blob.Manifest{}, fmt.Errorf("manifest %s: the server sent other content", hash)
	}
	return blob.DecodeManifest(data)
}

func (c *Client) getJSON(ctx context.Context, path string, v any) error {
	body, err := c.do(ctx, http.MethodGet, path, nil, "", false)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, contentType string, auth bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if auth && c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("%s %s: %s", method, path, apiErr.Error)
		}
		return nil, fmt.Errorf("%s %s: %s", method, path, resp.Status)
	}
	return data, nil
}
