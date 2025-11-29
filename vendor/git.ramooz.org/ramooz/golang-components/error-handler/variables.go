package error

// All custom error code should greater than 999 and start three numbers must use http status code like 400, 403, 500
// Custom error codes of business start from 400200 and messages are inserted to ErrorMessages by AddErrorMessages (in main function)
const (
	ERROR_HTTP_TOKEN_NOT_VALID                        = 401000
	ERROR_HTTP_FORBIDDEN                              = 403000
	ERROR_HTTP_FORBIDDEN_TO_SERVICE                   = 403001
	ERROR_HTTP_FORBIDDEN_TO_ACTION                    = 403002
	ERROR_HTTP_NOT_FOUND                              = 404000
	ERROR_HTTP_TOKEN_EXPIRED                          = 401001
	ERROR_HTTP_TOKEN_NOT_ISSUED                       = 401002
	ERROR_HTTP_TOKEN_NOT_VALID_YET                    = 401003
	ERROR_HTTP_TOKEN_INVALID_SEGMENTS                 = 401004
	ERROR_DEVICE_NOT_SELECETED                        = 401005
	ERROR_HTTP_INTERNAL_SERVER_ERROR                  = 500000
	ERROR_GLOBAL_MYSQL_ERROR                          = 500001
	ERROR_PARAMETERS_INVALID                          = 400000
	ERROR_INVALID_LOGIN                               = 400001
	ERROR_INVALID_USERNAME                            = 400002
	ERROR_INVALID_MOBILE                              = 400003
	ERROR_INVALID_EMAIL                               = 400004
	ERROR_INVALID_STATUS                              = 400006
	ERROR_INVALID_PASSWORD                            = 400005
	ERROR_INVALID_RESET_PASSWORD_CODE                 = 400007
	ERROR_RESET_PASSWORD_CODE_EXPIRED                 = 400008
	ERROR_RESET_PASSWORD_RESEND_TIME_NOT_ELAPSED      = 400009
	ERROR_DUPLICATE_USERNAME                          = 400010
	ERROR_DUPLICATE_MOBILE                            = 400011
	ERROR_DUPLICATE_EMAIL                             = 400012
	ERROR_DUPLICATE_KEY                               = 400013
	ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_DELETE_UPDATE = 400014
	ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_ADD_UPDATE    = 400015
	ERROR_SYSTEMIC_RECORD_UPDATE_DELETE               = 400016
	ERROR_USER_STATUS_IS_NOT_ACTIVE                   = 400017
	ERROR_USER_STATUS_IS_NOT_PENDING                  = 400018
	ERROR_UPLOADED_FILE_INCORRECT                     = 400021
	ERROR_UPLOADED_FILE_NOT_SUPPORTED                 = 400022
	ERROR_DOWNLOAD_FILE_INCORRECT                     = 400023
	ERROR_INVALID_TIME_RANGE                          = 400024
	ERROR_DELETE_ADMIN_USER                           = 400025
	ERROR_ADMIN_UPDATE_HIMSELF_STATUS                 = 400026
	ERROR_ADMIN_DELETE_HIMSELF_ADMIN_ROLE             = 400027
	ERROR_INVALID_TYPE                                = 400028
	ERROR_INVALID_USER_ROLE                           = 400029
	ERROR_INVALID_AMOUNT                              = 400030
	ERROR_UPDATE_ADMIN_PASSWORD                       = 400031
	ERROR_INVALID_CURRENT_PASSWORD                    = 400032
	ERROR_INVALID_NEW_PASSWORD                        = 400033
	ERROR_SYSTEMIC_SLIDER_DELETE                      = 400034
	ERROR_INVALID_TIME                                = 400035
	ERROR_USER_NOT_EXISTS                             = 400037
	ERROR_FIREBASE_TOKEN_DUPLICATED                   = 400038
	ERROR_LOGIN_INVALID                               = 400039
	ERROR_INVALID_MESSAGE_ACTION                      = 400041
	ERROR_INVALID_MESSAGE_OBJECT                      = 400042
	ERROR_DUPLICATE_STATUS                            = 400043
	ERROR_PRICE_CURRENCY_IS_INVALID                   = 400044
	ERROR_MONGO_DATABASE_VALIDATION_FAIELD            = 400045
	ErrorInvalidDBField                               = 400046
)

var (
	DefaultLanguage = LANG_EN

	errorServiceCode int

	// -->language:{code:message}<--
	I18nTranslates = map[string]map[int]string{
		LANG_EN: {
			ERROR_HTTP_TOKEN_NOT_VALID:                        "Invalid token",
			ERROR_HTTP_NOT_FOUND:                              "http not found",
			ERROR_HTTP_TOKEN_NOT_ISSUED:                       "token not issued",
			ERROR_HTTP_TOKEN_NOT_VALID_YET:                    "token isn't valid",
			ERROR_HTTP_TOKEN_INVALID_SEGMENTS:                 "token has invalid segments",
			ERROR_DEVICE_NOT_SELECETED:                        "device not selected",
			ERROR_HTTP_INTERNAL_SERVER_ERROR:                  "http internal server error",
			ERROR_HTTP_FORBIDDEN:                              "Forbidden access",
			ERROR_HTTP_FORBIDDEN_TO_SERVICE:                   "Forbidden access to %s service",
			ERROR_HTTP_FORBIDDEN_TO_ACTION:                    "Forbidden access to %s action",
			ERROR_HTTP_TOKEN_EXPIRED:                          "token expired",
			ERROR_PARAMETERS_INVALID:                          "invalid parameters",
			ERROR_GLOBAL_MYSQL_ERROR:                          "global mysql error",
			ERROR_INVALID_LOGIN:                               "invalid login data",
			ERROR_INVALID_MOBILE:                              "mobile is invalid",
			ERROR_INVALID_EMAIL:                               "email invalid",
			ERROR_INVALID_USERNAME:                            "username invalid",
			ERROR_INVALID_PASSWORD:                            "password invalid",
			ERROR_INVALID_STATUS:                              "status invalid",
			ERROR_INVALID_RESET_PASSWORD_CODE:                 "reset password code invalid",
			ERROR_RESET_PASSWORD_CODE_EXPIRED:                 "password code expired",
			ERROR_RESET_PASSWORD_RESEND_TIME_NOT_ELAPSED:      "password reset code time not elapsed",
			ERROR_DUPLICATE_USERNAME:                          "username duplicated",
			ERROR_DUPLICATE_MOBILE:                            "mobile duplicated",
			ERROR_DUPLICATE_EMAIL:                             "email duplicated",
			ERROR_DUPLICATE_KEY:                               "duplicate key",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_DELETE_UPDATE: "error in delete or update action",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_ADD_UPDATE:    "error in add new record",
			ERROR_SYSTEMIC_RECORD_UPDATE_DELETE:               "error systemic record",
			ERROR_USER_STATUS_IS_NOT_ACTIVE:                   "user isn't active",
			ERROR_USER_STATUS_IS_NOT_PENDING:                  "user in't pending",
			ERROR_UPLOADED_FILE_INCORRECT:                     "invalid file",
			ERROR_UPLOADED_FILE_NOT_SUPPORTED:                 "file not supported",
			ERROR_DOWNLOAD_FILE_INCORRECT:                     "download file invalid",
			ERROR_INVALID_TIME_RANGE:                          "invalid time range",
			ERROR_DELETE_ADMIN_USER:                           "delete admin user",
			ERROR_ADMIN_UPDATE_HIMSELF_STATUS:                 "update mimself status",
			ERROR_ADMIN_DELETE_HIMSELF_ADMIN_ROLE:             "can't delete your role",
			ERROR_INVALID_TYPE:                                "invalid type",
			ERROR_INVALID_USER_ROLE:                           "role isn't valid",
			ERROR_INVALID_AMOUNT:                              "amount isn't valid",
			ERROR_UPDATE_ADMIN_PASSWORD:                       "can't update admin password",
			ERROR_INVALID_CURRENT_PASSWORD:                    "invalid current password",
			ERROR_INVALID_NEW_PASSWORD:                        "invalid new password",
			ERROR_SYSTEMIC_SLIDER_DELETE:                      "can't delete systemic slider",
			ERROR_INVALID_TIME:                                "time invalid",
			ERROR_USER_NOT_EXISTS:                             "user not exists",
			ERROR_FIREBASE_TOKEN_DUPLICATED:                   "firebase token duplicated",
			ERROR_LOGIN_INVALID:                               "login invalid",
			ERROR_INVALID_MESSAGE_ACTION:                      "invalid message action",
			ERROR_INVALID_MESSAGE_OBJECT:                      "invalid message object",
			ERROR_DUPLICATE_STATUS:                            "status duplicated",
			ERROR_PRICE_CURRENCY_IS_INVALID:                   "invalid currency %s",
			ERROR_MONGO_DATABASE_VALIDATION_FAIELD:            "error in database",
			ErrorInvalidDBField:                               "invalid db field",
		},
		LANG_FA: {
			ERROR_HTTP_TOKEN_NOT_VALID:                        "توکن نامعتبر است",
			ERROR_HTTP_NOT_FOUND:                              "آدرس یافت نشد",
			ERROR_HTTP_TOKEN_NOT_ISSUED:                       "توکن صادر نشده است",
			ERROR_HTTP_TOKEN_NOT_VALID_YET:                    "توکن هنوز معتبر نیست",
			ERROR_HTTP_TOKEN_INVALID_SEGMENTS:                 "توکن دارای بخش‌های نامعتبر است",
			ERROR_DEVICE_NOT_SELECETED:                        "دستگاه انتخاب نشده است",
			ERROR_HTTP_INTERNAL_SERVER_ERROR:                  "خطای داخلی سرور",
			ERROR_HTTP_FORBIDDEN:                              "دسترسی ممنوع است",
			ERROR_HTTP_FORBIDDEN_TO_SERVICE:                   "دسترسی به سرویس %s ممنوع است",
			ERROR_HTTP_FORBIDDEN_TO_ACTION:                    "دسترسی به عملیات %s ممنوع است",
			ERROR_HTTP_TOKEN_EXPIRED:                          "توکن منقضی شده است",
			ERROR_PARAMETERS_INVALID:                          "پارامترها نامعتبر هستند",
			ERROR_GLOBAL_MYSQL_ERROR:                          "خطای عمومی MySQL",
			ERROR_INVALID_LOGIN:                               "اطلاعات ورود نامعتبر است",
			ERROR_INVALID_MOBILE:                              "شماره موبایل معتبر نمی‌باشد",
			ERROR_INVALID_EMAIL:                               "ایمیل نامعتبر است",
			ERROR_INVALID_USERNAME:                            "نام کاربری نامعتبر است",
			ERROR_INVALID_PASSWORD:                            "رمز عبور نامعتبر است",
			ERROR_INVALID_STATUS:                              "وضعیت نامعتبر است",
			ERROR_INVALID_RESET_PASSWORD_CODE:                 "کد بازیابی رمز عبور نامعتبر است",
			ERROR_RESET_PASSWORD_CODE_EXPIRED:                 "کد بازیابی رمز عبور منقضی شده است",
			ERROR_RESET_PASSWORD_RESEND_TIME_NOT_ELAPSED:      "زمان ارسال مجدد کد بازیابی رمز عبور نگذشته است",
			ERROR_DUPLICATE_USERNAME:                          "نام کاربری تکراری است",
			ERROR_DUPLICATE_MOBILE:                            "شماره موبایل تکراری است",
			ERROR_DUPLICATE_EMAIL:                             "ایمیل تکراری است",
			ERROR_DUPLICATE_KEY:                               "کلید تکراری است",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_DELETE_UPDATE: "خطا در عملیات حذف یا به‌روزرسانی",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_ADD_UPDATE:    "خطا در افزودن رکورد جدید",
			ERROR_SYSTEMIC_RECORD_UPDATE_DELETE:               "خطا در به‌روزرسانی یا حذف رکورد سیستمی",
			ERROR_USER_STATUS_IS_NOT_ACTIVE:                   "کاربر فعال نیست",
			ERROR_USER_STATUS_IS_NOT_PENDING:                  "کاربر در انتظار فعال‌سازی نیست",
			ERROR_UPLOADED_FILE_INCORRECT:                     "فایل نامعتبر است",
			ERROR_UPLOADED_FILE_NOT_SUPPORTED:                 "فایل پشتیبانی نمی‌شود",
			ERROR_DOWNLOAD_FILE_INCORRECT:                     "فایل دانلود نامعتبر است",
			ERROR_INVALID_TIME_RANGE:                          "بازه زمانی نامعتبر است",
			ERROR_DELETE_ADMIN_USER:                           "حذف کاربر مدیر",
			ERROR_ADMIN_UPDATE_HIMSELF_STATUS:                 "به‌روزرسانی وضعیت خود مدیر",
			ERROR_ADMIN_DELETE_HIMSELF_ADMIN_ROLE:             "نمی‌توانید نقش مدیر خود را حذف کنید",
			ERROR_INVALID_TYPE:                                "نوع نامعتبر است",
			ERROR_INVALID_USER_ROLE:                           "نقش کاربر نامعتبر است",
			ERROR_INVALID_AMOUNT:                              "مقدار نامعتبر است",
			ERROR_UPDATE_ADMIN_PASSWORD:                       "نمی‌توانید رمز عبور مدیر را به‌روزرسانی کنید",
			ERROR_INVALID_CURRENT_PASSWORD:                    "رمز عبور فعلی نامعتبر است",
			ERROR_INVALID_NEW_PASSWORD:                        "رمز عبور جدید نامعتبر است",
			ERROR_SYSTEMIC_SLIDER_DELETE:                      "نمی‌توانید اسلایدر سیستمی را حذف کنید",
			ERROR_INVALID_TIME:                                "زمان نامعتبر است",
			ERROR_USER_NOT_EXISTS:                             "کاربر وجود ندارد",
			ERROR_FIREBASE_TOKEN_DUPLICATED:                   "توکن Firebase تکراری است",
			ERROR_LOGIN_INVALID:                               "ورود نامعتبر است",
			ERROR_INVALID_MESSAGE_ACTION:                      "عملیات پیام نامعتبر است",
			ERROR_INVALID_MESSAGE_OBJECT:                      "شیء پیام نامعتبر است",
			ERROR_DUPLICATE_STATUS:                            "وضعیت تکراری است",
			ERROR_PRICE_CURRENCY_IS_INVALID:                   "واحد پولی %s نامعتبر است",
			ERROR_MONGO_DATABASE_VALIDATION_FAIELD:            "خطا در اعتبارسنجی پایگاه داده",
			ErrorInvalidDBField:                               "فیلد پایگاه داده نامعتبر است",
		},
		LANG_AR: {
			ERROR_HTTP_TOKEN_NOT_VALID:                        "الرمز غير صالح",
			ERROR_HTTP_NOT_FOUND:                              "الرابط غير موجود",
			ERROR_HTTP_TOKEN_NOT_ISSUED:                       "الرمز لم يتم إصداره بعد",
			ERROR_HTTP_TOKEN_NOT_VALID_YET:                    "الرمز غير صالح بعد",
			ERROR_HTTP_TOKEN_INVALID_SEGMENTS:                 "الرمز يحتوي على أقسام غير صالحة",
			ERROR_DEVICE_NOT_SELECETED:                        "لم يتم اختيار الجهاز",
			ERROR_HTTP_INTERNAL_SERVER_ERROR:                  "خطأ في الخادم الداخلي",
			ERROR_HTTP_FORBIDDEN:                              "حظر الوصول",
			ERROR_HTTP_FORBIDDEN_TO_SERVICE:                   "حظر الوصول إلى خدمة %s",
			ERROR_HTTP_FORBIDDEN_TO_ACTION:                    "حظر الوصول إلى عملية %s",
			ERROR_HTTP_TOKEN_EXPIRED:                          "انتهت صلاحية الرمز",
			ERROR_PARAMETERS_INVALID:                          "المعلمات غير صالحة",
			ERROR_GLOBAL_MYSQL_ERROR:                          "خطأ عام في MySQL",
			ERROR_INVALID_LOGIN:                               "بيانات تسجيل الدخول غير صالحة",
			ERROR_INVALID_MOBILE:                              "رقم الجوال غير صالح",
			ERROR_INVALID_EMAIL:                               "البريد الإلكتروني غير صالح",
			ERROR_INVALID_USERNAME:                            "اسم المستخدم غير صالح",
			ERROR_INVALID_PASSWORD:                            "كلمة المرور غير صالحة",
			ERROR_INVALID_STATUS:                              "الحالة غير صالحة",
			ERROR_INVALID_RESET_PASSWORD_CODE:                 "رمز إعادة تعيين كلمة المرور غير صالح",
			ERROR_RESET_PASSWORD_CODE_EXPIRED:                 "انتهت صلاحية رمز إعادة تعيين كلمة المرور",
			ERROR_RESET_PASSWORD_RESEND_TIME_NOT_ELAPSED:      "لم ينقض الوقت المطلوب لإعادة إرسال رمز إعادة تعيين كلمة المرور",
			ERROR_DUPLICATE_USERNAME:                          "اسم المستخدم مكرر",
			ERROR_DUPLICATE_MOBILE:                            "رقم الجوال مكرر",
			ERROR_DUPLICATE_EMAIL:                             "البريد الإلكتروني مكرر",
			ERROR_DUPLICATE_KEY:                               "مفتاح مكرر",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_DELETE_UPDATE: "خطأ في عملية الحذف أو التحديث",
			ERROR_FOREIGN_KEY_CONSTRAINT_FAILED_ADD_UPDATE:    "خطأ في إضافة سجل جديد",
			ERROR_SYSTEMIC_RECORD_UPDATE_DELETE:               "خطأ في تحديث أو حذف سجل نظامي",
			ERROR_USER_STATUS_IS_NOT_ACTIVE:                   "المستخدم غير نشط",
			ERROR_USER_STATUS_IS_NOT_PENDING:                  "المستخدم ليس في انتظار التفعيل",
			ERROR_UPLOADED_FILE_INCORRECT:                     "الملف غير صحيح",
			ERROR_UPLOADED_FILE_NOT_SUPPORTED:                 "الملف غير مدعوم",
			ERROR_DOWNLOAD_FILE_INCORRECT:                     "تنزيل الملف غير صحيح",
			ERROR_INVALID_TIME_RANGE:                          "نطاق الوقت غير صالح",
			ERROR_DELETE_ADMIN_USER:                           "حذف مستخدم المسؤول",
			ERROR_ADMIN_UPDATE_HIMSELF_STATUS:                 "تحديث حالة المسؤول نفسه",
			ERROR_ADMIN_DELETE_HIMSELF_ADMIN_ROLE:             "لا يمكن حذف دور المسؤول الخاص بك",
			ERROR_INVALID_TYPE:                                "نوع غير صالح",
			ERROR_INVALID_USER_ROLE:                           "الدور غير صالح",
			ERROR_INVALID_AMOUNT:                              "المبلغ غير صالح",
			ERROR_UPDATE_ADMIN_PASSWORD:                       "لا يمكن تحديث كلمة مرور المسؤول",
			ERROR_INVALID_CURRENT_PASSWORD:                    "كلمة المرور الحالية غير صالحة",
			ERROR_INVALID_NEW_PASSWORD:                        "كلمة المرور الجديدة غير صالحة",
			ERROR_SYSTEMIC_SLIDER_DELETE:                      "لا يمكن حذف المنزلق النظامي",
			ERROR_INVALID_TIME:                                "الوقت غير صحيح",
			ERROR_USER_NOT_EXISTS:                             "المستخدم غير موجود",
			ERROR_FIREBASE_TOKEN_DUPLICATED:                   "رمز Firebase مكرر",
			ERROR_LOGIN_INVALID:                               "تسجيل الدخول غير صالح",
			ERROR_INVALID_MESSAGE_ACTION:                      "عملية الرسالة غير صالحة",
			ERROR_INVALID_MESSAGE_OBJECT:                      "كائن الرسالة غير صالح",
			ERROR_DUPLICATE_STATUS:                            "الحالة مكررة",
			ERROR_PRICE_CURRENCY_IS_INVALID:                   "العملة غير صالحة %s",
			ERROR_MONGO_DATABASE_VALIDATION_FAIELD:            "خطأ في التحقق من قاعدة البيانات",
			ErrorInvalidDBField:                               "حقل قاعدة البيانات غير صالح",
		},
	}
)
