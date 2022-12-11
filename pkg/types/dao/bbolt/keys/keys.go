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

	// FullURL2
	// params: url 和 短链接 hash
	FullURL2 Key = "url/full/%s/%s"

	// ShortenedURL1
	// params: hash
	ShortenedURL1 Key = "url/shortened/%s"
)
