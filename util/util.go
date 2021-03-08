

package util

import (
	_ "github.com/busyfree/leaf-go/util/conf" // init conf
	"github.com/busyfree/leaf-go/util/redis"

	"github.com/busyfree/leaf-go/util/db"
	"github.com/busyfree/leaf-go/util/log"
	"github.com/busyfree/leaf-go/util/mc"
)

// GatherMetrics 收集一些被动指标
func GatherMetrics() {
	mc.GatherMetrics()
	redis.GatherMetrics()
	db.GatherMetrics()
}

// Reset all utils
func Reset() {
	log.Reset()
	db.ResetXORM()
	db.Reset()
	mc.Reset()

}

// Stop all utils
func Stop() {
}
