package server

import (
	"bd_service/internal/models"
	"bd_service/internal/storage"
	"context"
	"errors"

	todov1 "github.com/Mirwinli/proto_todo/gen/go/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage interface {
	CreateTask(ctx context.Context, title string, description string, uid int64) error
	DeleteTask(ctx context.Context, title string, uid int64) error
	DoneTask(ctx context.Context, title string, uid int64) error
	ListTasks(ctx context.Context, uid int64) ([]models.Task, error)
}

type Server struct {
	todov1.UnimplementedTodoServer
	todo Storage
}

func Register(grpcServer *grpc.Server, storage Storage) {
	todov1.RegisterTodoServer(grpcServer, &Server{todo: storage})
}

func (s *Server) CreateTask(ctx context.Context, req *todov1.CreateRequest) (*todov1.CreateResponse, error) {
	const op = "server.CreateTaskServer"

	if err := s.todo.CreateTask(ctx, req.GetTitle(), req.GetDescription(), req.GetUid()); err != nil {
		return nil, status.Error(codes.Internal, "error to create task")
	}
	return &todov1.CreateResponse{}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *todov1.DeleteRequest) (*todov1.DeleteResponse, error) {
	const op = "server.DeleteTask"

	if err := s.todo.DeleteTask(ctx, req.GetTitle(), req.GetUid()); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, "error to delete task")
	}
	return &todov1.DeleteResponse{}, nil
}

func (s *Server) DoneTask(ctx context.Context, req *todov1.DoneRequest) (*todov1.DoneResponse, error) {
	const op = "server.DoneTask"

	if err := s.todo.DoneTask(ctx, req.GetTitle(), req.GetUid()); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, "error to done task")
	}
	return &todov1.DoneResponse{}, nil
}

func (s *Server) ListTasks(ctx context.Context, req *todov1.ListRequest) (*todov1.ListResponse, error) {
	const op = "server.ListTasks"

	tasks, err := s.todo.ListTasks(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "You haven't tasks")
		}
		return nil, status.Error(codes.Internal, "error to list tasks")
	}

	grpcTasks := make([]*todov1.Task, 0, len(tasks))

	for _, task := range tasks {
		grpcTask := &todov1.Task{
			Title:       task.Title,
			Description: task.Description,
			CreateAt:    timestamppb.New(task.CreatedAt),
			Duration:    durationpb.New(*task.Duration),
			IsDone:      task.IsDone,
			DoneAt:      timestamppb.New(*task.DoneAt),
			TaskId:      task.TaskId,
		}
		grpcTasks = append(grpcTasks, grpcTask)
	}
	return &todov1.ListResponse{
		Tasks: grpcTasks,
	}, nil
}
