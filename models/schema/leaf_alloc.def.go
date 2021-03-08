package schema

import (
	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/timeutil"
)

type LeafAlloc struct {
	BizTag      string `xorm:"biz_tag VARCHAR(128) notnull pk" db:"biz_tag" json:"biz_tag"`
	MaxId       int64  `xorm:"max_id BIGINT(20) notnull default 0" db:"max_id" json:"max_id"`
	Step        int    `xorm:"step INT(11) notnull default 0" db:"step" json:"step"`
	Description string `xorm:"description VARCHAR(256) notnull" db:"description" json:"description"`
	CreatedAt   int64  `xorm:"created_at BIGINT(20) notnull default 0" db:"created_at" json:"created,omitempty"`
	DeletedAt   int64  `xorm:"deleted_at BIGINT(20) notnull default 0" db:"deleted_at" json:"deleted,omitempty"`
	UpdatedAt   int64  `xorm:"update_time BIGINT(20) notnull default 0" db:"updated_at" json:"updated_at"`
	Updated     string `xorm:"-" db:"-" json:"updated"`
}

func (p *LeafAlloc) TableName() string {
	prefix := conf.GetString("DB_DEFAULT_TABLE_PREFIX")
	if len(prefix) > 0 {
		return prefix + "_leaf_alloc"
	}
	return "leaf_alloc"
}

func (p *LeafAlloc) BeforeInsert() {
	p.CreatedAt = timeutil.MsTimestampNow()
	p.UpdatedAt = p.CreatedAt
	p.DeletedAt = 0
}

func (p *LeafAlloc) BeforeUpdate() {
	p.UpdatedAt = timeutil.MsTimestampNow()
}

func (p *LeafAlloc) AfterLoad() {
	p.UpdatedAt = timeutil.MsTimestampNow()
	p.Updated = timeutil.MsTimestamp2Time(p.UpdatedAt).Format("2006-01-02 15:04:05")
}
