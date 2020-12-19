package requestManager

import (
	"fmt"

	"github.com/rsheasby/gocrypt/protocol"
	"golang.org/x/crypto/bcrypt"
)

func validateRequest(req *protocol.Request) (err error) {
	if req.Cost < int32(bcrypt.MinCost) || req.Cost > int32(bcrypt.MaxCost) {
		return fmt.Errorf("invalid cost provided - cost must be between %d and %d, but cost of %d was provided", bcrypt.MinCost, bcrypt.MaxCost, req.Cost)
	}
	return nil
}
