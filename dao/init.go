package dao

import (
	"context"

	"github.com/busyfree/leaf-go/models/schema"
	"github.com/busyfree/leaf-go/util/db"
)

var (
	tableLeafAlloc = new(schema.LeafAlloc)
)

func SyncXORMTables() {
	ctx := context.Background()
	c := db.GetXORM(ctx, "default")
	_ = c.Sync2(tableLeafAlloc)
}
