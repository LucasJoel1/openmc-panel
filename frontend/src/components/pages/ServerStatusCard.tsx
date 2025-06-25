import { useEffect, useMemo, useState } from "react";
import RequestButton from "../RequestButton";
import {
    Card,
    CardHeader,
    CardDescription,
    CardTitle,
    CardContent,
    CardFooter,
} from "../ui/card";
import { type ServerStatusInterface } from "./StatusPage";
import { checkPermissions, formatTime, Permissions } from "@/lib/utils";

interface ServerStatusProps {
  data: ServerStatusInterface;
  permissions: bigint;
}

interface ParsedServerData {
  status: string;
  statusColour: string;
  uptime: string;
  memory: string;
  cpu: string;
}

export default function ServerStatusCard(props: ServerStatusProps) {
    const [parsedData, setparsedData] = useState<ParsedServerData | null>({
        status: "Loading...",
        statusColour: "bg-gray-500",
        uptime: "0",
        memory: "0",
        cpu: "",
    });

    const canStartServer = useMemo(() => 
        checkPermissions(props.permissions, Permissions.START_SERVER),
        [props.permissions]
    )

    const canStopServer = useMemo(() =>
        checkPermissions(props.permissions, Permissions.STOP_SERVER),
        [props.permissions]
    )

    const canRestartServer = useMemo(() =>
        checkPermissions(props.permissions, Permissions.RESTART_SERVER),
        [props.permissions]
    )

    useEffect(() => {
        const data: ParsedServerData = {
            status: "ERROR",
            statusColour: "bg-red-300",
            uptime: "-1",
            memory: "0",
            cpu: "0%",
        };

        switch (props.data.state) {
        case 0:
            data.status = "Online";
            data.statusColour = "bg-green-500";
            break;
        case 1:
            data.status = "Starting...";
            data.statusColour = "bg-yellow-300";
            break;
        case 2:
            data.status = "Offline";
            data.statusColour = "bg-gray-500";
            break;
        }

        data.uptime = formatTime(props.data.uptime);

        data.memory = props.data.RAMUsage.toFixed(2);

        data.cpu = props.data.CPU.toFixed(2) + "%";

        setparsedData(data);
    }, [props.data]);

    return (
        <Card className="bg-muted/50 rounded-xl">
            <CardHeader>
                <CardDescription>Server Status</CardDescription>
                <CardTitle className="text-3xl font-bold flex items-center gap-2">
                    <span
                        className={"w-6 h-6 rounded-full " + parsedData?.statusColour}
                    ></span>
                    {parsedData?.status}
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 grid grid-cols-2">
                <div>
                    <div>
                        <p className="text-sm text-muted-foreground">Uptime</p>
                        <p className="font-medium">{parsedData?.uptime} </p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">CPU Usage</p>
                        <p className="font-medium">{parsedData?.cpu}</p>
                    </div>
                </div>
                <div>
                    <div>
                        <p className="text-sm text-muted-foreground">TPS</p>
                        <p className="font-medium">{props.data.TPS}</p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">RAM Usage</p>
                        <p className="font-medium">
                            {parsedData?.memory}GB / {props.data.RamAlloc}GB
                        </p>
                    </div>
                </div>
            </CardContent>
            <CardFooter className={"flex flex-row gap-4"}>
                {canStartServer && <RequestButton contents="Start" req="/api/startServer" />}
                {canStopServer && <RequestButton contents="Stop" req="/api/stopServer" />}
                {canRestartServer && <RequestButton contents="Restart" req="/api/restartServer" />}
            </CardFooter>
        </Card>
    );
}
