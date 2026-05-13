package codes

import (
	"net/http"

	"github.com/downsized-devs/sdk-go/language"
)

type Message struct {
	StatusCode int
	TitleEN    string
	TitleID    string
	BodyEN     string
	BodyID     string
}

// HTTP message
var (
	// 4xx
	ErrMsgBadRequest = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusBadRequest),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusBadRequest),
		BodyEN:     "Invalid input. Please validate your input.",
		BodyID:     "Input data tidak valid. Mohon cek kembali input data anda.",
	}
	ErrMsgUnauthorized = Message{
		StatusCode: http.StatusUnauthorized,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusUnauthorized),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusUnauthorized),
		BodyEN:     "Unauthorized access. You are not authorized to access this resource.",
		BodyID:     "Akses ditolak. Anda tidak memiliki izin untuk mengakses laman ini.",
	}
	ErrMsgInvalidToken = Message{
		StatusCode: http.StatusUnauthorized,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusUnauthorized),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusUnauthorized),
		BodyEN:     "Invalid token. Please renew your session by reloading.",
		BodyID:     "Token tidak valid. Mohon perbarui sesi dengan memuat ulang laman.",
	}
	ErrMsgRefreshTokenExpired = Message{
		StatusCode: http.StatusUnauthorized,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusUnauthorized),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusUnauthorized),
		BodyEN:     "Session refresh token has expired. Please renew your session by reloading.",
		BodyID:     "Token pembaruan sudah tidak berlaku. Mohon perbarui sesi dengan memuat ulang laman.",
	}
	ErrMsgAccessTokenExpired = Message{
		StatusCode: http.StatusUnauthorized,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusUnauthorized),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusUnauthorized),
		BodyEN:     "Session access token has expired. Please renew your session by reloading.",
		BodyID:     "Token akses sudah tidak berlaku. Mohon perbarui sesi dengan memuat ulang laman.",
	}
	ErrMsgForbidden = Message{
		StatusCode: http.StatusForbidden,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusForbidden),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusForbidden),
		BodyEN:     "Forbidden. You don't have permission to access this resource.",
		BodyID:     "Terlarang. Anda tidak memiliki izin untuk mengakses laman ini.",
	}
	ErrMsgRevokeRefreshTokenFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusInternalServerError),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusInternalServerError),
		BodyEN:     "Failed revoking refresh token.",
		BodyID:     "Gagal mencabut refresh token.",
	}
	ErrMsgNotFound = Message{
		StatusCode: http.StatusNotFound,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusNotFound),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusNotFound),
		BodyEN:     "Record does not exist. Please validate your input or contact the administrator.",
		BodyID:     "Data tidak ditemukan. Mohon cek kembali input data anda atau hubungi administrator.",
	}
	ErrMsgContextTimeout = Message{
		StatusCode: http.StatusRequestTimeout,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusRequestTimeout),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusRequestTimeout),
		BodyEN:     "Request time has been exceeded.",
		BodyID:     "Waktu permintaan telah habis.",
	}
	ErrMsgConflict = Message{
		StatusCode: http.StatusConflict,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusConflict),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusConflict),
		BodyEN:     "Record has existed. Please validate your input or contact the administrator.",
		BodyID:     "Data sudah ada. Mohon cek kembali input data anda atau hubungi administrator.",
	}
	ErrMsgTooManyRequest = Message{
		StatusCode: http.StatusTooManyRequests,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusTooManyRequests),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusTooManyRequests),
		BodyEN:     "Too many requests. Please wait and try again after a few moments.",
		BodyID:     "Terlalu banyak permintaan. Mohon tunggu dan coba lagi setelah beberapa saat.",
	}

	// 5xx
	ErrMsgInternalServerError = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusInternalServerError),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusInternalServerError),
		BodyEN:     "Internal server error. Please contact the administrator.",
		BodyID:     "Terjadi kendala pada server. Mohon hubungi administrator.",
	}
	ErrMsgNotImplemented = Message{
		StatusCode: http.StatusNotImplemented,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusNotImplemented),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusNotImplemented),
		BodyEN:     "Not Implemented. Please contact the administrator.",
		BodyID:     "Layanan tidak tersedia. Mohon hubungi administrator.",
	}
	ErrMsgServiceUnavailable = Message{
		StatusCode: http.StatusServiceUnavailable,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusServiceUnavailable),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusServiceUnavailable),
		BodyEN:     "Service is unavailable. Please contact the administrator.",
		BodyID:     "Layanan sedang tidak tersedia. Mohon hubungi administrator.",
	}

	// Application specific messages
	ErrMsgResetPassword = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Failed To Reset Password",
		TitleID:    "Gagal Mengatur Ulang Kata Sandi",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgPasswordDoesNotMatch = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Entered Password Does Not Match",
		TitleID:    "Kata Sandi Yang Dimasukkan Salah",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgPasswordIsWeak = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Entered Password Combination Is Weak. Please Choose a Stronger Password",
		TitleID:    "Kombinasi Kata Sandi Yang Dimasukkan Lemah. Mohon Masukkan Password yang Lebih Kuat",
		BodyEN:     "Password must be at least 8 characters and include at least one uppercase letter, one lowercase letter, and one special character",
		BodyID:     "Kata sandi minimal 8 karakter dan harus mengandung minimal satu huruf besar, satu huruf kecil, dan satu karakter khusus",
	}

	ErrMsgResetTokenExpired = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Reset Password Token Has Expired",
		TitleID:    "Token Atur Ulang Kata Sandi Sudah Kadaluwarsa",
		BodyEN:     "Please create a new password reset request",
		BodyID:     "Silakan lakukan permintaan pengaturan ulang kata sandi kembali",
	}

	ErrMsgEmptyEmail = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Your Email is Empty",
		TitleID:    "Alamat Email Kosong",
		BodyEN:     "Please input your email",
		BodyID:     "Silakan isi email anda",
	}

	ErrMsgInvalidEmail = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Your Email is Invalid",
		TitleID:    "Alamat Email Tidak Valid",
		BodyEN:     "Please input a valid email",
		BodyID:     "Silakan isi email yang valid",
	}

	ErrMsgSameCurrentPassword = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Your new password cannot be the same with your current password",
		TitleID:    "Kata sandi baru anda tidak boleh sama dengan kata sandi saat ini",
		BodyEN:     "Please input a new password",
		BodyID:     "Silakan isi password yang baru",
	}

	ErrMsgResetTokenInvalid = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Reset Password Token Is Invalid",
		TitleID:    "Token Reset Pasword Tidak Valid",
		BodyEN:     "Please check your token",
		BodyID:     "Silakan periksa kembali token anda",
	}

	ErrMsgLockExist = Message{
		StatusCode: http.StatusTooManyRequests,
		TitleEN:    "Please wait for a while before requesting a new password",
		TitleID:    "Mohon tunggu sejenak sebelum meminta password baru",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgExecuteTemplateFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Failed to execute golang template",
		TitleID:    "Gagal mengeksekusi template golang",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgConvertMJMLToHTMLFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Failed to convert MJML to HTML",
		TitleID:    "Gagal mengkonversi MJML ke HTML",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgConvertGoTemplateToPDFFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Failed to convert (html/template) to PDF",
		TitleID:    "Gagal mengkonversi (html/template) ke PDF",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgGenerateToPDFFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Failed to generate file PDF",
		TitleID:    "Gagal menghasilkan file PDF",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgPasswordIsNotFilled = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Your new password and new confirmation password must be filled",
		TitleID:    "Kata sandi baru anda dan konfirmasi kata sandi baru harus diisi",
		BodyEN:     "Please input a new password and the confirmation password",
		BodyID:     "Silakan isi password dan konfirmasi password",
	}

	ErrMsgCodeParseHTMlTemplateFailed = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Failed to parse HTML Template",
		TitleID:    "Gagal menyusun template HTML",
		BodyEN:     "",
		BodyID:     "",
	}

	// Slack Alert
	ErrMsgSlackAlert = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Slack alerting failed",
		TitleID:    "Peringatan slack gagal",
		BodyEN:     "",
		BodyID:     "",
	}

	// Translator Error
	ErrMsgTranslatorlib = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Translator error",
		TitleID:    "Translator mengalami kegagalan",
	}

	// content too large common error
	ErrMsgContentTooLarge = Message{
		StatusCode: http.StatusRequestEntityTooLarge,
		TitleEN:    "File Content Too Large",
		TitleID:    "Konten File Terlalu Besar",
		BodyEN:     "",
		BodyID:     "",
	}

	// Image Too Big (10MB) error
	ErrMsgImageTooBig = Message{
		StatusCode: http.StatusRequestEntityTooLarge,
		TitleEN:    "Maximum Image size is 10 MB",
		TitleID:    "Maksimum ukuran gambar adalah 10 MB",
	}
)

// Application message
var (
	// Files upload
	SuccessDefault = Message{
		StatusCode: http.StatusOK,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusOK),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusOK),
		BodyEN:     "Request successful",
		BodyID:     "Request berhasil",
	}
	SuccessAccepted = Message{
		StatusCode: http.StatusAccepted,
		TitleEN:    language.HTTPStatusText(language.English, http.StatusAccepted),
		TitleID:    language.HTTPStatusText(language.Indonesian, http.StatusAccepted),
		BodyEN:     "Request accepted",
		BodyID:     "Request diterima",
	}

	// Image
	SuccessUploadImage = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Image Upload Success",
		TitleID:    "Anda Telah Berhasil Mengunggah Gambar",
		BodyEN:     "",
		BodyID:     "",
	}

	// Verif Token
	SuccessResetPassword = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Reset password success. Please check your email.",
		TitleID:    "Pengaturan ulang kata sandi berhasil. Silakan cek email anda.",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessInputResetPassword = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Access Granted",
		TitleID:    "Akses Diberikan",
		BodyEN:     "",
		BodyID:     "",
	}
)
