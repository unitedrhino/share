package caches

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
)

type UserOnline struct {
}

func (u *UserOnline) genKey(userID int64) string {
	return fmt.Sprintf("user:online:%v", userID)
}

func (u *UserOnline) SetUser(ctx context.Context, nodeID int64, userID int64) error {
	return store.SetexCtx(ctx, u.genKey(userID), cast.ToString(nodeID), 20)
}
func (u *UserOnline) DelUser(ctx context.Context, userID int64) error {
	_, err := store.DelCtx(ctx, u.genKey(userID))
	return err
}
