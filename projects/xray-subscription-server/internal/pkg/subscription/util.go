package subscription

import (
	"encoding/base64"
	"net/url"
	"strings"
)

func render(uris []*url.URL) []byte {
	var sb strings.Builder
	for _, u := range uris {
		sb.WriteString(u.String())
		sb.WriteByte('\n')
	}

	return []byte(base64.StdEncoding.EncodeToString([]byte(sb.String())))
}
