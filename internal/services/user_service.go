package services

import (
	"context"
	"git.ramooz.org/ramooz/golang-components/paginated-list/helper"
	"git.ramooz.org/ramooz/pb/apis-gen/imports/list"
	"github.com/mortezakhademan/user-service-sample/internal/domain"
	"github.com/mortezakhademan/user-service-sample/internal/repository"
	pbUser "github.com/mortezakhademan/user-service-sample/services/proto/apis-gen/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pbUser.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

func (s *UserService) InsertUser(ctx context.Context, req *pbUser.InsertUserRequest) (*pbUser.Id, error) {
	user := &domain.User{
		Name:  req.Name,
		Phone: req.Phone,
	}
	userId, err := s.userRepo.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pbUser.Id{
		Id: userId.Hex(),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pbUser.Id) (*pbUser.User, error) {
	id, _ := bson.ObjectIDFromHex(req.Id)
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pbUser.User{
		Id:    user.ID.Hex(),
		Name:  user.Name,
		Phone: user.Phone,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pbUser.UpdateUserRequest) (*emptypb.Empty, error) {
	userId, _ := bson.ObjectIDFromHex(req.Id)
	user := &domain.User{
		ID:    userId,
		Name:  req.Name,
		Phone: req.Phone,
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pbUser.Id) (*emptypb.Empty, error) {
	userId, _ := bson.ObjectIDFromHex(req.Id)
	if err := s.userRepo.Delete(ctx, userId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUserList(ctx context.Context, req *list.PaginatedListRequest) (*pbUser.UserListResponse, error) {
	paginatedList := helpers.NewListFromGrpcListRequest(req)
	users, err := s.userRepo.List(ctx, paginatedList)

	if err != nil {
		return nil, err
	}
	pbUsers := []*pbUser.UserListItem{}
	for _, user := range users {
		pbUsers = append(pbUsers, &pbUser.UserListItem{
			Id:    user.ID.Hex(),
			Name:  user.Name,
			Phone: user.Phone,
		})
	}
	return &pbUser.UserListResponse{
		Data:     pbUsers,
		Response: helpers.ToProtoResponse(paginatedList),
	}, nil
}
