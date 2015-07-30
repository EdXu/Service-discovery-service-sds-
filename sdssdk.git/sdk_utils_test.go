package sdssdk

import "testing"

func TestFilteContainStr(t *testing.T) {
	arrIn := make([]string, 0)
	arrIn = append(arrIn, "/root/A")
	arrIn = append(arrIn, "/root/A/B")
	arrIn = append(arrIn, "/root/C")
	arrIn = append(arrIn, "/root/D")
	arrIn = append(arrIn, "/root/F")
	arrIn = append(arrIn, "/root/D/A")
	parr := FilteContainStr(&arrIn)
	if len(*parr) != 2 {
		t.Errorf("TestFilteContainStr() failed.got %d, expected 2.", len(*parr))
	}
	for i := 0; i < len(*parr); i++ {
		if !((*parr)[i] == 1 || (*parr)[i] == 5) {
			t.Errorf("TestFilteContainStr() failed.got %d, expected 1,5.", (*parr)[i])
		}
	}
}

func TestFilteContainStr1(t *testing.T) {
	arrIn := make([]string, 0)
	arrIn = append(arrIn, "/root/A")
	arrIn = append(arrIn, "/root/A/B")
	arrIn = append(arrIn, "/root/C")
	arrIn = append(arrIn, "/root/D")
	arrIn = append(arrIn, "/root/F")
	arrIn = append(arrIn, "/root/D/A")
	arrIn = append(arrIn, "/root/")
	parr := FilteContainStr(&arrIn)
	if len(*parr) != 6 {
		t.Errorf("TestFilteContainStr() failed.got %d, expected 6.", len(*parr))
	}
	for i := 0; i < len(*parr); i++ {
		if !((*parr)[i] >= 0 || (*parr)[i] <= 5) {
			t.Errorf("TestFilteContainStr() failed.got %d, expected 0-5.", (*parr)[i])
		}
	}
}

func TestFilteContainStr2(t *testing.T) {
	arrIn := make([]string, 0)
	arrIn = append(arrIn, "/root/A")
	arrIn = append(arrIn, "/root/A/B")
	arrIn = append(arrIn, "/root/C")
	arrIn = append(arrIn, "/root/D/")
	arrIn = append(arrIn, "/root/F")
	arrIn = append(arrIn, "/root/D/A")
	parr := FilteContainStr(&arrIn)
	if len(*parr) != 2 {
		t.Errorf("TestFilteContainStr() failed.got %d, expected 2.", len(*parr))
	}
	for i := 0; i < len(*parr); i++ {
		if !((*parr)[i] >= 0 || (*parr)[i] <= 5) {
			t.Errorf("TestFilteContainStr() failed.got %d, expected 1,5.", (*parr)[i])
		}
	}
}
