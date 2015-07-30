package sdssdk

import (
	"testing"
	"time"
)
import dzhyun "gw.com.cn/dzhyun/dzhyun.git"

func TestCacheInit(t *testing.T) {
	//	pCache := DataCache{}
	//	pCache.Init(nil)
	//	mSocketGroup()
}

func TestUpdateSdsTopic(t *testing.T) {
	pCache := DataCache{}
	pCache.Init(nil)
	pm := &pCache.msdsTopic

	pCache.UpdateSdsTopic("/root/A", true)
	{
		value, ok := (*pm)["/root/A"]
		if !ok {
			t.Errorf("TestUpdateSdsTopic(/root/A) failed.got false, expected ok.")
		}
		if value != 1 {
			t.Errorf("TestUpdateSdsTopic(/root/A) failed.got %d, expected 1.", value)
		}
	}
	pCache.UpdateSdsTopic("/root/B", true)
	{
		value, ok := (*pm)["/root/B"]
		if !ok {
			t.Errorf("TestUpdateSdsTopic(/root/B) failed.got false, expected ok.")
		}
		if value != 1 {
			t.Errorf("TestUpdateSdsTopic(/root/B) failed.got %d, expected 1.", value)
		}
	}
	pCache.UpdateSdsTopic("/root", true)
	{
		value, ok := (*pm)["/root"]
		if !ok {
			t.Errorf("TestUpdateSdsTopic(/root) failed.got false, expected ok.")
		}
		if value != 1 {
			t.Errorf("TestUpdateSdsTopic(/root) failed.got %d, expected 1.", value)
		}
	}
	pCache.UpdateSdsTopic("/root/B", true)
	{
		value, ok := (*pm)["/root/B"]
		if !ok {
			t.Errorf("TestUpdateSdsTopic(/root/B) failed.got false, expected ok.")
		}
		if value != 2 {
			t.Errorf("TestUpdateSdsTopic(/root/B) failed.got %d, expected 2.", value)
		}
	}
}

func TestHaveSdsTopic(t *testing.T) {
	pCache := DataCache{}
	pCache.Init(nil)
	//	pm := &pCache.msdsTopic
	pCache.UpdateSdsTopic("/root/A", true)
	{
		ok := pCache.HaveSdsTopic("/root/A")
		if !ok {
			t.Errorf("TestHaveSdsTopic(/root/A) failed.got false, expected ok.")
		}
	}
	{
		ok := pCache.HaveSdsTopic("/root/B")
		if ok {
			t.Errorf("TestHaveSdsTopic(/root/B) failed.got true, expected false.")
		}
	}
	pCache.UpdateSdsTopic("/root/B", true)
	{
		ok := pCache.HaveSdsTopic("/root/B")
		if !ok {
			t.Errorf("TestHaveSdsTopic(/root/B) failed.got false, expected true.")
		}
	}
	pCache.UpdateSdsTopic("/root/B", true)
	{
		ok := pCache.HaveSdsTopic("/root/B")
		if !ok {
			t.Errorf("TestHaveSdsTopic(/root/B) failed.got false, expected true.")
		}
	}
}

func TestGetSdsTopic(t *testing.T) {
	pCache := DataCache{}
	pCache.Init(nil)
	//	pm := &pCache.msdsTopic
	{
		pCache.UpdateSdsTopic("/root/A", true)
		pCache.UpdateSdsTopic("/root/A", true)
		pCache.UpdateSdsTopic("/root", true)
		pCache.UpdateSdsTopic("/root/B", true)
		pCache.UpdateSdsTopic("/root/A", true)
	}
	pstrs := pCache.GetSdsTopic()
	for i := 0; i < len(*pstrs); i++ {
		if !((*pstrs)[i] == "/root/A" || (*pstrs)[i] == "/root/B" || (*pstrs)[i] == "/root") {
			t.Errorf("TestGetSdsTopic() failed.got %s, expected other.", (*pstrs)[i])
		}
	}
	if len(*pstrs) != 3 {
		t.Errorf("TestGetSdsTopic() failed.got %d, expected 3.", len(*pstrs))
	}
}

func TestInsertSocket(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	pm := &pCache.mSocketGroup

	{
		nodeName := "sds_1"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}

		pCache.insertSocket(ppoint)

		if len(*pm) != 1 {
			t.Errorf("TestInsertSocket() failed.got %d, expected 1.", len(*pm))
		}
	}
	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(*pm) != 2 {
			t.Errorf("TestInsertSocket() failed.got %d, expected 2.", len(*pm))
		}
	}

	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(*pm) != 2 {
			t.Errorf("TestInsertSocket() failed.got %d, expected 2.", len(*pm))
		}
	}
	{
		nodeName := "sds_3"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 UB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(*pm) != 2 {
			t.Errorf("TestInsertSocket() failed.got %d, expected 2.", len(*pm))
		}
	}

	time.Sleep(time.Second * 5)

}

func TestGetSdsAddr(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	{
		_, _, err := pCache.getSdsAddr("REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301")
		if err != nil {
			t.Errorf("TestGetSdsAddr() failed.got err= %S.", err.Error())
		}
	}
	{
		_, _, err := pCache.getSdsAddr("REP:tcp://10.15.144.105:10300 UB:tcp://10.15.144.105:10301")
		if err == nil {
			t.Errorf("TestGetSdsAddr() failed.got noerr, expected err.")
		}
	}
	{
		_, _, err := pCache.getSdsAddr("ReP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301")
		if err == nil {
			t.Errorf("TestGetSdsAddr() failed.got noerr, expected err.")
		}
	}

}

func TestUpdatePoints(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)

	{ //空表
		nodeName := "a_2"
		busiPath := "/root/A"
		state := int32(0)
		loding := float32(0.1)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		if pCache.mpointGroup.pointSize != 1 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 1.", pCache.mpointGroup.pointSize)
		}
		if pCache.mpointGroup.nodeSize != 2 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 2.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 0 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 0.", pCache.mpointGroup.chanSize)
		}
	}

	{ //删除不存在的节点
		nodeName := "a_1"
		busiPath := "/root/A"
		state := int32(1)
		loding := float32(0.1)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		if pCache.mpointGroup.pointSize != 1 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 1.", pCache.mpointGroup.pointSize)
		}
		if pCache.mpointGroup.nodeSize != 2 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 2.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 0 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 0.", pCache.mpointGroup.chanSize)
		}
	}

	{ //插入两节点
		nodeName := "b_1"
		busiPath := "/root/B"
		state := int32(0)
		loding := float32(0.1)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		nodeName1 := "b_2"
		busiPath1 := "/root/B"
		state1 := int32(0)
		loding1 := float32(0.1)
		ppoint1 := &dzhyun.SDSEndpoint{NodeName: &nodeName1, BusiPath: &busiPath1, State: &state1, Loading: &loding1}
		ps = append(ps, ppoint1)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		if pCache.mpointGroup.pointSize != 3 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 3.", pCache.mpointGroup.pointSize)
		}
		if pCache.mpointGroup.nodeSize != 3 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 3.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 0 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 0.", pCache.mpointGroup.chanSize)
		}
	}

	{ //删除两节点，一个存在另一个不存在
		nodeName1 := "b_2"
		busiPath1 := "/root/B"
		state1 := int32(1)
		loding1 := float32(0.1)
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ppoint1 := &dzhyun.SDSEndpoint{NodeName: &nodeName1, BusiPath: &busiPath1, State: &state1, Loading: &loding1}
		ps = append(ps, ppoint1)

		nodeName2 := "b_2"
		busiPath2 := "/root/D"
		state2 := int32(1)
		loding2 := float32(0.1)

		ppoint2 := &dzhyun.SDSEndpoint{NodeName: &nodeName2, BusiPath: &busiPath2, State: &state2, Loading: &loding2}
		ps = append(ps, ppoint2)
		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		if pCache.mpointGroup.pointSize != 2 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 2.", pCache.mpointGroup.pointSize)
		}
		if pCache.mpointGroup.nodeSize != 3 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 3.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 0 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 0.", pCache.mpointGroup.chanSize)
		}
	}

	{ //更新一个存在节点
		nodeName := "b_1"
		busiPath := "/root/B"
		state := int32(0)
		loding := float32(0.6)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)
		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		if pCache.mpointGroup.pointSize != 2 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 2.", pCache.mpointGroup.pointSize)
		}
		if pCache.mpointGroup.nodeSize != 3 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 3.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 0 {
			t.Errorf("TestUpdatePoints() failed.got %d, expected 0.", pCache.mpointGroup.chanSize)
		}
	}

}

func TestRefleshPoints(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)

	{ //空表也没新数据
		subName1 := "/root"
		Servicepath := "/root"
		pres1 := &dzhyun.SDSResponse{SubName: &subName1, Endpoints: nil}
		pCache.RefleshPoints(Servicepath, pres1)
		if pCache.mpointGroup.pointSize != 0 {
			t.Errorf("TestRefleshPoints() failed.got %d, expected 0.", pCache.mpointGroup.pointSize)
		}
	}
	{ //空表
		Servicepath := "/root"
		nodeName := "b_1"
		busiPath := "/root/B"
		state := int32(0)
		loding := float32(0.6)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.RefleshPoints(Servicepath, pres)
		if pCache.mpointGroup.pointSize != 1 {
			t.Errorf("TestRefleshPoints() failed.got %d, expected 1.", pCache.mpointGroup.pointSize)
		}
	}

	{ //数据有删除及插入
		/////////////////////////////////////////////////////////////////
		ps1 := make([]*dzhyun.SDSEndpoint, 0)
		nodeName4 := "b_2"
		busiPath4 := "/root/B"
		state4 := int32(0)
		loding4 := float32(0.1)
		ppoint4 := &dzhyun.SDSEndpoint{NodeName: &nodeName4, BusiPath: &busiPath4, State: &state4, Loading: &loding4}
		ps1 = append(ps1, ppoint4)

		nodeName5 := "b_3"
		busiPath5 := "/root/B"
		state5 := int32(0)
		loding5 := float32(0.9)
		ppoint5 := &dzhyun.SDSEndpoint{NodeName: &nodeName5, BusiPath: &busiPath5, State: &state5, Loading: &loding5}

		ps1 = append(ps1, ppoint5)

		subName1 := "/root"
		Servicepath := "/root"
		pres1 := &dzhyun.SDSResponse{SubName: &subName1, Endpoints: ps1}
		pCache.RefleshPoints(Servicepath, pres1)

		if pCache.mpointGroup.pointSize != 2 {
			t.Errorf("TestRefleshPoints() failed.got %d, expected 2.", pCache.mpointGroup.pointSize)
		}

	}

	{ //数据有删除 有更新
		nodeName := "b_1"
		busiPath := "/root/B"
		state := int32(0)
		loding := float32(0.1)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		nodeName1 := "b_2"
		busiPath1 := "/root/B"
		state1 := int32(0)
		loding1 := float32(0.1)
		ppoint1 := &dzhyun.SDSEndpoint{NodeName: &nodeName1, BusiPath: &busiPath1, State: &state1, Loading: &loding1}
		ps = append(ps, ppoint1)

		nodeName3 := "b_3"
		busiPath3 := "/root/B"
		state3 := int32(0)
		loding3 := float32(0.1)
		ppoint3 := &dzhyun.SDSEndpoint{NodeName: &nodeName3, BusiPath: &busiPath3, State: &state3, Loading: &loding3}

		ps = append(ps, ppoint3)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		/////////////////////////////////////////////////////////////////
		ps1 := make([]*dzhyun.SDSEndpoint, 0)
		nodeName4 := "b_2"
		busiPath4 := "/root/B"
		state4 := int32(0)
		loding4 := float32(0.1)
		ppoint4 := &dzhyun.SDSEndpoint{NodeName: &nodeName4, BusiPath: &busiPath4, State: &state4, Loading: &loding4}
		ps1 = append(ps1, ppoint4)

		nodeName5 := "b_3"
		busiPath5 := "/root/B"
		state5 := int32(0)
		loding5 := float32(0.9)
		ppoint5 := &dzhyun.SDSEndpoint{NodeName: &nodeName5, BusiPath: &busiPath5, State: &state5, Loading: &loding5}

		ps1 = append(ps1, ppoint5)

		subName1 := "/root"
		Servicepath := "/root"
		pres1 := &dzhyun.SDSResponse{SubName: &subName1, Endpoints: ps1}
		pCache.RefleshPoints(Servicepath, pres1)

		if pCache.mpointGroup.pointSize != 2 {
			t.Errorf("TestRefleshPoints() failed.got %d, expected 2.", pCache.mpointGroup.pointSize)
		}

	}

	{ //原来有数据，现在全没有了
		subName1 := "/root"
		Servicepath := "/root"
		pres1 := &dzhyun.SDSResponse{SubName: &subName1, Endpoints: nil}
		pCache.RefleshPoints(Servicepath, pres1)
		if pCache.mpointGroup.pointSize != 0 {
			t.Errorf("TestRefleshPoints() failed.got %d, expected 0.", pCache.mpointGroup.pointSize)
		}
	}
}

func TestApplyChan(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	{
		pCache.ApplyChan("/root/A")
		if pCache.mpointGroup.nodeSize != 2 {
			t.Errorf("TestApplyChan() failed.got %d, expected 2.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 1 {
			t.Errorf("TestApplyChan() failed.got %d, expected 1.", pCache.mpointGroup.chanSize)
		}
	}
	{
		pCache.ApplyChan("/root/A/B")
		if pCache.mpointGroup.nodeSize != 3 {
			t.Errorf("TestApplyChan() failed.got %d, expected 3.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 2 {
			t.Errorf("TestApplyChan() failed.got %d, expected 2.", pCache.mpointGroup.chanSize)
		}
	}
	{
		pCache.ApplyChan("/root/A/B")
		if pCache.mpointGroup.nodeSize != 3 {
			t.Errorf("TestApplyChan() failed.got %d, expected 3.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 3 {
			t.Errorf("TestApplyChan() failed.got %d, expected 3.", pCache.mpointGroup.chanSize)
		}
	}
	{
		ch := pCache.ApplyChan("/root/D")
		if pCache.mpointGroup.nodeSize != 4 {
			t.Errorf("TestApplyChan() failed.got %d, expected 2.", pCache.mpointGroup.nodeSize)
		}
		if pCache.mpointGroup.chanSize != 4 {
			t.Errorf("TestApplyChan() failed.got %d, expected 1.", pCache.mpointGroup.chanSize)
		}
		nodeName := "d_1"
		busiPath := "/root/D"
		state := int32(0)
		loding := float32(0.6)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)
		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)
		k := 0
		for flag := true; flag; {
			select {
			case <-ch:
				k++
			default:
				flag = false
			}
		}
		if k != 1 {
			t.Errorf("TestApplyChan() failed.got %d, expected 1.", k)
		}
		k = 0
		busiPath = "/root"
		pCache.UpdatePoints(pres)
		for flag := true; flag; {
			select {
			case <-ch:
				k++
			default:
				flag = false
			}
		}
		if k != 0 {
			t.Errorf("TestApplyChan() failed.got %d, expected 0.", k)
		}
		k = 0
		busiPath = "/root/D/A"
		pCache.UpdatePoints(pres)
		for flag := true; flag; {
			select {
			case <-ch:
				k++
			default:
				flag = false
			}
		}
		if k != 1 {
			t.Errorf("TestApplyChan() failed.got %d, expected 1.", k)
		}

	}
}

func TestFindOne(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	{
		nodeName := "b_1"
		busiPath := "/root/B"
		state := int32(0)
		loding := float32(0.1)
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding}
		ps := make([]*dzhyun.SDSEndpoint, 0)
		ps = append(ps, ppoint)

		nodeName1 := "b_2"
		busiPath1 := "/root/B"
		state1 := int32(0)
		loding1 := float32(0.2)
		ppoint1 := &dzhyun.SDSEndpoint{NodeName: &nodeName1, BusiPath: &busiPath1, State: &state1, Loading: &loding1}
		ps = append(ps, ppoint1)

		nodeName3 := "b_3"
		busiPath3 := "/root/B"
		state3 := int32(0)
		loding3 := float32(0.2)
		ppoint3 := &dzhyun.SDSEndpoint{NodeName: &nodeName3, BusiPath: &busiPath3, State: &state3, Loading: &loding3}

		ps = append(ps, ppoint3)

		nodeName4 := "b_4"
		busiPath4 := "/root"
		state4 := int32(0)
		loding4 := float32(0.5)
		ppoint4 := &dzhyun.SDSEndpoint{NodeName: &nodeName4, BusiPath: &busiPath4, State: &state4, Loading: &loding4}
		ps = append(ps, ppoint4)

		nodeName5 := "b_5"
		busiPath5 := "/root/L"
		state5 := int32(0)
		loding5 := float32(0.9)
		ppoint5 := &dzhyun.SDSEndpoint{NodeName: &nodeName5, BusiPath: &busiPath5, State: &state5, Loading: &loding5}
		ps = append(ps, ppoint5)

		nodeName6 := "b_6"
		busiPath6 := "/root/L"
		state6 := int32(0)
		loding6 := float32(0.7)
		ppoint6 := &dzhyun.SDSEndpoint{NodeName: &nodeName6, BusiPath: &busiPath6, State: &state6, Loading: &loding6}
		ps = append(ps, ppoint6)

		nodeName7 := "b_7"
		busiPath7 := "/root/T"
		state7 := int32(0)
		loding7 := float32(0.8)
		ppoint7 := &dzhyun.SDSEndpoint{NodeName: &nodeName7, BusiPath: &busiPath7, State: &state7, Loading: &loding7}
		ps = append(ps, ppoint7)

		subName := "/root"
		pres := &dzhyun.SDSResponse{SubName: &subName, Endpoints: ps}
		pCache.UpdatePoints(pres)

		if _, ok := pCache.FindOne("/root/A/B/A"); ok {
			t.Errorf("TestFindOne() failed.got true, expected false.")
		}

		if _, ok := pCache.FindOne("/root"); !ok {
			t.Errorf("TestFindOne() failed.got false, expected true.")
		}
	}
}

func TestGetSocket(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	{
		nodeName := "sds_1"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}

		pCache.insertSocket(ppoint)

		if len(pCache.mSocketGroup) != 1 {
			t.Errorf("TestGetSocket() failed.got %d, expected 1.", len(pCache.mSocketGroup))
		}
	}
	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 2 {
			t.Errorf("TestGetSocket() failed.got %d, expected 2.", len(pCache.mSocketGroup))
		}
	}

	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 2 {
			t.Errorf("TestGetSocket() failed.got %d, expected 2.", len(pCache.mSocketGroup))
		}
	}
	{
		nodeName := "sds_3"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10301"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 3 {
			t.Errorf("TestGetSocket() failed.got %d, expected 2.", len(pCache.mSocketGroup))
		}
	}

	time.Sleep(time.Second * 5)

	ps := pCache.GetSocket()
	if ps == nil {
		t.Errorf("TestGetSocket() failed.got nil, expected !nil.")
	}
	if !ps.mreqOK || !ps.msubOK {
		t.Errorf("TestGetSocket() failed.got !true, expected 2*True.")
	}
}

func TestGetSocket2(t *testing.T) {
	pCache := &DataCache{}
	pzmq := &Zmq{}
	pzmq.Init(pCache, nil)
	pCache.Init(pzmq)
	ps := pCache.GetSocket()
	if ps != nil {
		t.Errorf("TestGetSocket() failed.got !nil, expected nil.")
	}
	{
		nodeName := "sds_1"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10318"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}

		pCache.insertSocket(ppoint)

		if len(pCache.mSocketGroup) != 1 {
			t.Errorf("TestGetSocket() failed.got %d, expected 1.", len(pCache.mSocketGroup))
		}
	}
	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10318"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 2 {
			t.Errorf("TestGetSocket() failed.got %d, expected 2.", len(pCache.mSocketGroup))
		}
	}

	{
		nodeName := "sds_2"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10318"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 2 {
			t.Errorf("TestGetSocket() failed.got %d, expected 2.", len(pCache.mSocketGroup))
		}
	}
	{
		nodeName := "sds_3"
		busiPath := "/root/sds"
		state := int32(0)
		loding := float32(0.1)
		inf := "REP:tcp://10.15.144.105:10300 PUB:tcp://10.15.144.105:10318"
		ppoint := &dzhyun.SDSEndpoint{NodeName: &nodeName, BusiPath: &busiPath, State: &state, Loading: &loding, Interface: &inf}
		pCache.insertSocket(ppoint)
		if len(pCache.mSocketGroup) != 3 {
			t.Errorf("TestGetSocket() failed.got %d, expected 3.", len(pCache.mSocketGroup))
		}
	}

	time.Sleep(time.Second * 3)

	ps2 := pCache.GetSocket()
	if ps2 != nil {
		t.Errorf("TestGetSocket() failed.got !nil, expected nil.")
	}
}
