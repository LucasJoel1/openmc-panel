import { formatTime } from "@/lib/utils";
import {
    Card,
    CardDescription,
    CardHeader,
    CardTitle,
    CardContent,
} from "../ui/card";
import {
    Table,
    TableHeader,
    TableRow,
    TableHead,
    TableBody,
    TableCell,
} from "../ui/table";
import type { PlayersInterface } from "./StatusPage";

interface PlayerStatusProps {
  data: PlayersInterface;
}

export default function PlayerStatusCard(props: PlayerStatusProps) {
    return (
        <Card className="bg-muted/50 rounded-xl">
            <CardHeader>
                <CardDescription>Players</CardDescription>
                <CardTitle className="text-3xl font-bold">
                    {props.data.online} / {props.data.max}
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[40%]">Username</TableHead>
                            <TableHead className="w-[30%]">Ping</TableHead>
                            <TableHead className="w-[30%]">Time</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {props.data.list?.map((player) => {
                            return (
                                <TableRow key={player.name}>
                                    <TableCell className="font-medium">{player.name}</TableCell>
                                    <TableCell>{player.ping}ms</TableCell>
                                    <TableCell>{formatTime(player.time)}</TableCell>
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
}
