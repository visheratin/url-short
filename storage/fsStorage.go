package storage

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/visheratin/url-short/log"
)

// FSStorage implements storing of input-code pairs in the directory
// on the local disk in separate files. dataPath specifies a full
// path to the directory with input-code data.
type FSStorage struct {
	dataPath string
}

// newFSStorage creates an instance of FSStorage.
func newFSStorage(dataPath string) (FSStorage, error) {
	storage := FSStorage{
		dataPath: dataPath,
	}
	if _, err := os.Stat(storage.dataPath); os.IsNotExist(err) {
		err = os.MkdirAll(storage.dataPath, 0777)
		if err != nil {
			log.Log().Error.Println(err)
			return FSStorage{}, err
		}
	}
	return storage, nil
}

// Store creates a file with the name equal to the code
// and contents equal to the input string.
func (storage FSStorage) Store(code, input string) error {
	filename := path.Join(storage.dataPath, code)
	f, err := os.Create(filename)
	if err != nil {
		log.Log().Error.Println(err)
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(input))
	if err != nil {
		log.Log().Error.Println(err)
		return err
	}
	return nil
}

// LoadAll scans through all files in the directory of
// FSStorage instance and extracts input-code pairs.
func (storage FSStorage) LoadAll() ([][2]string, error) {
	f, err := os.Open(storage.dataPath)
	if err != nil {
		log.Log().Error.Println(err)
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Log().Error.Println(err)
		return nil, err
	}
	counter := 0
	result := make([][2]string, len(files))
	for _, info := range files {
		_, code := filepath.Split(info.Name())
		filepath := path.Join(storage.dataPath, code)
		input, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Log().Error.Println(err)
			continue
		}
		pair := [2]string{code, string(input)}
		result[counter] = pair
		counter++
	}
	// If there were any problems with reading files of the
	// directory we need to shrink the result slice.
	if counter < (len(files) - 1) {
		result = result[0:counter]
	}
	return result, nil
}
