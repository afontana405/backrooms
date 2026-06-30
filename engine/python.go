package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var (
	pyCmd    *exec.Cmd
	pyStdin  io.WriteCloser
	pyReader *bufio.Reader
	pyMu     sync.Mutex
	pyReady  bool
)

type pyRequest struct {
	Module string      `json:"module"`
	Method string      `json:"method"`
	Args   interface{} `json:"args"`
}

type pyResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// StartPython spawns the Python API subprocess using system Python.
func StartPython(scriptPath string) error {
	return StartPythonWith("python", scriptPath)
}

// StartPythonWith spawns the Python API subprocess with a specific Python executable.
func StartPythonWith(pythonExe, scriptPath string) error {
	pyCmd = exec.Command(pythonExe, scriptPath)
	pyCmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	pyCmd.Stderr = os.Stderr // Python logs go to our stderr

	stdin, err := pyCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("python stdin pipe: %w", err)
	}
	pyStdin = stdin

	stdout, err := pyCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("python stdout pipe: %w", err)
	}
	pyReader = bufio.NewReader(stdout)

	if err := pyCmd.Start(); err != nil {
		return fmt.Errorf("python start: %w", err)
	}

	pyReady = true
	return nil
}

// StopPython closes the Python subprocess cleanly.
func StopPython() {
	if !pyReady {
		return
	}
	pyReady = false
	pyStdin.Close()
	pyCmd.Wait()
}

// CallPython sends a request to the Python subprocess and returns the response data.
// args should be a map[string]interface{} (kwargs) or []interface{} (positional).
// Calls are serialized — Python processes one request at a time.
func CallPython(module, method string, args interface{}) (interface{}, error) {
	if !pyReady {
		return nil, fmt.Errorf("python bridge not started")
	}

	pyMu.Lock()
	defer pyMu.Unlock()

	req := pyRequest{Module: module, Method: method, Args: args}
	line, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("python marshal: %w", err)
	}

	if _, err := fmt.Fprintf(pyStdin, "%s\n", line); err != nil {
		return nil, fmt.Errorf("python write: %w", err)
	}

	respLine, err := pyReader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("python read: %w", err)
	}

	var resp pyResponse
	if err := json.Unmarshal([]byte(respLine), &resp); err != nil {
		return nil, fmt.Errorf("python response parse: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%s", resp.Error)
	}

	return resp.Data, nil
}
