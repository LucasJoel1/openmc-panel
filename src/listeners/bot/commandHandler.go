package botListners

import (
	"bufio"
	"context"
	"fmt"
	"openmc-panel/src/globals"
	"openmc-panel/src/listeners/share"
	"openmc-panel/src/server"
)

func CommandHandler(message string) string {
	switch message {
		case "start":
            if globals.GetServerRunning() {
                return "Server already online"
            }
            
            _, err := server.StartServer()
            if err != nil {
                return fmt.Sprintf("Failed to start server: %v", err)
            }
            
            // Create context with cancel
            ctx := context.WithoutCancel(context.Background())
            
            go func() {
                sharedListeners.LogsListener(ctx)
            }()
            
            return "Server started successfully"

		case "stop":
            if !globals.GetServerRunning() {
                return "Server is not online"
            }
            writer := bufio.NewWriter(globals.GetPipes().Stdin)
            globals.SetServerRunning(false)
            writer.WriteString("stop" + "\n")
            writer.Flush()
            return "Server stopping"
			
		default:
			return "command not found"
	}
}