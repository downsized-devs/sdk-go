package language

import "net/http"

const (
	FileUploadBackDateErrorID                 = "Backdate, tanggal file tidak bisa diterima sebelum tanggal mulai siklus."
	FileUploadWrongFileErrorID                = "Mohon periksa kembali format dan template excel."
	FileUploadForwardDateErrorID              = "Forward date, tanggal file tidak bisa diterima setelah tanggal mulai siklus"
	FileUploadInvalidDateErrorID              = "Mohon periksa kembali tanggal yang ada pada kolom excel untuk tidak melebihi waktu saat ini atau kurang dari waktu doc pertama."
	FileUploadInvalidPondIsNotRegisterErrorID = "Kolam tidak terdaftar. Mohon periksa template Anda."
	FileUploadWrongDateFormatErrorID          = "Format tanggal salah"
	FileUploadUnregisteredPondErrorID         = "Kolam belum terdaftar"
	FileUploadIncorrectDataFormatErrorID      = "Format data salah"
	FileUploadDocDopErrorID                   = "Kalkulasi DOC/DOP tidak sesuai"
)

// Surprisingly, we have a kinda official translation for HTTP status in Indonesian.
// Check it here: https://id.wikipedia.org/wiki/Daftar_kode_status_HTTP
var statusTextId = map[int]string{
	http.StatusContinue:           "Lanjutkan",
	http.StatusSwitchingProtocols: "Beralih Protokol",
	http.StatusProcessing:         "Processing",
	http.StatusEarlyHints:         "Petunjuk Awal",

	http.StatusOK:                   "OK",
	http.StatusCreated:              "Dibuat",
	http.StatusAccepted:             "Diterima",
	http.StatusNonAuthoritativeInfo: "Informasi Non-Resmi",
	http.StatusNoContent:            "Tanpa Konten",
	http.StatusResetContent:         "Setel Ulang Konten",
	http.StatusPartialContent:       "Konten Sebagian",
	http.StatusMultiStatus:          "Multi-Status",
	http.StatusAlreadyReported:      "Sudah Dilaporkan",
	http.StatusIMUsed:               "IM Used",

	http.StatusMultipleChoices:   "Pilihan ganda",
	http.StatusMovedPermanently:  "Dipindahkan Permanen",
	http.StatusFound:             "Ditemukan",
	http.StatusSeeOther:          "Lihat Lainnya",
	http.StatusNotModified:       "Tidak dimodifikasi",
	http.StatusUseProxy:          "Gunakan proxy",
	http.StatusTemporaryRedirect: "Pengalihan Sementara",
	http.StatusPermanentRedirect: "Pengalihan Permanen",

	http.StatusBadRequest:                   "Bad Request",
	http.StatusUnauthorized:                 "Tidak diperbolehkan",
	http.StatusPaymentRequired:              "Payment Required",
	http.StatusForbidden:                    "Terlarang",
	http.StatusNotFound:                     "Tidak Ditemukan",
	http.StatusMethodNotAllowed:             "Metode Tidak Diizinkan",
	http.StatusNotAcceptable:                "Tidak dapat diterima",
	http.StatusProxyAuthRequired:            "Diperlukan Otentikasi Proksi",
	http.StatusRequestTimeout:               "Request Timeout",
	http.StatusConflict:                     "Konflik",
	http.StatusGone:                         "Hilang",
	http.StatusLengthRequired:               "Panjang Diperlukan",
	http.StatusPreconditionFailed:           "Precondition Failed",
	http.StatusRequestEntityTooLarge:        "Payload Terlalu Besar",
	http.StatusRequestURITooLong:            "URI Terlalu Panjang",
	http.StatusUnsupportedMediaType:         "Jenis Media yang Tidak Didukung",
	http.StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	http.StatusExpectationFailed:            "Expectation Failed",
	http.StatusTeapot:                       "I'm a teapot",
	http.StatusMisdirectedRequest:           "Misdirected Request",
	http.StatusUnprocessableEntity:          "Unprocessable Entity",
	http.StatusLocked:                       "Locked",
	http.StatusFailedDependency:             "Failed Dependency",
	http.StatusTooEarly:                     "Too Early",
	http.StatusUpgradeRequired:              "Diperlukan Peningkatan",
	http.StatusPreconditionRequired:         "Precondition Required",
	http.StatusTooManyRequests:              "Too Many Requests",
	http.StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	http.StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	http.StatusInternalServerError:           "Kesalahan peladen dalam",
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
