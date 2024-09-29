package lib

import (
	"fmt"

	"github.com/fatih/color"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnpackGRPCError unpacks a gRPC error and prints user-friendly messages based on the error details.
// It handles specific error cases, particularly for InvalidArgument errors.
func UnpackGRPCError(err error) {
	red := color.New(color.FgRed).SprintFunc()

	st := status.Convert(err)
	if st.Code() == codes.InvalidArgument {
		for _, detail := range st.Details() {
			switch t := detail.(type) { //nolint:gocritic
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					fmt.Printf("The %s field %s\n", red(violation.GetField()), red(violation.GetDescription()))
				}
			}
		}
	} else {
		fmt.Printf("Please try again: %s\n", red(st.Message()))
	}
}
