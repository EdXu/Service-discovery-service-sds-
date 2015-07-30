/*******************************************************************************
 *
 *  文      件:    sdk_event.go
 *
 *  概述：事件响应主程序
 *
 *  版本历史
 *      1.0    2014-10-09    xufengping    创建并实现
 *
*******************************************************************************/
package sdssdk

import (
	"errors"
	"sync"
	"time"

	log4 "github.com/alecthomas/log4go"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

/******************************************************************************
* 概述:
* 类型名:     Event
* 成员列表：   成员名            成员类型       取值范围         描述
*            mzmq             Zmq                          Context
*            mcache           DataCache                    缓存数据类对象
*            mrun         	  bool                         用于停止异常恢复go程
*            msdsServicepath  string[]                       sds path
*	         mnowSocket       *ZmqSocket                   当前工作sds
*            mrwmutex         sync.RWMutex                 锁mnowSocket
*
******************************************************************************/
type Event struct {
	mzmq            Zmq
	mnowSocket      *ZmqSocket
	mcache          DataCache
	msdsServicepath string
	mrwmutex        sync.RWMutex
	mrun            bool
}

/******************************************************************************
* 概述：     Event初始化函数
* 函数名：    Init
* 返回值：    bool           true,false
* 参数列表：  参数名          参数类型      取值范围     描述
*			sdscntstring    []string                 req-Ip:port
*		    sdsServicepath  string                   sds路径
*
*******************************************************************************/
func (this *Event) Init(sdscntstring []string, sdsServicepath string) error {
	this.msdsServicepath = sdsServicepath
	this.mrun = true
	this.mcache.Init(&this.mzmq)
	if err := this.mzmq.Init(&this.mcache, this); err != nil {
		return errors.New("init mzmq err:" + err.Error())
	}
	i := 0
	for ; i < len(sdscntstring); i++ {
		if len(sdscntstring[i]) == 0 {
			continue
		}
		zmqReq := &ZmqReq{}
		if err := zmqReq.Init(sdscntstring[i], &this.mzmq); err != nil {
			zmqReq.Close()
			return errors.New("init MzmqReq err:" + err.Error())
		}
		err1, res := zmqReq.Req(sdsServicepath)
		zmqReq.Close()
		if err1 != nil || res == nil {
			log4.Debug("req sds From(%s) err:%s!So try next sdsAddr", sdscntstring[i], err1.Error())
			continue
		} else {
			this.mcache.UpdatePoints(res)
			break
		}
	}
	if i >= len(sdscntstring) {
		return errors.New("all sds req failed...")
	}
	if this.mcache.EmptySocketGroup() {
		return errors.New("no sds Socket create...")
	}
	log4.Debug("event init with sds num=%d", len(this.mcache.mSocketGroup))
	//	time.Sleep(time.Second * DEFAULT_TIMEWAIT)
	this.mcache.WaitSocketsState()

	if qs := this.mcache.GetSocket(); qs != nil {
		this.mnowSocket = qs
		this.mcache.UpdateSdsTopic(this.msdsServicepath, true)

		if err := this.mnowSocket.Register(this.msdsServicepath); err != nil {
			return errors.New("zmqSub register error:" + err.Error())
		}
	} else {
		return errors.New("nofind state ok of sds Socket...")
	}
	return nil
}

/******************************************************************************
* 概述：     Event停止函数
* 函数名：    Stop
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Event) Stop() {
	this.mrun = false
	//	this.mrwmutex.Lock()
	this.mnowSocket = nil
	this.mcache.Close()
	//	this.mrwmutex.UnLock()
	this.mzmq.Close()
}

/******************************************************************************
* 概述：     获取当前工作socket
* 函数名：    GetWorkSocket
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Event) GetWorkSocket() *ZmqSocket {
	this.mrwmutex.RLock()
	defer this.mrwmutex.RUnlock()
	return this.mnowSocket
}

/******************************************************************************
* 概述：     获取最佳节点信息，必要时请求数据
* 函数名：    GetServiceInfo
* 返回值：    error                       nil,error
			*dzhyun.SDSEndpoint          nil,*dzhyun.SDSEndpoint
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string                  路径
*
*******************************************************************************/
func (this *Event) GetServiceInfo(servicepath string) (error, *dzhyun.SDSEndpoint) {
	if !this.mcache.HaveSdsTopic(servicepath) {
		if err := this.reqZmqRegister(servicepath); err != nil {
			rerr := errors.New("reqZmqRegister(" + servicepath + ") err:" + err.Error())
			log4.Error(rerr.Error())
			return rerr, nil
		}
	}
	if point, ok := this.mcache.FindOne(servicepath); ok {
		return nil, point
	} else {
		rerr := errors.New("nofind(" + servicepath + ") point in cache")
		log4.Warn(rerr.Error())
		return rerr, nil
	}
}

/******************************************************************************
* 概述：     获取所有节点信息，必要时请求数据
* 函数名：    GetAllServiceInfo
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string                  路径
*
*******************************************************************************/
func (this *Event) GetAllServiceInfo(servicepath string) (error, *[]dzhyun.SDSEndpoint) {
	if !this.mcache.HaveSdsTopic(servicepath) {
		if err := this.reqZmqRegister(servicepath); err != nil {
			rerr := errors.New("reqZmqRegister(" + servicepath + ") err:" + err.Error())
			log4.Error(rerr.Error())
			return rerr, nil
		}
	}

	if points, ok := this.mcache.FindAll(servicepath); ok {
		return nil, points
	} else {
		rerr := errors.New("nofind(" + servicepath + ") point in cache")
		log4.Warn(rerr.Error())
		return rerr, nil
	}
}

/******************************************************************************
* 概述：     注册sdk订阅信息，必要时注册sds订阅信息
* 函数名：    RegistServiceInfo
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string                  路径
*
*******************************************************************************/
func (this *Event) RegistServiceInfo(servicepath string) (error, chan dzhyun.SDSEndpoint) {
	if !this.mcache.HaveSdsTopic(servicepath) {
		s := this.GetWorkSocket()
		if s == nil {
			log4.Error("no workSocket err")
			return errors.New("no workSocket err"), nil
		}
		if err := s.Register(servicepath); err != nil {
			rerr := errors.New("Regist(" + servicepath + ")failed, " + err.Error())
			log4.Error(rerr.Error())
			return rerr, nil
		}
		this.mcache.UpdateSdsTopic(servicepath, true)
	}
	return nil, this.mcache.ApplyChan(servicepath)
}

/******************************************************************************
* 概述：     注销sdk订阅信息
* 函数名：    UnRegistServiceInfo
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string                  路径
*           ch             chan dzhyun.SDSEndpoint
*
*******************************************************************************/
func (this *Event) UnRegistServiceInfo(servicepath string, ch chan dzhyun.SDSEndpoint) {
	this.mcache.DelChan(servicepath, ch)
}

/******************************************************************************
* 概述：     请求数据并订阅sds
* 函数名：    reqZmqRegister
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*           servicepath    string                  路径
*
*******************************************************************************/
func (this *Event) reqZmqRegister(servicepath string) error {
	workSocket := this.GetWorkSocket()
	if workSocket == nil {
		return errors.New("worksocket null")
	}
	err, res := workSocket.Req(servicepath)
	if err != nil || res == nil {
		return errors.New("reqzmq falied, " + err.Error())
	}
	this.mcache.UpdatePoints(res)
	if err := workSocket.Register(res.GetSubName()); err != nil {
		return errors.New("Register failed, " + err.Error())
	}
	this.mcache.UpdateSdsTopic(servicepath, true)

	return nil
}

/******************************************************************************
* 概述：     重新选择合适的worksocket(worksocket断开连接时调用)
* 函数名：    UpdateWorkSocket
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Event) UpdateWorkSocket() {
	this.mrwmutex.Lock()
	this.mnowSocket = nil
	this.mrwmutex.Unlock()
	var s *ZmqSocket
	go func() {
		for {
			s = this.mcache.GetSocket()
			if s != nil {
				break
			}
			if !this.mrun {
				return //结束
			}
			log4.Debug("wana find ok sdsSocket, but nofind")
			time.Sleep(time.Second * DEFAULT_TIMEWAIT4)
		}
		topics := this.mcache.GetSdsTopic()
		for i := 0; topics != nil && i < len(*topics); i++ {
			if err0 := s.Register((*topics)[i]); err0 != nil {
				log4.Error("register(%s) failed, %s ", (*topics)[i], err0.Error())
			}
			err, res := s.Req((*topics)[i])
			if err != nil || res == nil {
				log4.Error("req(%s) falied, %s", (*topics)[i], err.Error())
				continue
			}
			this.mcache.RefleshPoints(res.GetSubName(), res)
		}

		this.mrwmutex.Lock()
		defer this.mrwmutex.Unlock()
		this.mnowSocket = s

	}()
}
