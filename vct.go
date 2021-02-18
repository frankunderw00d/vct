package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	// version control structure
	VC struct {
		Version string `json:"version"`
	}
)

const (
	VersionFileName    = "version.json"
	TmpVersionFileName = "versionTemporary.txt"
	BaseVersion        = "1.0.0"
)

var (
	setVersion = flag.String("v", "", "set version number for current version")
)

func (vc *VC) VersionFileExists() bool {
	_, err := os.Stat(VersionFileName)
	if err != nil {
		return false
	}
	return true
}

func (vc *VC) CreateBasicVersionFile() error {
	file, err := os.OpenFile(VersionFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	vc.Version = BaseVersion

	data, err := json.Marshal(vc)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (vc *VC) LoadVersionFile() error {
	data, err := ioutil.ReadFile(VersionFileName)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, vc)
}

func (vc *VC) Upgrade() string {
	subs := strings.Split(vc.Version, ".")
	d, _ := strconv.Atoi(subs[2])
	subs[2] = strconv.Itoa(d + 1)

	vc.Version = strings.Join(subs, ".")

	return vc.Version
}

func (vc *VC) Store() error {
	file, err := os.OpenFile(VersionFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return err
	}

	data, err := json.Marshal(vc)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func init() {
	flag.Parse()
}

func main() {
	vc := VC{}
	currentVersion := BaseVersion

	if *setVersion != "" { // 设置版本号
		vc.Version = *setVersion
		currentVersion = *setVersion
		if err := vc.Store(); err != nil {
			log.Fatal(err)
		}
	} else { // 不设置版本号
		if !vc.VersionFileExists() {
			if err := vc.CreateBasicVersionFile(); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := vc.LoadVersionFile(); err != nil {
				log.Fatal(err)
			}
			currentVersion = vc.Upgrade()
			if err := vc.Store(); err != nil {
				log.Fatal(err)
			}
		}
	}

	// 存储临时版本文件
	if err := StoreTemporary(currentVersion); err != nil {
		log.Println(err.Error())
	}
}

func StoreTemporary(version string) error {
	file, err := os.OpenFile(TmpVersionFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return err
	}

	_, err = file.Write([]byte(version))
	return err
}