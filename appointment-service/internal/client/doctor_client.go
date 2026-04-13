package client

import (
	"context"

	"github.com/syndaly1/ap2-assignment2/appointment-service/internal/usecase"
	doctorpb "github.com/syndaly1/ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorClient struct {
	client doctorpb.DoctorServiceClient
}

func NewDoctorClient(client doctorpb.DoctorServiceClient) *DoctorClient {
	return &DoctorClient{client: client}
}

func (c *DoctorClient) GetDoctor(ctx context.Context, id string) error {
	_, err := c.client.GetDoctor(ctx, &doctorpb.GetDoctorRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return usecase.ErrDoctorUnavailable
		}

		switch st.Code() {
		case codes.NotFound:
			return usecase.ErrDoctorNotFound
		case codes.Unavailable:
			return usecase.ErrDoctorUnavailable
		default:
			return usecase.ErrDoctorUnavailable
		}
	}

	return nil
}
