package sdssdk

import (
	"fmt"
	"testing"
	"time"
	//	"time"
)

//import dzhyun "gw.com.cn/dzhyun/dzhyun.git"

func TestEventInit(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestInit() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestInit() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) < 1 {
		t.Errorf("TestInit() failed.got %d, expected >=1.", len(pEvent.mcache.mSocketGroup))
	}
}

func TestEventInit2(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10346")
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestInit2() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestInit2() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestInit2() failed.got %d, expected !0.", len(pEvent.mcache.mSocketGroup))
	}
}

func TestEventInit3(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10346")
	addrs = append(addrs, "tcp://10.15.144.105:10349")
	err := pEvent.Init(addrs, "/root/sds")
	if err == nil {
		t.Errorf("TestInit3() failed.got nil, expected err.")
	}
	if len(pEvent.mcache.msdsTopic) != 0 {
		t.Errorf("TestInit3() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) != 0 {
		t.Errorf("TestInit3() failed.got %d, expected 0.", len(pEvent.mcache.mSocketGroup))
	}
}

func TestEventStop(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestStop() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestStop() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestStop() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}

	time.Sleep(time.Second * 10)
	pEvent.Stop()
}

func TestGetWorkSocket(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestGetWorkSocket() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestGetWorkSocket() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestGetWorkSocket() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}
	s := pEvent.GetWorkSocket()
	if s == nil {
		t.Errorf("TestGetWorkSocket() failed.got nil, expected !nil.")
	}
}

func TestUpdateWorkSocket(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestGetWorkSocket() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestGetWorkSocket() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestGetWorkSocket() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}
	pEvent.mnowSocket = nil

	pEvent.UpdateWorkSocket()
	time.Sleep(time.Second * 10)
	if pEvent.mnowSocket == nil {
		t.Errorf("TestGetWorkSocket() failed.got nil, expected !nil.")
	}
}

func TestUpdateWorkSocket2(t *testing.T) {
	//	return
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestGetWorkSocket2() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestGetWorkSocket2() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestGetWorkSocket2() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}
	time.Sleep(time.Second * 10)
	pEvent.mnowSocket = nil

	for _, value := range pEvent.mcache.mSocketGroup {
		value.mok = 1
		value.mreqOK = false
		value.mChoose = false
	}

	pEvent.UpdateWorkSocket()
	time.Sleep(time.Second * 10)
	if pEvent.mnowSocket != nil {
		t.Errorf("TestGetWorkSocket2() failed.got !nil, expected nil.")
	}
}

func TestReqZmqRegister(t *testing.T) {
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestReqZmqRegister() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestReqZmqRegister() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestReqZmqRegister() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}
	if pEvent.mcache.mpointGroup.pointSize == 0 {
		t.Errorf("TestReqZmqRegister() failed.got 0, expected !0.", (pEvent.mcache.mpointGroup.pointSize))
	}

	if err := pEvent.reqZmqRegister("/root"); err != nil {
		t.Errorf("TestReqZmqRegister() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 2 {
		t.Errorf("TestReqZmqRegister() failed.got %d, expected 2.", len(pEvent.mcache.msdsTopic))
	}
	if pEvent.mcache.mpointGroup.pointSize == 0 {
		t.Errorf("TestReqZmqRegister() failed.got 0, expected !0.")
	}
}

func TestEventGetServiceInfo(t *testing.T) {
	pEvent := Event{}
	addrs := make([]string, 0)
	addrs = append(addrs, "tcp://10.15.144.105:10300")
	err := pEvent.Init(addrs, "/root/sds")
	if err != nil {
		t.Errorf("TestGetServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	if len(pEvent.mcache.msdsTopic) != 1 {
		t.Errorf("TestGetServiceInfo() failed.got %d, expected 1.", len(pEvent.mcache.msdsTopic))
	}
	if len(pEvent.mcache.mSocketGroup) == 0 {
		t.Errorf("TestGetServiceInfo() failed.got %d, expected >0.", len(pEvent.mcache.mSocketGroup))
	}
	if pEvent.mcache.mpointGroup.pointSize == 0 {
		t.Errorf("TestGetServiceInfo() failed.got 0, expected >0.")
	}
	fmt.Println("++++++++++++++++++++++W+++++++++++++++++++++++++")
	if err, _ := pEvent.GetServiceInfo("/root/test"); err != nil {
		t.Errorf("TestGetServiceInfo() failed.got %s, expected nil.", err.Error())
	}

	if err, _ := pEvent.GetServiceInfo("/root/test"); err != nil {
		t.Errorf("TestGetServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
	time.Sleep(time.Second * 10)
}
