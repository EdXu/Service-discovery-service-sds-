/*******************************************************************************
 *
 *  文      件:    sdk_cache.go
 *
 *  概述： 缓存数据
 *
 *  版本历史
 *      1.0    2014-10-09    xufengping    创建并实现
 *
*******************************************************************************/
package sdssdk

import (
	"errors"
	"time"
	log4 "github.com/alecthomas/log4go"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

import "math/rand"

import "sync"

/******************************************************************************
* 概述:       缓存数据类
* 类型名:     DataCache
* 成员列表：   成员名            成员类型      取值范围      描述
*            mpointGroup       Tree                     节点树及chan树
*            msdsTopic         map[string]int           sds主题
*            mSocketGroup      map[string]*ZmqSocket    所有sds连接
*			 mzmq              *Zmq
*	         msocketSerial     int                      给sds编号
*	         mtopicRWMutex     *sync.RWMutex            msdsTopic
*	         msocMutex         *sync.Mutex			    mSocketGroup
*	         mregRWMutex       sync.RWMutex             向用户推送及关闭chan时加锁
*
******************************************************************************/
type DataCache struct {
	mpointGroup     Tree
	msdsTopic       map[string]int //servicepath = subname
	mSocketGroup    map[string]*ZmqSocket
	mSocketsStateCh chan bool
	mtopicRWMutex   sync.RWMutex
	msocMutex       sync.Mutex
	mregRWMutex     sync.RWMutex
	mzmq            *Zmq
	msocketSerial   int
}

/******************************************************************************
* 概述：     DataCache初始化函数
* 函数名：    Init
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *DataCache) Init(pzmq *Zmq) {
	this.mzmq = pzmq
	rand.Seed(time.Now().Unix())
	this.mpointGroup.Init()
	this.msdsTopic = make(map[string]int, DEFAULT_InitTopicSize)
	this.mSocketGroup = make(map[string]*ZmqSocket, DEFAULT_InitSocketSize)
	this.mSocketsStateCh = make(chan bool, DEFAULT_ChanSocStateSize)
	this.msocketSerial = 0
	log4.Debug("init datacache ok")
}

/******************************************************************************
* 概述：     DataCache停止函数，关闭socket和销毁树
* 函数名：    Close
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *DataCache) Close() {
	this.msocMutex.Lock()
	defer this.msocMutex.Unlock()
	for key, s := range this.mSocketGroup {
		delete(this.mSocketGroup, key)
		s.Close()
	}
	this.mregRWMutex.Lock()
	defer this.mregRWMutex.Unlock()
	this.mpointGroup.Destory()
	close(this.mSocketsStateCh)
}

/******************************************************************************
* 概述：     DataCache查找某一个endpoint
* 函数名：    FindOne
* 返回值：    *dzhyun.SDSEndpoint       *Endpoint, nil
*           bool                       true,false
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string
*
*******************************************************************************/
func (this *DataCache) FindOne(servicepath string) (*dzhyun.SDSEndpoint, bool) {
	pp := this.mpointGroup.SearchBest(servicepath, DEFAULT_MinNodeNameNum)
	if pp == nil {
		return nil, false
	}
	return pp, true
}

/******************************************************************************
* 概述：     DataCache查找所有endpoint
* 函数名：    FindAll
* 返回值：    *[]dzhyun.SDSEndpoint       *[]Endpoint, nil
*           bool                       true,false
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string
*
*******************************************************************************/
func (this *DataCache) FindAll(servicepath string) (*[]dzhyun.SDSEndpoint, bool) {
	pArr := this.mpointGroup.SearchAll(servicepath)
	if pArr == nil {
		return nil, false
	}
	return pArr, true
}

/******************************************************************************
* 概述：     请求chan并订阅
* 函数名：    ApplyChan
* 返回值：    chan dzhyun.SDSEndpoint       chan Endpoint
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string
*
*******************************************************************************/
func (this *DataCache) ApplyChan(servicepath string) chan dzhyun.SDSEndpoint {
	ch := make(chan dzhyun.SDSEndpoint, DEFAULT_ChanBufferSize)
	this.mpointGroup.InsertChan(servicepath, ch)
	return ch
}

/******************************************************************************
* 概述：     删除chan并取消订阅
* 函数名：    DelChan
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string
*
*******************************************************************************/
func (this *DataCache) DelChan(servicepath string, ch chan dzhyun.SDSEndpoint) {
	if ch == nil {
		return
	}
	this.mpointGroup.DelChan(servicepath, ch)
	this.mregRWMutex.Lock()
	defer this.mregRWMutex.Unlock()
	close(ch)
}

/******************************************************************************
* 概述：     更新树的部分节点，如果是sds的就创建socket
* 函数名：    UpdatePoints
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*            res            *dzhyun.SDSResponse
*
*******************************************************************************/
func (this *DataCache) UpdatePoints(res *dzhyun.SDSResponse) {
	if res.Endpoints == nil {
		return
	}
	ps := &res.Endpoints
	subName := res.GetSubName()
	mainregChans := make([]chan dzhyun.SDSEndpoint, 0)

	this.mregRWMutex.RLock()
	defer this.mregRWMutex.RUnlock()
	this.mpointGroup.GetRegChans(&mainregChans, subName)

	for i := 0; i < len(*ps); i++ {
		subregChans := make([]chan dzhyun.SDSEndpoint, 0)
		this.insertSocket((*ps)[i])
		if (*ps)[i].GetState() == DEFAULT_PointStateOk {
			this.mpointGroup.InsertPoint((*ps)[i], &subregChans, subName)
		} else {
			this.mpointGroup.DelPoint((*ps)[i], &subregChans, subName)
		}
		for j := 0; j < len(mainregChans); j++ {
			select {
			case mainregChans[j] <- *(*ps)[i]:
			default:
				{
					log4.Error("chan full, so lose data")
				}
			}
		}
		for k := 0; k < len(subregChans); k++ {
			select {
			case subregChans[k] <- *(*ps)[i]:
			default:
				{
					log4.Error("chan full, so lose data")
				}
			}
		}
	}
}

/******************************************************************************
* 概述：     刷新树的一个分支，要与sds一致
* 函数名：    RefleshPoints
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*            Servicepath    string
*            res            *dzhyun.SDSResponse
*
*******************************************************************************/
func (this *DataCache) RefleshPoints(Servicepath string, res *dzhyun.SDSResponse) {
	oldppoint := this.mpointGroup.SearchAll(Servicepath)
	var newppoint *[]*dzhyun.SDSEndpoint
	if res.Endpoints == nil {
		newppoint = nil
	} else {
		newppoint = &res.Endpoints
	}

	tmpMap := make(map[string]*dzhyun.SDSEndpoint)
	for j := 0; newppoint != nil && j < len(*newppoint); j++ {
		this.insertSocket((*newppoint)[j])
		tmpMap[(*newppoint)[j].GetBusiPath()+"/"+(*newppoint)[j].GetNodeName()] = (*newppoint)[j]
	}

	for i := 0; oldppoint != nil && i < len(*oldppoint); i++ {
		str := (*oldppoint)[i].GetBusiPath() + "/" + (*oldppoint)[i].GetNodeName()
		if value, ok := tmpMap[str]; ok {
			if value.GetLoading() == (*oldppoint)[i].GetLoading() && value.GetState() == (*oldppoint)[i].GetState() {
				delete(tmpMap, str)
			}
		} else {
			if (*oldppoint)[i].State == nil {
				(*oldppoint)[i].State = new(int32)
			}
			*(*oldppoint)[i].State = DEFAULT_PointStateDel
			tmpMap[str] = &(*oldppoint)[i]
		}
	}

	subName := res.GetSubName()
	mainregChans := make([]chan dzhyun.SDSEndpoint, 0)

	this.mregRWMutex.RLock()
	defer this.mregRWMutex.RUnlock()
	this.mpointGroup.GetRegChans(&mainregChans, subName)

	for _, value := range tmpMap {
		subregChans := make([]chan dzhyun.SDSEndpoint, 0)
		if value.GetState() == DEFAULT_PointStateOk {
			this.mpointGroup.InsertPoint(value, &subregChans, subName)
		} else {
			this.mpointGroup.DelPoint(value, &subregChans, subName)
		}
		for j := 0; j < len(mainregChans); j++ {
			select {
			case mainregChans[j] <- *value:
			default:
				{
					log4.Error("chan full, so lose data")
				}
			}

		}
		for k := 0; k < len(subregChans); k++ {
			select {
			case subregChans[k] <- *value:
			default:
				{
					log4.Error("chan full, so lose data")
				}
			}
		}
	}

}

/******************************************************************************
* 概述：     更新sds topic
* 函数名：    UpdateSdsTopic
* 返回值：
* 参数列表：  参数名          参数类型      取值范围               描述
*           topics         string                            sdstopic
*           insert         bool        true(插入),false(删除)
*
*******************************************************************************/
func (this *DataCache) UpdateSdsTopic(topics string, insert bool) {
	this.mtopicRWMutex.Lock()
	defer this.mtopicRWMutex.Unlock()
	num, ok := this.msdsTopic[topics]
	if insert {
		if ok {
			this.msdsTopic[topics] = (num + 1)
		} else {
			this.msdsTopic[topics] = 1
		}
	} else {
		if ok {
			if num--; num > 0 {
				this.msdsTopic[topics] = num
			} else {
				delete(this.msdsTopic, topics)
			}
		}
	}
}

/******************************************************************************
* 概述：     获取所有topic
* 函数名：    GetSdsTopic
* 返回值：
* 参数列表：  参数名          参数类型      取值范围               描述
*
*******************************************************************************/
func (this *DataCache) GetSdsTopic() *[]string {
	arr := make([]string, 0)
	this.mtopicRWMutex.RLock()
	defer this.mtopicRWMutex.RUnlock()
	for key, _ := range this.msdsTopic {
		arr = append(arr, key)
	}

	return &arr
}

/******************************************************************************
* 概述：     是否已经有topic
* 函数名：    HaveSdsTopic
* 返回值：
* 参数列表：  参数名          参数类型      取值范围               描述
*
*******************************************************************************/
func (this *DataCache) HaveSdsTopic(path string) bool {
	this.mtopicRWMutex.RLock()
	defer this.mtopicRWMutex.RUnlock()
	if _, ok := this.msdsTopic[path]; ok {
		return true
	}
	return false
}

/******************************************************************************
* 概述：     判断是否有sds存在
* 函数名：    EmptySocketGroup
* 返回值：
* 参数列表：  参数名          参数类型      取值范围               描述
*
*******************************************************************************/
func (this *DataCache) EmptySocketGroup() bool {
	this.msocMutex.Lock()
	defer this.msocMutex.Unlock()
	if len(this.mSocketGroup) == 0 {
		return true
	}
	return false
}

/******************************************************************************
* 概述：     随机获取可用的sds
* 函数名：    GetSocket
* 返回值：
* 参数列表：  参数名          参数类型      取值范围               描述
*
*******************************************************************************/
func (this *DataCache) GetSocket() *ZmqSocket {
	this.msocMutex.Lock()
	defer this.msocMutex.Unlock()
	arrSize := len(this.mSocketGroup)
	if arrSize == 0 {
		return nil
	}
	randPos := rand.Int() % arrSize
	log4.Debug("rand pos %d arrSize %d", randPos, arrSize)

	//从randpos开始查找合适的socket，直到末尾
	i := 0
	for _, s := range this.mSocketGroup {
		if i != randPos {
			i++
			continue
		}
		if s.SetChoose() {
			return s
		}
	}

	//从开始位置开始查找合适的socket，直到randpos
	i = 0
	for _, s := range this.mSocketGroup {
		if i == randPos {
			break
		}
		if s.SetChoose() {
			return s
		}
	}
	return nil
}

/******************************************************************************
* 概述：     获取一组req和pub地址
* 函数名：   getSdsAddr
* 返回值：   string                                  reqaddr
*         	string					                subaddr
* 			error
* 参数列表：  参数名          参数类型      取值范围     描述
*            iface          string
*
*******************************************************************************/
func (this *DataCache) getSdsAddr(iface string) (string, string, error) {
	reqAddr, subAddr := "", ""

	var ok bool
	if reqAddr, ok = GetSubStr(iface, "REP:", " "); !ok {
		return "", "", errors.New("GetREPStr from Interface err")
	}
	if subAddr, ok = GetSubStr(iface, "PUB:", " "); !ok {
		return "", "", errors.New("GetPUBStr from Interface err")
	}
	return reqAddr, subAddr, nil
}

/******************************************************************************
* 概述：     判断并创建一个新的sds socket
* 函数名：    insertSocket
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*            pp             *dzhyun.SDSEndpoint
*
*******************************************************************************/
func (this *DataCache) insertSocket(pp *dzhyun.SDSEndpoint) {
	if pp == nil || pp.GetBusiPath() != DEFAULT_SdsServicepath {
		return
	}
	this.msocMutex.Lock()
	defer this.msocMutex.Unlock()
	if _, ok := this.mSocketGroup[(*pp).GetNodeName()]; !ok {
		newZmqSocket := &ZmqSocket{}
		if reqUrl, subUrl, err := this.getSdsAddr((*pp).GetInterface()); err == nil {
			this.msocketSerial++
			if err := newZmqSocket.Init(reqUrl, subUrl, this.mzmq, this.msocketSerial, this.mSocketsStateCh); err == nil {
				this.mSocketGroup[(*pp).GetNodeName()] = newZmqSocket
				log4.Debug("insertSocket size=%d nodename=%s %s %s", len(this.mSocketGroup), (*pp).GetNodeName(), reqUrl, subUrl)
			}
		}

	}
}

func (this *DataCache) WaitSocketsState() {
	timer := time.NewTimer(time.Second * DEFAULT_TIMEWAIT)
	for realNum := 0; true; {
		select {
		case <-this.mSocketsStateCh:
			realNum++
			if (realNum) == len(this.mSocketGroup) {
				return
			}
		case <-timer.C:
			return
		}
	}
}
