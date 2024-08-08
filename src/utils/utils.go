package utils

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/afistapratama12/micli/constants"
	"github.com/afistapratama12/micli/src/model"
)

func ReadMessage(b []byte, dest interface{}) error {
	err := json.Unmarshal(b, dest)
	return err
}

func RunCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 3 scenario
// if file exist and data exist, return
// if file not exist, create file and add data, return
// if file exist but data not exist, add data, return
func ReadFileCache(path string) (res model.CacheData, err error) {
	var f *os.File

	f, err = os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(path)
			if err != nil {
				return res, err
			}
		} else {
			return res, err
		}
	}

	f.Close()

	f, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return res, err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return res, err
	}

	if len(string(b)) < 1 {
		// write file
		data := strings.Join(constants.DEFAULT_PAIR, ";")

		_, err = f.Write([]byte(data))
		if err != nil {
			return res, err
		}

		res.ListPair = constants.DEFAULT_PAIR
	} else {
		res.ListPair = strings.Split(string(b), ";")
	}

	return res, nil
}

func CreateMapPairWS(listPair []string) map[string]string {
	var mapPairWS = make(map[string]string)

	for _, pair := range listPair {
		mapPairWS[strings.Replace(pair, "_", "", 1)] = pair
	}

	return mapPairWS
}
