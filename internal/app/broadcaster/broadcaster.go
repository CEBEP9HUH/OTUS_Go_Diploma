package broadcaster

import "context"

type Broadcaster interface {
	Run(context.Context) error
}
