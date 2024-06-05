package language

import "net/http"

const (
	FileUploadBackDateErrorEN                 = "Backdate, file date cannot be accepted before cycle start date"
	FileUploadWrongFileErrorEN                = "Please check excel template and format."
	FileUploadForwardDateErrorEN              = "Forward date, file date cannot be accepted beyond cycle start date"
	FileUploadInvalidDateErrorEN              = "Please check the date column in excel to not exceed the current time or before the start of doc"
	FileUploadInvalidPondIsNotRegisterErrorEN = "Pond is not registered. Please check the template file"
	FileUploadWrongDateFormatErrorEN          = "Incorrect date format"
	FileUploadUnregisteredPondErrorEN         = "Unregistered pond"
	FileUploadIncorrectDataFormatErrorEN      = "Incorrect data format"
	FileUploadDocDopErrorEN                   = "Incorrcet DOC/DOP calculation"
)

var statusTextEn = map[int]string{
	http.StatusContinue:           "Continue",
	http.StatusSwitchingProtocols: "Switching Protocols",
	http.StatusProcessing:         "Processing",
	http.StatusEarlyHints:         "Early Hints",

	http.StatusOK:                   "OK",
	http.StatusCreated:              "Created",
	http.StatusAccepted:             "Accepted",
	http.StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	http.StatusNoContent:            "No Content",
	http.StatusResetContent:         "Reset Content",
	http.StatusPartialContent:       "Partial Content",
	http.StatusMultiStatus:          "Multi-Status",
	http.StatusAlreadyReported:      "Already Reported",
	http.StatusIMUsed:               "IM Used",

	http.StatusMultipleChoices:   "Multiple Choices",
	http.StatusMovedPermanently:  "Moved Permanently",
	http.StatusFound:             "Found",
	http.StatusSeeOther:          "See Other",
	http.StatusNotModified:       "Not Modified",
	http.StatusUseProxy:          "Use Proxy",
	http.StatusTemporaryRedirect: "Temporary Redirect",
	http.StatusPermanentRedirect: "Permanent Redirect",

	http.StatusBadRequest:                   "Bad Request",
	http.StatusUnauthorized:                 "Unauthorized",
	http.StatusPaymentRequired:              "Payment Required",
	http.StatusForbidden:                    "Forbidden",
	http.StatusNotFound:                     "Not Found",
	http.StatusMethodNotAllowed:             "Method Not Allowed",
	http.StatusNotAcceptable:                "Not Acceptable",
	http.StatusProxyAuthRequired:            "Proxy Authentication Required",
	http.StatusRequestTimeout:               "Request Timeout",
	http.StatusConflict:                     "Conflict",
	http.StatusGone:                         "Gone",
	http.StatusLengthRequired:               "Length Required",
	http.StatusPreconditionFailed:           "Precondition Failed",
	http.StatusRequestEntityTooLarge:        "Request Entity Too Large",
	http.StatusRequestURITooLong:            "Request URI Too Long",
	http.StatusUnsupportedMediaType:         "Unsupported Media Type",
	http.StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	http.StatusExpectationFailed:            "Expectation Failed",
	http.StatusTeapot:                       "I'm a teapot",
	http.StatusMisdirectedRequest:           "Misdirected Request",
	http.StatusUnprocessableEntity:          "Unprocessable Entity",
	http.StatusLocked:                       "Locked",
	http.StatusFailedDependency:             "Failed Dependency",
	http.StatusTooEarly:                     "Too Early",
	http.StatusUpgradeRequired:              "Upgrade Required",
	http.StatusPreconditionRequired:         "Precondition Required",
	http.StatusTooManyRequests:              "Too Many Requests",
	http.StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	http.StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	http.StatusInternalServerError:           "Internal Server Error",
	http.StatusNotImplemented:                "Not Implemented",
	http.StatusBadGateway:                    "Bad Gateway",
	http.StatusServiceUnavailable:            "Service Unavailable",
	http.StatusGatewayTimeout:                "Gateway Timeout",
	http.StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	http.StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	http.StatusInsufficientStorage:           "Insufficient Storage",
	http.StatusLoopDetected:                  "Loop Detected",
	http.StatusNotExtended:                   "Not Extended",
	http.StatusNetworkAuthenticationRequired: "Network Authentication Required",
}
