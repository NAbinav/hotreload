package watcher

import "testing"

func TestShouldIgnore(t *testing.T) {

	cases := []struct {
		path string
		want bool
	}{
		// obvious ignores
		{".git/config", true},
		{"node_modules/lodash/index.js", true},
		{"bin/server", true},
		{"./bin/app", true},

		//break when using strings.Contains
		{"binary/main.go", false},
		{"notgit/file", false},
		{"my_node/index.js", false},

		// regular paths
		{"main.go", false},
		{"internal/server/handler.go", false},
		{"cmd/app/main.go", false},
	}

	for _, c := range cases {

		got := shouldIgnore(c.path)

		if got != c.want {
			t.Errorf("shouldIgnore(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}
