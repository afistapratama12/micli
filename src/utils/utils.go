package utils

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/afistapratama12/micli/src/constants"
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

func CheckOrCreateFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	defer f.Close()
	return nil
}

func ModifyFileCache(f *os.File, data []byte) error {
	var err error

	// replace all data in file
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// 3 scenario
// if file exist and data exist, return
// if file not exist, create file and add data, return
// if file exist but data not exist, add data, return
func ReadFileCache(path string) (f *os.File, res model.CacheData, err error) {
	err = CheckOrCreateFile(path)
	if err != nil {
		return
	}

	f, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, res, err
	}

	// defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, res, err
	}

	if len(string(b)) < 1 {
		// write file
		data := strings.Join(constants.DEFAULT_PAIR, ";")

		_, err = f.Write([]byte(data))
		if err != nil {
			return nil, res, err
		}

		res.ListPair = constants.DEFAULT_PAIR
	} else {
		res.ListPair = strings.Split(string(b), ";")
	}

	return f, res, nil
}

func CreateMapPairWS(listPair []string) map[string]string {
	var mapPairWS = make(map[string]string)

	for _, pair := range listPair {
		mapPairWS[strings.Replace(pair, "_", "", 1)] = pair
	}

	return mapPairWS
}

func RemoveItems[V comparable](orig []V, itemsToRemove []V) []V {
	/// change itemsToRemove to map, for O(n) complexity
	itemsMap := make(map[V]struct{})

	for _, item := range itemsToRemove {
		itemsMap[item] = struct{}{}
	}

	result := []V{}

	// check every item in orig
	for _, item := range orig {
		if _, ok := itemsMap[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// return true if item not exist in orig
func CompareData[V comparable](orig []V, itemsToCheck []V) (bool, []V) {
	var isItemNotExist = false
	var itemsNotExist []V

	origMap := make(map[V]struct{})

	for _, item := range orig {
		origMap[item] = struct{}{}
	}

	for _, item := range itemsToCheck {
		if _, ok := origMap[item]; !ok {
			isItemNotExist = true
			itemsNotExist = append(itemsNotExist, item)
		}
	}

	return isItemNotExist, itemsNotExist
}
