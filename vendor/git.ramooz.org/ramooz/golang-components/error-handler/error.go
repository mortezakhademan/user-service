package error

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "git.ramooz.org/ramooz/pb/apis-gen/imports/errors/v2"
	integrations "google.golang.org/genproto/googleapis/cloud/integrations/v1alpha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"runtime"
	"strconv"
)

type Error struct {
	ServiceCode int
	ErrorCode   int
	Message     string
	Details     []string
	TraceLines  []string
}

func SetServiceCode(serviceCode int) {
	errorServiceCode = serviceCode
	addServiceCodeToMainErrors()

}

func SetDefaultLanguage(language string) {
	DefaultLanguage = language
}
func New(errorCode int, details []string, params ...interface{}) *Error {
	errorCode, _ = strconv.Atoi(strconv.Itoa(errorServiceCode) + strconv.Itoa(errorCode))
	e := &Error{
		ServiceCode: errorServiceCode,
		ErrorCode:   errorCode,
		Message:     GetErrorMessage(errorCode, params...),
		Details:     details,
	}
	if withTraceDetail {
		e.TraceLines = getLastTraceLines(5)
	}
	return e
}
func getLastTraceLines(count int) []string {
	var traceLines []string

	// Use runtime.Callers to retrieve multiple stack frames
	pcs := make([]uintptr, count)
	n := runtime.Callers(3, pcs) // Start from the calling function
	frames := runtime.CallersFrames(pcs[:n])

	// Iterate through the frames and build trace lines
	for {
		frame, more := frames.Next()
		traceLines = append(traceLines, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return traceLines
}

// Deprecated: NewError
func NewError(errorCode int, params ...interface{}) *Error {
	responseErrorCode, _ := strconv.Atoi(strconv.Itoa(errorServiceCode) + strconv.Itoa(errorCode))
	e := &Error{
		ErrorCode: responseErrorCode,
		Message:   GetErrorMessage(errorCode, params...),
	}
	if withTraceDetail {
		e.TraceLines = getLastTraceLines(5)
	}
	return e
}

func GetErrorMessageByLanguage(errorCode int, language string, params ...interface{}) string {
	if _, ok := I18nTranslates[language][errorCode]; ok || I18nTranslates[language][errorCode] != "" {
		return fmt.Sprintf(I18nTranslates[language][errorCode], params...)
	}
	if _, ok := I18nTranslates[DefaultLanguage][errorCode]; ok {
		return fmt.Sprintf(I18nTranslates[DefaultLanguage][errorCode], params...)
	}
	return ""
}

func GetHttpHeaderCode(statusCode int) int {
	if statusCode == 0 {
		return http.StatusOK
	}
	errorCodeStr := strconv.Itoa(statusCode)
	switch len(errorCodeStr) {
	case 5, 6:
		errorCodeStr = errorCodeStr[:3]
	case 7:
		errorCodeStr = errorCodeStr[1:4]
	case 8:
		errorCodeStr = errorCodeStr[2:5]
	case 9:
		errorCodeStr = errorCodeStr[3:6]
	}
	code, _ := strconv.Atoi(errorCodeStr)
	return code
}

func (e *Error) GetSplitCode() (serviceCode int, httpCode int, customCode int) {
	errorCodeStr := strconv.Itoa(e.ErrorCode)
	switch len(errorCodeStr) {
	case 5, 6:
		httpCode, _ = strconv.Atoi(errorCodeStr[:3])
		customCode, _ = strconv.Atoi(errorCodeStr[3:])
	case 7:
		serviceCode, _ = strconv.Atoi(errorCodeStr[:1])
		httpCode, _ = strconv.Atoi(errorCodeStr[1:4])
		customCode, _ = strconv.Atoi(errorCodeStr[4:])
	case 8:
		serviceCode, _ = strconv.Atoi(errorCodeStr[:2])
		httpCode, _ = strconv.Atoi(errorCodeStr[2:5])
		customCode, _ = strconv.Atoi(errorCodeStr[5:])
	}
	return
}
func (e *Error) GetErrorCodeWithoutServiceCode() int {
	errorCodeStr := strconv.Itoa(e.ErrorCode)
	code := e.ErrorCode
	switch len(errorCodeStr) {
	case 7:
		code, _ = strconv.Atoi(errorCodeStr[1:])
	case 8:
		code, _ = strconv.Atoi(errorCodeStr[2:])
	}
	return code
}
func (e *Error) GRPCStatus() *status.Status {
	var code codes.Code
	code = codes.Code(uint32(e.ErrorCode))
	statusObj := status.New(code, e.GetMessage())
	for _, detail := range e.Details {
		statusObj, _ = statusObj.WithDetails(&integrations.ErrorDetail{ErrorMessage: detail})
	}
	if withTraceDetail {
		statusObj, _ = statusObj.WithDetails(&errors2.TraceDetail{Lines: e.TraceLines})
	}
	return statusObj
}

func (e *Error) Error() string {
	if e.Message == "" {
		e.Message = e.GetMessage()
	}
	message := "error code: (" + strconv.Itoa(e.ErrorCode) + ")\n" + e.GetMessage() + "\n"
	message += "details:\n"
	for _, detail := range e.Details {
		message += detail + "\n"
	}
	return message
}

func (e *Error) GetMessage() string {
	if e.Message != "" {
		return e.Message
	}
	errorCode := e.ErrorCode
	errorCodeStr := strconv.Itoa(errorCode)
	if errorServiceCode != 0 {
		errorCode, _ = strconv.Atoi(errorCodeStr[len(strconv.Itoa(errorServiceCode)):])
	}

	return GetErrorMessage(errorCode)
}

func GetErrorMessage(errorCode int, params ...interface{}) string {
	if message, ok := I18nTranslates[DefaultLanguage][errorCode]; ok {
		return fmt.Sprintf(message, params...)
	}

	errorCode, _ = strconv.Atoi(fmt.Sprintf("%d%d", errorServiceCode, errorCode))
	if message, ok := I18nTranslates[DefaultLanguage][errorCode]; ok {
		return fmt.Sprintf(message, params...)
	}

	return ""
}

func addServiceCodeToMainErrors() {
	newI18nTranslates := map[string]map[int]string{}
	for lang, errMessages := range I18nTranslates {
		if _, ok := newI18nTranslates[lang]; !ok {
			newI18nTranslates[lang] = map[int]string{}
		}
		for errCode, errText := range errMessages {
			errCode, _ = strconv.Atoi(fmt.Sprintf("%v%v", errorServiceCode, errCode))
			newI18nTranslates[lang][errCode] = errText
		}
	}
	I18nTranslates = newI18nTranslates
}

func AddErrorMessages(i18nErrors map[string]map[int]string) error {
	for lang, errMessages := range i18nErrors {

		if _, ok := I18nTranslates[lang]; !ok {
			I18nTranslates[lang] = map[int]string{}
		}

		for errCode, errText := range errMessages {
			errCode, _ = strconv.Atoi(fmt.Sprintf("%v%v", errorServiceCode, errCode))
			if _, ok := I18nTranslates[lang][errCode]; ok {
				return errors.New("error code already defined! error language:" + lang + " code:" + strconv.Itoa(errCode))
			} else {
				I18nTranslates[lang][errCode] = errText
			}
		}
	}
	return nil
}

func ConvertRecoverToError(r interface{}) error {
	switch x := r.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		// Fallback err (per specs, error strings should be lowercase w/o punctuation
		return errors.New(fmt.Sprint(x))
	}
}

func ExportErrorsAsJson() string {
	list, err := json.Marshal(I18nTranslates)
	if err != nil {
		panic("there is error on list of errors")
	}
	return string(list)
}
