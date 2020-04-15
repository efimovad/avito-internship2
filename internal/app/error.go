package app

import (
	"net/http"
	"strconv"
	"time"
)

func get429ErrorText(limit int) string {
	return "<html>" +
		"\n\t<head>" +
		"\n\t\t<title>Too Many Requests</title>" +
		"\n\t</head>" +
		"\n\t<body>" +
		"\n\t\t<h1>Too Many Requests</h1>" +
		"\n\t\t<p>I only allow " + strconv.Itoa(limit) + " requests per period to this Web site per net. Try again soon.</p>" +
		"\n\t</body>" +
		"\n</html>"
}

func Error (w http.ResponseWriter, error string, code int, wait time.Duration) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Retry-After", wait.String())
	w.WriteHeader(code)
	_, _ = w.Write([]byte(error))
}
