/*******************************************************************************
 *
 *  文      件:    sdk_socket.go
 *
 *  概述：         网络通信
 *
 *  版本历史
 *      1.0    2014-10-09    xufengping    创建并实现
 *
*******************************************************************************/
package sdssdk

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	log4 "github.com/alecthomas/log4go"
	zmq4 "github.com/pebbe/zmq4"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
	common "gw.com.cn/dzhyun/utils.git/sdsutils"
)

/***************************Zmq************************************************
* 概述:       Zmq类，所有zmq共用一个该对象
* 类型名:     Zmq
* 成员列表：   成员名            成员类型       取值范围       描述
*			 Mcontext         *zmq4.Context
*            MdataCache       *DataCache
*	         Mevent           *Event
*
******************************************************************************/
type Zmq struct {
	Mcontext   *zmq4.Context
	MdataCache *DataCache
	Mevent     *Event
}

/******************************************************************************
* 概述：     Zmq初始化
* 函数名：    Init
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Zmq) Init(dataCache *DataCache, event *Event) error {
	this.MdataCache = dataCache
	this.Mevent = event
	if this.Mcontext != nil {
		return nil
	}
	context, err := zmq4.NewContext()
	if err != nil {
		return err
	}
	this.Mcontext = context
	return nil
}

/******************************************************************************
* 概述：     Zmq关闭
* 函数名：    Close
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Zmq) Close() bool {
	if this.Mcontext != nil {
		if err := this.Mcontext.Term(); err != nil {
			switch zmq4.AsErrno(err) {
			case zmq4.Errno(syscall.EINTR):
				this.Mcontext = nil
				return false
			case zmq4.ETERM:
				this.Mcontext = nil
				return true
			case zmq4.Errno(syscall.EFAULT):
				this.Mcontext = nil
				return true
			}
		}
		this.Mcontext = nil
	}
	log4.Info("conetxt term ok")
	return true
}

/***************************Zmq************************************************
* 概述:       event初始化时调用
* 类型名:     ZmqReq
* 成员列表：   成员名            成员类型       取值范围       描述
			 mzmq              *zmq4.Context
			 mreqSocket        *zmq4.Socket
*
******************************************************************************/
type ZmqReq struct {
	mzmq       *Zmq
	mreqSocket *zmq4.Socket
}

/******************************************************************************
* 概述：     ZmqReq初始化
* 函数名：    Init
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqReq) Init(requrl string, zmq *Zmq) error {
	this.mzmq = zmq
	if socket, err := this.mzmq.Mcontext.NewSocket(zmq4.REQ); err != nil {
		return err
	} else {
		this.mreqSocket = socket
		if err := this.mreqSocket.SetRcvtimeo(DEFAULT_RECVTIMEOUT1); err != nil {
			rerr := errors.New("SetRcvtimeo " + err.Error())
			return rerr
		}
		if err := this.mreqSocket.SetLinger(0); err != nil {
			rerr := errors.New("SetLinger " + err.Error())
			return rerr
		}
	}
	if err := this.mreqSocket.Connect(requrl); err != nil {
		rerr := errors.New("Connect " + err.Error())
		return rerr
	}
	return nil
}

/******************************************************************************
* 概述：     ZmqReq关闭
* 函数名：    Close
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqReq) Close() {
	if this.mreqSocket != nil {
		this.mreqSocket.Close()
		this.mreqSocket = nil
	}
}

/******************************************************************************
* 概述：     ZmqReq 请求
* 函数名：    Req
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqReq) Req(servicepath string) (error, *dzhyun.SDSResponse) {
	stream, err0 := common.WrapRequest(DEFAULT_Version, servicepath, true)
	if err0 != nil {
		return err0, nil
	}
	if _, er1 := this.mreqSocket.SendBytes(stream, zmq4.DONTWAIT); er1 != nil {
		return er1, nil
	}
	stream2, err2 := this.mreqSocket.RecvBytes(0) //
	if err2 != nil {
		log4.Debug("zmqreq %s", err2.Error())
		return err2, nil
	}
	if len(stream2) == 0 {
		log4.Debug("zmqreq timeout")
		return errors.New("req timeout"), nil
	}

	frm, err3 := common.UnwrapBaseProto(stream2)
	if err3 != nil {
		return err3, nil
	}
	mid := common.UnOffset(*frm.GetBody().Mid, 4)
	if mid[4-1] == 200 {
		return errors.New("failed protocal"), nil
	}
	res, err4 := common.UnwrapResponse(frm.GetBody().Mdata)
	if err4 != nil {
		return err4, nil
	}
	return nil, res
}

/*************************ZmqReq***********************************************
* 概述:       ZmqReq类
* 类型名:     ZmqReq
* 成员列表：   成员名            成员类型       取值范围       描述
*            mzmq              *Zmq
*            mreqSocket        *zmq4.Socket
*            msubSocket        *zmq4.Socket
*            mok               int           0，1，2       网络已连接的个数
*            mreqOK            bool                        req连接是否正常
*          	 msubOK            bool                        sub连接是否正常
*            mChoose           bool                        是否被选为工作socket
*            mreqUrl           string
*	         msubUrl           string
*	         mcmdSoc           *zmq4.Socket  Stop及topic   订阅退订及go退出
*            mpoller           *zmq4.Poller
*	         mmons1            *zmq4.Socket                MonitorReq
*	         mmons2            *zmq4.Socket                MonitorSub
*	         mreqMutex         sync.Mutex                  Req
*	         mcmdMutex         sync.Mutex                  mcmdSoc
*	         mchooseMutex      sync.RWMutex                mok及mChoose
*	         mserial           int                         编号，monitor时有用
*
******************************************************************************/
type ZmqSocket struct {
	mzmq         *Zmq
	mreqSocket   *zmq4.Socket
	msubSocket   *zmq4.Socket
	mok          int
	mreqOK       bool
	msubOK       bool
	mChoose      bool
	mreqUrl      string
	msubUrl      string
	mcmdSoc      *zmq4.Socket
	mpoller      *zmq4.Poller
	mmons1       *zmq4.Socket
	mmons2       *zmq4.Socket
	mreqMutex    sync.Mutex
	mcmdMutex    sync.Mutex
	mchooseMutex sync.RWMutex
	mserial      int
	mstateCh     chan bool
	mstateChFlag bool
}

/******************************************************************************
* 概述：     ZmqSocket初始化
* 函数名：    Init
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) Init(requrl, subUrl string, zmq *Zmq, serial int, stateCh chan bool) error {
	this.mzmq = zmq
	this.mreqUrl = requrl
	this.msubUrl = subUrl
	this.mserial = serial
	this.mstateCh = stateCh
	this.mstateChFlag = true
	this.mpoller = zmq4.NewPoller()

	this.mok = 0
	this.mChoose = false
	this.mreqOK = false
	this.msubOK = false
	if socket, err := this.mzmq.Mcontext.NewSocket(zmq4.REQ); err != nil {
		log4.Error("NewSocket Failed,%s", err.Error())
		return err
	} else {
		this.mreqSocket = socket
		if err1 := this.mreqSocket.SetRcvtimeo(DEFAULT_RECVTIMEOUT1); err1 != nil {
			rerr := errors.New("ReqSocket SetRcvtimeo " + err1.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err2 := this.mreqSocket.SetLinger(0); err2 != nil {
			rerr := errors.New("ReqSocket SetLinger " + err2.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err3 := this.mreqSocket.SetSndtimeo(DEFAULT_RECVTIMEOUT1); err3 != nil {
			rerr := errors.New("ReqSocket SetSndtimeo " + err3.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err4 := this.mreqSocket.SetReconnectIvl(time.Second * DEFAULT_TIMEWAIT3); err4 != nil {
			rerr := errors.New("ReqSocket SetReconnectIvl " + err4.Error())
			log4.Error(rerr.Error())
			return rerr
		}
	}

	if socket, err := this.mzmq.Mcontext.NewSocket(zmq4.SUB); err != nil {
		log4.Error("NewSocket Failed,%s", err.Error())
		return err
	} else {
		this.msubSocket = socket
		if err1 := this.msubSocket.SetRcvtimeo(0); err1 != nil {
			rerr := errors.New("SubSocket SetRcvtimeo " + err1.Error())
			log4.Error(rerr.Error())
			return rerr
		}

		if err2 := this.msubSocket.SetLinger(0); err2 != nil {
			rerr := errors.New("SubSocket SetLinger " + err2.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err4 := this.msubSocket.SetReconnectIvl(time.Second * DEFAULT_TIMEWAIT3); err4 != nil {
			rerr := errors.New("ReqSocket SetReconnectIvl " + err4.Error())
			log4.Error(rerr.Error())
			return rerr
		}
	}

	if s, err := this.mzmq.Mcontext.NewSocket(zmq4.PAIR); err != nil {
		log4.Error("NewSocket Failed,%s", err.Error())
		return err
	} else {
		this.mcmdSoc = s
		if err1 := this.mcmdSoc.SetRcvtimeo(-1); err1 != nil {
			rerr := errors.New("Pair SetRcvtimeo " + err1.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err1 := this.mcmdSoc.SetSndtimeo(-1); err1 != nil {
			rerr := errors.New("Pair SetSndtimeo " + err1.Error())
			log4.Error(rerr.Error())
			return rerr
		}
		if err2 := this.mcmdSoc.SetLinger(0); err2 != nil {
			rerr := errors.New("Pair SetLinger " + err2.Error())
			log4.Error(rerr.Error())
			return rerr
		}
	}
	if err := this.Monitor(); err != nil {
		log4.Debug("monitor err:%s", err.Error())
		return err
	}
	if err := this.mreqSocket.Connect(requrl); err != nil {
		rerr := errors.New("Req Connect " + err.Error())
		log4.Error(rerr.Error())
		return rerr
	}
	if err := this.msubSocket.Connect(subUrl); err != nil {
		rerr := errors.New("Sub Connect " + err.Error())
		log4.Error(rerr.Error())
		return rerr
	}

	if err := this.run(); err != nil {
		return err
	}
	log4.Debug("init socket ok %s %s", this.mreqUrl, this.msubUrl)
	return nil
}

/******************************************************************************
* 概述：    ZmqSocket停止函数
* 函数名：    close()
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) Close() {
	//	log4.Debug("start zmqSocket Close******")
	this.mok = 0
	this.mreqOK = false
	this.msubOK = false
	this.mChoose = false
	this.mcmdMutex.Lock()
	if this.mcmdSoc != nil {
		this.mcmdSoc.Send("stop ", 0)
		this.mcmdSoc.Recv(0)
		this.mcmdSoc.Close()
		this.mcmdSoc = nil
	}
	this.mcmdMutex.Unlock()

	if this.mreqSocket != nil {
		this.mreqSocket.Close()
		this.mreqSocket = nil
	}
	if this.msubSocket != nil {
		this.msubSocket.Close()
		this.msubSocket = nil
	}
	if this.mmons1 != nil {
		this.mmons1.Close()
		this.mmons1 = nil
	}
	if this.mmons2 != nil {
		this.mmons2.Close()
		this.mmons2 = nil
	}
}

/******************************************************************************
* 概述：     事件监听初始化
* 函数名：    Monitor
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) Monitor() error {

	addr1 := fmt.Sprintf("inproc://monitor.rep%d_%d", os.Getpid(), this.mserial)
	if err := this.mreqSocket.Monitor(addr1, zmq4.EVENT_CLOSED|zmq4.EVENT_DISCONNECTED|zmq4.EVENT_CONNECTED); err != nil {
		log4.Error("Req Monitor(%s) Failed,%s", addr1, err.Error())
		return err
	}
	s1, err0 := this.mzmq.Mcontext.NewSocket(zmq4.PAIR)
	if err0 != nil {
		log4.Error("NewSocket Failed,%s", err0.Error())
		return err0
	}
	if err := s1.SetRcvtimeo(0); err != nil {
		log4.Error("SetRcvtimeo Failed,%s", err.Error())
		return err
	}
	if err := s1.Connect(addr1); err != nil {
		log4.Error("Connect Failed,%s", err.Error())
		return err
	}
	this.mmons1 = s1
	if !strings.Contains(runtime.GOOS, "windows") {
		this.mpoller.Add(this.mmons1, zmq4.POLLIN)
	} else {
		this.mreqOK = true
	}

	addr2 := fmt.Sprintf("inproc://monitor.sub%d_%d", os.Getpid(), this.mserial)
	if err := this.msubSocket.Monitor(addr2, zmq4.EVENT_CLOSED|zmq4.EVENT_DISCONNECTED|zmq4.EVENT_CONNECTED); err != nil {
		log4.Error("Sub Monitor(%s) Failed,%s", addr2, err.Error())
		return err
	}
	s2, err1 := this.mzmq.Mcontext.NewSocket(zmq4.PAIR)
	if err1 != nil {
		log4.Error("NewSocket Failed,%s", err1.Error())
		return err1
	}
	if err := s2.SetRcvtimeo(0); err != nil {
		log4.Error("SetRcvtimeo Failed,%s", err.Error())
		return err
	}
	if err := s2.Connect(addr2); err != nil {
		log4.Error("Connect Failed,%s", err.Error())
		return err
	}
	this.mmons2 = s2
	if !strings.Contains(runtime.GOOS, "windows") {
		this.mpoller.Add(this.mmons2, zmq4.POLLIN)
	} else {
		this.msubOK = true
	}

	return nil
}

/******************************************************************************
* 概述：     ZmqSocket的请求函数
* 函数名：    Req
* 返回值：    error
*           *dzhyun.SDSResponse                   响应
* 参数列表：  参数名          参数类型      取值范围     描述
*			servicepath   string                  节点路径
*
*******************************************************************************/
func (this *ZmqSocket) Req(servicepath string) (error, *dzhyun.SDSResponse) {
	this.mchooseMutex.RLock()
	if !this.mChoose {
		this.mchooseMutex.RUnlock()
		return errors.New("ZmqSocket not workSocket"), nil
	}
	this.mchooseMutex.RUnlock()

	stream, err0 := common.WrapRequest(DEFAULT_Version, servicepath, true)
	if err0 != nil {
		return err0, nil
	}

	this.mreqMutex.Lock()
	if _, er1 := this.mreqSocket.SendBytes(stream, zmq4.DONTWAIT); er1 != nil {
		this.mreqMutex.Unlock()
		return er1, nil
	}
	stream2, err2 := this.mreqSocket.RecvBytes(0) //
	if err2 != nil {
		this.mreqMutex.Unlock()
		log4.Debug("req timeout %s", err2.Error())
		return err2, nil
	}
	if len(stream2) == 0 {
		this.mreqMutex.Unlock()
		log4.Debug("req timeout")
		return errors.New("req timeout"), nil
	}
	this.mreqMutex.Unlock()

	frm, err3 := common.UnwrapBaseProto(stream2)
	if err3 != nil {
		return err3, nil
	}
	mid := common.UnOffset(*frm.GetBody().Mid, 4)
	if mid[4-1] == 200 {
		return errors.New("failed protocal"), nil
	}
	res, err4 := common.UnwrapResponse(frm.GetBody().Mdata)
	if err4 != nil {
		return err4, nil
	}
	return nil, res
}

/******************************************************************************
* 概述：    ZmqSocket订阅函数
* 函数名：    Register
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*           topic          string                  sds topic
*
*******************************************************************************/
func (this *ZmqSocket) Register(topic string) error {

	this.mchooseMutex.RLock()
	if !this.mChoose {
		this.mchooseMutex.RUnlock()
		return errors.New("ZmqSocket not workSocket")
	}
	this.mchooseMutex.RUnlock()

	this.mcmdMutex.Lock()
	defer this.mcmdMutex.Unlock()
	if this.mcmdSoc == nil {
		return nil
	}
	if _, err := this.mcmdSoc.Send("r"+topic, 0); err != nil {
		return err
	}
	//	log4.Debug("cmdsoc send ok")
	return nil
}

/******************************************************************************
* 概述：    ZmqSocket注销订阅函数
* 函数名：    UnRegister
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*           topic          string                  sds topic
*
*******************************************************************************/
func (this *ZmqSocket) UnRegister(topic string) error {
	this.mchooseMutex.RLock()
	if !this.mChoose {
		this.mchooseMutex.RUnlock()
		return errors.New("ZmqSocket not workSocket")
	}
	this.mchooseMutex.RUnlock()

	this.mcmdMutex.Lock()
	defer this.mcmdMutex.Unlock()
	if this.mcmdSoc == nil {
		return nil
	}
	if _, err := this.mcmdSoc.Send("u"+topic, 0); err != nil {
		return err
	}
	return nil
}

/******************************************************************************
* 概述：    订阅退订，事件监控，接收推送数据
* 函数名：    run
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) run() error {
	addr0 := fmt.Sprintf("inproc://mon.ipc%d_%d", os.Getpid(), this.mserial)
	cmdSocket, err := this.mzmq.Mcontext.NewSocket(zmq4.PAIR)
	if err != nil {
		log4.Error("newsocket failed, %s", err.Error())
		return err
	}
	if err := cmdSocket.SetRcvtimeo(0); err != nil {
		log4.Error("SetRecvTimeO falied, %s", err.Error())
		return err
	}
	if err := cmdSocket.Bind(addr0); err != nil {
		log4.Error("Bind falied, %s", err.Error())
		return err
	}
	if err := this.mcmdSoc.Connect(addr0); err != nil {
		log4.Error("Connect falied, %s", err.Error())
		return err
	}

	this.mpoller.Add(cmdSocket, zmq4.POLLIN)
	this.mpoller.Add(this.msubSocket, zmq4.POLLIN)
	go func() {
		defer cmdSocket.Close()
		msgHeadFlag := true
		for {
			sockets, _ := this.mpoller.Poll(-1)
			for _, socket := range sockets {
				switch s := socket.Socket; s {
				case cmdSocket:
					flag, _ := this.dealCmd(s)
					if flag == 1 {
						log4.Debug("zmqSocket stop")
						return
					}
				case this.mmons1:
					this.dealMonitorEvent(s, true)
				case this.mmons2:
					this.dealMonitorEvent(s, false)
				case this.msubSocket:
					this.dealSub(s, &msgHeadFlag)
				}
			}
		}
	}()
	return nil
}

/******************************************************************************
* 概述：    订阅退订
* 函数名：    dealCmd
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) dealCmd(cmdSocket *zmq4.Socket) (int, error) {
	log4.Debug("start deal cmd...")
	for { //半包处理
		cmdsStr, err0 := cmdSocket.Recv(zmq4.DONTWAIT)

		if err0 != nil {
			errno1 := zmq4.AsErrno(err0)
			switch errno1 {
			case zmq4.Errno(syscall.EAGAIN):
				return 0, nil
			case zmq4.Errno(syscall.EINTR):
				continue
			default:
				log4.Debug("zmq req Get err %v, %d!", errno1, errno1)
			}
		}
		if len(cmdsStr) == 0 {
			log4.Debug("deal cmd return")
			return 0, nil
		}
		ss := strings.Split(cmdsStr, " ")
		log4.Debug("recv cmd %s", cmdsStr)
		for i := 0; i < len(ss); i++ {
			if len(ss[i]) == 0 {
				continue
			}
			if ss[i] == "stop" {
				//				log4.Debug("recv cmd will stop %s %s", this.mreqUrl, this.msubUrl)
				cmdSocket.Send("0", 0)
				return 1, nil
			}
			if !this.mChoose {
				//				log4.Debug("recv cmd ,but notChoose so return")
				return 0, nil
			}
			if ss[i][0] == 'r' {
				if err := this.msubSocket.SetSubscribe(ss[i][1:]); err != nil {
					log4.Error("SetSubscribe(%s) falied, %s", ss[i][1:], err.Error())
					return 0, err
				}
				log4.Debug("setSubscribe ok %s", ss[i][1:])
				continue
			}
			if ss[i][0] == 'u' {
				if err := this.msubSocket.SetUnsubscribe(ss[i][1:]); err != nil {
					log4.Error("SetUnSubscribe(%s) falied, %s", ss[i][1:], err.Error())
					return 0, err
				}
				log4.Debug("setUnSubscribe ok %s", ss[i][1:])
				continue
			}
		}
	}
	return 0, nil
}

/******************************************************************************
* 概述：    接收处理推送信息
* 函数名：    dealSub
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) dealSub(soc *zmq4.Socket, msgHeadFlag *bool) error {
	//	log4.Debug("start deal sub0...")
	//	this.mchooseMutex.RLocker()
	if !this.mChoose {
		//	this.mchooseMutex.RUnlock()
		return nil
	}
	log4.Debug("start deal sub...")
	//	this.mchooseMutex.RUnlock()
	for {
		if *msgHeadFlag {
			topic, err := soc.Recv(zmq4.DONTWAIT)
			if err != nil {
				errno1 := zmq4.AsErrno(err)
				switch errno1 {
				case zmq4.Errno(syscall.EAGAIN):
					return nil
				case zmq4.Errno(syscall.EINTR):
					continue
				default:
					log4.Debug("zmq req Get err %v, %d!", errno1, errno1)
				}
			}
			if len(topic) == 0 {
				//				log4.Debug("sub return1")
				return nil
			}
			//			log4.Debug("recv sub head =", topic)
			*msgHeadFlag = false
		} else {
			stream, err2 := soc.RecvBytes(zmq4.DONTWAIT)
			if err2 != nil {
				errno1 := zmq4.AsErrno(err2)
				switch errno1 {
				case zmq4.Errno(syscall.EAGAIN):
					return nil
				case zmq4.Errno(syscall.EINTR):
					continue
				default:
					log4.Debug("zmq req Get err %v, %d!", errno1, errno1)
				}
			}
			if len(stream) == 0 {
				//				log4.Debug("sub return2")
				return nil
			}
			*msgHeadFlag = true

			frm, err3 := common.UnwrapBaseProto(stream)
			if err3 != nil {
				log4.Error("UnwrapBaseProto falied, %s", err3.Error())
				return err3
			}

			mid := common.UnOffset(*frm.GetBody().Mid, 4)
			if mid[4-1] == 200 {
				log4.Error("sdssdk zmqsub get mid == 200 err")
				continue
			}
			res, err4 := common.UnwrapResponse(frm.GetBody().Mdata)
			if err4 != nil {
				log4.Error("sdssdk sub UnwrapResponse error:", err4)
				continue
			}
			this.mzmq.MdataCache.UpdatePoints(res)
		}
	}
	return nil
}

/******************************************************************************
* 概述：     事件处理
* 函数名：    dealMonitorEvent
* 返回值：    error
* 参数列表：  参数名          参数类型      取值范围     描述
*            isReq          bool                     是否为req事件
*
*******************************************************************************/
func (this *ZmqSocket) dealMonitorEvent(s *zmq4.Socket, isReq bool) error {
	log4.Debug("start dealMonitorEvent...")
	for {
		a, b, c, err := s.RecvEvent(0)
		//if runtime.GOOS == "windows" && len(b) == 0 && c == 0 {
		//	log4.Debug("monitor eagan windows")
		//	return nil
		//}
		if err != nil {
			errno1 := zmq4.AsErrno(err)
			switch errno1 {
			case zmq4.Errno(syscall.EAGAIN):
				log4.Debug("monitor eagan ")
				return nil
			case zmq4.Errno(syscall.EINTR):
				log4.Debug("monitor EINTR")
				continue
			default:
				log4.Debug("zmq req Get err %v, %d!", errno1, errno1)
			}
		}
		if a == 0 {
			//			log4.Debug("monitor return")
			return nil
		}
		switch a {
		case zmq4.EVENT_CONNECTED:
			log4.Info("sub or req monitor event CONNECTED, url:%s %s", this.mreqUrl, this.mreqUrl)
			this.mchooseMutex.Lock()
			defer this.mchooseMutex.Unlock()
			if isReq {
				this.mreqOK = true
			} else {
				this.msubOK = true
			}
			if this.mstateChFlag && this.mreqOK && this.msubOK {
				this.mstateChFlag = false
				select {
				case this.mstateCh <- true:
				default:
				}
			}
			//			this.mok++
		case zmq4.EVENT_DISCONNECTED:
			log4.Error("sub or req monitor event DISCONNECTED, url:%s %s", this.mreqUrl, this.mreqUrl)
			this.mchooseMutex.Lock()
			defer this.mchooseMutex.Unlock()
			if isReq {
				this.mreqOK = false
			} else {
				this.msubOK = false
			}
			//			this.mok--
			if this.mChoose {
				this.mChoose = false
				topics := this.mzmq.MdataCache.GetSdsTopic()
				for i := 0; i < len(*topics); i++ {
					if err := this.msubSocket.SetUnsubscribe((*topics)[i]); err != nil {
						log4.Error("SetUnSubscribe(%s) falied, %s", (*topics)[i], err.Error())
					}
				}
				this.mzmq.Mevent.UpdateWorkSocket()
			}
		case zmq4.EVENT_CLOSED:
			//			this.mchooseMutex.Lock()
			//			this.mok
			//			if this.mChoose{
			//				this.mChoose = false
			//				this.mevent.UpdateWorkSocket()
			//			log4.Error("sdssdk zmqreq monitor event CLOSED", b, c)
			//			}
			//			this.mchooseMutex.UnLock()
		default:
			log4.Debug("zmqreq monitor unknow event err", a, b, c, err)
		}
	}
	return nil

}

/******************************************************************************
* 概述：    选中为工作socket
* 函数名：    SetChoose
* 返回值：    bool           true被选中，false未符合条件
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *ZmqSocket) SetChoose() bool {
	this.mchooseMutex.Lock()
	defer this.mchooseMutex.Unlock()
	//	if this.mok == 2 {
	if this.mreqOK && this.msubOK {
		this.mChoose = true
		return true
	}
	return false
}
