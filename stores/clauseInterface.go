package stores

import "gorm.io/gorm/clause"

const AuthModify = "authModify:%v"

type Opt int64

const (
	Create Opt = iota + 1
	Update
	Delete
	Select
)

type ClauseInterface struct {
}

func (sd ClauseInterface) Name() string {
	return ""
}

func (sd ClauseInterface) Build(clause.Builder) {

}

func (sd ClauseInterface) MergeClause(*clause.Clause) {

}
