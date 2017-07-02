package histories

import "time"

type Synchronizer interface {
	SyncRecent() (int, error)
	SyncAll() (int, error)
	GetLastTime() time.Time
}
