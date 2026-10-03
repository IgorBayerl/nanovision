package blob

import (
	"bytes"
	"testing"
)

// Every decoder must reject bad input with an error, never a panic or a huge
// allocation, and whatever it accepts must encode back to the same bytes: one
// value has exactly one encoding, which keeps hashes stable.

func FuzzDecodeCoverage(f *testing.F) {
	f.Add(EncodeCoverage(sampleCoverage()))
	f.Add(EncodeCoverage(Coverage{Reports: 1, Lines: []Line{{Number: 1, Hits: 2, ReportHits: []int{2}}}}))
	f.Fuzz(func(t *testing.T, data []byte) {
		c, err := DecodeCoverage(data)
		if err != nil {
			return
		}
		if again := EncodeCoverage(c); !bytes.Equal(again, data) {
			t.Fatalf("decode/encode changed the bytes\n in: %x\nout: %x", data, again)
		}
	})
}

func FuzzDecodeAnalysis(f *testing.F) {
	f.Add(EncodeAnalysis(Analysis{
		TotalLines: 10,
		Functions:  []Function{{Name: "f", StartLine: 1, EndLine: 9, Complexity: intPtr(3)}},
		Statements: []Statement{{StartLine: 2, EndLine: 2, Type: "return"}},
	}))
	f.Fuzz(func(t *testing.T, data []byte) {
		a, err := DecodeAnalysis(data)
		if err != nil {
			return
		}
		if again := EncodeAnalysis(a); !bytes.Equal(again, data) {
			t.Fatalf("decode/encode changed the bytes\n in: %x\nout: %x", data, again)
		}
	})
}

func FuzzDecodeManifest(f *testing.F) {
	f.Add(EncodeManifest(sampleManifest()))
	f.Fuzz(func(t *testing.T, data []byte) {
		m, err := DecodeManifest(data)
		if err != nil {
			return
		}
		if again := EncodeManifest(m); !bytes.Equal(again, data) {
			t.Fatalf("decode/encode changed the bytes\n in: %x\nout: %x", data, again)
		}
	})
}

func FuzzDecodeDiff(f *testing.F) {
	f.Add(EncodeDiff(Diff{Files: []DiffFile{{Path: "a.go", Kind: DiffModified, Added: []int{1, 2, 7}, Modified: []int{9}}}}))
	f.Fuzz(func(t *testing.T, data []byte) {
		d, err := DecodeDiff(data)
		if err != nil {
			return
		}
		if again := EncodeDiff(d); !bytes.Equal(again, data) {
			t.Fatalf("decode/encode changed the bytes\n in: %x\nout: %x", data, again)
		}
	})
}
