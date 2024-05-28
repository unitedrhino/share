package stores

import (
	"fmt"
	"gorm.io/gorm"
)

type Tree struct {
	ID       int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	ParentID int64  `gorm:"column:parent_id;type:bigint;NOT NULL"`      // 上级区域ID(雪花ID)
	IDPath   string `gorm:"column:id_path;type:varchar(1024);NOT NULL"` // 1-2-3-的格式记录顶级区域到当前区域的路径
}

type TreeWithName struct {
	ID       int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	ParentID int64  `gorm:"column:parent_id;type:bigint;NOT NULL"`        // 上级区域ID(雪花ID)
	IDPath   string `gorm:"column:id_path;type:varchar(1024);NOT NULL"`   // 1-2-3-的格式记录顶级区域到当前区域的路径
	NamePath string `gorm:"column:name_path;type:varchar(1024);NOT NULL"` // 1-2-3-的格式记录顶级区域到当前区域的路径
}

type IDPath struct {
	IDPath string `gorm:"column:id_path;type:varchar(1024);index"` // 1-2-3-的格式记录顶级到当前的路径
	ID     int64  `gorm:"column:id;type:bigint(20);index"`         //2是未分类,未使用的,1是根节点
}

type IDPathFilter struct {
	IDPath string
	ID     int64
}

func (i *IDPathFilter) Filter(db *gorm.DB, prefix string) *gorm.DB {
	if i == nil {
		return db
	}
	if i.ID != 0 {
		col := "id"
		if prefix != "" {
			col = Col(prefix + "_id")
		}
		db = db.Where(fmt.Sprintf("%s = ?", col), i.ID)
	}
	if i.IDPath != "" {
		col := "id_path"
		if prefix != "" {
			col = Col(prefix + "_id_path")
		}
		db = db.Where(fmt.Sprintf("%s like ?", col), i.IDPath+"%")
	}
	return db
}
