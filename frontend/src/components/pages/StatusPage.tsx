import { useEffect, useRef, useState } from "react";
import { Button } from "../ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "../ui/card";
import { Input } from "../ui/input";
import ServerStatusCard from "./ServerStatusCard";
import PlayerStatusCard from "./PlayerStatusCard";
import ServerPropertiesCard from "./ServerPropertiesCard";

interface StatusPageProps {
  visible: boolean;
}

export interface ServerData {
  Status: ServerStatusInterface;
  Players: PlayersInterface;
  Properties: ServerPropertiesInterface;
}

export interface ServerStatusInterface {
  state: number;
  uptime: number;
  TPS: number;
  CPU: number;
  RAMUsage: number;
  RamAlloc: number;
}

export interface PlayersInterface {
  online: number;
  max: number;
  list: Player[] | null;
}

export interface Player {
  ping: number;
  time: number;
  name: string;
}

export interface ServerPropertiesInterface {
  version: string;
  GameMode: number;
  Difficulty: number;
  IP: string;
  Port: number;
  ModLoader: string;
  ModLoaderVersion: string;
}

export default function StatusPage(props: StatusPageProps) {

    const [serverStatusData, setServerStatusData] =
    useState<ServerStatusInterface | null>(null);
    const [playerData, setPlayersData] = useState<PlayersInterface | null>(null);
    const [serverPropertiesData, setServerPropertiesData] =
    useState<ServerPropertiesInterface | null>(null);

    const ws = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (!props.visible) return;

        ws.current = new WebSocket(
            "ws://" +
        window.location.href.split("#")[0].split("://")[1] +
        "api/ws/serverInfo"
        );

        ws.current.onopen = () => {
            console.log("connected to ws");
        };

        ws.current.onmessage = (event) => {
            const data: ServerData = JSON.parse(event.data);

            setServerStatusData(data.Status);
            setPlayersData(data.Players);
            setServerPropertiesData(data.Properties);
        };

        return () => {
            ws.current?.close();
        };
    }, [props.visible]);

    useEffect(() => {});

    return (
        <main className={`${props.visible !== true ? "hidden" : "block"} p-4`}>
            <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
                {serverStatusData && <ServerStatusCard data={ serverStatusData } />}
                {playerData && <PlayerStatusCard data={playerData} />}
                {serverPropertiesData && <ServerPropertiesCard data={serverPropertiesData} />}
            </div>
            <Card className="mt-4">
                <CardHeader className="flex flex-row items-center justify-between">
                    <div>
                        <CardDescription>Server Logs</CardDescription>
                        <CardTitle className="text-2xl font-bold">
              Recent Activity
                        </CardTitle>
                    </div>
                </CardHeader>
                <CardContent>
                    
                </CardContent>
                <CardFooter>
                    <div className="flex w-full items-center space-x-2">
                        <Input placeholder="Enter a command..." type="text" />
                        <Button variant={"outline"}>Send</Button>
                    </div>
                </CardFooter>
            </Card>
        </main>
    );
}
