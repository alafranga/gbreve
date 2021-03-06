package usl

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/purell"

	"github.com/alaturka/gbreve/text/textutil"
)

const (
	fallbackScheme = "https"
)

var (
	supportedSchemes = newSupported(
		"https",
		"http",
		"ssh",
		"git",
		"git+ssh",
		"ftp",
		"ftps",
		"file",
	)

	supportedProviders = newSupported(
		"bitbucket.com",
		"github.com",
		"gitlab.com",
		"salsa.debian.org",
	)

	supportedClasses = newSupported(
		"git",
		"tar.bz2",
		"tar.gz",
		"tar.xz",
		"tgz",
		"zip",
	)
)

type supported struct {
	support map[string]struct{}
	list    []string
}

func newSupported(items ...string) *supported {
	s := &supported{
		support: map[string]struct{}{},
	}

	for _, item := range items {
		s.support[item] = struct{}{}
		s.list = append(s.list, item)
	}

	return s
}

func (s *supported) contains(support string) bool {
	_, ok := s.support[support]

	return ok
}

// USL should be commented
type USL struct {
	Class    string // Source class
	Domain   string // url.URL Host without port
	Fragment string // url.URL Fragment
	BasePath string // url.URL Path without leading and trailing slashes
	Host     string // url.URL Host
	ID       string // Source identifier
	InPath   string // Relative path after root source
	Name     string // Name of the source in relative path form
	Password string // url.Userinfo Password
	Path     string // url.URL Port
	Port     string // url.URL Port
	Ref      string // Git reference (i.e. branch, tag, commit)
	Scheme   string // url.URL Scheme
	Source   string // Transport string
	Username string // url.Userinfo Username
}

func newFromURL(u *url.URL) *USL {
	var username, password, domain, port string

	if u.User != nil {
		username = u.User.Username()

		if p, ok := u.User.Password(); ok {
			password = p
		}
	}

	domain = u.Host

	if u.Host != "" {
		var err error
		if domain, port, err = net.SplitHostPort(u.Host); err != nil {
			domain = u.Host
			port = ""
		}
	}

	return &USL{
		Domain:   domain,
		Fragment: u.Fragment,
		Host:     u.Host,
		Password: password,
		Path:     u.Path,
		Port:     port,
		Scheme:   u.Scheme,
		Username: username,
	}
}

// Parse should be commented
func Parse(rawurl string) (*USL, error) {
	u, err := parse(rawurl)
	if err != nil {
		return nil, err
	}

	us := newFromURL(u)
	if err = us.compute(); err != nil {
		return nil, err
	}

	return us, nil
}

// Parse should be commented
func ParseMayLocalPath(rawurl string) (*USL, error) {
	in, err := reduceLocal(rawurl)
	if err != nil {
		return nil, err
	}

	return Parse(in)
}

// Map should be commented
func (us *USL) Map() (map[string]string, []string) {
	m := map[string]string{}

	e := reflect.ValueOf(us).Elem()

	for i := 0; i < e.NumField(); i++ {
		if e.Field(i).CanInterface() {
			k := strings.ToLower(e.Type().Field(i).Name)
			v := reflect.ValueOf(e.Field(i).Interface()).String()

			m[k] = v
		}
	}

	ks := make([]string, 0, len(m))

	for k := range m {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	return m, ks
}

// MapCustom should be commented
func (us *USL) MapCustom(templateMap map[string]string) (map[string]string, []string) {
	m, ks := us.Map()

	for k, v := range templateMap {
		m[k] = textutil.RenderString(v, m)

		ks = append(ks, k)
	}

	sort.Strings(ks)

	return m, ks
}

func IsLocal(in string) bool {
	for _, p := range []string{"/", "./", "../"} {
		if strings.HasPrefix(in, p) {
			return true
		}
	}

	return false
}

func (us *USL) String() string {
	return us.Source
}

// Private methods

func (us *USL) compute() error {
	if strings.HasSuffix(us.Scheme, "+ssh") {
		us.Scheme = "ssh" //nolint:goconst
	}

	if path, ref, ok := parseRef(us.Path); ok {
		us.Path = path
		us.Ref = ref
	}

	if before, class, after, ok := parseClass(us.Path); ok && supportedClasses.contains(class) {
		us.Path = before
		us.Name = relPath(before)
		us.InPath = relPath(after)
		us.Class = class
	}

	us.BasePath = relPath(us.Path)

	if supportedProviders.contains(us.Host) {
		if us.Class == "" {
			us.Class = "git" //nolint:goconst
		}

		if us.Name == "" {
			parts := strings.Split(relPath(us.Path), "/")
			if len(parts) < 2 {
				return fmt.Errorf("incomplete repository path %q for provider %q", us.Path, us.Host)
			}

			us.Name = strings.Join(parts[:2], "/")
			us.InPath = strings.Join(parts[2:], "/")
		}
	}

	if us.Name == "" {
		us.Name = us.BasePath
	}

	if us.Ref != "" && us.Class != "git" {
		return fmt.Errorf("reference found for non git source: %q", us.Ref)
	}

	us.Source = us.source()
	us.ID = us.id()

	return nil
}

func (us *USL) id() string {
	s := us.Source

	if us.Ref != "" {
		var buf strings.Builder

		buf.WriteString(s)
		buf.WriteByte('@')
		buf.WriteString(us.Ref)

		s = buf.String()
	}

	return url.PathEscape(s)
}

func reduceLocal(rawurl string) (string, error) {
	if !IsLocal(rawurl) {
		return rawurl, nil
	}

	var buf strings.Builder

	buf.WriteString("file://")

	abspath, err := filepath.Abs(rawurl)
	if err != nil {
		return "", err
	}

	buf.WriteString(abspath)

	return buf.String(), nil
}

func (us *USL) source() string { //nolint:funlen
	var buf strings.Builder

	if us.Scheme == "file" {
		buf.WriteString(us.Scheme)
		buf.WriteString("://")

		if us.Host != "" {
			buf.WriteString(us.Host)
		}

		buf.WriteString(us.Path)

		if us.Class != "" && us.Class != "git" {
			buf.WriteByte('.')
			buf.WriteString(us.Class)
		}

		return buf.String()
	}

	if us.Scheme == "ssh" && us.Port == "" { //nolint:goconst
		if us.Username != "" {
			buf.WriteString(us.Username)
			buf.WriteByte('@')
		}

		buf.WriteString(us.Host)
		buf.WriteByte(':')
		buf.WriteString(us.Name)

		if us.Class != "" {
			buf.WriteByte('.')
			buf.WriteString(us.Class)
		}

		return buf.String()
	}

	buf.WriteString(us.Scheme)
	buf.WriteString("://")

	if us.Username != "" {
		buf.WriteString(us.Username)

		if us.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(us.Password)
		}

		buf.WriteByte('@')
	}

	buf.WriteString(us.Host)

	if us.Class == "" {
		if us.BasePath != "" {
			buf.WriteByte('/')
			buf.WriteString(us.BasePath)
		}
	} else {
		buf.WriteByte('/')
		buf.WriteString(us.Name)
		buf.WriteByte('.')
		buf.WriteString(us.Class)
	}

	return buf.String()
}

// Helpers

func cut(s string, c string) (string, string) {
	i := strings.Index(s, c)

	if i < 0 {
		return s, ""
	}

	return s[:i], s[i+len(c):]
}

func groupPatternFromSlice(group string, ss []string) string {
	escaped := make([]string, 0, len(ss))

	for _, s := range ss {
		escaped = append(escaped, regexp.QuoteMeta(s))
	}

	return `(?P<` + group + `>` + strings.Join(escaped, "|") + `)`
}

func namedMatches(re *regexp.Regexp, in string) (map[string]string, bool) {
	match := re.FindStringSubmatch(in)

	if match == nil {
		return nil, false
	}

	result := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result, true
}

var reClass = regexp.MustCompile(
	`^(?P<before>.*?)[.]` + groupPatternFromSlice("class", supportedClasses.list) + `(?P<after>/.*)?$`,
)

func parseClass(path string) (string, string, string, bool) {
	if m, ok := namedMatches(reClass, path); ok {
		return m["before"], m["class"], m["after"], true
	}

	return path, "", "", false
}

var reRef = regexp.MustCompile(`^(?P<before>.+)@(?P<ref>[^@]*)$`)

func parseRef(path string) (string, string, bool) {
	if m, ok := namedMatches(reRef, path); ok {
		return m["before"], m["ref"], true
	}

	return path, "", false
}

func parse(rawurl string) (*url.URL, error) {
	in := rawurl

	if IsLocal(in) {
		return nil, fmt.Errorf("local file paths not allowed %q", rawurl)
	}

	if scheme, remaining := cut(in, "://"); remaining == "" {
		if m, ok := matchSpecial(in); ok {
			return parseSpecial(in, m)
		}

		if m, ok := matchSSH(in); ok {
			return parseSSH(in, m)
		}

		in = fallbackScheme + "://" + in
	} else {
		scheme = strings.ToLower(scheme)

		if !supportedSchemes.contains(scheme) {
			return nil, fmt.Errorf("unsupported scheme %q", scheme)
		}
	}

	return parseUsual(in, nil)
}

var reSSH = regexp.MustCompile(
	`^((?P<user>[a-zA-Z0-9_]+)@)?(?P<host>[a-zA-Z0-9._-]+):(?P<path>.*)$`,
)

func matchSSH(in string) (map[string]string, bool) {
	return namedMatches(reSSH, in)
}

func parseSSH(_ string, match map[string]string) (*url.URL, error) {
	return &url.URL{
		Host:   match["host"],
		User:   url.User(match["user"]),
		Path:   filepath.Clean(match["path"]),
		Scheme: "ssh",
	}, nil
}

var reSpecial = regexp.MustCompile(
	`^((?P<user>[a-zA-Z0-9_.-]+)@)?` + groupPatternFromSlice("provider", supportedProviders.list) + `(?P<sep>[/:])` + `(?P<path>.*)?$`, //nolint:lll
)

func matchSpecial(in string) (map[string]string, bool) {
	return namedMatches(reSpecial, in)
}

func parseSpecial(_ string, match map[string]string) (*url.URL, error) {
	user := ""
	scheme := "https"

	if match["sep"] == ":" {
		scheme = "ssh"
		user = match["user"]

		if user == "" {
			user = "git"
		} else if user != "git" {
			return nil, fmt.Errorf("user must be 'git' where found %q", match["user"])
		}
	}

	host := match["provider"]
	path := filepath.Clean(match["path"])

	return &url.URL{
		Host:   host,
		User:   url.User(user),
		Path:   path,
		Scheme: scheme,
	}, nil
}

func parseUsual(in string, _ map[string]string) (*url.URL, error) {
	normurl, err := purell.NormalizeURLString(
		in, purell.FlagsUsuallySafeGreedy|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveFragment,
	)
	if err != nil {
		return nil, err
	}

	return url.Parse(normurl)
}

func relPath(path string) string {
	return strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
}
