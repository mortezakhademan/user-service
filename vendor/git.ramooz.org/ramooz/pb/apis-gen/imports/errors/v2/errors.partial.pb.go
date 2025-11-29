package errors

//
//import (
//	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
//)
//
//func (x *GetErrorsRequest) GetErrorResponse(clientsErr ...map[int64]string) (*GetErrorsResponse, error) {
//	resp := &GetErrorsResponse{}
//
//	if len(x.LanguageCode) == 0 {
//		x.LanguageCode = "en"
//	}
//
//	resp.Language = x.LanguageCode
//	resp.ErrorMessage = getErrorsByLanguageCode(x.LanguageCode, clientsErr...)
//
//	return resp, nil
//}
//
//func getErrorsByLanguageCode(languageCode string, clientsErr ...map[int64]string) map[int64]string {
//	errs := make(map[int64]string)
//
//	for lang, errors := range componentsError.I18nTranslates {
//		if len(languageCode) == 0 {
//			languageCode = "en"
//		}
//
//		if lang == languageCode {
//			for code, msg := range errors {
//				errs[int64(code)] = msg
//			}
//		}
//	}
//
//	if len(clientsErr) != 0 {
//		for _, err := range clientsErr {
//			for code, msg := range err {
//				errs[code] = msg
//			}
//		}
//	}
//
//	return errs
//}
