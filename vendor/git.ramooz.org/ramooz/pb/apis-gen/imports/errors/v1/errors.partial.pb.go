package errors

import (
	"encoding/json"
	"fmt"
	error2 "git.ramooz.org/ramooz/golang-components/error-handler"
)

func (x *GetErrorsRequest) GetErrorResponse() (*GetErrorsResponse, error) {
	if x.ResultType == GetErrorsRequest_GRPC {
		grpcResponse := GrpcErrors{
			Errors: convertErrorMapToGrpcResponse(),
		}
		return &GetErrorsResponse{Errors: &GetErrorsResponse_Grpc{Grpc: &grpcResponse}}, nil
	} else if x.ResultType == GetErrorsRequest_JSON {
		errors, err := json.Marshal(error2.I18nTranslates)
		if err != nil {
			panic(err)
		}
		errorsString := string(errors)
		return &GetErrorsResponse{Errors: &GetErrorsResponse_Json{Json: errorsString}}, nil
	}
	return nil, fmt.Errorf("please select result type")
}

func convertErrorMapToGrpcResponse() []*GrpcError {
	grpcMap := []*GrpcError{}
	for lang, errors := range error2.I18nTranslates {
		grpcMap = append(grpcMap, &GrpcError{
			Language:     lang,
			ErrorMessage: convertMapInToInt64(errors),
		})
	}
	return grpcMap
}

func convertMapInToInt64(mapInt map[int]string) map[int64]string {
	map64 := map[int64]string{}
	for k, v := range mapInt {
		map64[int64(k)] = v
	}
	return map64
}
