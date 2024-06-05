package codes

import (
	"math"

	"github.com/downsized-devs/sdk-go/language"
	"github.com/downsized-devs/sdk-go/operator"
)

type Code uint32

type AppMessage map[Code]Message

type DisplayMessage struct {
	StatusCode int    `json:"statusCode"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

// * Important: For new codes, please add them to the bottom of corresponding list to avoid changing existing codes and potentially breaking existing flow

const NoCode Code = math.MaxUint32

const (
	// Success code
	CodeSuccess = Code(iota + 10)
	CodeAccepted
)

const (
	// cycle file upload codes
	CodeSuccessUpload = Code(iota + 50)
	CodeSuccessUploadSowExcel
	CodeSuccessUploadDailyMonitoring
	CodeSuccessUploadSampling
	CodeSuccessUploadPartialHarvest
	CodeSuccessUploadTotalHarvest
	CodeSuccessUploadProductionPlan
)

const (
	// cycle codes
	CodeSuccessCreateCycle = Code(iota + 100)
	CodeSuccessUpdateCycle
	CodeSuccessDeleteCycle
	CodeSuccessEndCycle
)

const (
	// farms codes
	CodeSuccessCreateFarm = Code(iota + 150)
	CodeSuccessUpdateFarm
	CodeSuccessDeleteFarm
)

const (
	// module codes
	CodeSuccessCreateModule = Code(iota + 200)
	CodeSuccessUpdateModule
	CodeSuccessDeleteModule
)

const (
	// pond codes
	CodeSuccessCreatePond = Code(iota + 250)
	CodeSuccessUpdatePond
	CodeSuccessDeletePond
)

const (
	// ticket codes
	CodeSuccessCreateTicket = Code(iota + 300)
	CodeSuccessUpdateTicket
	CodeSuccessDeleteTicket
)

const (
	// user codes
	CodeSuccessCreateUser = Code(iota + 350)
	CodeSuccessUpdateUser
	CodeSuccessDeleteUser
	CodeSuccessReadUser
	CodeSuccessUploadImage
	CodeSuccessAssignRole
)

const (
	// verif token codes
	CodeSuccessResetPassword = Code(iota + 400)
	CodeSuccessInputResetPassword
)

const (
	// faq codes
	CodeSuccessCreateFaqFeedback = Code(iota + 450)
	CodeFailInsertFeedbackWeeklyLimit
)

const (
	// pond disease codes
	CodeSuccessCreatePondDisease = Code(iota + 500)
	CodeSuccessUpdatePondDisease
	CodeSuccessDeletePondDisease
)

const (
	// bulk upload order
	CodeSuccessUploadBulkOrder = Code(iota + 550)
)

const (
	// common errors
	CodeInvalidValue = Code(iota + 1000)
	CodeContextDeadlineExceeded
	CodeContextCanceled
	CodeInternalServerError
	CodeServerUnavailable
	CodeNotImplemented
	CodeBadRequest
	CodeNotFound
	CodeConflict
	CodeUnauthorized
	CodeTooManyRequest
	CodeMarshal
	CodeUnmarshal
)

const (
	// SQL errors
	CodeSQL = Code(iota + 1300)
	CodeSQLInit
	CodeSQLBuilder
	CodeSQLTxBegin
	CodeSQLTxCommit
	CodeSQLTxRollback
	CodeSQLTxExec
	CodeSQLPrepareStmt
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLRecordDoesNotExist
	CodeSQLUniqueConstraint
	CodeSQLConflict
	CodeSQLNoRowsAffected
)

const (
	// third party/client errors
	CodeClient = Code(iota + 1500)
	CodeClientMarshal
	CodeClientUnmarshal
	CodeClientErrorOnRequest
	CodeClientErrorOnReadBody
)

const (
	// general file I/O errors
	CodeFile = Code(iota + 1600)
	CodeFilePathOpenFailed
	CodeFileTooBig
)

const (
	// auth errors
	CodeAuth = Code(iota + 1700)
	CodeAuthRefreshTokenExpired
	CodeAuthAccessTokenExpired
	CodeAuthFailure
	CodeAuthInvalidToken
	CodeForbidden
	CodeAuthRevokeRefreshTokenFailed
)

const (
	// JSON encoding errors
	CodeJSONSchema = Code(iota + 1900)
	CodeJSONSchemaInvalid
	CodeJSONSchemaNotFound
	CodeJSONStructInvalid
	CodeJSONRawInvalid
	CodeJSONValidationError
	CodeJSONMarshalError
	CodeJSONUnmarshalError
)

const (
	// XML encoding errors
	CodeXMLSchema = Code(iota + 1950)
	CodeXMLMarshalError
	CodeXMLUnmarshalError
)

const (
	// Excel Errors
	CodeExcelFailedParsing = Code(iota + 2000)
	CodeExcelInvalidType
	CodeExcelFailedToSaveFile
)

const (
	// Storage Errors
	CodeStorage = Code(iota + 2100)
	CodeStorageS3Upload
	CodeStorageS3Download
	CodeStorageS3Delete
)

const (
	// data conversion error
	CodeConvert = Code(iota + 2200)
	CodeConvertTime
)

const (
	// File Upload Errors
	CodeFileUploadWrongFile = Code(iota + 3000)
	CodeFileUploadSowWrongFile
	CodeFileUploadDailyWrongFile
	CodeFileUploadHarvestWrongFile
	CodeFileUploadSamplingWrongFile
	CodeFileUploadTreatmentWrongFile
	CodeFileUploadProductionPlanWrongFile
	CodeFileUploadSowBackdate
	CodeFileUploadDailyBackdate
	CodeFileUploadHarvestBackdate
	CodeFileUploadSamplingBackdate
	CodeFileUploadTreatmentBackdate
	CodeFileUploadBackdate
	CodeFileProductionPlanNotUploaded
	CodeFileSowNotUploaded
	CodeFileUploadInvalidExtension
	CodeFileUploadSowForwardDate
	CodeFileUploadDailyForwardDate
	CodeFileUploadHarvestForwardDate
	CodeFileUploadSamplingForwardDate
	CodeFileUploadTreatmentForwardDate
	CodeFileUploadForwardDate
	CodeFileInvalidateDateColumn
	CodeFileInvalidPondIsNotRegister
	CodeFileUploadBulkOrderWrongFile
	CodeFileUploadWrongDateFormat
	CodeFileUploadUnregisteredPond
	CodeFileUploadIncorrectDataFormat
	CodeFileUploadDocDopMiscalculation
)

const (
	// Farm Errors
	CodeFarmError = Code(iota + 3100)
	CodeFarmNotFound
	CodeFarmNoActivePond
)

const (
	// Module Errors
	CodeModuleError = Code(iota + 3200)
	CodeModuleNotFound
)

const (
	// Pond Errors
	CodePondError = Code(iota + 3300)
	CodePondNotFound
)

const (
	// Cycle Errors
	CodeCycleError = Code(iota + 3400)
	CodeCycleNotFound
	CodeCycleEndFailedNoActivePond
	CodeCycleEndFailedNoSelectedPond
	CodeCycleStartFailedNameTooLong
	CodeCycleStartFailedOngoingCycle
)

const (
	// Metrics Errors
	CodeMetricsError = Code(iota + 3500)
	CodeMetricsProductivityNoDaily
)

const (
	// SES Errors
	CodeSendEmailFailed = Code(iota + 3700)
)

const (
	// Reset Password Error
	CodePasswordDoesNotMatch = Code(iota + 3800)
	CodeFailedResetPassword
	CodeResetPasswordTokenExpired
	CodeEmptyEmail
	CodeInvalidEmail
	CodeSameCurrentPassword
	CodePasswordIsNotFilled
	CodeResetPasswordTokenInvalid
	CodePasswordIsWeak
)

const (
	// Redis Cache Error
	CodeRedisGet = Code(iota + 3900)
	CodeRedisSetex
	CodeFailedLock
	CodeFailedReleaseLock
	CodeLockExist
	CodeCacheMarshal
	CodeCacheUnmarshal
	CodeCacheGetSimpleKey
	CodeCacheSetSimpleKey
	CodeCacheDeleteSimpleKey
	CodeCacheGetHashKey
	CodeCacheSetHashKey
	CodeCacheDeleteHashKey
	CodeCacheSetExpiration
	CodeCacheDecode
	CodeCacheLockNotAcquired
	CodeCacheInvalidCastType
	CodeCacheNotFound
)

const (
	CodeErrorHttpNewRequest = Code(iota + 4000)
	CodeErrorHttpDo
	CodeErrorIoutilReadAll
	CodeHttpUnmarshal
	CodeHttpMarshal
)

const (
	// Code Feature Flag Retriever Errors
	CodeFeatureFlagRetrieverFailed = Code(iota + 4100)
)

const (
	// Code Go-html template errors
	CodeExecuteTemplateFailed = Code(iota + 4200)
	CodeConvertMJMLToHTMLFailed
	CodePDFToJSONFailed
	CodePDFGeneratorFromJSONFailed
	CodeGeneratePDFFailed
	CodeParseHTMlTemplateFailed
)

const (
	// Assigining Role Error
	CodeErrorUserAssignRole = Code(iota + 4300)
	CodeErrorWrongAssigneeRole
	CodeErrorWrongUser
)

const (
	// Deleting User Error
	CodeErrorDeleteUser = Code(iota + 4400)
	CodeErrorWrongUserToDelete
)

const (
	// Farm Recommendation
	CodeErrorHarvestPlanMinBigger = Code(iota + 4500)
	CodeErrorHarvestPlanDayIntervalLower
	CodeErrorHarvestPlanMinMaxInterval
	CodeErrorABWSizeError
)

const (
	// Slack Alert Error
	CodeErrorSlackAlert = Code(iota + 4600)
)

const (
	// Security Error
	CodeErrorSecurityInvalidChipper = Code(iota + 4700)
)

const (
	// Timelib error
	CodeErrorTimelib = Code(iota + 4800)
)

const (
	// Pond Disease Errors
	CodePondDiseaseError = Code(iota + 4900)
	CodePondDiseaseNotFound
)

const (
	// Translator Error
	CodeTranslatorError = Code(iota + 5000)
)

const (
	// Image Upload Error
	CodeImageUploadSizeTooBig = Code(iota + 5100)
)

// Error messages only
var ErrorMessages = AppMessage{
	CodeInvalidValue:            ErrMsgBadRequest,
	CodeContextDeadlineExceeded: ErrMsgContextTimeout,
	CodeContextCanceled:         ErrMsgContextTimeout,
	CodeInternalServerError:     ErrMsgInternalServerError,
	CodeServerUnavailable:       ErrMsgServiceUnavailable,
	CodeNotImplemented:          ErrMsgNotImplemented,
	CodeBadRequest:              ErrMsgBadRequest,
	CodeNotFound:                ErrMsgNotFound,
	CodeConflict:                ErrMsgConflict,
	CodeUnauthorized:            ErrMsgUnauthorized,
	CodeTooManyRequest:          ErrMsgTooManyRequest,
	CodeMarshal:                 ErrMsgBadRequest,
	CodeUnmarshal:               ErrMsgBadRequest,
	CodeJSONMarshalError:        ErrMsgBadRequest,
	CodeJSONUnmarshalError:      ErrMsgBadRequest,
	CodeJSONValidationError:     ErrMsgBadRequest,

	CodeSQL:                   ErrMsgInternalServerError,
	CodeSQLInit:               ErrMsgInternalServerError,
	CodeSQLBuilder:            ErrMsgInternalServerError,
	CodeSQLTxBegin:            ErrMsgInternalServerError,
	CodeSQLTxCommit:           ErrMsgInternalServerError,
	CodeSQLTxRollback:         ErrMsgInternalServerError,
	CodeSQLTxExec:             ErrMsgInternalServerError,
	CodeSQLPrepareStmt:        ErrMsgInternalServerError,
	CodeSQLRead:               ErrMsgInternalServerError,
	CodeSQLRowScan:            ErrMsgInternalServerError,
	CodeSQLRecordDoesNotExist: ErrMsgNotFound,
	CodeSQLUniqueConstraint:   ErrMsgConflict,
	CodeSQLConflict:           ErrMsgConflict,
	CodeSQLNoRowsAffected:     ErrMsgInternalServerError,

	CodeClient:                ErrMsgInternalServerError,
	CodeClientMarshal:         ErrMsgInternalServerError,
	CodeClientUnmarshal:       ErrMsgInternalServerError,
	CodeClientErrorOnRequest:  ErrMsgBadRequest,
	CodeClientErrorOnReadBody: ErrMsgBadRequest,

	CodeAuth:                         ErrMsgUnauthorized,
	CodeAuthRefreshTokenExpired:      ErrMsgRefreshTokenExpired,
	CodeAuthAccessTokenExpired:       ErrMsgAccessTokenExpired,
	CodeAuthFailure:                  ErrMsgUnauthorized,
	CodeAuthInvalidToken:             ErrMsgInvalidToken,
	CodeForbidden:                    ErrMsgForbidden,
	CodeAuthRevokeRefreshTokenFailed: ErrMsgRevokeRefreshTokenFailed,

	CodeExcelFailedParsing:    ErrMsgBadRequest,
	CodeExcelInvalidType:      ErrMsgFileUploadWrongFile,
	CodeExcelFailedToSaveFile: ErrMsgInternalServerError,

	CodeStorageS3Upload: ErrMsgBadRequest,

	CodeConvert:     ErrMsgInternalServerError,
	CodeConvertTime: ErrMsgInternalServerError,

	CodeFileUploadSowBackdate:          ErrMsgFileUploadSowBackdate,
	CodeFileUploadDailyBackdate:        ErrMsgFileUploadDailyBackdate,
	CodeFileUploadHarvestBackdate:      ErrMsgFileUploadHarvestBackdate,
	CodeFileUploadSamplingBackdate:     ErrMsgFileUploadSamplingBackdate,
	CodeFileUploadTreatmentBackdate:    ErrMsgFileUploadTreatmentBackdate,
	CodeFileUploadBackdate:             ErrMsgFileUploadBackdate,
	CodeFileSowNotUploaded:             ErrMsgFileUploadSowNotUploaded,
	CodeFileProductionPlanNotUploaded:  ErrMsgFileUploadProductionPlanNotUploaded,
	CodeFileUploadInvalidExtension:     ErrMsgFileUploadInvalidExtension,
	CodeFileUploadSowForwardDate:       ErrMsgFileUploadSowBackdate,
	CodeFileUploadDailyForwardDate:     ErrMsgFileUploadDailyBackdate,
	CodeFileUploadHarvestForwardDate:   ErrMsgFileUploadHarvestBackdate,
	CodeFileUploadSamplingForwardDate:  ErrMsgFileUploadSamplingBackdate,
	CodeFileUploadTreatmentForwardDate: ErrMsgFileUploadTreatmentBackdate,
	CodeFileUploadForwardDate:          ErrMsgFileUploadBackdate,
	CodeFileInvalidateDateColumn:       ErrMsgFileInvalidDate,
	CodeFileInvalidPondIsNotRegister:   ErrMsgFileInvalidPondIsNotRegister,
	CodeFileUploadWrongDateFormat:      ErrMsgFileUploadWrongDate,
	CodeFileUploadUnregisteredPond:     ErrMsgFileUploadUnregisteredPond,
	CodeFileUploadIncorrectDataFormat:  ErrMsgFileUploadIncorrectDataFormat,

	CodeFileUploadWrongFile:               ErrMsgFileUploadWrongFile,
	CodeFileUploadSowWrongFile:            ErrMsgFileUploadSowWrongFile,
	CodeFileUploadDailyWrongFile:          ErrMsgFileUploadDailyWrongFile,
	CodeFileUploadHarvestWrongFile:        ErrMsgFileUploadHarvestWrongFile,
	CodeFileUploadSamplingWrongFile:       ErrMsgFileUploadSamplingWrongFile,
	CodeFileUploadTreatmentWrongFile:      ErrMsgFileUploadTreatmentWrongFile,
	CodeFileUploadProductionPlanWrongFile: ErrMsgFileUploadProductionPlanWrongFile,
	CodeFileUploadBulkOrderWrongFile:      ErrMsgFileUploadBulkOrderWrongFile,

	CodeCycleEndFailedNoActivePond:   ErrMsgCycleEndFailedNoActivePond,
	CodeCycleEndFailedNoSelectedPond: ErrMsgCycleEndFailedNoSelectedPond,
	CodeCycleStartFailedNameTooLong:  ErrMsgCodeCycleStartFailedNameTooLong,
	CodeCycleNotFound:                ErrMsgCycleNotFound,
	CodeCycleStartFailedOngoingCycle: ErrMsgCycleStartFailedOngoingCycle,

	CodeMetricsProductivityNoDaily: ErrMsgMetricsProductivityNoDaily,

	CodeFailedResetPassword:       ErrMsgResetPassword,
	CodePasswordDoesNotMatch:      ErrMsgPasswordDoesNotMatch,
	CodeResetPasswordTokenExpired: ErrMsgResetTokenExpired,
	CodeEmptyEmail:                ErrMsgEmptyEmail,
	CodeInvalidEmail:              ErrMsgInvalidEmail,
	CodeSameCurrentPassword:       ErrMsgSameCurrentPassword,
	CodePasswordIsNotFilled:       ErrMsgPasswordIsNotFilled,
	CodeResetPasswordTokenInvalid: ErrMsgResetTokenInvalid,
	CodePasswordIsWeak:            ErrMsgPasswordIsWeak,

	CodeLockExist:            ErrMsgLockExist,
	CodeRedisGet:             ErrMsgInternalServerError,
	CodeRedisSetex:           ErrMsgInternalServerError,
	CodeFailedLock:           ErrMsgInternalServerError,
	CodeFailedReleaseLock:    ErrMsgInternalServerError,
	CodeCacheMarshal:         ErrMsgInternalServerError,
	CodeCacheUnmarshal:       ErrMsgInternalServerError,
	CodeCacheGetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheSetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheDeleteSimpleKey: ErrMsgInternalServerError,
	CodeCacheGetHashKey:      ErrMsgInternalServerError,
	CodeCacheSetHashKey:      ErrMsgInternalServerError,
	CodeCacheDeleteHashKey:   ErrMsgInternalServerError,
	CodeCacheSetExpiration:   ErrMsgInternalServerError,
	CodeCacheDecode:          ErrMsgInternalServerError,
	CodeCacheLockNotAcquired: ErrMsgInternalServerError,
	CodeCacheInvalidCastType: ErrMsgInternalServerError,
	CodeCacheNotFound:        ErrMsgInternalServerError,

	CodeFailInsertFeedbackWeeklyLimit: ErrMsgInsertFeedbackWeeklyLimit,

	CodeErrorHttpNewRequest: ErrMsgInternalServerError,
	CodeErrorHttpDo:         ErrMsgInternalServerError,
	CodeErrorIoutilReadAll:  ErrMsgInternalServerError,
	CodeHttpMarshal:         ErrMsgInternalServerError,
	CodeHttpUnmarshal:       ErrMsgInternalServerError,

	CodeFeatureFlagRetrieverFailed: ErrMsgInternalServerError,

	// GO-HTML Template
	CodeExecuteTemplateFailed:      ErrMsgExecuteTemplateFailed,
	CodeConvertMJMLToHTMLFailed:    ErrMsgConvertMJMLToHTMLFailed,
	CodePDFToJSONFailed:            ErrMsgConvertGoTemplateToPDFFailed,
	CodePDFGeneratorFromJSONFailed: ErrMsgConvertGoTemplateToPDFFailed,
	CodeGeneratePDFFailed:          ErrMsgGenerateToPDFFailed,
	CodeParseHTMlTemplateFailed:    ErrMsgCodeParseHTMlTemplateFailed,

	// User Error
	CodeErrorUserAssignRole:    ErrMsgAssignRole,
	CodeErrorWrongAssigneeRole: ErrMsgWrongAssigneeRole,
	CodeErrorWrongUser:         ErrMsgWrongUser,
	CodeErrorDeleteUser:        ErrMsgDeleteUser,
	CodeErrorWrongUserToDelete: ErrMsgDeleteUser,

	// Farm Recommendation Error
	CodeErrorHarvestPlanMinBigger:        ErrMsgHarvestPlanMinBigger,
	CodeErrorHarvestPlanDayIntervalLower: ErrMsgHarvestPlanIntervalLower,
	CodeErrorHarvestPlanMinMaxInterval:   ErrMsgHarvestPlanMinMaxInterval,

	// Slack Error
	CodeErrorSlackAlert: ErrMsgSlackAlert,

	// Timelib Error
	CodeErrorTimelib: ErrMsgTimelib,

	CodeFarmNoActivePond: ErrMsgFarmNoActivePond,

	CodeModuleError: ErrMsgModule,

	CodePondError: ErrMsgPond,

	CodePondDiseaseError: ErrMsgPondDisease,

	// File I/O error
	CodeFile:               ErrMsgInternalServerError,
	CodeFilePathOpenFailed: ErrMsgInternalServerError,
	CodeFileTooBig:         ErrMsgContentTooLarge,

	// Code Translator
	CodeTranslatorError: ErrMsgTranslatorlib,

	// Error Image
	CodeImageUploadSizeTooBig: ErrMsgImageTooBig,
}

// Successful messages only
var ApplicationMessages = AppMessage{
	// Other
	CodeAccepted: SuccessAccepted,

	// Files upload
	CodeSuccessUploadSowExcel:        SuccessUploadSow,
	CodeSuccessUploadDailyMonitoring: SuccessUploadDailyMonitoring,
	CodeSuccessUploadSampling:        SuccessUploadSampling,
	CodeSuccessUploadPartialHarvest:  SuccessUploadPartialHarvest,
	CodeSuccessUploadTotalHarvest:    SuccessUploadTotalHarvest,
	CodeSuccessUploadProductionPlan:  SuccessUploadProductionPlan,

	// Cycles
	CodeSuccessCreateCycle: SuccessCreateCycle,
	CodeSuccessUpdateCycle: SuccessUpdateCycle,
	CodeSuccessDeleteCycle: SuccessDeleteCycle,
	CodeSuccessEndCycle:    SuccessEndCycle,

	// Farm
	CodeSuccessCreateFarm: SuccessCreateFarm,
	CodeSuccessUpdateFarm: SuccessUpdateFarm,
	CodeSuccessDeleteFarm: SuccessDeleteFarm,

	// Module
	CodeSuccessCreateModule: SuccessCreateModule,
	CodeSuccessUpdateModule: SuccessUpdateModule,
	CodeSuccessDeleteModule: SuccessDeleteModule,

	// Pond
	CodeSuccessCreatePond: SuccessCreatePond,
	CodeSuccessUpdatePond: SuccessUpdatePond,
	CodeSuccessDeletePond: SuccessDeletePond,

	// Ticket
	CodeSuccessCreateTicket: SuccessCreateTicket,
	CodeSuccessUpdateTicket: SuccessUpdateTicket,
	CodeSuccessDeleteTicket: SuccessDeleteTicket,

	// User
	CodeSuccessCreateUser: SuccessCreateUser,
	CodeSuccessUpdateUser: SuccessUpdateUser,
	CodeSuccessDeleteUser: SuccessDeleteUser,

	// Verif token
	CodeSuccessResetPassword:      SuccesResetPassword,
	CodeSuccessInputResetPassword: SuccessInputResetPassword,

	// Faq
	CodeSuccessCreateFaqFeedback: SuccesCreateFaqFeedback,

	// Pond Disease
	CodeSuccessCreatePondDisease: SuccessCreatePondDisease,
	CodeSuccessUpdatePondDisease: SuccessUpdatePondDisease,
	CodeSuccessDeletePondDisease: SuccessDeletePondDisease,

	// Bulk Upload Order
	CodeSuccessUploadBulkOrder: SuccessBulkUploadOrder,
}

func Compile(code Code, lang string) DisplayMessage {
	if appMsg, ok := ApplicationMessages[code]; ok {
		return DisplayMessage{
			StatusCode: appMsg.StatusCode,
			Title:      operator.Ternary(lang == language.Indonesian, appMsg.TitleID, appMsg.TitleEN),
			Body:       operator.Ternary(lang == language.Indonesian, appMsg.BodyID, appMsg.BodyEN),
		}
	}

	return DisplayMessage{
		StatusCode: SuccessDefault.StatusCode,
		Title:      operator.Ternary(lang == language.Indonesian, SuccessDefault.TitleID, SuccessDefault.TitleEN),
		Body:       operator.Ternary(lang == language.Indonesian, SuccessDefault.BodyID, SuccessDefault.BodyEN),
	}
}
