import { useEffect, useMemo, useState } from "react"
import { Button } from "../ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card"
// import { ScrollArea } from "../ui/scroll-area"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table"
import { toast } from "sonner"
import { checkPermissions, Permissions } from "@/lib/utils"
import CreateUserForm from "./CreateUserForm"
import ModifyUserPopout from "./ModifyUserPopout"
import ChangePasswordDialog from "./ChangePasswordDialog"

interface UserManagementProps {
    permissions: bigint
    visible: boolean
}

export interface User {
    username: string
    isAdmin: boolean
}

export default function UserManagementPage(props: UserManagementProps) {
    const [users, setUsers] = useState<User[]>([])

    const canCreateUsers = useMemo(() =>
        checkPermissions(props.permissions, Permissions.CREATE_USER),
        [props.permissions]
    )

    const canDeleteUsers = useMemo(() =>
        checkPermissions(props.permissions, Permissions.DELETE_USER),
        [props.permissions]
    )

    const canModifyUsers = useMemo(() =>
        checkPermissions(props.permissions, Permissions.MODIFY_USER),
        [props.permissions]
    )

    useEffect(() => {
        fetch("/api/getUsers")
            .then((res) => res.json())
            .then((data) => {
                setUsers(data)
            })
            .catch(err => {
                console.error("error fetching users", err)
                toast(err)
            })
    }, [])

    return (
        <Card
            className={
                props.visible !== true
                    ? "hidden"
                    : "bg-muted/50 rounded-xl flex flex-col h-full"
            }
        >
            <CardHeader>
                <CardTitle>Manage Users</CardTitle>
                <CardDescription>Manage, delete, create and modify panel users.</CardDescription>
            </CardHeader>            <CardContent className="flex-1 min-h-0 relative flex flex-col">
                <div className="flex-1 min-h-0">
                    <Table>
                        <TableHeader>
                            <TableRow className="flex">
                                <TableHead className="flex-1">Username</TableHead>
                                <TableHead className="pr-26.5">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {
                                users.map(user => (
                                    <TableRow key={user.username} className="flex">
                                        <TableCell className="flex-1 flex items-center">{user.username}</TableCell>
                                        <TableCell className="flex-auto flex items-center justify-end">
                                            <div className="flex gap-2">
                                                {canModifyUsers && <ModifyUserPopout user={user} loggedInUserPermissionLevel={props.permissions} />}
                                                {canModifyUsers && <ChangePasswordDialog user={user} />}
                                                {canDeleteUsers && <Button variant="destructive" disabled={user.isAdmin}>Delete</Button>}
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                ))
                            }
                        </TableBody>
                    </Table>
                </div>
                {canCreateUsers &&
                    <div className="flex justify-end mt-auto pt-4 border-t">
                        <CreateUserForm setUsers={setUsers} />
                    </div>
                }
            </CardContent>
        </Card>
    )
}