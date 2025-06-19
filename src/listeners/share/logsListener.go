package sharedListeners

import (
	"bufio"
	"context"
	"fmt"
	"openmc-panel/src/globals"
	"regexp"
	"time"
)

func LogsListener(ctx context.Context) {
    pipes := globals.GetPipes()
    done := make(chan bool)

    joinMatcher, _ :=  regexp.Compile(`^\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: \w+\[\/\d{1,3}(?:\.\d{1,3}){3}:\d+\] logged in with entity id \d+ at \((-?\d+\.\d+, -?\d+\.\d+, -?\d+\.\d+)\)$`)
    leaveMatcher, _ := regexp.Compile(`^\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: \w+ left the game$`)
    chatMatcher, _ := regexp.Compile(`^\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: <\w+> [\s\S]*$`)
    onlineMatch, _ := regexp.Compile(`^\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO]: Done \(\d.\d{3}s\)! For help, type "help"$`)

    go func() {
        stopLogging := make(chan bool)
        go func() {
            globals.StartLogging()
            stopLogging <- true
        }()

        scanner := bufio.NewScanner(pipes.Stdout)
        for {
            select {
            case <-ctx.Done():
                globals.SetServerRunning(false)
                done <- true
                return
            default:
                if scanner.Scan() {
                    line := scanner.Text()
					globals.AppendBuffer(line + "\n")
					if joinMatcher.MatchString(line) {
                        OnConnect(line)
					} else if leaveMatcher.MatchString(line) {
                        OnDisconnect(line)
                    } else if chatMatcher.MatchString(line) {
                        OnGameMessage(line)
                    } else if onlineMatch.MatchString(line) {
                        globals.SetServerRunning(true)
                        globals.LockServer(false)
                        globals.SetStartTime(time.Now().Unix())
                    }
                } else {
                    if err := scanner.Err(); err != nil {
                        fmt.Printf("Error reading stdout: %v\n", err)
                    }
                    globals.SetServerRunning(false)
                    globals.SetStartTime(0)
                    done <- true
                    return
                }
            }
        }
    }()

    <-done
}