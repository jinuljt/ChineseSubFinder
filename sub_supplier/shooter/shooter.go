package shooter

import (
	"crypto/md5"
	"fmt"
	"github.com/allanpk716/ChineseSubFinder/common"
	"github.com/allanpk716/ChineseSubFinder/sub_supplier"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Supplier struct {

}

func NewSupplier() *Supplier {
	return &Supplier{}
}

func (s Supplier) GetSubListFromFile(filePath string, httpProxy string) ([]sub_supplier.SubInfo, error) {

	const qLan = "Chn"
	var outSubInfoList []sub_supplier.SubInfo
	var jsonList []SublistShooter

	hash, err := s.computeFileHash(filePath)
	if err != nil {
		return nil, err
	}
	if hash == "" {
		return nil, common.ShooterFileHashIsEmpty
	}

	fileName := filepath.Base(filePath)

	httpClient := common.NewHttpClient(httpProxy)

	_, err = httpClient.R().
		SetFormData(map[string]string{
			"filehash": hash,
			"pathinfo": fileName,
			"format": "json",
			"lang": qLan,
		}).
		SetResult(&jsonList).
		Post(common.SubShooterRootUrl)
	if err != nil {
		return nil, err
	}
	for _, shooter := range jsonList {
		for _, file := range shooter.Files {
			subExt := file.Ext
			if strings.Contains(file.Ext, ".") == false {
				subExt = "." + subExt
			}
			outSubInfoList = append(outSubInfoList, *sub_supplier.NewSubInfo(fileName, qLan, "", file.Link, 0, shooter.Delay, subExt))
		}
	}
	return outSubInfoList, nil
}

func (s Supplier) GetSubListFromKeyword(keyword string, httpProxy string) ([]sub_supplier.SubInfo, error) {
	panic("not implemented")
}

func (s Supplier) computeFileHash(filePath string) (string, error) {
	hash := ""
	fp, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	stat, err := fp.Stat()
	if err != nil {
		return "", err
	}
	size := float64(stat.Size())
	if size < 0xF000 {
		return "", common.VideoFileIsTooSmall
	}
	samplePositions := [4]int64{
		4 * 1024,
		int64(math.Floor(size / 3 * 2)),
		int64(math.Floor(size / 3)),
		int64(size - 8*1024)}
	var samples [4][]byte
	for i, position := range samplePositions {
		samples[i] = make([]byte, 4*1024)
		_, err = fp.ReadAt(samples[i], position)
		if err != nil {
			return "", err
		}
	}
	for _, sample := range samples {
		if len(hash) > 0 {
			hash += ";"
		}
		hash += fmt.Sprintf("%x", md5.Sum(sample))
	}

	return hash, nil
}


type FilesShooter struct {
	Ext  string `json:"ext"`
	Link string `json:"link"`
}
type SublistShooter struct {
	Desc  string         `json:"desc"`
	Delay int64          `json:"delay"`
	Files []FilesShooter `json:"files"`
}