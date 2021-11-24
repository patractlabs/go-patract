package rand_extension

import "fmt"

type EventRandomUpdated struct {
	New VecU8 `scale:"operator"`
}

func (e EventRandomUpdated) String() string {
	return fmt.Sprintf("event RandomUpdated: %v", e.New)
}
