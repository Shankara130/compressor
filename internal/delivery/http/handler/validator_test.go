package handler

import "testing"

func TestMaxSizeForMime(t *testing.T) {
	tests := []struct {
		mime string
		want int64
	}{
		{"image/jpeg", MaxImageSize},
		{"image/png", MaxImageSize},
		{"video/mp4", MaxVideoSize},
		{"video/x-matroska", MaxVideoSize},
		{"application/pdf", MaxPDFSize},
		{"application/octet-stream", 0},
		{"text/plain", 0},
		{"", 0},
	}
	for _, tc := range tests {
		if got := MaxSizeForMime(tc.mime); got != tc.want {
			t.Errorf("MaxSizeForMime(%q) = %d, want %d", tc.mime, got, tc.want)
		}
	}
}

func TestValidateFile(t *testing.T) {
	tests := []struct {
		name     string
		mime     string
		size     int64
		filename string
		wantErr  bool
	}{
		{"valid jpeg", "image/jpeg", 1 << 20, "photo.jpg", false},
		{"valid png", "image/png", 1 << 20, "photo.png", false},
		{"valid mp4", "video/mp4", 10 << 20, "clip.mp4", false},
		{"valid pdf", "application/pdf", 1 << 20, "doc.pdf", false},
		{"image wrong ext", "image/jpeg", 1 << 20, "photo.gif", true},
		{"image too large", "image/jpeg", MaxImageSize + 1, "photo.jpg", true},
		{"video too large", "video/mp4", MaxVideoSize + 1, "clip.mp4", true},
		{"pdf too large", "application/pdf", MaxPDFSize + 1, "doc.pdf", true},
		{"unsupported type", "text/plain", 100, "note.txt", true},
		{"zero-size valid type", "image/jpeg", 0, "photo.jpg", false},
		{"boundary image size", "image/png", MaxImageSize, "photo.png", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFile(tc.mime, tc.size, tc.filename)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateFile(%q, %d, %q) err = %v, wantErr %v",
					tc.mime, tc.size, tc.filename, err, tc.wantErr)
			}
		})
	}
}

func TestDetectMime(t *testing.T) {
	opaque := []byte{0x00, 0x01, 0x02, 0x03} // sniffs as application/octet-stream
	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xE0}    // real JPEG magic

	cases := []struct {
		name string
		sniff []byte
		file  string
		want  string
	}{
		{"mov fallback", opaque, "clip.mov", "video/quicktime"},
		{"mkv fallback", opaque, "clip.mkv", "video/x-matroska"},
		{"avi fallback", opaque, "clip.avi", "video/x-msvideo"},
		{"mp4 fallback", opaque, "clip.mp4", "video/mp4"},
		{"pdf fallback", opaque, "doc.pdf", "application/pdf"},
		{"sniff wins over extension", jpeg, "renamed.mov", "image/jpeg"},
		{"unknown stays opaque", opaque, "file.xyz", "application/octet-stream"},
	}
	for _, c := range cases {
		if got := detectMime(c.sniff, c.file); got != c.want {
			t.Errorf("%s: detectMime got %q, want %q", c.name, got, c.want)
		}
	}
}
