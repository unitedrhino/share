package def

import (
	sq "gitee.com/unitedrhino/squirrel"
	"gorm.io/gorm"
	"time"
)

type PageInfo2 struct {
	TimeStart int64     `json:"timeStart"`
	TimeEnd   int64     `json:"timeEnd"`
	Page      int64     `json:"page" form:"page"`       // 页码
	Size      int64     `json:"size" form:"size"`       // 每页大小
	Orders    []OrderBy `json:"orderBy" form:"orderBy"` // 排序信息
}

type TimeRange struct {
	Start int64 `json:"start,optional"` //开始时间 unix时间戳
	End   int64 `json:"end,optional"`   //结束时间 unix时间戳
}

type DateRange struct {
	Start string `json:"start,optional"` //开始时间 格式：yyyy-mm-dd
	End   string `json:"end,optional"`   //结束时间 格式：yyyy-mm-dd
}

// 排序结构体
type OrderBy struct {
	Filed string `json:"filed" form:"filed"` //要排序的字段名
	Sort  int64  `json:"sort" form:"sort"`   //排序的方式：0 OrderAsc、1 OrderDesc
}

func (p PageInfo2) GetLimit() int64 {
	//if p.Size == 0 {
	//	return 20000
	//}
	return p.Size
}
func (p PageInfo2) GetOffset() int64 {
	if p.Page == 0 {
		return 0
	}
	return p.Size * (p.Page - 1)
}
func (p PageInfo2) GetTimeStart() time.Time {
	return time.UnixMilli(p.TimeStart)
}
func (p PageInfo2) GetTimeEnd() time.Time {
	return time.UnixMilli(p.TimeEnd)
}

func (p PageInfo2) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.TimeStart != 0 {
		sql = sql.Where("ts>=?", p.GetTimeStart())
	}
	if p.TimeEnd != 0 {
		sql = sql.Where("ts<=?", p.GetTimeEnd())
	}
	if p.Size != 0 {
		sql = sql.Limit(uint64(p.GetLimit()))
		if p.Page != 0 {
			sql = sql.Offset(uint64(p.GetOffset()))
		}
	}
	return sql
}

func (p PageInfo2) FmtSql2(sql *gorm.DB) *gorm.DB {
	if p.TimeStart != 0 {
		sql = sql.Where("ts>=?", p.GetTimeStart())
	}
	if p.TimeEnd != 0 {
		sql = sql.Where("ts<=?", p.GetTimeEnd())
	}
	if p.Size != 0 {
		sql = sql.Limit(int(p.GetLimit()))
		if p.Page != 0 {
			sql = sql.Offset(int(p.GetOffset()))
		}
	}
	if len(p.Orders) == 0 {
		sql = sql.Order("ts desc")
	}
	return sql
}

func (p PageInfo2) FmtWhere(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.TimeStart != 0 {
		sql = sql.Where(sq.GtOrEq{"ts": p.GetTimeStart()})
	}
	if p.TimeEnd != 0 {
		sql = sql.Where(sq.LtOrEq{"ts": p.GetTimeEnd()})
	}
	return sql
}

func (t TimeRange) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if t.Start != 0 {
		sql = sql.Where("created_time>=?", time.Unix(t.Start, 0))
	}
	if t.End != 0 {
		sql = sql.Where("created_time<=?", time.Unix(t.End, 0))
	}
	return sql
}

func (t *TimeRange) ToGorm(db *gorm.DB, column string) *gorm.DB {
	if t == nil {
		return db
	}
	if t.Start != 0 {
		db = db.Where(column+">=?", time.Unix(t.Start, 0))
	}
	if t.End != 0 {
		db = db.Where(column+"<=?", time.Unix(t.End, 0))
	}
	return db
}
