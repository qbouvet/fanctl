package reader

import "os/exec"
import "io"
import "bytes"
import "strconv"

import "logging"


/*
 *	Helper around exec.Cmd
 */
type CmdLine struct {
	Cmd  string
	Args []string 
}

func NewCmdLine(s ...string) *CmdLine {
	if len(s) <1 {
		panic("Not enough arguments")
	}
	res := &CmdLine {Cmd: s[0]}
	if len(s) >1 {
		res.Args = s[1:len(s)]
	}
	return res
}

func (C *CmdLine) ToExec() *exec.Cmd {
	return exec.Command(C.Cmd, C.Args...)
}



/*
 * 	CmdReader reads a value from one command or several piped commands
 */
type CmdReader struct {
	CmdStrs []*CmdLine
	running []*exec.Cmd
}

func FromCmd(c ...*CmdLine) Reader {
	if len(c) <1 {
		panic("Too few arguments")
	}
	return &CmdReader {CmdStrs: c, running: make([]*exec.Cmd, len(c))}
}

func (R *CmdReader) Read () float64 {
	if len(R.CmdStrs) <=0 {
		panic("len(cmds) <=0")
	}
	var out []byte
	if len(R.CmdStrs) ==1 {
	// Simple case with single command
		logging.Debug("Executing: %v\n", R.CmdStrs[0])
		_out, err := R.CmdStrs[0].ToExec().CombinedOutput() 
		if err !=nil {
			panic("exec error: "+err.Error())
		}
		out = _out
	} else {
	// Several piped commands
		// Instanciate exec.Cmds
		for i,cmdstr := range(R.CmdStrs) {
			R.running[i] = cmdstr.ToExec()
		}
		// 1. Build piping
		for i,_ := range(R.running) {
			pipeRd, pipeWr := io.Pipe()
			R.running[i].Stdout = pipeWr
			if i < len(R.running)-1 {
				R.running[i+1].Stdin = pipeRd
			}
		}
		var stdout bytes.Buffer
		R.running[len(R.running)-1].Stdout = &stdout
		// 2. Run and wait
		for i,_ := range(R.running) {
			R.running[i].Start()
		}
		for i,_ := range(R.running) {
			R.running[i].Wait()
		 	if i > 0 {						
				// Close stdin (pipe) unless first command
				logging.Debug("Closing stdin of piped command %d\n", i) 
				R.running[i].Stdin.(*io.PipeReader).Close()		
		 	}
			if i < len(R.running)-1 { 
				// Close stdout (pipe) unless last command
				logging.Debug("Closing stdout of piped command %d\n", i) 
		 		R.running[i].Stdout.(*io.PipeWriter).Close()	
		 	}
		}
		out = stdout.Bytes()
	}
	// Parse output to produce float
	if len(out) <= 0 {
		panic("exec error: nothing returned")
	}
	if out[len(out)-1] == '\n' {	// strip trailing \n
		out = out [0:len(out)-1]
	}
	res, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		panic("Can't convert: "+string(out))
	}
	return res
}


