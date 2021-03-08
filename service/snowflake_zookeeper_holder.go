package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/busyfree/leaf-go/util/check"
	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/log"
	"github.com/busyfree/leaf-go/util/timeutil"
	"github.com/go-zookeeper/zk"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

//holder-init-zk_addr:/snowflake/com.sankuai.leaf.opensource.test/forever/169.254.121.237:2181-0000000000

var (
	PREFIX_ZK_PATH = "/snowflake/" + conf.GetString("LEAF_NAME")
	PROP_PATH      = filepath.Join(conf.GetConfigPath(), conf.GetString("LEAF_NAME")) + "/leafconf/{port}/workerID.toml"
	PATH_FOREVER   = PREFIX_ZK_PATH + "/forever" //保存所有数据持久的节点
	logger         = log.Get(context.Background())
)

type Endpoint struct {
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Timestamp int64  `json:"timestamp"`
}

type SnowFlakeZookeeperHolder struct {
	ZKAddressNode  string
	listenAddress  string
	WorkerId       int
	ip             string
	port           string
	connectionStr  string
	lastUpdateTime int64
}

func NewSnowFlakeZookeeperHolder(ip, port, connectionStr string) *SnowFlakeZookeeperHolder {
	s := new(SnowFlakeZookeeperHolder)
	s.ip = ip
	s.port = port
	s.listenAddress = ip + ":" + port
	s.connectionStr = connectionStr
	return s
}

func (s *SnowFlakeZookeeperHolder) watch(ev <-chan zk.Event) {
	for {
		select {
		case e := <-ev:
			logger.Infof("SnowFlakeZookeeperHolderE:%+v", e)
		}
	}
}

func (s *SnowFlakeZookeeperHolder) Init() bool {
	c, _, err := zk.Connect([]string{s.listenAddress}, time.Duration(6)*time.Second)
	if err != nil {
		panic(err)
	}
	boolExist, _, _ := c.Exists(PATH_FOREVER)
	if !boolExist {
		zkAddr := s.createNode(c)
		s.updateLocalWorkerID(s.WorkerId)
		go s.scheduledUploadData(c, zkAddr)
		return true
	}
	keys, _, err := c.Children(PATH_FOREVER)
	if err != nil {
		logger.Infof("c.ChildrenErr:%+v", err)
		return false
	}
	if len(keys) > 0 {
		nodeMap := make(map[string]int, 0)
		realNodeMap := make(map[string]string, 0)
		for _, node := range keys {
			nodeKey := strings.Split(node, "-")
			realNodeMap[nodeKey[0]] = node
			nodeMap[nodeKey[0]] = cast.ToInt(nodeKey[1])
		}
		if workId, ok := nodeMap[s.listenAddress]; ok {
			zkAddrNode := PATH_FOREVER + "/" + realNodeMap[s.listenAddress]
			s.WorkerId = workId
			if !s.checkInitTimeStamp(c, zkAddrNode) {
				viper.SetConfigName("workerID")
				viper.AddConfigPath(conf.GetConfigPath())
				err := viper.ReadInConfig()
				if err != nil {
					panic(fmt.Errorf("read file error:%v", err))
				}
				logger.Errorf("START FAILED ,use local node file properties workerID-{%d}", viper.GetInt("workerID"))
				return false
			}
			s.ZKAddressNode = zkAddrNode
			s.doService(c)
			s.updateLocalWorkerID(s.WorkerId)
			logger.Infof("[Old NODE]find forever node have this endpoint ip-{%s} port-{%s} workid-{%d} childnode and start SUCCESS", s.ip, s.port, s.WorkerId)
		} else {
			newNode := s.createNode(c)
			s.ZKAddressNode = newNode
			nodeKeyArr := strings.Split(newNode, "-")
			s.WorkerId = cast.ToInt(nodeKeyArr[1])
			s.doService(c)
			s.updateLocalWorkerID(workId)
			logger.Infof("[New NODE]can not find node on forever node that endpoint ip-{%s} port-{%s} workid-{%d},create own node on forever node and start SUCCESS", s.ip, s.port, s.WorkerId)
		}
	}
	return true
}

func (s *SnowFlakeZookeeperHolder) doService(client *zk.Conn) {
	go s.scheduledUploadData(client, s.ZKAddressNode)
}

func (s *SnowFlakeZookeeperHolder) scheduledUploadData(client *zk.Conn, zkAddrNode string) {
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	for {
		select {
		case <-ticker.C:
			s.updateNewData(client, zkAddrNode)
		}
	}
}

func (s *SnowFlakeZookeeperHolder) createNode(client *zk.Conn) string {
	path := PATH_FOREVER + "/" + s.listenAddress + "-0000000000"
	zkNodePaths := strings.Split(path, "/")
	var root = "/"
	if len(zkNodePaths) > 0 {
		for _, p := range zkNodePaths {
			if len(p) == 0 {
				continue
			}
			root += p
			_, err := client.Create(root, []byte{}, zk.FlagTTL, zk.WorldACL(zk.PermAll))
			root += "/"
			if err != nil {
				continue
			}
		}
	}
	root = strings.TrimRight(root, "/")
	_, err := client.Create(root, s.buildData(), zk.FlagTTL, zk.WorldACL(zk.PermAll))
	if err != nil {
		logger.Infof("zk CreateValErr:%+v", err)
		return path
	}
	return path
}

func (s *SnowFlakeZookeeperHolder) updateNewData(client *zk.Conn, path string) {
	if timeutil.MsTimestampNow() < s.lastUpdateTime {
		return
	}
	_, stat, _ := client.Get(path)
	var version int32 = 0
	if stat != nil {
		version = stat.Version
	}
	_, err := client.Set(path, s.buildData(), version)
	if err != nil {
		return
	}
	s.lastUpdateTime = timeutil.MsTimestampNow()
	return
}

func (s *SnowFlakeZookeeperHolder) buildData() []byte {
	endPoint := new(Endpoint)
	endPoint.IP = s.ip
	endPoint.Port = s.port
	endPoint.Timestamp = timeutil.MsTimestampNow()
	encodeArr, _ := json.Marshal(endPoint)
	return encodeArr
}

func (s *SnowFlakeZookeeperHolder) deBuildData(val []byte) *Endpoint {
	endPoint := new(Endpoint)
	_ = json.Unmarshal(val, endPoint)
	return endPoint
}

func (s *SnowFlakeZookeeperHolder) checkInitTimeStamp(client *zk.Conn, zkAddrNode string) bool {
	data, _, _ := client.Get(zkAddrNode)
	endpoint := s.deBuildData(data)
	return !(endpoint.Timestamp > timeutil.MsTimestampNow())
}

func (s *SnowFlakeZookeeperHolder) updateLocalWorkerID(workId int) {
	filePath := strings.Replace(PROP_PATH, "{port}", s.port, -1)
	if !check.CheckFileExist(filePath) {
		dirs := filepath.Dir(filePath)
		err := os.MkdirAll(dirs, os.ModePerm)
		if err != nil {
			return
		}
	}
	err := ioutil.WriteFile(filePath, []byte(fmt.Sprintf("workerID=%d", workId)), os.ModePerm)
	if err != nil {
		logger.Infof("%+v", err)
		return
	}
	return
}

func (s *SnowFlakeZookeeperHolder) GetWorkerId() int {
	return s.WorkerId
}
