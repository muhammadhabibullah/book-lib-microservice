package grpc

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func ParseErrorStatus(err error) error {
	errStatus, ok := status.FromError(err)
	if ok {
		errMessage := errStatus.Message()
		errCode := errStatus.Code().String()

		return fmt.Errorf("[%s] %s", errCode, errMessage)
	}

	return nil
}
