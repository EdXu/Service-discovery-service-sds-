package sdssdk

import (
	"fmt"
	"testing"
	"time"

	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
)

func TestSDKStart(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start(""); err == nil {
		t.Errorf("TestStart() failed.got nil, expected err.")
	}
}

func TestSDKStart2(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start(" tcp://10.15.144.105:10300 ; tcp://10.15.144.105:10300 "); err != nil {
		t.Errorf("TestStart2() failed.got %s, expected nil.", err.Error())
	}
}

func TestSDKStart3(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start(" ;;  tcp://10.15.144.105:10357 ; tcp://10.15.144.105:10300 "); err != nil {
		t.Errorf("TestStart3() failed.got %s, expected nil.", err.Error())
	}
}

func TestSDKStart4(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10357 ; tcp://10.15.144.105:10378 "); err == nil {
		t.Errorf("TestStart4() failed.got nil, expected err.")
	}
}

func TestSDKStart5(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300 ; tcp://10.15.144.105:10300 "); err != nil {
		t.Errorf("TestStart5() failed.got %s, expected nil.", err.Error())
	}
}

func TestSDKStart6(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestStart6() failed.got %s, expected nil.", err.Error())
	}
}

func TestSDKStop(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestStop() failed.got %s, expected nil.", err.Error())
	}
	if err := pSDK.Stop(); err != nil {
		t.Errorf("TestStop() failed.got %s, expected nil.", err.Error())
	}
}

func TestSdkStop2(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10345"); err == nil {
		t.Errorf("TestStop2() failed.got nil, expected err.")
	}
	if err := pSDK.Stop(); err != nil {
		t.Errorf("TestStop() failed.got %s, expected nil.", err.Error())
	}
}

func TestSdkGetServiceInfo(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestGetServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	if err, point := pSDK.GetServiceInfo("/root/sds"); err != nil {
		t.Errorf("TestGetServiceInfo() failed.got %s, expected nil.", err.Error())
	} else {
		fmt.Println("point= ", point)
	}
	if err, _ := pSDK.GetServiceInfo("/root/A/B/Dhg"); err == nil {
		t.Errorf("TestGetServiceInfo() failed.got nil, expected err.")
	}
}

func TestSdkGetAllServiceInfo(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestGetAllServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	ppoints0 := pSDK.GetAllServiceInfo("/root/sds")
	if ppoints0 == nil || len(ppoints0) == 0 {
		t.Errorf("TestGetAllServiceInfo() failed.got nil, expected !0.")
	}
	ppoints := pSDK.GetAllServiceInfo("/root/A/B/Dhg")
	if ppoints != nil && len(ppoints) != 0 {
		t.Errorf("TestGetAllServiceInfo() failed.got %d, expected 0.", len(ppoints))
	}
}

func TestSdkGetAllServiceInfo2(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10389"); err == nil {
		t.Errorf("TestGetAllServiceInfo2() failed.got nil, expected err.")
	}
	if ppoints := pSDK.GetAllServiceInfo("/root/sds"); ppoints != nil {
		t.Errorf("TestGetAllServiceInfo2() failed.got !nil, expected nil.")
	}
	if ppoints := pSDK.GetAllServiceInfo("/root/A/B/Dhg"); ppoints != nil {
		t.Errorf("TestGetAllServiceInfo2() failed.got !nil, expected nil.")
	}
}

func TestSdkRegistServiceInfo(t *testing.T) {
	//	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestRegistServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	{
		err, ch := pSDK.RegistServiceInfo("/root")
		if err != nil {
			t.Errorf("TestRegistServiceInfo() failed.got %s, expected nil.", err.Error())
		}
		err2, ch2 := pSDK.RegistServiceInfo("/root")
		if err2 != nil {
			t.Errorf("TestRegistServiceInfo() failed.got %s, expected nil.", err2.Error())
		}
		time.Sleep(time.Second * 20)
		k := 0
		q := 0
		for flag := true; flag; {
			select {
			case <-ch:
				k++
				//				fmt.Println("/root info=", p.GetNodeName(), p.GetLoading())
			case <-ch2:
				q++
				//				fmt.Println("/root info2=", p.GetNodeName(), p.GetLoading())
			default:
				flag = false
			}
		}
		if k != q {
			t.Errorf("TestRegistServiceInfo() failed.got k!=q, expected k=q.")
		}
	}

	{
		err, ch := pSDK.RegistServiceInfo("/root/A/A/B/D")
		if err != nil {
			t.Errorf("TestRegistServiceInfo() failed.got %s, expected nil.", err.Error())
		}
		time.Sleep(time.Second * 20)
		k := 0
		for flag := true; flag; {
			select {
			case <-ch:
				k++
				//				fmt.Println("/root info=", p.GetNodeName(), p.GetLoading())
			default:
				flag = false
			}
		}
		if k > 0 {
			t.Errorf("TestRegistServiceInfo() failed.got !0, expected 0.")
		}
	}

}

func TestSdkUnRegistServiceInfo(t *testing.T) {
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		t.Errorf("TestUnRegistServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	ch0 := make(chan dzhyun.SDSEndpoint, 1024)
	pSDK.UnRegistServiceInfo("/root", ch0)

	err, ch := pSDK.RegistServiceInfo("/root")
	if err != nil {
		t.Errorf("TestUnRegistServiceInfo() failed.got %s, expected nil.", err.Error())
	}
	err2, ch2 := pSDK.RegistServiceInfo("/root")
	if err2 != nil {
		t.Errorf("TestUnRegistServiceInfo() failed.got %s, expected nil.", err2.Error())
	}
	time.Sleep(time.Second * 20)
	k := 0
	q := 0
	for flag := true; flag; {
		select {
		case <-ch:
			k++
			//				fmt.Println("/root info=", p.GetNodeName(), p.GetLoading())
		case <-ch2:
			q++
			//				fmt.Println("/root info2=", p.GetNodeName(), p.GetLoading())
		default:
			flag = false
		}
	}
	if k != q {
		t.Errorf("TestUnRegistServiceInfo() failed.got k!=q, expected k=q.")
	}

	pSDK.UnRegistServiceInfo("/root", ch2)
	time.Sleep(time.Second * 20)
	k = 0
	for flag := true; flag; {
		select {
		case <-ch:
			k++
			//				fmt.Println("/root info=", p.GetNodeName(), p.GetLoading())
		default:
			flag = false
		}
	}
	if k == 0 {
		t.Errorf("TestUnRegistServiceInfo() failed.got 0, expected !0.")
	}
}

func BenchmarkSdkGetServiceInfo(b *testing.B) {
	//	return
	b.StopTimer()
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		return
	}
	err0, ch0 := pSDK.RegistServiceInfo("/root/test")
	if err0 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err0.Error())
		return
	}

	err1, ch1 := pSDK.RegistServiceInfo("/root")
	if err1 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err1.Error())
		return
	}

	err2, ch2 := pSDK.RegistServiceInfo("/root")
	if err2 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err2.Error())
		return
	}

	go func() {
		for {
			select {
			case <-ch0:

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch1:

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch2:

			}
		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pSDK.GetServiceInfo("/root")
	}
}

func BenchmarkSdkGetAllServiceInfo(b *testing.B) {
	//	return
	b.StopTimer()
	pSDK := &SDSSDK{}
	if err := pSDK.Start("tcp://10.15.144.105:10300"); err != nil {
		return
	}
	err0, ch0 := pSDK.RegistServiceInfo("/root/test")
	if err0 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err0.Error())
		return
	}

	err1, ch1 := pSDK.RegistServiceInfo("/root")
	if err1 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err1.Error())
		return
	}

	err2, ch2 := pSDK.RegistServiceInfo("/root")
	if err2 != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err2.Error())
		return
	}

	go func() {
		for {
			select {
			case <-ch0:

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch1:

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch2:

			}
		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pSDK.GetAllServiceInfo("/root")
	}
}
