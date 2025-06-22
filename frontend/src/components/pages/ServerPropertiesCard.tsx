import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "../ui/card";
import type { ServerPropertiesInterface } from "./StatusPage";

interface ServerPropertiesProps {
  data: ServerPropertiesInterface;
}

export default function ServerPropertiesCard(props: ServerPropertiesProps) {
    return (
        <Card className="bg-muted/50 rounded-xl">
            <CardHeader>
                <CardDescription>Server Properties</CardDescription>
                <CardTitle className="text-3xl font-bold">1.21.5</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 grid grid-cols-2">
                <div>
                    <div>
                        <p className="text-sm text-muted-foreground">Gamemode</p>
                        <p className="font-medium">
                            {String(props.data.GameMode).charAt(0).toUpperCase() +
                String(props.data.GameMode).slice(1)}
                        </p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">IP Address</p>
                        <span className="font-medium">{props.data.IP}</span>
                        <span className="text-muted-foreground">:{props.data.Port}</span>
                    </div>
                </div>
                <div>
                    <div>
                        <p className="text-sm text-muted-foreground">Difficulty</p>
                        <p className="font-medium">
                            {String(props.data.Difficulty).charAt(0).toUpperCase() +
                String(props.data.Difficulty).slice(1)}
                        </p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">Modloader</p>
                        <span className="font-medium">{String(props.data.ModLoader)} </span>
                        <span className="text-muted-foreground">
                            {String(props.data.ModLoaderVersion)}
                        </span>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
