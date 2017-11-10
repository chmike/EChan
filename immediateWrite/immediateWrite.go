package immediateWrite

import (
	"github.com/chmike/EChan"
)

// New directly forwards in to out
func New() echan.Interface {
	return echan.ChanForward
}
