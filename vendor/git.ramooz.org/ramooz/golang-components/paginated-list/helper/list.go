package helpers

import (
	"bytes"
	componentsList "git.ramooz.org/ramooz/golang-components/paginated-list"
	"git.ramooz.org/ramooz/pb/apis-gen/imports/list"
)

const (
	DEFAULT_PAGE_SIZE          = 20
	RESULT_TYPE_PAGINATED_JSON = 0
	RESULT_TYPE_EXCEL          = 1
)

func ToProtoResponse(filterList *componentsList.List) (response *list.PaginatedListResponse) {
	if filterList != nil {
		response = &list.PaginatedListResponse{
			Page:            int32(filterList.PageNo),
			TotalItemsCount: int64(filterList.TotalItemsCount),
			Sort:            filterList.Sort,
			Filters:         filterList.Filters,
			PageSize:        int32(filterList.PageSize),
			ResultType:      int32(filterList.ResultType),
		}
		if filterList.ResultType == RESULT_TYPE_EXCEL {
			var buffer bytes.Buffer

			filterList.ExcelFile.Write(&buffer)
			response.FileData = buffer.Bytes()
		}
	}

	return response
}

func NewListFromGrpcListRequest(request *list.PaginatedListRequest) *componentsList.List {
	if request == nil {
		request = &list.PaginatedListRequest{
			Page:     1,
			PageSize: DEFAULT_PAGE_SIZE,
		}
	}
	list := componentsList.NewList()
	list.Filters = request.Filters
	list.Sort = request.Sort
	if request.Page > 0 {
		list.PageNo = int(request.Page)
	}
	if request.PageSize > 0 {
		list.PageSize = int(request.PageSize)
	}
	list.ResultType = int(request.ResultType)
	return list
}

//func ConvertListToGrpcList(ctx iris.Context, authToken string) *grpc.PaginatedListRequest {
//	paginationRequest := &grpc.PaginatedListRequest{
//		Auth:     &grpc.AuthRequest{Token: authToken},
//		PageSize: DEFAULT_PAGE_SIZE,
//		PageNo:   1,
//	}
//	list := &componentsList2.List{}
//	global.FormDecoder.Decode(list, ctx.FormValues())
//	if list.PageNo > 0 {
//		paginationRequest.PageNo = int32(list.PageNo)
//	}
//	if list.PageSize > 0 {
//		paginationRequest.PageSize = int32(list.PageSize)
//	}
//	paginationRequest.Sort = list.Sort
//	paginationRequest.Filters = list.Filters
//	paginationRequest.ResultType = int32(list.ResultType)
//	return paginationRequest
//}

func ConvertListToGrpcListFromService(filters map[string]string) *list.PaginatedListRequest {
	return &list.PaginatedListRequest{
		PageSize: DEFAULT_PAGE_SIZE,
		Page:     1,
		Filters:  filters,
	}
}
