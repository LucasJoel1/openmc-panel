import { ScrollArea } from "../ui/scroll-area";
import { Button } from "../ui/button";
import {
    Card,
    CardHeader,
    CardDescription,
    CardTitle,
    CardContent,
    CardFooter,
} from "../ui/card";
import { Input } from "../ui/input";
import { useEffect, useRef, useState } from "react";
import { Separator } from "../ui/separator";
import React from "react";
import LogsHistory from "./LogsHistory";

interface ConsolePageProps {
    visible: boolean;
    isPage: boolean;
}

// interface LoadedLog {
//     date: Date
//     index: number
//     data: string
// }

export default function ConsolePage(props: ConsolePageProps) {
    const [logs, setLogs] = useState<string[]>([]);
    const [loadedLog, setLoadedLog] = useState<string>("")
    const [useLive, setUseLive] = useState<boolean>(true)
    const [autoscroll, setAutoscroll] = useState<boolean>(true);
    const [command, setCommand] = useState<string>("")
    const ws = useRef<WebSocket | null>(null);
    const lastElement = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        ws.current = new WebSocket(
            "ws://" +
            window.location.href.split("#")[0].split("://")[1] +
            "api/ws/serverLogs"
        );

        ws.current.onopen = () => {
            console.log("connected to server logs ws");
        };

        ws.current.onmessage = (event) => {
            setLogs((prevLogs) => [...prevLogs, event.data]);
        };

        return () => {
            ws.current?.close();
        };
    }, []);

    useEffect(() => {
        if (autoscroll && logs.length > 0 && lastElement.current !== null) {
            lastElement.current.scrollIntoView({ behavior: "smooth", block: "end" });
        }
    }, [logs.length, autoscroll]);

    const executeCommand = () => {
        if (!command.trim()) return;

        fetch(`/api/executeCommand?command=${encodeURIComponent(command)}`, {
            method: "POST"
        })
            .then(response => {
                if (!response.ok) {
                    console.error('Failed to execute command');
                }
            })
            .catch(error => {
                console.error('Error executing command:', error);
            });

        setCommand("");
    }

    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            executeCommand();
        }
    }

    return (
        <Card
            className={
                props.visible !== true
                    ? "hidden"
                    : "bg-muted/50 rounded-xl flex flex-col h-full"
            }
        >
            <CardHeader className="flex flex-row items-center justify-between">
                <div className="w-full">
                    <CardDescription>Server Logs</CardDescription>
                    <div className="flex justify-between w-full">
                        <CardTitle className="text-2xl font-bold">
                            Recent Activity
                        </CardTitle>
                        <div className="flex items-center gap-3">
                            {!useLive && <Button
                                variant="default"
                                onClick={() => setUseLive(true)}
                            >
                                Show Live
                            </Button>
                            }
                            <Button
                                variant={autoscroll ? "default" : "outline"}
                                onClick={() => setAutoscroll(!autoscroll)}
                            >
                                Autoscroll
                            </Button>
                            {props.isPage && <LogsHistory setLoadedLogs={setLoadedLog} setUseLive={setUseLive} />}
                        </div>
                    </div>
                </div>{" "}
            </CardHeader>{" "}
            <CardContent className="flex-1 min-h-0 relative">
                <div className="absolute inset-0 px-6">
                    <ScrollArea className="h-full w-full rounded-md border">
                        <div className="p-4">
                            {useLive ? logs.map((log, index) => (
                                <React.Fragment key={log}>
                                    <pre className="text-sm whitespace-pre-wrap">{log}</pre>
                                    <Separator
                                        className="my-2"
                                        ref={index === logs.length - 1 ? lastElement : null}
                                    />
                                </React.Fragment>
                            ))
                                :
                                <pre className="text-sm whitespace-pre-wrap">{loadedLog}</pre>
                            }
                        </div>
                    </ScrollArea>
                </div>
            </CardContent>            <CardFooter>
                <div className="flex w-full items-center space-x-2">
                    <Input
                        placeholder="Enter a command..."
                        type="text"
                        value={command}
                        onChange={(e) => setCommand(e.target.value)}
                        onKeyDown={handleKeyPress}
                    />
                    <Button variant={"outline"} onClick={() => executeCommand()}>Send</Button>
                </div>
            </CardFooter>
        </Card>
    );
}
