package acl

import componentsError "git.ramooz.org/ramooz/golang-components/error-handler"

var ErrorMessages = map[string]map[int]string{
	componentsError.LANG_FA: {
		ERROR_AUTHORIZATION_NOT_FOUND:       "توکن یافت نشد",
		ERROR_MISSING_BEARER_IN_HEADER:      "در هدر یافت نشد Bearer",
		ERROR_NO_HEADER_IN_REQUEST:          "هدری در درخواست وجود ندارد",
		ERROR_PRIVATE_SECRET_KEY_IS_EMPTY:   "کد امنیتی توکن خالی است",
		ERROR_CREDENTIAL_NOT_VALID:          "اطلاعات سطح دسترسی معتبر نمی باشد",
		ERROR_ACL_DATA_NOT_VALID:            "اطلاعات اعتبارسنجی معتبر نمی باشد",
		ERROR_API_KEY_FUNC_INVALID:          "تابع get api key info تعریف نشده است",
		ERROR_SECRET_KEY_IS_EMPTY:           "کلید امنیتی توکن خالی است",
		ERROR_SERVICE_CODE_HAS_BEEN_NOT_SET: "کد سرویس تنظیم نشده است",
		ERROR_JWT_TOKEN_IS_INVALID:          "توکن نامعتبر می باشد",
		ERR_EXTRACT_ACL_FROM_CONTEXT:        "نمی توانیم اطلاعات acl را از کانتکست استخراج کنیم",
	},
	componentsError.LANG_EN: {
		ERROR_AUTHORIZATION_NOT_FOUND:       "not found jwt token in header",
		ERROR_MISSING_BEARER_IN_HEADER:      "missing Bearer prefix in Authorization header",
		ERROR_NO_HEADER_IN_REQUEST:          "no headers in request",
		ERROR_CREDENTIAL_NOT_VALID:          "credential not valid",
		ERROR_ACL_DATA_NOT_VALID:            "acl data not valid",
		ERROR_API_KEY_FUNC_INVALID:          "api key function option is empty",
		ERROR_SECRET_KEY_IS_EMPTY:           "secret key is empty",
		ERROR_SERVICE_CODE_HAS_BEEN_NOT_SET: "service code has been not set",
		ERROR_PRIVATE_SECRET_KEY_IS_EMPTY:   "private secret key is empty",
		ERROR_JWT_TOKEN_IS_INVALID:          "jwt token is invalid",
		ERR_EXTRACT_ACL_FROM_CONTEXT:        "can't extract acl from context",
		ERR_METHOD_NOT_IMPLEMENTED:          "method not implemented",
	},
	componentsError.LANG_AR: {
		ERROR_AUTHORIZATION_NOT_FOUND:       "يرجى تسجيل الدخول في البداية",
		ERROR_MISSING_BEARER_IN_HEADER:      "يرجى تسجيل الدخول في البداية",
		ERROR_NO_HEADER_IN_REQUEST:          "يرجى تسجيل الدخول في البداية",
		ERROR_CREDENTIAL_NOT_VALID:          "يرجى تسجيل الدخول في البداية",
		ERROR_ACL_DATA_NOT_VALID:            "يرجى تسجيل الدخول في البداية",
		ERROR_API_KEY_FUNC_INVALID:          "خيار وظيفة المفتاح فارغ",
		ERROR_SECRET_KEY_IS_EMPTY:           "رقم السري الخاص فارغ",
		ERROR_SERVICE_CODE_HAS_BEEN_NOT_SET: "لم يتم تعيين رمز الخدمة",
		ERROR_PRIVATE_SECRET_KEY_IS_EMPTY:   "رقم السري الخاص فارغ",
		ERROR_JWT_TOKEN_IS_INVALID:          "رمز غير صالح",
		ERR_EXTRACT_ACL_FROM_CONTEXT:        "لا يمكن استخراج قائمة التحكم من السياق",
		ERR_METHOD_NOT_IMPLEMENTED:          "لم يتم تنفيذ الطريقة",
	},
}

const (
	ERROR_AUTHORIZATION_NOT_FOUND       = 401100
	ERROR_MISSING_BEARER_IN_HEADER      = 401101
	ERROR_NO_HEADER_IN_REQUEST          = 401102
	ERROR_CREDENTIAL_NOT_VALID          = 401103
	ERROR_PRIVATE_SECRET_KEY_IS_EMPTY   = 400101
	ERROR_ACL_DATA_NOT_VALID            = 400102
	ERROR_API_KEY_FUNC_INVALID          = 400103
	ERROR_SECRET_KEY_IS_EMPTY           = 400104
	ERROR_SERVICE_CODE_HAS_BEEN_NOT_SET = 400105
	ERROR_JWT_TOKEN_IS_INVALID          = 400106
	ERR_EXTRACT_ACL_FROM_CONTEXT        = 500101
	ERR_METHOD_NOT_IMPLEMENTED          = 500102
)
