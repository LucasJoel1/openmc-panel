import { useState, useRef, type FormEvent } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { checkPassword } from "@/lib/authVerification";
import { toast } from "sonner";
import type { User } from "./UserManagementPage";

interface ChangePasswordDialogProps {
    user: User;
}

export default function ChangePasswordDialog(props: ChangePasswordDialogProps) {
    const [newPassword, setNewPassword] = useState<string>("")
    const [confirmPassword, setConfirmPassword] = useState<string>("")
    const closeButtonRef = useRef<HTMLButtonElement>(null)

    const submitChangePasswordForm = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()

        if (newPassword !== confirmPassword) {
            toast.error("Passwords do not match")
            return
        }

        const passwordCheckOutput = checkPassword(newPassword)
        if (passwordCheckOutput !== "") {
            toast.error(passwordCheckOutput)
            return
        }

        // TODO: Implement password change API call
        fetch(`/api/changePassword`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: props.user.username,
                newPassword: newPassword,
                confirmPassword: confirmPassword,
            })
        })
            .then(async (res) => {
                if (res.status !== 200) {
                    const errorText = await res.text();
                    throw new Error(errorText || "An unexpected error has occurred");
                }
                toast.success(`Password changed for ${props.user.username}`)
                setNewPassword("")
                setConfirmPassword("")
                closeButtonRef.current?.click() // Close dialog on success
            })
            .catch((err) => {
                console.error("error changing password: " + err)
                toast.error(err.message || "An unknown error occurred")
            })
    }

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="outline" size="sm" disabled={props.user.isAdmin && localStorage.getItem("username") !== props.user.username}>
                    Change Password
                </Button>
            </DialogTrigger>
            <DialogContent className="space-y-6">
                <form onSubmit={(e) => submitChangePasswordForm(e)} className="space-y-6">
                    <DialogHeader>
                        <DialogTitle>Change Password</DialogTitle>
                        <DialogDescription>Change the password for {props.user.username}</DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-6">
                        <div className="grid gap-3">
                            <Label htmlFor="newPassword">New Password</Label>
                            <Input
                                id="newPassword"
                                type="password"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                            />
                        </div>
                        <div className="grid gap-3">
                            <Label htmlFor="confirmPassword">Confirm Password</Label>
                            <Input
                                id="confirmPassword"
                                type="password"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <DialogClose asChild>
                            <Button ref={closeButtonRef} variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button type="submit">Change Password</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
