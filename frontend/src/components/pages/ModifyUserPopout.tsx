import { Permissions } from "@/lib/utils";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { ScrollArea } from "../ui/scroll-area";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger } from "../ui/sheet";
import { Checkbox } from "../ui/checkbox";
import { useState } from "react";
import type { User } from "./UserManagementPage";

interface ModifyUserPopoutProps {
    user: User;
    loggedInUserPermissionLevel: bigint;
}

export default function ModifyUserPopout(props: ModifyUserPopoutProps) {
    const permissionEntries = Object.entries(Permissions);
    const [permissionStates, setPermissionStates] = useState<boolean[]>(
        new Array(permissionEntries.length).fill(false)
    );

    const handlePermissionChange = (index: number, checked: boolean) => {
        const newStates = [...permissionStates];
        newStates[index] = checked;
        setPermissionStates(newStates);
    };

    return (
        <Sheet>
            <SheetTrigger>
                <Button variant="outline" disabled={props.user.isAdmin && localStorage.getItem("username") !== props.user.username}>Permissions</Button>
            </SheetTrigger>
            <SheetContent className="sm:max-w-md flex flex-col h-full">
                <SheetHeader>
                    <SheetTitle className="text-xl">Manage Permissions</SheetTitle>
                    <SheetDescription>Modify permissions for {props.user.username}</SheetDescription>
                </SheetHeader>

                <div className="flex-1 min-h-0 relative p-4">
                    <div className="absolute inset-0 px-4 pb-20">
                        <SheetTitle className="text-lg mb-3">Permissions</SheetTitle>
                        <ScrollArea className="h-full w-full border rounded-md p-4">
                            <div className="space-y-3">
                                {permissionEntries.map(([perm, value], index) => (
                                    <div key={String(value)} className="flex items-center space-x-3">
                                        <Checkbox
                                            id={`perm-${value}`}
                                            checked={permissionStates[index]}
                                            onCheckedChange={(checked) => handlePermissionChange(index, checked as boolean)}
                                        />
                                        <Label
                                            htmlFor={`perm-${value}`}
                                            className="text-sm font-normal cursor-pointer flex-1"
                                        >
                                            {perm.replace(/_/g, " ").toLowerCase().replace(/\b\w/g, l => l.toUpperCase())}
                                        </Label>
                                    </div>
                                ))}
                            </div>
                        </ScrollArea>
                    </div>
                </div>

                <div className="p-4 border-t">
                    <Button className="w-full">Update Permissions</Button>
                </div>
            </SheetContent>
        </Sheet>
    )
}