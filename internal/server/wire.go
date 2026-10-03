package server

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// The upload protocol sends hash lists as raw 16-byte hashes, and blobs as a
// gzip stream of frames: the hash, the length as a varint, then the bytes.
// Both are much smaller than JSON, and a manifest of 30,000 files has three
// hashes per file.

// EncodeHashes writes hashes back to back.
func EncodeHashes(hashes []blob.Hash) []byte {
	out := make([]byte, 0, len(hashes)*blob.HashSize)
	for _, h := range hashes {
		out = append(out, h[:]...)
	}
	return out
}

// DecodeHashes reads a list written by EncodeHashes.
func DecodeHashes(data []byte) ([]blob.Hash, error) {
	if len(data)%blob.HashSize != 0 {
		return nil, fmt.Errorf("hash list of %d bytes is not a multiple of %d", len(data), blob.HashSize)
	}
	hashes := make([]blob.Hash, len(data)/blob.HashSize)
	for i := range hashes {
		copy(hashes[i][:], data[i*blob.HashSize:])
	}
	return hashes, nil
}

// EncodeBlobs writes blobs as one gzip stream of frames.
func EncodeBlobs(w io.Writer, blobs []store.Blob) error {
	zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return err
	}
	var size [binary.MaxVarintLen64]byte
	for _, b := range blobs {
		if _, err := zw.Write(b.Hash[:]); err != nil {
			return err
		}
		n := binary.PutUvarint(size[:], uint64(len(b.Data)))
		if _, err := zw.Write(size[:n]); err != nil {
			return err
		}
		if _, err := zw.Write(b.Data); err != nil {
			return err
		}
	}
	return zw.Close()
}

// the largest blob a frame may carry; the manifest of a run is the biggest
const maxFrame = 64 << 20

// DecodeBlobs reads frames written by EncodeBlobs and hands each blob to fn.
func DecodeBlobs(r io.Reader, fn func(store.Blob) error) error {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("blob stream: %w", err)
	}
	br := bufio.NewReader(zr)
	for {
		var h blob.Hash
		if _, err := io.ReadFull(br, h[:]); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("blob stream: %w", err)
		}
		size, err := binary.ReadUvarint(br)
		if err != nil {
			return fmt.Errorf("blob stream: %w", err)
		}
		if size > maxFrame {
			return fmt.Errorf("blob %s: %d bytes is too large", h, size)
		}
		data := make([]byte, size)
		if _, err := io.ReadFull(br, data); err != nil {
			return fmt.Errorf("blob stream: %w", err)
		}
		if err := fn(store.Blob{Hash: h, Data: data}); err != nil {
			return err
		}
	}
}
