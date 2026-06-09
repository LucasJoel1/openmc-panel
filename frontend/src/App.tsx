import { useEffect, useState } from "react";
import { ThemeProvider } from "./components/theme-provider";
import { SidebarInset, SidebarProvider } from "./components/ui/sidebar";
import AppSidebar from "./components/AppSidebar";
import AppHeader from "./components/AppHeader";
import * as Icons from "lucide-react";
import type { NavSection } from "./types";
import PageManager from "./components/PageManager";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import LoginPage from "./components/pages/Login";
import { checkPermissions, Permissions } from "./lib/utils";

function App() {
    const [selection, setSelection] = useState<number>(0);
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
    const [username, setUsername] = useState<string>("");
    const [permissions, setPermissions] = useState<bigint>(0n)
    const [navSections, setNavSections] = useState<NavSection[]>([])
    const queryClient = new QueryClient();

    useEffect(() => {
        const localStorageUsername = localStorage.getItem("username")
        const permissions = localStorage.getItem("permissions")
        if (localStorageUsername !== null) {
            setIsLoggedIn(true)
            setUsername(localStorageUsername)
            if (permissions !== null) {
                setPermissions(BigInt(permissions))
            } else {
                setPermissions(0n)
            }
        }
    }, [isLoggedIn])

    useEffect(() => {
        setNavSections([])
        const tempNavSections: NavSection[] = []

        if (permissions !== 0n) {
            if (checkPermissions(permissions, Permissions.VIEW_SERVER_STATUS)) {
                tempNavSections.push({
                    id: 0,
                    icon: <Icons.ChartBar size={48} color="#ffffff" strokeWidth={1.5} />,
                    name: "Status",
                    sideBar: true
                })
            }
            if (checkPermissions(permissions, Permissions.READ_SERVER_LOGS)) {
                tempNavSections.push({
                    id: 1,
                    icon: <Icons.Terminal size={48} color="#ffffff" strokeWidth={1.5} />,
                    name: "Console",
                    sideBar: true
                })
            }
            if (checkPermissions(permissions, Permissions.VIEW_USERS)) {
                tempNavSections.push({
                    id: 100,
                    icon: <Icons.Users size={48} color="#ffffff" strokeWidth={1.5} />,
                    name: "User Management",
                    sideBar: false
                })
            }

            tempNavSections.push(
            {
                id: 2,
                icon: <Icons.User size={48} color="#ffffff" strokeWidth={1.5} />,
                name: "Player Management",
                sideBar: true
            },
            {
                id: 3,
                icon: <Icons.Pencil size={48} color="#ffffff" strokeWidth={1.5} />,
                name: "Edit Properties",
                sideBar: true
            },
            {
                id: 4,
                icon: <Icons.FolderOpen size={48} color="#ffffff" strokeWidth={1.5} />,
                name: "File Manager",
                sideBar: true
            })
        }
        setNavSections(tempNavSections)

        // Set selection to first available section if current selection is not valid
        if (tempNavSections.length > 0 && !tempNavSections.find(section => section.id === selection)) {
            setSelection(tempNavSections[0].id)
        }
    }, [permissions])

    useEffect(() => {
        const handleHashChange = () => {
            const hash = window.location.hash.replace("#", "");
            const match = navSections.find(
                (page) => page.name.replace(" ", "-") === hash
            );

            if (match) {
                setSelection(match.id);
            } else if (navSections.length > 0) {
                // If no match found but we have sections, select the first one
                setSelection(navSections[0].id);
            }
        };

        handleHashChange();

        window.addEventListener("hashchange", handleHashChange);

        return () => {
            window.removeEventListener("hashchange", handleHashChange);
        };
    }, [navSections]);

    return (
        <QueryClientProvider client={queryClient}>
            <ThemeProvider defaultTheme="dark">
                {isLoggedIn ?
                    <SidebarProvider>
                        <AppSidebar
                            navSections={navSections}
                            selectedItem={selection}
                            setSelectedItem={setSelection}
                            username={username}
                        />
                        <SidebarInset>
                            <AppHeader section={navSections.find(section => section.id === selection)?.name || navSections.length <= 0 && "No Pages Allowed" || "Loading..."} />
                            <PageManager page={selection} permissions={permissions} />
                        </SidebarInset>
                    </SidebarProvider> :
                    <LoginPage setIsLoggedIn={setIsLoggedIn} />
                }
                <Toaster
                    theme={
                        ["dark", "light", "system"].includes(
                            localStorage.getItem("theme") || ""
                        )
                            ? (localStorage.getItem("theme") as "dark" | "light" | "system")
                            : "system"
                    }
                />
            </ThemeProvider>
        </QueryClientProvider>
    );
}

export default App;
