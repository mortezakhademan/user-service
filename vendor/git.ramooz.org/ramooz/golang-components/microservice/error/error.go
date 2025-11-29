package error

import (
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
)

var ErrorMessages = map[string]map[int]string{
	componentsError.LANG_EN: {
		ERR_VALIDATE_FIELDS:     "fields validation failed",
		ErrorInvalidModel:       "can't convert model to proto model",
		ErrorInvalidProtoModel:  "can't convert proto model to model",
		ErrorHttpMethodNotFound: "http method not found in context meta data",
		ErrorModelIsNil:         "model is nil",
	},
	componentsError.LANG_AR: {
		ERR_VALIDATE_FIELDS:     "فشل التحقق من الحقول",
		ErrorInvalidModel:       "لا يمكن تحويل النموذج إلى نموذج أولي",
		ErrorInvalidProtoModel:  "لا يمكن تحويل النموذج إلى نموذج أولي",
		ErrorHttpMethodNotFound: "لم يتم العثور على طريقة اتش تي تي بي في بيانات السياق الفرعي",
		ErrorModelIsNil:         "النموذج فارغ",
	},
	componentsError.LANG_FA: {
		ERR_VALIDATE_FIELDS:     "اعتبارسنجی فیلدهاناموفق بود",
		ErrorInvalidModel:       "نمی توان مدل را به مدل پروتو تبدیل کرد",
		ErrorInvalidProtoModel:  "نمی توان مدل پروتو را به مدل تبدیل کرد",
		ErrorHttpMethodNotFound: "متد اچ تی تی پی در داده های متا کانتکست یافت نشد",
		ErrorModelIsNil:         "مدل خالی است",
	},
}

const (
	ERR_VALIDATE_FIELDS     = 400120
	ErrorInvalidModel       = 400121
	ErrorInvalidProtoModel  = 400122
	ErrorHttpMethodNotFound = 500023
	ErrorModelIsNil         = 400123
)
