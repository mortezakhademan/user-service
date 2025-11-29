package serviceInfo

import (
	"context"
	"encoding/json"
	service "git.ramooz.org/ramooz/pb/apis-gen/services/user/v2"
	"google.golang.org/grpc"
	"io/ioutil"

	"git.ramooz.org/ramooz/golang-components/microservice/helper"
)

/*
*
service_info used for component internal functions
*/
type ServiceInfo struct {
	Name        string       `bson:"name" json:"name"`
	Title       string       `bson:"title" json:"title"`
	Code        int32        `bson:"code" json:"code"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}

// GetFromFile get service from file
func GetFromFile(jsonFilePath string) *ServiceInfo {
	serviceInfo := readJsonServiceInfoFile(jsonFilePath)
	setPermissionsInUserService(serviceInfo)

	return serviceInfo
}

// GetFromEmbed get service from embedded file
func GetFromEmbed(svcInfo string) *ServiceInfo {
	serviceInfo := readJsonServiceInfo(svcInfo)
	setPermissionsInUserService(serviceInfo)

	return serviceInfo
}

func SetEnvFromFile(filePath ...string) error {
	if err := helper.SetEnvFile(filePath...); err != nil {
		return err
	}
	return nil
}

func SetEnvFromMap(envData map[string]string) error {
	if err := helper.SetEnvFromMap(envData); err != nil {
		return err
	}
	return nil
}

func (serviceInfo *ServiceInfo) SyncServiceInfo(userServiceConnection *grpc.ClientConn, forceUpdate bool) {
	SetServiceInfoClientConnection(userServiceConnection)
	serviceInfoClient := GetServiceInfoClient()

	serviceInfoProto := service.ServiceInfo{}
	helper.MapToMap(serviceInfo, &serviceInfoProto)
	res, err := serviceInfoClient.SyncServiceInfo(context.Background(), &service.ServiceInfoRequest{
		ServiceInfo: &serviceInfoProto,
		ForceUpdate: forceUpdate,
	})

	if err != nil {
		helper.Log.Errorf("unable to Sync service info to user service, %v", err)
		return
	}

	helper.Log.Infof("service info Sync result, %v", res)
}

func readJsonServiceInfoFile(filePath string) *ServiceInfo {
	jsonData, errFile := ioutil.ReadFile(filePath)
	if errFile != nil {
		helper.Log.Panic(errFile)
	}

	serviceInfo := ServiceInfo{}
	err := json.Unmarshal(jsonData, &serviceInfo)
	if err != nil {
		helper.Log.Panic(err)
	}

	return &serviceInfo
}

func readJsonServiceInfo(infoData string) *ServiceInfo {
	serviceInfo := ServiceInfo{}
	err := json.Unmarshal([]byte(infoData), &serviceInfo)
	if err != nil {
		helper.Log.Panic(err)
	}

	return &serviceInfo
}

func setPermissionsInUserService(serviceInfo *ServiceInfo) {
	//todo: implement
}

func (serviceInfo *ServiceInfo) GetPermissionNamesByPermissionCodes(permissionCodes []int32) (permissionNames []string) {
	permissionNames = make([]string, len(permissionCodes))
	for codeIndex, code := range permissionCodes {
		for _, permission := range serviceInfo.Permissions {
			if permission.Code == code {
				permissionNames[codeIndex] = permission.Name
				break
			}
		}
	}
	return permissionNames
}

func (serviceInfo *ServiceInfo) GetPermissionCode(name string) int32 {

	permissions := serviceInfo.Permissions
	for _, permission := range permissions {
		if permission.Name == name {
			return permission.Code
		}
	}
	return 0
}
