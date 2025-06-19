import { useEffect, useMemo, useState } from "react";
import { ThemeProvider } from "./components/theme-provider";
import { SidebarInset, SidebarProvider } from "./components/ui/sidebar";
import AppSidebar from "./components/AppSidebar";
import AppHeader from "./components/AppHeader";
import * as Icons from 'lucide-react'
import type { NavSection } from "./types";
import PageManager from "./components/PageManager";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";

function App() {
    const [selection, setSelection] = useState<number>(0);
    const queryClient = new QueryClient()

    const navSections: NavSection[] = useMemo(() => [
        {
            id: 0,
            icon: <Icons.ChartBar size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "Status",
        },
        {
            id: 1,
            icon: <Icons.Server size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "Server Manager",
        },
        {
            id: 2,
            icon: <Icons.Terminal size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "Console",
        },
        {
            id: 3,
            icon: <Icons.User size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "Player Management",
        },
        {
            id: 4,
            icon: <Icons.Pencil size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "Edit Properties",
        },
        {
            id: 5,
            icon: <Icons.FolderOpen size={48} color="#ffffff" strokeWidth={1.5} />,
            name: "File Manager",
        },
    ], []);

    useEffect(() => {
        const handleHashChange = () => {
            const hash = window.location.hash.replace("#", "")
            const match = navSections.find((page) => page.name.replace(" ", "-") === hash);

            if (match) {
                setSelection(match.id);
            }
        }

        handleHashChange()

        window.addEventListener("hashchange", handleHashChange)

        return () => {
            window.removeEventListener("hashchange", handleHashChange)
        }
    }, [navSections])

    return (
        <QueryClientProvider client={queryClient}>
            <ThemeProvider defaultTheme="dark">
                <SidebarProvider>
                    <AppSidebar navSections={navSections} selectedItem={selection} setSelectedItem={setSelection} />
                    <SidebarInset>
                        <AppHeader section={navSections[selection].name} />
                        <PageManager page={selection} />
                    </SidebarInset>
                </SidebarProvider>
                <Toaster
                    theme={
                        (["dark", "light", "system"].includes(localStorage.getItem("theme") || "")
                            ? (localStorage.getItem("theme") as "dark" | "light" | "system")
                            : "system")
                    }
                />
            </ThemeProvider>
        </QueryClientProvider>
    );
}

export default App;
