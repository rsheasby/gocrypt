package requestManager

import (
	"fmt"

	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/protocol"
	"golang.org/x/crypto/bcrypt"
)

func validateRequest(req *protocol.Request) (err error) {
	// Ensure the request type is valid
	if req.RequestType != protocol.Request_HASHPASSWORD && req.RequestType != protocol.Request_VERIFYPASSWORD {
		return fmt.Errorf("invalid request type provided - should be either HASHPASSWORD or VERIFYPASSWORD but received invalid int instead: %d", req.RequestType)
	}

	// Input validation for all request types
	if len(req.ResponseKey) < config.MinResponseKeyLength {
		return fmt.Errorf("response key is too short - should be %d characters at a minimum, but provided key had a length of %d", config.MinResponseKeyLength, len(req.ResponseKey))
	}
	if len(req.Password) == 0 {
		return fmt.Errorf("password field is empty")
	}

	// Input validation for HASHPASSWORD request
	if req.RequestType == protocol.Request_HASHPASSWORD {
		if req.Cost < int32(bcrypt.MinCost) || req.Cost > int32(bcrypt.MaxCost) {
			return fmt.Errorf("invalid cost provided - cost must be between %d and %d, but cost of %d was provided", bcrypt.MinCost, bcrypt.MaxCost, req.Cost)
		}
	}

	// Input validation for VERIFYPASSWORD request
	if req.RequestType == protocol.Request_VERIFYPASSWORD {
		if len(req.Hash) == 0 {
			return fmt.Errorf("hash field is empty")
		}
	}
	return nil
}
