package usl

import (
	"os"
	"testing"
)

type testParse struct {
	in  string
	out map[string]string
}

func init() { //nolint:gochecknoinits
	err := os.Chdir("/tmp")
	if err != nil {
		panic(err)
	}
}

//nolint:funlen
func TestParse(t *testing.T) {
	t.Parallel()

	tests := map[string][]testParse{
		"Schemeless usual": {
			{
				"github.com/user/repo", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"id":     `https:%2F%2Fgithub.com%2Fuser%2Frepo.git`,
					"inpath": "",
					"name":   "user/repo",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo@unstable", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"id":     `https:%2F%2Fgithub.com%2Fuser%2Frepo.git@unstable`,
					"inpath": "",
					"name":   "user/repo",
					"ref":    "unstable",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo/a/b@unstable", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"id":     `https:%2F%2Fgithub.com%2Fuser%2Frepo.git@unstable`,
					"inpath": "a/b",
					"name":   "user/repo",
					"ref":    "unstable",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo.git/a/b", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"id":     `https:%2F%2Fgithub.com%2Fuser%2Frepo.git`,
					"inpath": "a/b",
					"name":   "user/repo",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"example.com", map[string]string{
					"source": "https://example.com",

					"class":  "",
					"id":     `https:%2F%2Fexample.com`,
					"inpath": "",
					"name":   "",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"example.com/", map[string]string{
					"source": "https://example.com",

					"class":  "",
					"id":     `https:%2F%2Fexample.com`,
					"inpath": "",
					"name":   "",
					"ref":    "",
					"scheme": "https",
				},
			},
		},
		"Schemeless SSH": {
			{
				"github.com:user/repo", map[string]string{
					"source": "git@github.com:user/repo.git",

					"class":    "git",
					"id":       `git@github.com:user%2Frepo.git`,
					"inpath":   "",
					"name":     "user/repo",
					"ref":      "",
					"scheme":   "ssh",
					"username": "git",
				},
			},
			{
				"git@github.com:user/repo", map[string]string{
					"source": "git@github.com:user/repo.git",

					"class":    "git",
					"id":       `git@github.com:user%2Frepo.git`,
					"inpath":   "",
					"name":     "user/repo",
					"ref":      "",
					"scheme":   "ssh",
					"username": "git",
				},
			},
			{
				"user@example.com:a/b", map[string]string{
					"source": "user@example.com:a/b",

					"class":    "",
					"id":       `user@example.com:a%2Fb`,
					"inpath":   "",
					"name":     "a/b",
					"path":     "a/b",
					"ref":      "",
					"scheme":   "ssh",
					"username": "user",
				},
			},
		},
		"File scheme": {
			{
				"file:///path/to/repo.git/a/b", map[string]string{
					"source": "file:///path/to/repo",

					"class":  "git",
					"id":     `file:%2F%2F%2Fpath%2Fto%2Frepo`,
					"inpath": "a/b",
					"name":   "path/to/repo",
					"ref":    "",
					"scheme": "file",
				},
			},
			{
				"file:///path/to/repo.git/a/b@next", map[string]string{
					"source": "file:///path/to/repo",

					"class":  "git",
					"id":     `file:%2F%2F%2Fpath%2Fto%2Frepo@next`,
					"inpath": "a/b",
					"name":   "path/to/repo",
					"ref":    "next",
					"scheme": "file",
				},
			},
			{
				"file://path/to/repo.git/a/b", map[string]string{
					"source": "file://path/to/repo",

					"class":  "git",
					"id":     `file:%2F%2Fpath%2Fto%2Frepo`,
					"inpath": "a/b",
					"name":   "to/repo",
					"ref":    "",
					"scheme": "file",
				},
			},
			{
				"file:///path/to/repo.zip/a/b", map[string]string{
					"source": "file:///path/to/repo.zip",

					"class":  "zip",
					"id":     `file:%2F%2F%2Fpath%2Fto%2Frepo.zip`,
					"inpath": "a/b",
					"name":   "path/to/repo",
					"ref":    "",
					"scheme": "file",
				},
			},
		},
		"HTTPS scheme": {
			{
				"https://github.com/user/repo", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"domain": "github.com",
					"id":     `https:%2F%2Fgithub.com%2Fuser%2Frepo.git`,
					"inpath": "",
					"name":   "user/repo",
					"port":   "",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"https://example.com:8080/user/repo.git/a/b", map[string]string{
					"source": "https://example.com:8080/user/repo.git",

					"class":  "git",
					"domain": "example.com",
					"id":     `https:%2F%2Fexample.com:8080%2Fuser%2Frepo.git`,
					"inpath": "a/b",
					"name":   "user/repo",
					"port":   "8080",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"https://example.com", map[string]string{
					"source": "https://example.com",

					"class":  "",
					"id":     `https:%2F%2Fexample.com`,
					"inpath": "",
					"name":   "",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"https://example.com/", map[string]string{
					"source": "https://example.com",

					"class":  "",
					"id":     `https:%2F%2Fexample.com`,
					"inpath": "",
					"name":   "",
					"ref":    "",
					"scheme": "https",
				},
			},
		},
		"SSH scheme": {
			{
				"ssh://git@github.com/user/repo", map[string]string{
					"source": "git@github.com:user/repo.git",

					"class":    "git",
					"id":       `git@github.com:user%2Frepo.git`,
					"inpath":   "",
					"name":     "user/repo",
					"ref":      "",
					"scheme":   "ssh",
					"username": "git",
				},
			},
			{
				"ssh://user:pass@example.com/a/b", map[string]string{
					"source": "user@example.com:a/b",

					"class":    "",
					"id":       `user@example.com:a%2Fb`,
					"inpath":   "",
					"name":     "a/b",
					"password": "pass",
					"path":     "/a/b",
					"ref":      "",
					"scheme":   "ssh",
					"username": "user",
				},
			},
			{
				"ssh://user:pass@example.com:22/a/b", map[string]string{
					"source": "ssh://user:pass@example.com:22/a/b",

					"class":    "",
					"domain":   "example.com",
					"host":     "example.com:22",
					"inpath":   "",
					"name":     "a/b",
					"password": "pass",
					"path":     "/a/b",
					"port":     "22",
					"ref":      "",
					"scheme":   "ssh",
					"username": "user",
				},
			},
		},
	}

	for name, ts := range tests {
		ts := ts // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range ts {
				got, err := Parse(tc.in)

				if err != nil {
					t.Errorf("Parse(%q) = unexpected err %q", tc.in, err)
					continue
				}

				m, _ := got.Map()

				for ke, ve := range tc.out {
					if va, ok := m[ke]; ok {
						if ve != va {
							t.Errorf("\t%40s    %-12s\twant: %-12s\tgot:  %-12s", tc.in, ke, ve, va)
						}
					}
				}
			}
		})
	}
}

//nolint:funlen
func TestParseMayLocalPath(t *testing.T) {
	t.Parallel()

	tests := map[string][]testParse{
		"Schemeless file": {
			{
				"/a/b", map[string]string{
					"source": "file:///a/b",

					"class":    "",
					"basepath": "a/b",
					"scheme":   "file",
				},
			},
			{
				"./a/b", map[string]string{
					"source": "file:///tmp/a/b",

					"class":    "",
					"basepath": "tmp/a/b",
					"scheme":   "file",
				},
			},
			{
				"../a/b", map[string]string{
					"source": "file:///a/b",

					"class":    "",
					"basepath": "a/b",
					"scheme":   "file",
				},
			},
			{
				"./a/b.git/x@next", map[string]string{
					"source": "file:///tmp/a/b",

					"class":    "git",
					"basepath": "tmp/a/b",
					"inpath":   "x",
					"ref":      "next",
					"scheme":   "file",
				},
			},
			{
				"./a/b.zip/x", map[string]string{
					"source": "file:///tmp/a/b.zip",

					"class":    "zip",
					"basepath": "tmp/a/b",
					"inpath":   "x",
					"scheme":   "file",
				},
			},
		},
	}

	for name, ts := range tests {
		ts := ts // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range ts {
				got, err := ParseMayLocalPath(tc.in)

				if err != nil {
					t.Errorf("Parse(%q) = unexpected err %q", tc.in, err)
					continue
				}

				m, _ := got.Map()

				for ke, ve := range tc.out {
					if va, ok := m[ke]; ok {
						if ve != va {
							t.Errorf("\t%40s    %-12s\twant: %-12s\tgot:  %-12s", tc.in, ke, ve, va)
						}
					}
				}
			}
		})
	}
}

type testTemplate struct {
	in     string
	custom map[string]string
	out    map[string]string
}

//nolint:funlen
func TestTemplate(t *testing.T) {
	t.Parallel()

	tests := map[string][]testTemplate{
		"Schemeless usual": {
			{
				"github.com/user/repo", map[string]string{
					"custom": `{{ .source | pathescape }}`,
				}, map[string]string{
					"source": "https://github.com/user/repo.git",
					"custom": `https:%2F%2Fgithub.com%2Fuser%2Frepo.git`,
				},
			},
		},
	}

	for name, ts := range tests {
		ts := ts // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range ts {
				got, err := Parse(tc.in)

				if err != nil {
					t.Errorf("Parse(%q) = unexpected err %q", tc.in, err)
					continue
				}

				m, _ := got.MapCustom(tc.custom)

				for ke, ve := range tc.out {
					if va, ok := m[ke]; ok {
						if ve != va {
							t.Errorf("\t%40s    %-12s\twant: %-12s\tgot:  %-12s", tc.in, ke, ve, va)
						}
					}
				}
			}
		})
	}
}
