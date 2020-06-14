package reader

import "io/ioutil"
import "strconv"


type FileReader struct {
	path string
}

func FromFile(path string) Reader {
	return &FileReader{path: path}
}

func (R *FileReader) Read() float64 {
	out, err := ioutil.ReadFile(R.path)
    if err!=nil {
		panic("read error on "+R.path+": "+err.Error())
	}
	if out[len(out)-1] == '\n' {	// strip trailing \n
		out = out [0:len(out)-1]
	}
	res, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		panic("Convertion error on "+string(out)+": "+err.Error())
	}	
	return res
}