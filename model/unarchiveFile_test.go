package model

import (
	"testing"
)

func TestUnArchiveFile(t *testing.T) {

	desRoot := "C:\\Tmp"
	//file := "C:\\Tmp\\123.zip"
	//file := "C:\\Tmp\\456.rar"
	file := "C:\\Tmp\\789.tar"
	err := UnArchiveFile(file, desRoot)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnArr(t *testing.T) {
	desRoot := "C:\\Tmp"
	//file := "C:\\Tmp\\[subhd]_0_162236051219240.zip"
	file := "C:\\Tmp\\123.zip"
	//file := "C:\\Tmp\\Tmp.7z"
	//file := "C:\\Tmp\\[zimuku]_0_[zmk.pw]奥斯陆.Oslo.[WEB.1080P]中英文字幕.zip"
	err := unArr7z(file, desRoot)
	if err != nil {
		t.Fatal(err)
	}
}