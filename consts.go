package meater

var API_SERVER = "https://public-api.cloud.meater.com/"

var API_VERSION = "v1"

var apiURLs = map[string]string{
	"API_LOGIN":   "/login",
	"API_DEVICES": "/devices",
}

var apiStatuses = map[string]string{
	"API_STATUS_OK": "OK", // 200
}

var apiErrors = map[string]string{
	"API_ERROR_NOT_FOUND":             "Not Found",             // 404
	"API_ERROR_UNAUTHORIZED":          "Unauthorized",          // 401
	"API_ERROR_BAD_REQUEST":           "Bad Request",           // 400
	"API_ERROR_TOO_MANY_REQUESTS":     "Too Many Requests",     // 429
	"API_ERROR_INTERNAL_SERVER_ERROR": "Internal Server Error", // 500
}

const (
	GET  = "GET"
	POST = "POST"
)
