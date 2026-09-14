package entities

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestClassifyOutboundMedia(t *testing.T) {
	cases := []struct {
		mime     string
		size     int
		wantType string
		wantErr  bool
	}{
		{"image/png", 1024, MediaTypeImage, false},
		{"image/jpeg; charset=binary", 1024, MediaTypeImage, false},
		{"application/pdf", MaxAttachmentBytes, MediaTypeDocument, false},
		{"video/mp4", 2048, MediaTypeVideo, false},
		{"audio/ogg", 2048, MediaTypeAudio, false},
		{"image/png", MaxImageBytes + 1, "", true},
		{"application/pdf", MaxAttachmentBytes + 1, "", true},
		{"application/x-msdownload", 10, "", true},
		{"image/gif", 10, "", true},
		{"application/pdf", 0, "", true},
	}
	for _, tc := range cases {
		got, err := ClassifyOutboundMedia(tc.mime, tc.size)
		if tc.wantErr {
			if err == nil || !errors.Is(err, ErrMediaNotAllowed) {
				t.Fatalf("%s/%d: expected ErrMediaNotAllowed, got %v", tc.mime, tc.size, err)
			}
			continue
		}
		if err != nil || got != tc.wantType {
			t.Fatalf("%s/%d: got %q %v", tc.mime, tc.size, got, err)
		}
	}
}

func TestChatMediaKey(t *testing.T) {
	now := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
	key := ChatMediaKey(26, "application/pdf", "factura.pdf", "abc", now)
	if key != "whatsapp/26/2026/09/abc.pdf" {
		t.Fatalf("unexpected key %q", key)
	}
	if got := ChatMediaKey(26, "application/octet-stream", "notas.ZIP", "x", now); !strings.HasSuffix(got, ".zip") {
		t.Fatalf("extension from filename expected, got %q", got)
	}
	if got := ChatMediaKey(26, "", "", "x", now); !strings.HasSuffix(got, ".bin") {
		t.Fatalf("fallback extension expected, got %q", got)
	}
}

func TestSafeMediaFilename(t *testing.T) {
	if got := SafeMediaFilename("../../etc/passwd"); got != "passwd" {
		t.Fatalf("path traversal must be stripped, got %q", got)
	}
	if got := SafeMediaFilename("C:\\Users\\ana\\foto.png"); got != "foto.png" {
		t.Fatalf("windows path must be stripped, got %q", got)
	}
}
