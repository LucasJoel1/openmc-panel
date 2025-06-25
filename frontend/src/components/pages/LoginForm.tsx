import { useState, type FormEvent } from "react";
import { Button } from "../ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { toast } from "sonner";

interface LoginFormProps {
    setIsLoggedIn: (value: boolean) => void;
}

export default function LoginForm(props: LoginFormProps) {
    const [username, setUsername] = useState<string>("")
    const [password, setPassword] = useState<string>("")

    const handleLogin = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        fetch("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        })
            .then((res) => res.json())
            .then((data) => {
                localStorage.setItem("username", data.username)
                localStorage.setItem("permissions", BigInt(data.permissions).toString())
                props.setIsLoggedIn(true)
                toast(`Logged in as ${data.username}`)
            })
            .catch(error => {
                console.error(error)
                toast(error)
            })
    }

    return (
        <div className={"flex flex-col gap-6"}>
            <Card>
                <CardHeader className="text-center">
                    <CardTitle className="text-xl">Login</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={(e) => handleLogin(e)}>
                        <div className="grid gap-6">
                            <div className="grid gap-6">
                                <div className="grid gap-3">
                                    <Label htmlFor="email">Username</Label>
                                    <Input
                                        id="username"
                                        type="text"
                                        required
                                        onChange={(e) => setUsername(e.target.value)}
                                    />
                                </div>
                                <div className="grid gap-3">
                                    <div className="flex items-center">
                                        <Label htmlFor="password">Password</Label>
                                    </div>
                                    <Input id="password" type="password" onChange={(e) => setPassword(e.target.value)} required />
                                </div>
                                <Button type="submit" className="w-full">
                                    Login
                                </Button>
                            </div>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    )
}