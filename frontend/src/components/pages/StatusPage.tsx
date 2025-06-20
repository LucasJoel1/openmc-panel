import { useEffect, useRef, useState } from "react";
import ServerStatusCard from "./ServerStatusCard";
import PlayerStatusCard from "./PlayerStatusCard";
import ServerPropertiesCard from "./ServerPropertiesCard";
import ConsolePage from "./ConsolePage";

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
            console.log("connected to server status ws");
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
        <div
            className={`${
                props.visible !== true ? "hidden" : "flex flex-col gap-4 h-full"
            }`}
        >
            <div className="grid auto-rows-min gap-4 md:grid-cols-3">
                {serverStatusData && <ServerStatusCard data={serverStatusData} />}
                {playerData && <PlayerStatusCard data={playerData} />}
                {serverPropertiesData && (
                    <ServerPropertiesCard data={serverPropertiesData} />
                )}
            </div>{" "}
            <div className="flex-1 min-h-0">
                <ConsolePage visible={true} isPage={false} />
            </div>
        </div>
    );
}
