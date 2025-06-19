package server

import (
	"fmt"
	"openmc-panel/src/globals"
	"os/exec"

	"github.com/shirou/gopsutil/v4/process"
)

func StartServer() (*exec.Cmd, error) {
    globals.LockServer(true)
    data := globals.GetServerSettings()

    if data == nil {
        return nil, fmt.Errorf("server settings cannot be retrieved")
    }

    args := []string{"-Xmx" + data.Memory[1], "-Xms" + data.Memory[0], "-jar", data.JarName}

    if !data.GUI {
        args = append(args, "nogui")
    }

    cmd := exec.Command("java", args...)
	cmd.Dir = data.Path
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("stdout pipe not available: %v", err)
    }

    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("stdin pipe not available: %v", err)
    }

    globals.SetPipes(stdout, stdin)

    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("failed to start server: %v", err)
    }

    proc, err := process.NewProcess(int32(cmd.Process.Pid))

    if err != nil {
        return nil, fmt.Errorf("error retrieving pid")
    }

    globals.SetProcess(proc)

    globals.SetServerVersions(globals.GetServerType())

    return cmd, nil
}