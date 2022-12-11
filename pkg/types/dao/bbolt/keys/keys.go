package keys

import "fmt"

// Key 键
type Key string

// Format format
func (k Key) Format(params ...interface{}) []byte {
	return []byte(fmt.Sprintf(string(k), params...))
}

var (
	URLBucket = Key("urls")

	// FullURL1
	// params: url
	FullURL1 Key = "url/full/%s"

	// ShortenedURL1
	// params: hash
	ShortenedURL1 Key = "url/shortened/%s"
)
