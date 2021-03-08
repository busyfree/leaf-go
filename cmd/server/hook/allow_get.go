package hook

import (
	"context"

	"github.com/bilibili/twirp"
)

func NeAllowGet() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			ctx = twirp.WithAllowGET(ctx, true)
			return ctx, nil
		},
	}
}
