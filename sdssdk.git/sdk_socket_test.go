package sdssdk

import (
	"fmt"
	"testing"
	"time"
)

//	"time"

//import dzhyun "gw.com.cn/dzhyun/dzhyun.git"

func TestSocketInit1(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestSocketInit1() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestSocketInit1() failed.got false, expected 2*true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestSocketInit1() failed.got true, expected false.")
	}
}

func TestSocketInit2(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10349", pzmq, 1); err != nil {
		t.Errorf("TestSocketInit2() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || pZmqSocket.msubOK {
		t.Errorf("TestSocketInit2() failed.got first false or second true, expected true and false.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestSocketInit2() failed.got true, expected false.")
	}
}

func TestSocketInit3(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10323", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestSocketInit3() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestSocketInit3() failed.got first true or second false, expected false and true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestSocketInit3() failed.got true, expected false.")
	}
}

func TestSocketInit4(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10323", "tcp://10.15.144.105:10335", pzmq, 1); err != nil {
		t.Errorf("TestSocketInit4() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if pZmqSocket.mreqOK || pZmqSocket.msubOK {
		t.Errorf("TestSocketInit4() failed.got one true, expected 2*false.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestSocketInit4() failed.got true, expected false.")
	}
}

func TestSocketClose(t *testing.T) {
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestSocketClose() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestSocketClose() failed.got one false, expected 2*true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestSocketClose() failed.got true, expected false.")
	}
	pZmqSocket.Close()
}

func TestReq(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestReq() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestReq() failed.got one false, expected 2*true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestReq() failed.got true, expected false.")
	}
	pZmqSocket.SetChoose()

	err, res := pZmqSocket.Req("/root/D/A")
	if err != nil {
		t.Errorf("TestReq() failed.got %s, expected nil.", err.Error())
	}
	if !(res == nil || len(res.Endpoints) == 0) {
		t.Errorf("TestReq() failed.got !0, expected 0.")
	}

	err1, res1 := pZmqSocket.Req("/root")
	if err1 != nil {
		t.Errorf("TestReq() failed.got %s, expected nil.", err1.Error())
	}
	if res1 == nil || len(res1.Endpoints) == 0 {
		t.Errorf("TestReq() failed.got 0, expected !0.")
	}
}

func TestReq2(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10345", "tcp://10.15.144.105:10356", pzmq, 1); err != nil {
		t.Errorf("TestReq2() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if pZmqSocket.mreqOK || pZmqSocket.msubOK {
		t.Errorf("TestReq2() failed.got one true, expected 2*false.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestReq2() failed.got true, expected false.")
	}
	pZmqSocket.mok = 2
	pZmqSocket.SetChoose()

	err, res := pZmqSocket.Req("/root")
	if err == nil {
		t.Errorf("TestReq2() failed.got %s, expected nil.", err.Error())
	}
	fmt.Println(err.Error())
	if !(res == nil || len(res.Endpoints) == 0) {
		t.Errorf("TestReq2() failed.got !0, expected 0.")
	}
}

func TestRegister(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestRegister() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestRegister() failed.got one false, expected 2*true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestRegister() failed.got true, expected false.")
	}
	pZmqSocket.SetChoose()

	err := pZmqSocket.Register("/root")
	if err != nil {
		t.Errorf("TestRegister() failed.got %s, expected nil.", err.Error())
	}

	//	time.Sleep(time.Second * 100)
}

func TestUnRegister(t *testing.T) {
	//	return
	pzmq := &Zmq{}
	pEvent := &Event{}
	pEvent.Init(nil, "/root/sds")
	pzmq.Init(&pEvent.mcache, pEvent)

	pZmqSocket := &ZmqSocket{}
	if err := pZmqSocket.Init("tcp://10.15.144.105:10300", "tcp://10.15.144.105:10301", pzmq, 1); err != nil {
		t.Errorf("TestUnRegister() failed.got %s, expected nil.", err.Error())
	}
	time.Sleep(time.Second * 10)
	if !pZmqSocket.mreqOK || !pZmqSocket.msubOK {
		t.Errorf("TestUnRegister() failed.got one false, expected 2*true.")
	}
	if pZmqSocket.mChoose {
		t.Errorf("TestUnRegister() failed.got true, expected false.")
	}

	pZmqSocket.SetChoose()

	err := pZmqSocket.Register("/root")
	if err != nil {
		t.Errorf("TestUnRegister() failed.got %s, expected nil.", err.Error())
	}

	time.Sleep(time.Second * 10)
	fmt.Println("start unregister /root")
	err2 := pZmqSocket.UnRegister("/root")
	if err2 != nil {
		t.Errorf("TestUnRegister() failed.got %s, expected nil.", err2.Error())
	}

	//	time.Sleep(time.Second * 100)
}
