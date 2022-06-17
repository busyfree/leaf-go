package job

import (
	"context"
	"fmt"

	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/db"
	"github.com/busyfree/leaf-go/util/log"
)

func init() {
	// go run main.go job once test args ,will trigger this func
	manual("cfg", func(ctx context.Context) error {
		path := conf.GetConfigPath()
		var (
			logger = log.Get(ctx)
			ts     int64
		)
		logger.Infof("conf.GetConfigPath:%s", path)
		defaultDSN := conf.GetString("DB_DEFAULT_DSN")
		logger.Infof("defaultDSN:%s", defaultDSN)
		c := db.Get(ctx, "default")
		sqlSelect := fmt.Sprintf("SELECT unix_timestamp(now()) AS ts")
		q := db.SQLSelect("1", sqlSelect)
		row := c.QueryRowContext(ctx, q)
		err := row.Scan(&ts)
		if err != nil {
			return err
		}
		logger.Infof("unix_timestamp(now()) IS :%d", ts)
		return nil
	})
}
