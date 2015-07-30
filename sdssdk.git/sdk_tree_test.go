package sdssdk

import (
	"fmt"
	"testing"
)
import dzhyun "gw.com.cn/dzhyun/dzhyun.git"

//import "fmt"
//func TestAdd1(t *testing.T) {
//	r := Add(1, 2)
//	if r != 2 {
//		t.Errorf("Add(1,2) failed.got %d, expected 3.", r)
//	}
//}

//func BenchmarkAdd1(b *testing.B) {
//	b.StopTimer()

//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		Add(1, 2)
//	}
//}

func TestTreeInit(t *testing.T) {
	//	return
	pTree := &Tree{}
	pTree.Init()
	ch1 := make(chan dzhyun.SDSEndpoint, 1024)
	myTestInsertChan(pTree, "/root/sds", ch1, t)
	ch2 := make(chan dzhyun.SDSEndpoint, 1024)
	myTestInsertChan(pTree, "/root", ch2, t)
	ch3 := make(chan dzhyun.SDSEndpoint, 1024)
	myTestInsertChan(pTree, "/root", ch3, t)
	ch4 := make(chan dzhyun.SDSEndpoint, 1024)
	myTestInsertChan(pTree, "/root/ddd", ch4, t)
	myTestGetRegChans(pTree, "/root/sds", 3, 3, t)
	myTestGetRegChans(pTree, "/root", 2, 3, t)

	myTestInsertPoint(pTree, "n1", "/root/sds", "/root", int32(0), float32(0.1), 1, 3, t)
	//	myTesDelPoint(pTree, "n1", "/root/sds", "/root", int32(0), float32(0.1), 1, 3, t)
	myTestInsertPoint(pTree, "n2", "/root/sfr", "/root", int32(0), float32(0.1), 0, 4, t)
	//	myTesDelPoint(pTree, "n2", "/root/sfr", "/root", int32(0), float32(0.1), 0, 4, t)
	myTestInsertPoint(pTree, "n3", "/root", "/root", int32(0), float32(0.1), 0, 4, t)
	//	myTesDelPoint(pTree, "n3", "/root", "/root", int32(0), float32(0.1), 0, 4, t)
	myTestGetRegChans(pTree, "/root", 2, 4, t)
	ch5 := make(chan dzhyun.SDSEndpoint, 1024)
	myTestInsertChan(pTree, "/root", ch5, t)
	myTestGetRegChans(pTree, "/root", 3, 4, t)
	myTestDelChan(pTree, "/root", ch5, t)
	myTestGetRegChans(pTree, "/root", 2, 4, t)
	pTree.Destory()
	myTestGetRegChans(pTree, "/root", 0, 4, t)
	myTestGetRegChans(pTree, "/root/sds", 0, 4, t)
	myTestGetRegChans(pTree, "/root/sfr", 0, 4, t)
	myTestGetRegChans(pTree, "/root/ddd", 0, 4, t)

	//myTestDelChan(pTree, "/root", &ch, t)
	//myTestGetRegChans(pTree, "/root/sds", 1, 2, t)
	//myTestDelChan(pTree, "/root/sds", &ch, t)
	//myTestGetRegChans(pTree, "/root/sds", 0, 2, t)
	return

	myTestInsertPoint(pTree, "n1", "/root/sfr", "/root", int32(0), float32(0.1), 0, 2, t)
	myTestInsertPoint(pTree, "n2", "/root/sfr2", "/root", int32(0), float32(0.1), 0, 3, t)
	myTestInsertPoint(pTree, "n3", "/root", "/root", int32(0), float32(0.1), 0, 3, t)
	myTestGetRegChans(pTree, "/root", 0, 3, t)

	myTestGetRegChans(pTree, "/root/sds", 0, 3, t)
	myTestGetRegChans(pTree, "/root/sfr", 0, 3, t)
	myTestGetRegChans(pTree, "/root/sfr2", 0, 3, t)

	myTestSearchAll(pTree, "/root", 3, t)
	myTestInsertPoint(pTree, "n4", "/root/sfr", "/root", int32(0), float32(0.1), 0, 3, t)
	myTestSearchAll(pTree, "/root/sfr", 2, t)
	myTestSearchAll(pTree, "/root/sds", 0, t)

	myTesDelPoint(pTree, "n0", "/root/sfr", "/root", int32(0), float32(0.1), 0, 3, t)
	myTestSearchAll(pTree, "/root", 4, t)
	myTesDelPoint(pTree, "n1", "/root/sds", "/root", int32(0), float32(0.1), 0, 3, t)
	myTestSearchAll(pTree, "/root", 4, t)
	myTestSearchAll(pTree, "/root/sds", 0, t)
	myTesDelPoint(pTree, "n4", "/root/sfr", "/root/sds", int32(0), float32(0.1), 0, 3, t)
	myTestSearchAll(pTree, "/root", 3, t)
	myTestSearchAll(pTree, "/root/sfr", 1, t)

}

func myTestSearchAll(pTree *Tree, path string, pointsNum int, t *testing.T) {
	//	fmt.Println("++++++++++++++---myTestSearchAll------------------------------", path)
	pp := pTree.SearchAll(path)
	if pp == nil && pointsNum != 0 {
		t.Errorf("myTestSearchAll(%s) failed.got nil, expected %d.", path, pointsNum)
	}
	if pp == nil {
		return
	}
	if len(*pp) != pointsNum {
		t.Errorf("myTestSearchAll(%s) failed.got %d, expected %d.", path, len(*pp), pointsNum)
	}
	//	fmt.Println("------myTestSearchAll--------", path)
	for i := 0; i < len(*pp); i++ {
		fmt.Println("nodeInfo:", (*pp)[i].GetBusiPath(), (*pp)[i].GetNodeName(), (*pp)[i].GetLoading())
	}
	//	fmt.Println("-----------------------------")
}

func myTestGetRegChans(pTree *Tree, path string, chanNum int, nodeNum int, t *testing.T) {
	//	fmt.Println("+++++++++-----myTestGetRegChans------------------------------", path)
	RChans := make([]chan dzhyun.SDSEndpoint, 0)
	pTree.GetRegChans(&RChans, path)
	if len(RChans) != chanNum {
		t.Errorf("myTestGetRegChans(%s) failed.got %d, expected %d.", path, len(RChans), chanNum)
	}
	if pTree.nodeSize != nodeNum {
		t.Errorf("myTestGetRegChans(%s) failed.got %d, expected %d.", path, pTree.nodeSize, nodeNum)
	}
}

func myTestInsertPoint(pTree *Tree, nodeName, busiPath, startPath string, state int32, loding float32, chanNum, nodeNum int, t *testing.T) {
	//	fmt.Println("+++++++--------myTestInsertPoint------------------------------", busiPath)
	ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
	RChans := make([]chan dzhyun.SDSEndpoint, 0)
	pTree.InsertPoint(ppoint, &RChans, startPath)
	if pTree.nodeSize != nodeNum {
		t.Errorf("MyTestInsertPoint(%s) failed.got %d, expected %d.", busiPath, pTree.nodeSize, nodeNum)
	}
	if len(RChans) != chanNum {
		t.Errorf("MyTestInsertPoint(%s) failed.got %d, expected %d.", busiPath, len(RChans), chanNum)
	}
}

func myTesDelPoint(pTree *Tree, nodeName, busiPath, startPath string, state int32, loding float32, chanNum, nodeNum int, t *testing.T) {
	//	fmt.Println("++++++++++-----myTesDelPoint------------------------------", busiPath)
	ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
	RChans := make([]chan dzhyun.SDSEndpoint, 0)
	pTree.DelPoint(ppoint, &RChans, startPath)
	if pTree.nodeSize != nodeNum {
		t.Errorf("myTesDelPoint(%s) failed.got %d, expected %d.", busiPath, pTree.nodeSize, nodeNum)
	}
	if len(RChans) != chanNum {
		t.Errorf("myTesDelPoint(%s) failed.got %d, expected %d.", busiPath, len(RChans), chanNum)
	}
}

func myTestInsertChan(pTree *Tree, Path string, ch chan dzhyun.SDSEndpoint, t *testing.T) {
	pTree.InsertChan(Path, ch)
}

func myTestDelChan(pTree *Tree, Path string, ch chan dzhyun.SDSEndpoint, t *testing.T) {
	pTree.DelChan(Path, ch)
}
