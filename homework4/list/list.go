package list

import (
	"homework4/model"
	"io/ioutil"
	"path/filepath"
)

type FindFolder struct {
	Dir string
}


func (s FindFolder) GetList(extension string)([]model.File, error) {
	files, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}

	stFile := make([]model.File, 0)

	for _, f := range files{
		name:= f.Name()
		filext := filepath.Ext(name)
		size := f.Size()
		if filext == extension || extension == ""{
			stFile = append(stFile, model.File{Name: name, Extension: filext, Size: size})

		}	
		return stFile, nil

	}

	

	return stFile, nil
	
	
}
