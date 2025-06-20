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

interface ConsolePageProps {
  visible: boolean;
  isPage: boolean;
}

export default function ConsolePage(props: ConsolePageProps) {
    const [logs, setLogs] = useState<string[]>([]);
    const [autoscroll, setAutoscroll] = useState<boolean>(true);
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
                            <Button
                                variant={autoscroll ? "default" : "outline"}
                                onClick={() => setAutoscroll(!autoscroll)}
                            >
            Autoscroll
                            </Button>
                            <Button>Logs History</Button>
                        </div>
                    </div>
                </div>{" "}
            </CardHeader>{" "}
            <CardContent className="flex-1 min-h-0 relative">
                <div className="absolute inset-0 px-6">
                    <ScrollArea className="h-full w-full rounded-md border">
                        <div className="p-4">
                            {logs.map((log, index) => (
                                <React.Fragment key={log}>
                                    <pre className="text-sm whitespace-pre-wrap">{log}</pre>
                                    <Separator
                                        className="my-2"
                                        ref={index === logs.length - 1 ? lastElement : null}
                                    />
                                </React.Fragment>
                            ))}
                        </div>
                    </ScrollArea>
                </div>
            </CardContent>
            <CardFooter>
                <div className="flex w-full items-center space-x-2">
                    <Input placeholder="Enter a command..." type="text" />
                    <Button variant={"outline"}>Send</Button>
                </div>
            </CardFooter>
        </Card>
    );
}
