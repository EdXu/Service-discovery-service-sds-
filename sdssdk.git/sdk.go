/*******************************************************************************
 *
 *  文      件:    sdk.go
 *
 *  概述：         sdk接口
 *
 *  版本历史
 *      1.0    2014-10-09    xufengping    创建并实现
 *
*******************************************************************************/
package sdssdk

import (
	"errors"
	"strings"

	log4 "github.com/alecthomas/log4go"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

/******************************************************************************
* 概述:       sdk类
* 类型名:     SDSSDK
* 成员列表：   成员名            成员类型       取值范围         描述
*            mevent           Event
*
*******************************************************************************/
type SDSSDK struct {
	mevent Event
}

/******************************************************************************
* 概述：     sdssdk启动函数
* 函数名：    Start
* 返回值：    error          error,nil
* 参数列表：  参数名          参数类型      取值范围     描述
*			sdscntstrings   string                  req-Ip:port
*
*******************************************************************************/
func (sdk *SDSSDK) Start(sdscntstrings string) error {
	if len(sdscntstrings) == 0 {
		return errors.New("nil sdscntstrings")
	}
	sdscntstrings = strings.Replace(sdscntstrings, " ", "", len(sdscntstrings))
	addrs := strings.Split(sdscntstrings, ";")

	if err := sdk.mevent.Init(addrs, DEFAULT_SdsServicepath); err != nil {
		log4.Error("Start failed, %s", err.Error())
		return errors.New("Start failed " + err.Error())
	}
	log4.Info("start sdssdk ok...")
	return nil
}

/******************************************************************************
* 概述：     sdssdk停止函数
* 函数名：    Stop
* 返回值：    error          error,nil
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (sdk *SDSSDK) Stop() error {
	sdk.mevent.Stop()
	log4.Info("stop sdssdk ok...")
	return nil
}

/******************************************************************************
* 概述：     获取该路径下最优的一个节点
* 函数名：    GetServiceInfo
* 返回值：    error                      error,nil
*           dzhyun.SDSEndpoint
* 参数列表：  参数名          参数类型      取值范围     描述
*			servicepath    string                  路径
*
*******************************************************************************/
func (sdk *SDSSDK) GetServiceInfo(servicepath string) (error, dzhyun.SDSEndpoint) {
	if len(servicepath) == 0 || servicepath[0] != '/' {
		return errors.New("wrong servicepath"), dzhyun.SDSEndpoint{}
	}
	err, ppoint := sdk.mevent.GetServiceInfo(servicepath)
	if ppoint == nil {
		return err, dzhyun.SDSEndpoint{}
	}
	return err, *ppoint

}

/******************************************************************************
* 概述：     获取该路径下所有节点
* 函数名：    GetAllServiceInfo
* 返回值：    []dzhyun.SDSEndpoint
* 参数列表：  参数名          参数类型      取值范围     描述
*			servicepath    string                  路径
*
*******************************************************************************/
func (sdk *SDSSDK) GetAllServiceInfo(servicepath string) []dzhyun.SDSEndpoint {
	if len(servicepath) == 0 || servicepath[0] != '/' {
		log4.Error("wrong servicepath")
		return nil
	}
	err, pArr := sdk.mevent.GetAllServiceInfo(servicepath)
	if err != nil || pArr == nil {
		return nil
	}
	return *pArr
}

/******************************************************************************
* 概述：     注册订阅信息
* 函数名：    RegistServiceInfo
* 返回值：    error                      error,nil
            chan dzhyun.SDSEndpoint     chan,nil
* 参数列表：  参数名          参数类型      取值范围     描述
*			servicepath    string                  路径
*
*******************************************************************************/
func (sdk *SDSSDK) RegistServiceInfo(servicepath string) (error, chan dzhyun.SDSEndpoint) {
	if len(servicepath) == 0 || servicepath[0] != '/' {
		return errors.New("wrong servicepath"), nil
	}
	err, ch := sdk.mevent.RegistServiceInfo(servicepath)
	return err, ch
}

/******************************************************************************
* 概述：      注销订阅信息
* 函数名：    UnRegistServiceInfo
* 返回值：    error                      error,nil
* 参数列表：  参数名          参数类型      取值范围     描述
*			servicepath    string                    路径
*           ch             chan dzhyun.SDSEndpoint   将关闭的chan
*
*******************************************************************************/
func (sdk *SDSSDK) UnRegistServiceInfo(servicepath string, ch chan dzhyun.SDSEndpoint) error {
	if len(servicepath) == 0 || servicepath[0] != '/' {
		return errors.New("wrong servicepath")
	}
	if ch == nil {
		return errors.New("nil chan")
	}
	sdk.mevent.UnRegistServiceInfo(servicepath, ch)
	return nil
}
