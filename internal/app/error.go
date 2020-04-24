package app

import (
	"net/http"
)

func getResponseText(title string, body string) string {
	return "<html>" +
		"\n\t<head>" +
		"\n\t\t<title>" + title + "</title>" +
		"\n\t</head>" +
		"\n\t<body>" +
		"\n\t\t<h1>" + title + "</h1>" +
		"\n\t\t<p>"+ body +"</p>" +
		"\n\t</body>" +
		"\n</html>"
}

const (
	TOO_MANY_REQ_ERROR = "Too Many Requests"
	WRONG_IP_ERROR = "Wrong IP"
	NO_ACCESS_ERROR = "No access to admin request"
)

func Respond(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(body))
}