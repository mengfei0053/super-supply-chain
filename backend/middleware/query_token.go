package middleware

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var tokenQueryPattern = regexp.MustCompile(`(?i)(token=)[^&\s]*`)

// RedactSensitiveQuery removes the secret value of a token query parameter.
func RedactSensitiveQuery(raw string) string {
	if raw == "" || !strings.Contains(strings.ToLower(raw), "token=") {
		return raw
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return "token=REDACTED"
	}
	if _, ok := values["token"]; ok {
		values.Set("token", "REDACTED")
	}
	return values.Encode()
}

// CaptureAndRedactQueryToken keeps ?token= available to auth, then removes it
// from the request URL so later logs cannot record the secret.
func CaptureAndRedactQueryToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.Request.URL.RawQuery
		if strings.Contains(strings.ToLower(raw), "token=") {
			values, err := url.ParseQuery(raw)
			if err != nil {
				c.Request.URL.RawQuery = ""
				c.Request.RequestURI = c.Request.URL.Path
			} else if token := strings.TrimSpace(values.Get("token")); token != "" {
				c.Set(contextQueryToken, token)
				values.Del("token")
				c.Request.URL.RawQuery = values.Encode()
				if c.Request.URL.RawQuery == "" {
					c.Request.RequestURI = c.Request.URL.Path
				} else {
					c.Request.RequestURI = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
				}
			}
		}
		c.Next()
	}
}

func scrubSensitiveDump(dump string) string {
	lines := strings.Split(dump, "\n")
	for i, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		switch {
		case strings.HasPrefix(lower, "authorization:"),
			strings.HasPrefix(lower, "x-api-key:"),
			strings.HasPrefix(lower, "cookie:"):
			name := strings.SplitN(line, ":", 2)[0]
			lines[i] = name + ": REDACTED"
		default:
			if strings.Contains(lower, "token=") {
				lines[i] = tokenQueryPattern.ReplaceAllString(line, "${1}REDACTED")
			}
		}
	}
	return strings.Join(lines, "\n")
}
