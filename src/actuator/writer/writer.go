package writer

import "strconv"
import "io/ioutil"

import "logging"

/*		Writer is analogous to Reader: 
 *		- Convert our standard units to those of the actual sink (nothing to 
 *		  do yet)
 *		- Abstract away concrete writing details (e.g. need to "enable" pwm 
 *		  output)
 *
 *		We only have 1 type of sinks, so currently Writer is a struct and 
 *		not an interface. It is a concrete "LinuxPwmWriter"
 */


type Writer struct {
	path string 
}

func LinuxPwmWriter(path string) *Writer {
	W := &Writer {path: path}
	err := ioutil.WriteFile(W.path+"_enable", []byte{'0'}, 0644)
	if err !=nil {
		panic("Could not enable Linux PWM Writer: "+err.Error())
	}
	return W
}

func (W *Writer) Write (value int) {
	valuestr := strconv.Itoa(value)
	err := ioutil.WriteFile(W.path+"_enable", []byte{'0'}, 0644)
	if err !=nil {
		panic("Write error: Could not enable Linux PWM Writer: "+err.Error())
	}
	logging.Debug("Writing %s <- %s\n", W.path, valuestr)
	err = ioutil.WriteFile(W.path, []byte(valuestr), 0644)
	if err !=nil {
		panic("Write error: "+err.Error())
	}
}