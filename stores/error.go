package stores

import (
	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
	"strings"
)

func ErrFmt(err error) error {
	if err == nil {
		return nil
	}
	if err.Error() == "redis: nil" {
		return errors.NotFind
	}
	if _, ok := err.(*errors.CodeError); ok {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFind.WithStack(1)
	}
	if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
		return errors.Duplicate.AddDetail(err)
	}
	return errors.Database.AddDetail(err)
}
