import LoginForm from "./LoginForm";
import logo from '../../assets/logo.svg'

interface LoginPageProps {
    setIsLoggedIn: (value: boolean) => void;
}

export default function LoginPage({ setIsLoggedIn }: LoginPageProps) {
    return (
        <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
            <div className="flex w-full max-w-sm flex-col gap-6">
                <a href="#" className="flex items-center gap-2 self-center font-medium">
                    <div className="size-6">
                        <img src={logo} className="size-6" />
                    </div>
                    OPENMC PANEL
                </a>
                <LoginForm setIsLoggedIn={setIsLoggedIn} />
            </div>
        </div>
    )
}