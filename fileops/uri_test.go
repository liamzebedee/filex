package fileops

import "testing"

// ParseURIList and UriToPath are tested here because they're pure logic
// that was moved out of the ui package (which requires GTK) for testability.

func TestParseURIList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "single URI",
			input: "file:///home/user/file.txt\r\n",
			want:  []string{"file:///home/user/file.txt"},
		},
		{
			name:  "multiple URIs",
			input: "file:///home/user/a.txt\r\nfile:///home/user/b.txt\r\n",
			want:  []string{"file:///home/user/a.txt", "file:///home/user/b.txt"},
		},
		{
			name:  "with comments",
			input: "# comment\r\nfile:///home/user/file.txt\r\n",
			want:  []string{"file:///home/user/file.txt"},
		},
		{
			name:  "empty lines ignored",
			input: "\r\nfile:///tmp/a\r\n\r\nfile:///tmp/b\r\n",
			want:  []string{"file:///tmp/a", "file:///tmp/b"},
		},
		{
			name:  "unix line endings",
			input: "file:///tmp/a\nfile:///tmp/b\n",
			want:  []string{"file:///tmp/a", "file:///tmp/b"},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "only whitespace",
			input: "  \r\n  \r\n",
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseURIList(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("ParseURIList() = %v (len %d), want %v (len %d)", got, len(got), tt.want, len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ParseURIList()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestURIToPath(t *testing.T) {
	tests := []struct {
		uri  string
		want string
	}{
		{"file:///home/user/file.txt", "/home/user/file.txt"},
		{"file:///tmp/a b c.txt", "/tmp/a b c.txt"},
		{"/already/a/path", "/already/a/path"},
		{"file://", "file://"},   // too short, returned as-is
		{"http://example.com", "http://example.com"}, // not file://, returned as-is
	}
	for _, tt := range tests {
		got := URIToPath(tt.uri)
		if got != tt.want {
			t.Errorf("URIToPath(%q) = %q, want %q", tt.uri, got, tt.want)
		}
	}
}
