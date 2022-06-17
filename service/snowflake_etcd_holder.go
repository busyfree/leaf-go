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

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client/v3"

	"github.com/busyfree/leaf-go/util/check"
	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/timeutil"
)

type SnowFlakeEtcdHolder struct {
	endpoints       []string
	etcdAddressNode string
	ip              string
	port            string
	listenAddress   string
	connectionStr   string
	lastUpdateTime  int64
	dialTimeout     int
	WorkerId        int
}

func NewSnowFlakeEtcdHolder(ip, port string, endpoints []string, dialTimeout int) *SnowFlakeEtcdHolder {
	s := new(SnowFlakeEtcdHolder)
	s.ip = ip
	s.port = port
	s.endpoints = endpoints
	if dialTimeout == 0 {
		dialTimeout = 10
	}
	s.dialTimeout = dialTimeout
	s.listenAddress = ip + ":" + port
	s.Init()
	return s
}

func (s *SnowFlakeEtcdHolder) Init() bool {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   s.endpoints,
		DialTimeout: time.Duration(s.dialTimeout) * time.Second,
	})
	if err != nil {
		panic(err)
	}
	var retryCount = 0
RETRY:
	resp, err := c.Get(context.Background(), PATH_FOREVER)
	if err != nil {
		panic(err)
	}
	switch true {
	case resp == nil:
	case resp != nil && len(resp.Kvs) == 0:
		zkAddr := s.createNode(c)
		s.updateLocalWorkerID(s.WorkerId)
		go s.scheduledUploadData(c, zkAddr)
		return true
	}
	resp, err1 := c.Get(context.Background(), PATH_FOREVER)
	if err1 != nil {
		panic(err1)
	}
	if len(resp.Kvs) == 0 {
		if retryCount > 3 {
			return false
		}
		retryCount++
		goto RETRY
	}
	nodeMap := make(map[string]int, 0)
	realNodeMap := make(map[string]string, 0)
	for _, node := range resp.Kvs {
		nodeKey := strings.Split(string(node.Value), "-")
		realNodeMap[nodeKey[0]] = string(node.Value)
		nodeMap[nodeKey[0]] = cast.ToInt(nodeKey[1])
	}
	if workId, ok := nodeMap[s.listenAddress]; ok {
		etcdAddrNode := PATH_FOREVER + "/" + realNodeMap[s.listenAddress]
		s.WorkerId = workId
		if !s.checkInitTimeStamp(c, etcdAddrNode) {
			viper.SetConfigName("workerID")
			viper.AddConfigPath(conf.GetConfigPath())
			err := viper.ReadInConfig()
			if err != nil {
				panic(fmt.Errorf("read file error:%v", err))
			}
			logger.Errorf("START FAILED ,use local node file properties workerID-{%d}", viper.GetInt("workerID"))
			return false
		}
		s.etcdAddressNode = etcdAddrNode
		s.doService(c)
		s.updateLocalWorkerID(s.WorkerId)
		logger.Infof("[Old NODE]find forever node have this endpoint ip-{%s} port-{%s} workid-{%d} childnode and start SUCCESS", s.ip, s.port, s.WorkerId)
	} else {
		newNode := s.createNode(c)
		s.etcdAddressNode = newNode
		nodeKeyArr := strings.Split(newNode, "-")
		s.WorkerId = cast.ToInt(nodeKeyArr[1])
		s.doService(c)
		s.updateLocalWorkerID(workId)
		logger.Infof("[New NODE]can not find node on forever node that endpoint ip-{%s} port-{%s} workid-{%d},create own node on forever node and start SUCCESS", s.ip, s.port, s.WorkerId)
	}
	return true
}

func (s *SnowFlakeEtcdHolder) doService(client *clientv3.Client) {
	go s.scheduledUploadData(client, s.etcdAddressNode)
}

func (s *SnowFlakeEtcdHolder) scheduledUploadData(client *clientv3.Client, zkAddrNode string) {
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	for {
		select {
		case <-ticker.C:
			s.updateNewData(client, zkAddrNode)
		}
	}
}

func (s *SnowFlakeEtcdHolder) createNode(client *clientv3.Client) string {
	path := PATH_FOREVER + "/" + s.listenAddress + "-0000000000"
	_, err := client.Put(context.Background(), PATH_FOREVER, s.listenAddress+"-0000000000")
	if err != nil {
		return ""
	}
	_, err = client.Put(context.Background(), path, string(s.buildData()))
	if err != nil {
		logger.Infof("etcd CreateValErr:%+v", err)
		return ""
	}
	return path
}

func (s *SnowFlakeEtcdHolder) updateNewData(client *clientv3.Client, path string) {
	if timeutil.MsTimestampNow() < s.lastUpdateTime {
		return
	}
	_, err := client.Put(context.Background(), path, string(s.buildData()))
	if err != nil {
		return
	}
	s.lastUpdateTime = timeutil.MsTimestampNow()
	return
}

func (s *SnowFlakeEtcdHolder) buildData() []byte {
	endPoint := new(Endpoint)
	endPoint.IP = s.ip
	endPoint.Port = s.port
	endPoint.Timestamp = timeutil.MsTimestampNow()
	encodeArr, _ := json.Marshal(endPoint)
	return encodeArr
}

func (s *SnowFlakeEtcdHolder) deBuildData(val []byte) *Endpoint {
	endPoint := new(Endpoint)
	_ = json.Unmarshal(val, endPoint)
	return endPoint
}

func (s *SnowFlakeEtcdHolder) checkInitTimeStamp(client *clientv3.Client, zkAddrNode string) bool {
	resp, err := client.Get(context.Background(), zkAddrNode)
	if err != nil {
		return false
	}
	endpoint := s.deBuildData(resp.Kvs[0].Value)
	return !(endpoint.Timestamp > timeutil.MsTimestampNow())
}

func (s *SnowFlakeEtcdHolder) updateLocalWorkerID(workId int) {
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

func (s *SnowFlakeEtcdHolder) GetWorkerId() int {
	return s.WorkerId
}
