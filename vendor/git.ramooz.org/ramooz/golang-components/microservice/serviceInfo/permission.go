package serviceInfo

type Permission struct {
	Name           string          `bson:"name" json:"name"`
	Title          string          `bson:"title" json:"title"`
	Code           int32           `bson:"code" json:"code"`
	SubPermissions []SubPermission `bson:"sub_permissions" json:"sub_permissions"`
}

type SubPermission struct {
	ServiceName     string  `bson:"service_name" json:"service_name"`
	PermissionCodes []int32 `bson:"permission_codes" json:"permission_codes"`
}
