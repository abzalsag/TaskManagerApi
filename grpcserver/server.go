package grpcserver

import (
	"context"
	"errors"

	"TaskManagerApi/proto"
	"TaskManagerApi/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskGRPCServer struct {
	proto.UnimplementedTaskServiceServer
	taskService *service.TaskService
}

func NewTaskGRPCServer(taskService *service.TaskService) *TaskGRPCServer {
	return &TaskGRPCServer{
		taskService: taskService,
	}
}

func (s *TaskGRPCServer) GetTask(
	ctx context.Context,
	req *proto.GetTaskRequest,
) (*proto.TaskResponse, error) {
	taskID := int(req.GetId())
	userID := 1

	task, err := s.taskService.GetTaskById(
		ctx,
		taskID,
		userID,
	)

	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"task not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"internal server error",
		)
	}

	return &proto.TaskResponse{
		Id:          int32(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}, nil
}
