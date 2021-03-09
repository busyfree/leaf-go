package service

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/busyfree/leaf-go/models"
	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/timeutil"
	"github.com/spf13/cast"
)

var (
	workerIdBits       = 10
	maxWorkerId        = ^(-1 << workerIdBits)
	sequenceBits       = 12
	workerIdShift      = sequenceBits
	timestampLeftShift = sequenceBits + workerIdBits
	sequenceMask       = ^(-1 << sequenceBits)
)

type SnowFlakeIdGenImpl struct {
	twepoch       int64
	workerId      int64
	sequence      int64
	lastTimestamp int64
}

func NewSnowFlakeIdGenImpl(zkAddr string, port int, twepoch int64) *SnowFlakeIdGenImpl {
	s := new(SnowFlakeIdGenImpl)
	s.twepoch = twepoch
	if !(timeutil.MsTimestampNow() > twepoch) {
		panic("Snowflake not support twepoch gt currentTime")
	}
	zkEnable := conf.GetBool("LEAF_SNOWFLAKE_ENABLE_ZK")
	if zkEnable {
		ip := s.getHostAddress(conf.GetString("LEAF_SNOWFLAKE_ETHER"))
		holder := NewSnowFlakeZookeeperHolder(ip, fmt.Sprintf("%d", port), zkAddr)
		logger.Infof("twepoch:{%d} ,ip:{%s} ,zkAddress:{%s} port:{%d}", twepoch, ip, zkAddr, port)
		if !holder.Init() {
			panic("Snowflake Id Gen is not init ok")
		}
		s.workerId = int64(holder.GetWorkerId())
	} else {
		s.workerId = conf.GetInt64("LEAF_SNOWFLAKE_WORKER_ID")
	}
	logger.Infof("START SUCCESS USE ZK WORKERID-{%d}", s.workerId)
	if !(s.workerId >= 0 && s.workerId <= int64(maxWorkerId)) {
		panic("workerID must gte 0 and lte 1023")
	}
	return s
}

func (s *SnowFlakeIdGenImpl) Init(ctx context.Context) bool {
	return true
}

func (s *SnowFlakeIdGenImpl) Get(ctx context.Context, key string) models.Result {
	var ts = timeutil.MsTimestampNow()
	if ts > s.lastTimestamp {
		offset := s.lastTimestamp - ts
		if offset <= 5 {
			time.Sleep(time.Duration(offset<<1) * time.Millisecond)
			ts = timeutil.MsTimestampNow()
			if ts < s.lastTimestamp {
				return models.Result{Id: -1, Status: models.EXCEPTION}
			}
		} else {
			return models.Result{Id: -3, Status: models.EXCEPTION}
		}
	}
	if ts == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & int64(sequenceMask)
		if s.sequence == 0 {
			s.sequence = int64(rand.Intn(100))
			ts = s.tilNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequence = int64(rand.Intn(100))
	}
	s.lastTimestamp = ts
	id := ((ts - s.twepoch) << timestampLeftShift) | (s.workerId << workerIdShift) | s.sequence
	return models.Result{Id: id, Status: models.SUCCESS}
}

func (s *SnowFlakeIdGenImpl) tilNextMillis(lastTimestamp int64) int64 {
	var ts = timeutil.MsTimestampNow()
	for ts <= lastTimestamp {
		ts = timeutil.MsTimestampNow()
	}
	return ts
}

func (s *SnowFlakeIdGenImpl) getHostAddress(interfaceName string) string {
	ips, err := s.ips()
	if err != nil {
		return ""
	}
	if len(interfaceName) > 0 {
		if val, ok := ips[interfaceName]; ok {
			return val
		}
	} else {
		for _, ip := range ips {
			return ip
		}
	}
	return ""
}

func (s *SnowFlakeIdGenImpl) ips() (map[string]string, error) {
	ips := make(map[string]string)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ips[byName.Name] = ip.String()
			case *net.IPAddr:
				ip = v.IP
				ips[byName.Name] = ip.String()
			}
		}
	}
	return ips, nil
}

func (s *SnowFlakeIdGenImpl) DecodeSnowflakeId(idStr string) map[string]interface{} {
	var out = make(map[string]interface{}, 0)
	var snowflakeId = cast.ToInt64(idStr)
	originTimestamp := (snowflakeId >> 22) + s.twepoch
	out["timestamp"] = fmt.Sprintf("%d (%s)", originTimestamp, timeutil.MsTimestamp2Time(originTimestamp).Format("2006-01-02 15:04:05.000"))
	workerId := (snowflakeId >> 12) ^ (snowflakeId >> 22 << 10)
	sequence := snowflakeId ^ (snowflakeId >> 12 << 12)
	out["workerId"] = workerId
	out["sequenceId"] = sequence
	return out
}
