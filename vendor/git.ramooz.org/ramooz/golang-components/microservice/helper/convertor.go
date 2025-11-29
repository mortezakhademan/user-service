package helper

import (
	"encoding/json"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/encoding/protojson"
)

func MapToMap(fromStruct interface{}, toStruct interface{}) error {
	jsonBinary, err := json.Marshal(fromStruct)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBinary, toStruct)
	if err != nil {
		return err
	}

	return nil
}

func MapToPBStruct(m map[string]interface{}) (*structpb.Struct, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func PBStructToMap(pbStruct *structpb.Struct) map[string]interface{} {
	modelMap := make(map[string]interface{})
	b, _ := json.Marshal(pbStruct)
	_ = json.Unmarshal(b, &modelMap)
	return modelMap
}

func StringToBoolean(flag string) bool {
	switch flag {
	case "true", "1":
		return true
	case "false", "0":
		return false
	default:
		return true
	}
}
