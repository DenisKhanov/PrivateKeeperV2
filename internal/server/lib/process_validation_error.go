package lib

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProcessValidationError processes validation errors and returns a structured error.
func ProcessValidationError(msg string, report map[string]string) error {
	st := status.New(codes.InvalidArgument, msg)
	details := make([]*errdetails.BadRequest_FieldViolation, 0, len(report))
	for field, message := range report {
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: message,
		})
	}
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, details...)
	st, err := st.WithDetails(br)
	if err != nil {
		logrus.WithError(err).Error("Internal error: failed to set details")
		return status.Error(codes.Internal, "internal error")
	}
	return st.Err()
}
