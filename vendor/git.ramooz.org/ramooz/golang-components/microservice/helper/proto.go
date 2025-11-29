package helper

import (
	"encoding/json"
	"errors"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"

	error3 "git.ramooz.org/ramooz/golang-components/microservice/error"
)

// ConvertProtoToModel convert proto structure to model
func ConvertProtoToModel[T any](proto proto.Message) (model T, err error) {
	m, ok := reflect.New(reflect.TypeOf(model).Elem()).Interface().(T)
	if !ok {
		return model, errors.New("model type is invalid")
	}
	b, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(proto)
	if err != nil {
		return m, componentsError.New(error3.ErrorInvalidProtoModel, []string{err.Error(), string(b)})
	}
	if err = json.Unmarshal(b, m); err != nil {
		return m, componentsError.New(error3.ErrorInvalidProtoModel, []string{err.Error(), string(b)})
	}
	return m, nil
}

// ConvertModelToProto model structure to proto
func ConvertModelToProto[T proto.Message](model any) (protoModel T, err error) {
	m, ok := reflect.New(reflect.TypeOf(protoModel).Elem()).Interface().(T)
	if !ok {
		return protoModel, errors.New("proto type is invalid")
	}
	if reflect.ValueOf(model).IsZero() {
		return m, componentsError.New(error3.ErrorModelIsNil, nil)
	}
	b, err := json.Marshal(model)
	if err != nil {
		return m, componentsError.New(error3.ErrorInvalidModel, []string{err.Error(), string(b)})
	}
	if err = (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(b, m); err != nil {
		return m, componentsError.New(error3.ErrorInvalidModel, []string{err.Error(), string(b)})
	}
	return m, nil
}

/*// ConvertProtoArrayToModels convert proto list to models
func ConvertProtoArrayToModels[T any](protos []proto.Message) (models []T, err error) {
	for _, proto := range protos {
		m, err := ConvertProtoToModel[T](proto)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}*/
