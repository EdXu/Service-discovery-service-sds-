package sdssdk

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

/***************************Zmq************************************************
* 概述:
* 类型名:     Tree
* 成员列表：   成员名            成员类型       取值范围       描述
*             proot            *TreeNode
*             nodeSize         int                 树有效孩子结点数据(不包括第一个冗余结点)
*             pointSize        int                 point节点的数量
*             chanSize         int                 chan的数量
*             mmutex           sync.Mutex          整棵树的
*
******************************************************************************/
type Tree struct {
	proot     *TreeNode //第一个节点为冗余节点，无用
	nodeSize  int
	pointSize int
	chanSize  int
	mmutex    sync.Mutex
}

/***************************Zmq************************************************
* 概述:       树结点结构
* 类型名:     TreeNode
* 成员列表：   成员名            成员类型       取值范围               描述
*             name             string                             节点名称
*         	  regChans         map[chan dzhyun.SDSEndpoint]int    注册的chan
*	          nodes            map[string]*dzhyun.SDSEndpoint     nodes存储
*	          childs           map[string]*TreeNode               孩子们
*
******************************************************************************/
type TreeNode struct {
	name     string
	regChans map[chan dzhyun.SDSEndpoint]int
	nodes    map[string]*dzhyun.SDSEndpoint
	childs   map[string]*TreeNode
}

/******************************************************************************
* 概述：      初始化
* 函数名：    Init
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) Init() {
	this.nodeSize = 0
	this.pointSize = 0
	this.chanSize = 0
	rand.Seed(time.Now().UnixNano())
	this.proot = &TreeNode{name: "/", regChans: nil, nodes: nil, childs: nil}
	this.proot.childs = make(map[string]*TreeNode)
	this.proot.nodes = make(map[string]*dzhyun.SDSEndpoint)
	this.proot.regChans = make(map[chan dzhyun.SDSEndpoint]int)
}

/******************************************************************************
* 概述：      销毁chans
* 函数名：    Destory
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) Destory() {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	pTreeNode := this.proot
	this.subDestory(pTreeNode)
}

/******************************************************************************
* 概述：      递归销毁chans
* 函数名：    subDestory
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) subDestory(pTreeNode *TreeNode) {
	if pTreeNode == nil {
		return
	}
	if pTreeNode.regChans != nil {
		for key, _ := range pTreeNode.regChans {
			delete(pTreeNode.regChans, key)
			close(key)
			this.chanSize--
		}
	}

	if pTreeNode.childs == nil {
		return
	}
	for _, value := range pTreeNode.childs {
		this.subDestory(value)
	}
}

/******************************************************************************
* 概述：      查找目标路径，然后递归获取point，最后选取一个
* 函数名：    SearchBest
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*            path           string
*            num            int                      在load最小的num个point中选取
*
*******************************************************************************/
func (this *Tree) SearchBest(path string, num int) *dzhyun.SDSEndpoint {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(path, "/")
	pTreeNode := this.proot
	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			return nil
		}
		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
		} else {
			return nil
		}
	}

	arr := make([]*dzhyun.SDSEndpoint, num, num)
	arrSize := 0
	this.subSearchSome(pTreeNode, &arr, &arrSize, num)
	if arrSize == 0 {
		return nil
	}

	randPos := rand.Intn(num)
	//	fmt.Println("rand pos =", randPos)

	randPos = randPos % arrSize
	p := *arr[randPos]
	return &p
}

/******************************************************************************
* 概述：      递归获取point(排序)
* 函数名：    subSearchSome
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) subSearchSome(pTreeNode *TreeNode, arr *[]*dzhyun.SDSEndpoint, arrSize *int, num int) {
	if pTreeNode == nil {
		return
	}
	if pTreeNode.nodes != nil {
		for _, value := range pTreeNode.nodes {
			SortMinPoints(arr, arrSize, num, value)
		}
	}
	if pTreeNode.childs == nil {
		return
	}
	for _, value := range pTreeNode.childs {
		this.subSearchSome(value, arr, arrSize, num)
	}
}

/******************************************************************************
* 概述：      查找目标路径，然后递归获取point
* 函数名：    SearchAll
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) SearchAll(path string) *[]dzhyun.SDSEndpoint {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(path, "/")
	pTreeNode := this.proot
	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			return nil
		}
		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
		} else {
			return nil
		}
	}

	arr := make([]dzhyun.SDSEndpoint, 0)
	this.subSearchAll(pTreeNode, &arr)
	return &arr
}

/******************************************************************************
* 概述：      递归获取point
* 函数名：    subSearchAll
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) subSearchAll(pTreeNode *TreeNode, arr *[]dzhyun.SDSEndpoint) {
	if pTreeNode == nil {
		return
	}
	if pTreeNode.nodes != nil {
		for _, value := range pTreeNode.nodes {
			*arr = append(*arr, *value)
		}
	}
	if pTreeNode.childs == nil {
		return
	}
	for _, value := range pTreeNode.childs {
		this.subSearchAll(value, arr)
	}
}

/******************************************************************************
* 概述：      插入point，并获取从startPath（不包含）开始的所有chan
* 函数名：    InsertPoint
* 返回值：
* 参数列表：  参数名          参数类型      取值范围         描述
*            ppoint         *dzhyun.SDSEndpoint
*            pRChans        *[]chan dzhyun.SDSEndpoint  返回的chans
*            startPath      qstring                     从startPath（不包含）开始的所有chan
*
*******************************************************************************/
func (this *Tree) InsertPoint(ppoint *dzhyun.SDSEndpoint, pRChans *[]chan dzhyun.SDSEndpoint, startPath string) {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(ppoint.GetBusiPath(), "/")
	pTreeNode := this.proot
	startRegPaths := strings.Split(startPath, "/")
	j := len(startRegPaths) - 1

	for i := 0; i < len(paths); i++ { //顺便推送
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			pTreeNode.childs = make(map[string]*TreeNode)
			newOne := &TreeNode{name: paths[i], regChans: nil, nodes: nil, childs: nil}
			pTreeNode.childs[paths[i]] = newOne
			pTreeNode = newOne
			this.nodeSize++
			continue
		}

		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
			if pRChans != nil && i > j && pTreeNode.regChans != nil {
				//for key, _ := range pTreeNode.regChans {
				//	*key <- *ppoint
				//}
				for key, _ := range pTreeNode.regChans {
					*pRChans = append(*pRChans, key)
				}
			}
		} else {
			newOne := &TreeNode{name: paths[i], regChans: nil, nodes: nil, childs: nil}
			pTreeNode.childs[paths[i]] = newOne
			pTreeNode = newOne
			this.nodeSize++
		}
	}
	if pTreeNode.nodes == nil {
		pTreeNode.nodes = make(map[string]*dzhyun.SDSEndpoint)
	}
	if _, ok := pTreeNode.nodes[ppoint.GetNodeName()]; !ok {
		this.pointSize++
	}

	pTreeNode.nodes[ppoint.GetNodeName()] = ppoint
}

/******************************************************************************
* 概述：      删除point，并获取从startPath（不包含）开始的所有chan
* 函数名：    DelPoint
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) DelPoint(ppoint *dzhyun.SDSEndpoint, pRChans *[]chan dzhyun.SDSEndpoint, startPath string) {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(ppoint.GetBusiPath(), "/")
	pTreeNode := this.proot

	startRegPaths := strings.Split(startPath, "/")
	j := len(startRegPaths) - 1

	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			return
		}

		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
			if pRChans != nil && i > j && pTreeNode.regChans != nil {
				//for key, _ := range pTreeNode.regChans {
				//	*key <- *ppoint
				//}
				for key, _ := range pTreeNode.regChans {
					*pRChans = append(*pRChans, key)
				}
			}
		} else {
			return
		}
	}
	if pTreeNode.nodes != nil {
		if _, ok := pTreeNode.nodes[ppoint.GetNodeName()]; ok {
			this.pointSize--
		}
		delete(pTreeNode.nodes, ppoint.GetNodeName())
	}
}

/******************************************************************************
* 概述：      获取根到endPath（包含）之间的所有chan
* 函数名：    GetRegChans
* 返回值：
* 参数列表：  参数名          参数类型      取值范围   描述
*            pRChans                               *[]chan dzhyun.SDSEndpoint
*            endPath        string                 根到endPath（包含）之间的所有chan
*
*******************************************************************************/
func (this *Tree) GetRegChans(pRChans *[]chan dzhyun.SDSEndpoint, endPath string) {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(endPath, "/")
	pTreeNode := this.proot

	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			return
		}

		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
			if pTreeNode.regChans != nil {
				for key, _ := range pTreeNode.regChans {
					*pRChans = append(*pRChans, key)
				}
			}
		} else {
			return
		}
	}
}

/******************************************************************************
* 概述：     插入chan 并订阅
* 函数名：    InsertChan
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*
*******************************************************************************/
func (this *Tree) InsertChan(path string, ch chan dzhyun.SDSEndpoint) {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	if len(path) < 1 || path[0] != '/' {
		return
	}

	paths := strings.Split(path, "/")
	pTreeNode := this.proot
	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			pTreeNode.childs = make(map[string]*TreeNode)
			newOne := &TreeNode{name: paths[i], regChans: nil, nodes: nil, childs: nil}
			pTreeNode.childs[paths[i]] = newOne
			pTreeNode = newOne
			this.nodeSize++
			continue
		}
		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
		} else {
			newOne := &TreeNode{name: paths[i], regChans: nil, nodes: nil, childs: nil}
			pTreeNode.childs[paths[i]] = newOne
			pTreeNode = newOne
			this.nodeSize++
		}
	}
	if pTreeNode.regChans == nil {
		pTreeNode.regChans = make(map[chan dzhyun.SDSEndpoint]int)
	}
	if _, ok := pTreeNode.regChans[ch]; !ok {
		this.chanSize++
	}
	pTreeNode.regChans[ch] = 0
}

/******************************************************************************
* 概述：     删除chan 并取消订阅
* 函数名：    InsertChan
* 返回值：
* 参数列表：  参数名          参数类型      取值范围     描述
*            path           string
*            ch             chan dzhyun.SDSEndpoint  与订阅时传入的ch相同
*
*******************************************************************************/
func (this *Tree) DelChan(path string, ch chan dzhyun.SDSEndpoint) {
	this.mmutex.Lock()
	defer this.mmutex.Unlock()
	paths := strings.Split(path, "/")
	pTreeNode := this.proot
	for i := 0; i < len(paths); i++ {
		if len(paths[i]) == 0 {
			continue
		}
		if pTreeNode.childs == nil {
			return
		}
		if value, ok := pTreeNode.childs[paths[i]]; ok {
			pTreeNode = value
		} else {
			return
		}
	}
	if pTreeNode.regChans == nil {
		return
	}
	if _, ok := pTreeNode.regChans[ch]; ok {
		this.chanSize--
	}
	delete(pTreeNode.regChans, ch)

}

///////////////////////////subnameTree,未测试//////////////////////////
//type TopicTree struct {
//	proot    *TopicTreeNode //第一个节点为冗余节点，无用
//	regSize  int
//	nodeSize int
//}

//type TopicTreeNode struct {
//	name    string
//	regFlag bool
//	childs  map[string]*TopicTreeNode
//}

//func (this *TopicTree) Init() {
//	this.regSize = 0
//	this.nodeSize = 0
//	this.proot = &TopicTreeNode{name: "/", regFlag: false, childs: nil}
//	this.proot.childs = make(map[string]*TopicTreeNode)
//}

//func (this *TopicTree) Destory() {

//}

//func (this *TopicTree) Insert(path string) bool {
//	paths := strings.Split(path, "/")
//	pTreeNode := this.proot

//	for i := 0; i < len(paths); i++ {
//		if len(paths[i]) == 0 {
//			continue
//		}
//		if pTreeNode.childs == nil {
//			pTreeNode.childs = make(map[string]*TopicTreeNode)
//			newOne := &TopicTreeNode{name: paths[i], regFlag: false, childs: nil}
//			pTreeNode.childs[paths[i]] = newOne
//			pTreeNode = newOne
//			this.nodeSize++
//			continue
//		}

//		if value, ok := pTreeNode.childs[paths[i]]; ok {
//			pTreeNode = value
//			if pTreeNode.regFlag {
//				return false
//			}
//		} else {
//			newOne := &TopicTreeNode{name: paths[i], regFlag: false, childs: nil}
//			pTreeNode.childs[paths[i]] = newOne
//			pTreeNode = newOne
//			this.nodeSize++
//		}
//	}
//	pTreeNode.regFlag = true
//	return true
//}

//func (this *TopicTree) Search(path string) bool {
//	paths := strings.Split(path, "/")
//	pTreeNode := this.proot
//	for i := 0; i < len(paths); i++ {
//		if len(paths[i]) == 0 {
//			continue
//		}
//		if pTreeNode.childs == nil {
//			return false
//		}
//		if value, ok := pTreeNode.childs[paths[i]]; ok {
//			pTreeNode = value
//			if pTreeNode.regFlag {
//				return true
//			}
//		} else {
//			return false
//		}
//	}
//	if pTreeNode.regFlag {
//		return true
//	} else {
//		return false
//	}

//}

//func (this *TopicTree) GetAllTopics() *[]string {
//	pTreeNode := this.proot
//	arr := make([]string, 0)
//	headpath := ""
//	this.subGetAllTopics(pTreeNode, headpath, &arr)
//	return &arr
//}

//func (this *TopicTree) subGetAllTopics(pTreeNode *TopicTreeNode, headpath string, arr *[]string) {
//	if pTreeNode == nil {
//		return
//	}
//	if pTreeNode.name != "/" {
//		headpath = headpath + "/" + pTreeNode.name
//		if pTreeNode.regFlag {
//			(*arr) = append((*arr), headpath)
//			return
//		}
//	}

//	if pTreeNode.childs == nil {
//		return
//	}
//	for _, value := range pTreeNode.childs {
//		this.subGetAllTopics(value, headpath, arr)
//	}
//}
