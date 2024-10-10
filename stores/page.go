package stores

import (
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
)

type Order = int64

const (
	OrderAsc  Order = iota + 1 //从久到近排序
	OrderDesc                  //时间从近到久排序
)

var orderMap = map[int64]string{
	OrderAsc:  "asc",
	OrderDesc: "desc",
}

type PageInfo struct {
	Page   int64     `json:"page" form:"page"`         // 页码
	Size   int64     `json:"pageSize" form:"pageSize"` // 每页大小
	Orders []OrderBy `json:"orderBy" form:"orderBy"`   // 排序信息
}

// 排序结构体
type OrderBy struct {
	Field string `json:"field" form:"field"` //要排序的字段名
	Sort  Order  `json:"sort" form:"sort"`   //排序的方式：1 OrderAsc、2 OrderDesc
}

func (p *PageInfo) GetLimit() int64 {
	if p == nil || p.Size == 0 {
		return 2000
	}
	return p.Size
}
func (p *PageInfo) GetOffset() int64 {
	if p == nil || p.Page == 0 {
		return 0
	}
	return p.Size * (p.Page - 1)
}

// 获取排序参数
func (p *PageInfo) getOrders() (arr []string) {
	if p != nil && len(p.Orders) > 0 {
		for _, o := range p.Orders {
			arr = append(arr, fmt.Sprintf("%s %s", Col(utils.CamelCaseToUdnderscore(o.Field)), orderMap[o.Sort]))
		}
	}
	return
}
func (p *PageInfo) WithDefaultOrder(in ...OrderBy) *PageInfo {
	if p == nil {
		p = &PageInfo{}
	}
	if len(p.Orders) == 0 {
		p.Orders = in
	}
	return p
}

func (p *PageInfo) WithOrder(in ...OrderBy) *PageInfo {
	if p == nil {
		p = &PageInfo{}
	}
	p.Orders = in
	return p
}

func (p *PageInfo) ToGorm(db *gorm.DB) *gorm.DB {
	if p == nil {
		return db.Limit(1000)
	}
	if p.Size != 0 {
		db = db.Limit(int(p.GetLimit()))
		if p.Page != 0 {
			db = db.Offset(int(p.GetOffset()))
		}
	}

	if len(p.Orders) != 0 {
		orders := p.getOrders()
		for _, o := range orders {
			db = db.Order(o)
		}
	}
	return db
}
