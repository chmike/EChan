// Package immediateWrite provides a echan which synchronously writes all items from in to out.
package immediateWrite

import (
	"github.com/chmike/EChan"
)

// New directly forwards in to out
func New() echan.Interface {
	return echan.ChanForward
}
