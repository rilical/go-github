package specfetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FetchResult struct {
	Data      []byte
	ETag      string
	Unchanged bool
}

type Fetcher interface {
	Fetch(ctx context.Context, url, lastETag string) (FetchResult, error)
}

type fetcher struct {
	client *http.Client
}

func New(client *http.Client) Fetcher {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	return fetcher{client: client}
}

const maxSpecBytes = 64 << 20

func (f fetcher) Fetch(ctx context.Context, url, lastETag string) (FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return FetchResult{}, err
	}
	if lastETag != "" {
		req.Header.Set("If-None-Match", lastETag)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return FetchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return FetchResult{ETag: lastETag, Unchanged: true}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return FetchResult{}, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxSpecBytes+1))
	if err != nil {
		return FetchResult{}, err
	}
	if len(data) > maxSpecBytes {
		return FetchResult{}, fmt.Errorf("fetch %s: spec exceeds %d-byte cap", url, maxSpecBytes)
	}

	// Some origins serve the spec without an ETag. Synthesize a strong validator
	// from the body so downstream change-detection and the emit window marker have
	// a stable per-content identity instead of an empty string.
	etag := resp.Header.Get("ETag")
	if etag == "" {
		sum := sha256.Sum256(data)
		etag = "sha256:" + hex.EncodeToString(sum[:])
	}

	return FetchResult{Data: data, ETag: etag}, nil
}
