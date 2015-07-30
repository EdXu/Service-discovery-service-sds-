package sdssdk

import (
	"fmt"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	return
	pSDK := &SDSSDK{}
	if err := pSDK.Start(" tcp://10.15.144.105:10300 ; tcp://10.15.144.105:10300 "); err != nil {
		fmt.Println("TestMutex() failed.got %s, expected nil.", err.Error())
		return
	}
	flag := true

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
		f2 := true
		for flag {
			select {
			case <-ch0:
				f2 = true
				fmt.Println("==========1")
			default:
				if f2 {
					f2 = false
				} else {
					//	flag = false
				}
				fmt.Println("==========12")
				time.Sleep(time.Second * 20)
			}
		}
	}()

	go func() {
		f2 := true
		for flag {
			select {
			case <-ch1:
				f2 = true
				fmt.Println("==========2")
			default:
				if f2 {
					f2 = false
				} else {
					//	flag = false
				}
				fmt.Println("==========22")
				time.Sleep(time.Second * 20)
			}
		}
	}()

	go func() {
		f2 := true
		for flag {
			select {
			case <-ch2:
				f2 = true
				fmt.Println("==========3")
			default:
				if f2 {
					f2 = false
				} else {
					//	flag = false
				}
				fmt.Println("==========33")
				time.Sleep(time.Second * 20)
			}
		}
	}()

	go func() {
		for flag {
			err4, ch4 := pSDK.RegistServiceInfo("/root")
			if err4 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err4.Error())
				flag = false

			}
			time.Sleep(time.Second * 2)
			err5 := pSDK.UnRegistServiceInfo("/root", ch4)
			if err5 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err5.Error())
				flag = false
			}
			fmt.Println("++++++++++++++1")
		}

	}()

	go func() {
		for flag {
			err4, ch4 := pSDK.RegistServiceInfo("/root/test")
			if err4 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err4.Error())
				flag = false
			}
			time.Sleep(time.Second * 1)
			err5 := pSDK.UnRegistServiceInfo("/root/test", ch4)
			if err5 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err5.Error())
				flag = false
			}
			fmt.Println("++++++++++++++2")
		}
	}()

	go func() {
		for flag {
			err4, ch4 := pSDK.RegistServiceInfo("/root/A/B")
			if err4 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err4.Error())
				flag = false
			}
			time.Sleep(time.Second * 2)
			err5 := pSDK.UnRegistServiceInfo("/root/A/B", ch4)
			if err5 != nil {
				fmt.Println("TestMutex() failed.got %s, expected nil.", err5.Error())
				flag = false
			}
			fmt.Println("++++++++++++++3")
		}
	}()

	go func() {
		for flag {
			err, po := pSDK.GetServiceInfo("/root")
			if err == nil {
				fmt.Println("---1 getServiceInfo ", po.GetNodeName())
			} else {
				fmt.Println("---1 getServiceInfo err=", err.Error())
			}
			time.Sleep(time.Second * 2)
			fmt.Println("--------------------1")
		}

	}()

	go func() {
		for flag {
			err, _ := pSDK.GetServiceInfo("/root/test")
			if err != nil {

			}
			time.Sleep(time.Second * 1)
			fmt.Println("--------------------2")
		}

	}()

	go func() {
		for flag {
			err, _ := pSDK.GetServiceInfo("/root")
			if err != nil {

			}
			time.Sleep(time.Second * 1)
			fmt.Println("--------------------3")
		}

	}()

	go func() {
		for flag {
			pSDK.GetAllServiceInfo("/root/test")

			time.Sleep(time.Second * 1)
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@1")
		}

	}()

	go func() {
		for flag {
			pSDK.GetAllServiceInfo("/root/test/A")

			time.Sleep(time.Second * 1)
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@2")
		}

	}()
	go func() {
		for flag {
			pSDK.GetAllServiceInfo("/root")

			time.Sleep(time.Second * 1)
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@3")
		}

	}()

	for flag {
		time.Sleep(time.Second * 60)
	}
	fmt.Println("end because flag == false")
	pSDK.Stop()
}
