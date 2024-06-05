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
	ErrMsgFileUploadSowBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Stocking Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Data Tebar Gagal Diunggah, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadDailyBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Daily Monitoring Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Data Monitoring Harian Gagal Diunggah, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadHarvestBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Harvest Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Data Panen Gagal Diunggah, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadSamplingBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Sampling Monitoring Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Data Sampling Gagal Diunggah, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadTreatmentBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Treatment Upload Failed",
		TitleID:    "Data Treatment Gagal Diunggah",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadBackdate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Data Gagal Diunggah, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadBackDateErrorEN,
		BodyID:     language.FileUploadBackDateErrorID,
	}

	ErrMsgFileUploadDocDopMiscalculation = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, file DOC/DOP is miscalculated",
		TitleID:    "Data Gagal Diunggah, kalkulasi DOC/DOP tidak bersesuaian",
		BodyEN:     language.FileUploadDocDopErrorEN,
		BodyID:     language.FileUploadDocDopErrorID,
	}

	ErrMsgFileUploadWrongDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, the date format is wrong.",
		TitleID:    "Data Gagal Diunggah, format tanggal salah.",
		BodyEN:     language.FileUploadWrongDateFormatErrorEN,
		BodyID:     language.FileUploadWrongDateFormatErrorID,
	}

	ErrMsgFileUploadUnregisteredPond = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, unregistered pond",
		TitleID:    "Data Gagal Diunggah, kolam belum terdaftar",
		BodyEN:     language.FileUploadUnregisteredPondErrorEN,
		BodyID:     language.FileUploadUnregisteredPondErrorID,
	}

	ErrMsgFileUploadIncorrectDataFormat = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, incorrect data format",
		TitleID:    "Data Gagal Diunggah, format data salah",
		BodyEN:     language.FileUploadIncorrectDataFormatErrorEN,
		BodyID:     language.FileUploadIncorrectDataFormatErrorID,
	}

	ErrMsgFileUploadSowForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Stocking Upload Failed",
		TitleID:    "Data Tebar Gagal Diunggah",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileUploadDailyForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Daily Monitoring Upload Failed, file date can't be older than cycle start date.",
		TitleID:    "Monitoring Harian Gagal, tanggal berkas harus lebih baru dari tanggal mulai siklus.",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileUploadHarvestForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Harvest Upload Failed",
		TitleID:    "Data Panen Gagal Diunggah",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileUploadSamplingForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Sampling Upload Failed",
		TitleID:    "Data Sampling Gagal Diunggah",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileUploadTreatmentForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Treatment Upload Failed",
		TitleID:    "Data Treatment Gagal Diunggah",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileUploadForwardDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed",
		TitleID:    "Data Gagal Diunggah",
		BodyEN:     language.FileUploadForwardDateErrorEN,
		BodyID:     language.FileUploadForwardDateErrorID,
	}

	ErrMsgFileInvalidDate = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed",
		TitleID:    "Data Gagal Diunggah",
		BodyEN:     language.FileUploadInvalidDateErrorEN,
		BodyID:     language.FileUploadInvalidDateErrorID,
	}

	ErrMsgFileInvalidPondIsNotRegister = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed",
		TitleID:    "Gagal Mengunggah Berkas",
		BodyEN:     language.FileUploadInvalidPondIsNotRegisterErrorEN,
		BodyID:     language.FileUploadInvalidPondIsNotRegisterErrorID,
	}

	ErrMsgFileUploadWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed, wrong file selected.",
		TitleID:    "Data Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadSowWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Stocking Upload Failed, wrong file selected.",
		TitleID:    "Data Tebar Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadDailyWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Daily Monitoring Upload Failed, wrong file selected.",
		TitleID:    "Data Monitoring Harian Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadHarvestWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Harvest Upload Failed, wrong file selected.",
		TitleID:    "Data Panen Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadSamplingWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Sampling Monitoring Upload Failed, wrong file selected.",
		TitleID:    "Data Sampling Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadTreatmentWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Treatment Upload Failed, wrong file selected.",
		TitleID:    "Data Treatment Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadProductionPlanWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Production Plan Upload Failed, wrong file selected.",
		TitleID:    "Rencana Produksi Gagal Diunggah, berkas terpilih salah.",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadBulkOrderWrongFile = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Bulk Order Upload Failed",
		TitleID:    "Order Bulk Gagal Diunggah",
		BodyEN:     language.FileUploadWrongFileErrorEN,
		BodyID:     language.FileUploadWrongFileErrorID,
	}

	ErrMsgFileUploadSowNotUploaded = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Stocking File Not Uploaded",
		TitleID:    "Data Tebar Belum Diunggah",
		BodyEN:     "Please upload stocking data first before uploading other data.",
		BodyID:     "Mohon unggah data tebar terlebih dahulu sebelum mengunggah data lainnya.",
	}

	ErrMsgFileUploadProductionPlanNotUploaded = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Production Plan File Not Uploaded",
		TitleID:    "File Rencana Produksi Belum Diunggah",
		BodyEN:     "Please upload production plan data first before uploading sow data.",
		BodyID:     "Mohon unggah rencana produksi terlebih dahulu sebelum mengunggah data tebar.",
	}

	ErrMsgFileUploadInvalidExtension = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "File Upload Failed",
		TitleID:    "Data Gagal Diunggah",
		BodyEN:     "File format or file extension not valid",
		BodyID:     "Format file atau ekstensi file tidak valid",
	}

	ErrMsgCycleEndFailedNoActivePond = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "End Cycle Failed",
		TitleID:    "Gagal Mengakhiri Siklus",
		BodyEN:     "Cycle doesn't have active ponds.",
		BodyID:     "Belum ada kolam aktif di dalam siklus.",
	}

	ErrMsgCycleEndFailedNoSelectedPond = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "End Cycle Failed",
		TitleID:    "Gagal Mengakhiri Siklus",
		BodyEN:     "No pond is selected.",
		BodyID:     "Tidak ada kolam yang dipilih.",
	}

	ErrMsgMetricsProductivityNoDaily = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Productivity Data Not Available",
		TitleID:    "Data Produktivitas Tidak Tersedia",
		BodyEN:     "Productivity data is not available daily.",
		BodyID:     "Data produktivitas tidak dapat ditampilkan harian.",
	}

	ErrMsgCodeCycleStartFailedNameTooLong = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Create Cycle Failed, cycle name too long.",
		TitleID:    "Gagal Membuat Siklus, nama siklus terlalu panjang.",
		BodyEN:     "Name is too long.",
		BodyID:     "Nama terlalu panjang.",
	}

	ErrMsgResetPassword = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Failed To Reset Password",
		TitleID:    "Gagal Mengatur Ulang Kata Sandi",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgCycleNotFound = Message{
		StatusCode: http.StatusNotFound,
		TitleEN:    "Cycle Not Found",
		TitleID:    "Siklus Tidak Ditemukan",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgCycleStartFailedOngoingCycle = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Create Cycle Failed, ongoing cycle already exists.",
		TitleID:    "Gagal Membuat Siklus, sudah ada siklus berjalan.",
		BodyEN:     "Please end current cycle before creating new ones",
		BodyID:     "Mohon akhiri siklus yang sedang berjalan sebelum membuat siklus baru",
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

	ErrMsgInsertFeedbackWeeklyLimit = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Cannot Insert Feedback",
		TitleID:    "Saran Tidak Dapat Dikirim",
		BodyEN:     "Feedback can only be given once a week",
		BodyID:     "Saran hanya dapat diberikan sekali setiap minggu",
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

	ErrMsgFarmNoActivePond = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "End Cycle Failed. No Active Ponds.",
		TitleID:    "Tidak Dapat Akhiri Siklus. Tidak Ada Kolam Aktif.",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgModule = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Internal Error. Module Data Unavailable.",
		TitleID:    "Kesalahan Internal. Data Modul Tidak Tersedia.",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgPond = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Internal Error. Pond Data Unavailable.",
		TitleID:    "Kesalahan Internal. Data Kolam Tidak Tersedia.",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgPondDisease = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Internal Error. Pond Disease Data Unavailable.",
		TitleID:    "Kesalahan Internal. Data Penyakit Tambak Tidak Tersedia.",
		BodyEN:     "",
		BodyID:     "",
	}

	// User Assign Role
	ErrMsgWrongAssigneeRole = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "User has no authority to change selected user",
		TitleID:    "Pengguna tidak berhak mengganti pengguna yang dipilih",
		BodyEN:     "Please use authorized user",
		BodyID:     "Gunakan pengguna yang diperbolehkan",
	}

	ErrMsgWrongUser = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "User has no authority to change role",
		TitleID:    "Pengguna tidak berhak mengganti peran",
		BodyEN:     "Please use authorized user",
		BodyID:     "Gunakan pengguna yang diperbolehkan",
	}

	ErrMsgAssignRole = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Assigning role failed",
		TitleID:    "Penggantian Pengguna Gagal",
		BodyEN:     "",
		BodyID:     "",
	}

	// User Delete Role
	ErrMsgWrongUserToDelete = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "User unauthorized to delete the selected user",
		TitleID:    "Pengguna tidak berhak menghapus pengguna yang dipilih",
		BodyEN:     "Please use authorized user",
		BodyID:     "Gunakan pengguna yang diperbolehkan",
	}

	ErrMsgDeleteUser = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Deleting role failed",
		TitleID:    "Penghapusan Pengguna Gagal",
		BodyEN:     "",
		BodyID:     "",
	}

	// Farm Recommendation
	ErrMsgHarvestPlanMinBigger = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "End size value is bigger than start size value",
		TitleID:    "End size melebihi nilai start size",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgHarvestPlanIntervalLower = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "Harvest interval minimum is 6 days",
		TitleID:    "Minimal jarak antar panen 6 hari",
		BodyEN:     "",
		BodyID:     "",
	}

	ErrMsgHarvestPlanMinMaxInterval = Message{
		StatusCode: http.StatusBadRequest,
		TitleEN:    "The interval for the value is 60 - 110",
		TitleID:    "Nilai yang dapat digunakan hanya 60 - 110",
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

	// TImelib Error
	ErrMsgTimelib = Message{
		StatusCode: http.StatusInternalServerError,
		TitleEN:    "Timelib error",
		TitleID:    "Timelib mengalami kegagalan",
	}

	// TImelib Error
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
	SuccessUploadSow = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Stocking Upload Success",
		TitleID:    "Data Tebar Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUploadDailyMonitoring = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Daily Monitoring Upload Successs",
		TitleID:    "Data Monitoring Harian Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUploadSampling = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Sampling Upload Success",
		TitleID:    "Data Sampling Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUploadPartialHarvest = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Partial Harvest Upload Success",
		TitleID:    "Data Panen Parsial Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUploadTotalHarvest = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Total Harvest Upload Success",
		TitleID:    "Data Panen Total Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUploadProductionPlan = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Production Plan Upload Success",
		TitleID:    "Rencana Produksi Berhasil Diunggah",
		BodyEN:     "",
		BodyID:     "",
	}

	// Cycles
	SuccessCreateCycle = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You Have Successfully Created A New Cycle",
		TitleID:    "Anda Telah Berhasil Membuat Siklus Baru",
		BodyEN:     "Input data by uploading files",
		BodyID:     "Masukan data dengan mengunggah file",
	}
	SuccessUpdateCycle = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Cycle",
		TitleID:    "Anda Telah Berhasil Memperbarui Siklus",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeleteCycle = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Cycle",
		TitleID:    "Anda Telah Berhasil Menghapus Siklus",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessEndCycle = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Ended The Cycle",
		TitleID:    "Anda Telah Berhasil Mengakhiri Siklus",
		BodyEN:     "",
		BodyID:     "",
	}

	// Farms
	SuccessCreateFarm = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You Have Successfully Created A New Farm",
		TitleID:    "Anda Telah Berhasil Mendaftarkan Tambak Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdateFarm = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Farm",
		TitleID:    "Anda Telah Berhasil Memperbarui Tambak",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeleteFarm = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Farm",
		TitleID:    "Anda Telah Berhasil Menghapus Tambak",
		BodyEN:     "",
		BodyID:     "",
	}

	// Modules
	SuccessCreateModule = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Created A New Module",
		TitleID:    "Anda Telah Berhasil Membuat Modul Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdateModule = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Module",
		TitleID:    "Anda Telah Berhasil Memperbarui Modul",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeleteModule = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Module",
		TitleID:    "Anda Telah Berhasil Menghapus Modul",
		BodyEN:     "",
		BodyID:     "",
	}

	// Ponds
	SuccessCreatePond = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Created A New Pond",
		TitleID:    "Anda Telah Berhasil Membuat Kolam Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdatePond = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Pond",
		TitleID:    "Anda Telah Berhasil Memperbarui Kolam",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeletePond = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Pond",
		TitleID:    "Anda Telah Berhasil Menghapus Kolam",
		BodyEN:     "",
		BodyID:     "",
	}

	//Pond Diseases
	SuccessCreatePondDisease = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Created A New Pond Disease",
		TitleID:    "Anda Telah Berhasil Membuat Penyakit Tambak Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdatePondDisease = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Pond Disease",
		TitleID:    "Anda Telah Berhasil Memperbarui Penyakit Tambak",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeletePondDisease = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Pond Disease",
		TitleID:    "Anda Telah Berhasil Menghapus Penyakit Tambak",
		BodyEN:     "",
		BodyID:     "",
	}

	// Tickets
	SuccessCreateTicket = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Created A New Ticket",
		TitleID:    "Anda Telah Berhasil Membuat Tiket Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdateTicket = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The Ticket",
		TitleID:    "Anda Telah Berhasil Memperbarui Tiket",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeleteTicket = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The Ticket",
		TitleID:    "Anda Telah Berhasil Menghapus Tiket",
		BodyEN:     "",
		BodyID:     "",
	}

	// User Codes
	SuccessCreateUser = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Created A New User",
		TitleID:    "Anda Telah Berhasil Membuat Pengguna Baru",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessUpdateUser = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Updated The User",
		TitleID:    "Anda Telah Berhasil Memperbarui Pengguna",
		BodyEN:     "",
		BodyID:     "",
	}
	SuccessDeleteUser = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "You Have Successfully Deleted The User",
		TitleID:    "Anda Telah Berhasil Menghapus Pengguna",
		BodyEN:     "",
		BodyID:     "",
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
	SuccesResetPassword = Message{
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

	// Faq Token
	SuccesCreateFaqFeedback = Message{
		StatusCode: http.StatusOK,
		TitleEN:    "Create feedback success.",
		TitleID:    "Pembuatan saran berhasil.",
		BodyEN:     "",
		BodyID:     "",
	}

	// Bulk Upload Order
	SuccessBulkUploadOrder = Message{
		StatusCode: http.StatusCreated,
		TitleEN:    "You have Successfully Upload Bulk Order",
		TitleID:    "Anda Telah Berhasil Mengunggah Bulk Order",
		BodyEN:     "",
		BodyID:     "",
	}
)
