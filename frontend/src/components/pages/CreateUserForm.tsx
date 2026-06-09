import { useState, useRef, type FormEvent } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { checkPassword, checkUsername } from "@/lib/authVerification";
import { toast } from "sonner";
import type { User } from "./UserManagementPage";

interface CreateUserFormProps {
    setUsers: React.Dispatch<React.SetStateAction<User[]>>
}

export default function CreateUserForm(props: CreateUserFormProps) {
    const [username, setUsername] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [confirmPassword, setConfirmPassword] = useState<string>("")
    const closeButtonRef = useRef<HTMLButtonElement>(null)

    const submitCreateUserForm = (e: FormEvent<HTMLFormElement>) => {
        console.log("submit")
        e.preventDefault()
        if (password !== confirmPassword) {
            toast.error("Passwords do not match")
            return
        }

        const usernameCheckOutput = checkUsername(username)
        if (usernameCheckOutput !== "") {
            toast.error(usernameCheckOutput)
            return
        }

        const passwordCheckOutput = checkPassword(password)
        if (passwordCheckOutput !== "") {
            toast.error(passwordCheckOutput)
            return
        }

        fetch("/api/registerUser", {
            method: "POST",
            body: JSON.stringify({
                username: username,
                password: password,
                confirmPassword: confirmPassword,
            })
        })
            .then(async (res) => {
                if (res.status !== 200) {
                    const errorText = await res.text();
                    throw new Error(errorText || "An unexpected error has occured");
                }                fetch("/api/getUsers")
                    .then((res) => res.json())
                    .then((data) => {
                        toast(`User ${username} created`)
                        props.setUsers(data)
                        closeButtonRef.current?.click() // Close dialog on success
                    })
                    .catch(err => {
                        console.error("error fetching users", err)
                        toast.error(typeof err === "string" ? err : "an unknown error has occured fetching users")
                    })
            }).catch((err) => {
                console.error("error creating user: " + err)
                toast.error(err.message || "An unknown error occurred")
            })
    }

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="default">Create User</Button>
            </DialogTrigger>
            <DialogContent>
                <form onSubmit={(e) => submitCreateUserForm(e)}>
                    <DialogHeader className="mb-4">
                        <DialogTitle>Create New User</DialogTitle>
                        <DialogDescription>Create a new user for the panel.</DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4">
                        <div className="grid gap-3">
                            <Label htmlFor="username">Username</Label>
                            <Input id="username" type="text" onChange={(e) => setUsername(e.target.value)} />
                        </div>
                        <div className="grid gap-3">
                            <Label htmlFor="password">Password</Label>
                            <Input id="password" type="password" onChange={(e) => setPassword(e.target.value)} />
                        </div>
                        <div className="grid gap-3">
                            <Label htmlFor="confirmPassword">Confirm Password</Label>
                            <Input id="confirmPassword" type="password" onChange={(e) => setConfirmPassword(e.target.value)} />
                        </div>
                    </div>                    
                    <DialogFooter className="mt-4">
                        <DialogClose asChild>
                            <Button ref={closeButtonRef} variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button type="submit">Save changes</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}