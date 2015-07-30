/*******************************************************************************
 *
 *  文      件:    sdk_utils.go
 *
 *  概述：         通用函数
 *
 *  版本历史
 *      1.0    2014-10-09    xufengping    创建并实现
 *
*******************************************************************************/
package sdssdk

import (
	"strings"

	log4 "github.com/alecthomas/log4go"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

/******************************************************************************
* 概述：     异常恢复函数
* 函数名：   RecoverPanic
* 返回值：
* 参数列表：  参数名          参数类型      取值范围       描述
*
*******************************************************************************/
func RecoverPanic() {
	if r := recover(); r != nil {
		log4.Debug("Runtime error caught: %v, then recover", r)
	}
}

/******************************************************************************
* 概述：     按序插入数组
* 函数名：   SortMinPoints
* 返回值：
* 参数列表：  参数名          参数类型                取值范围       描述
*           arr            []*dzhyun.SDSEndpoint               操作的数组
*           arrSize        *int                                数组当前元素个数
*           arrCap         int                                 数组的总容量大小
*           insertData     *dzhyun.SDSEndpoint                 被插入的元素
*
*******************************************************************************/
func SortMinPoints(arr *[]*dzhyun.SDSEndpoint, arrSize *int, arrCap int, insertData *dzhyun.SDSEndpoint) {
	pos := *arrSize
	for i := *arrSize - 1; i >= 0; i-- {
		if insertData.GetLoading() < (*arr)[i].GetLoading() {
			if i+1 < arrCap {
				(*arr)[i+1] = (*arr)[i]
			}
			pos = i
		} else {
			break
		}
	}

	if pos <= arrCap-1 {
		(*arr)[pos] = insertData
	}
	if (*arrSize) < arrCap {
		(*arrSize)++
	}
}

/******************************************************************************
* 概述：     获取子串
* 函数名：    GetSubAddr
* 返回值：    string                     “”,addr
*           bool                       true,false
* 参数列表：  参数名          参数类型      取值范围     描述
*           str            string
*           startstr       string
*           endstr         string
*
*******************************************************************************/
func GetSubStr(str, startstr, endstr string) (string, bool) {
	var i, j int
	if len(str) == 0 {
		return "", false
	}
	if i = strings.Index(str, startstr); i == -1 {
		return "", false
	}
	if j = strings.Index(str[i:], endstr); j == -1 {
		j = len(str)
	}
	return str[i+len(startstr) : j], true
}

/******************************************************************************
* 概述：     将包含的字窜剔除
* 函数名：    FilteContainStr
* 返回值：    *[]int                     被剔除的字窜pos
*
* 参数列表：  参数名          参数类型      取值范围     描述
*            arrIn          *[]string
*
*******************************************************************************/
func FilteContainStr(arrIn *[]string) *[]int {
	arrOut := make([]int, 0)
	for i := 0; i < len(*arrIn); i++ {
		newStr := (*arrIn)[i]

		if newStr[len(newStr)-1] != '/' {
			newStr = newStr + "/"
		}
		for j := 0; j < len(arrOut); j++ {
			if i == (arrOut)[j] {

				goto HERE1
			}
		}
		for j := 0; j < len(*arrIn); j++ {
			if j == i {

				continue
			}
			for z := 0; z < len(arrOut); z++ {
				if j == (arrOut)[z] {

					goto HERE2
				}
			}
			if strings.Contains((*arrIn)[j], newStr) {
				arrOut = append(arrOut, j)
			}
		HERE2:
		}
	HERE1:
	}
	return &arrOut
}
