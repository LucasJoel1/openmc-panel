import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "../ui/sheet";
import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";
import { Separator } from "../ui/separator";
import { FileText, Download, Trash2, RefreshCw } from "lucide-react";
import { useEffect, useState } from "react";

interface historicalLog {
    id: number,
    date: string,
    time: string,
    size: string
}


// retireved from api
interface logType {
    name: string
    size: number
}

export default function LogsHistory() {
    const [historicalLogs, setHistoricalLogs] = useState<historicalLog[]>([]);

    const getLogs = () => {
        fetch("/api/savedLogs", {
            method: "GET"
        }).then((res) => res.json())
            .then((data) => {
                const logs: historicalLog[] = []
                const items = data.logs
                items.forEach((element: logType) => {
                    const dateObj = new Date(Number(element.name))
                    const date = dateObj.toLocaleDateString(undefined, {
                        year: "numeric",
                        month: "short",
                        day: "2-digit"
                    })
                    const time = dateObj.toLocaleTimeString(undefined, {
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit"
                    })
                    const item: historicalLog = {
                        id: parseInt(element.name),
                        date,
                        time,
                        size: String((element.size / 1000).toFixed(2)) + "kb"
                    }
                    logs.push(item)
                });
                setHistoricalLogs(logs)
            })
    }

    useEffect(() => {
        getLogs()
    }, [])

    const handleLogLoad = (logId: number) => {
        console.log("Loading log:", logId);
    };

    const handleLogDelete = (logId: number) => {
        console.log("Deleting log:", logId);
    };

    return (
        <Sheet>
            <SheetTrigger asChild>
                <Button>Logs History</Button>
            </SheetTrigger>
            <SheetContent className="w-[90vw] sm:w-[400px] md:w-[500px] flex flex-col">                <SheetHeader>
                <SheetTitle>Logs History</SheetTitle>
                <SheetDescription>Select a historical log to view.</SheetDescription>
            </SheetHeader>
            <div className="px-6 pb-3">
                <Button
                    size="sm"
                    variant="outline"
                    className="w-full"
                    onClick={() => getLogs()}
                >
                    <RefreshCw className="h-4 w-4 mr-2" />
                        Refresh
                </Button>
            </div>
            <ScrollArea className="flex-1 h-full w-full">
                <div className="p-4 pb-12">
                    {historicalLogs.map((log, index) => (
                        <div key={log.id}>
                            <div className="flex items-center justify-between py-3 px-2 rounded-md hover:bg-muted/50 transition-colors">
                                <div className="flex items-center space-x-3">
                                    <FileText className="h-4 w-4 text-muted-foreground" />
                                    <div>
                                        <p className="font-medium text-sm">{log.date}</p>
                                        <p className="text-xs text-muted-foreground">{log.time} • {log.size}</p>
                                    </div>
                                </div>
                                <div className="flex gap-2">
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        onClick={() => handleLogLoad(log.id)}
                                        className="h-7 px-3"
                                    >
                                        <Download className="h-3 w-3 mr-1" />
                                            Load
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="destructive"
                                        onClick={() => handleLogDelete(log.id)}
                                        className="h-7 px-3"
                                    >
                                        <Trash2 className="h-3 w-3 mr-1" />
                                            Delete
                                    </Button>
                                </div>
                            </div>
                            {index < historicalLogs.length - 1 && <Separator />}
                        </div>
                    ))}
                </div>
            </ScrollArea>
            </SheetContent>
        </Sheet>
    );
}