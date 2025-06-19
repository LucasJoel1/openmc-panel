import type { Dispatch, SetStateAction } from "react";
import type { NavSection } from "../types";
import {
    Sidebar,
    SidebarGroup,
    SidebarGroupContent,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "./ui/sidebar";
import logo from '../assets/logo.svg'

interface AppSidebarProps {
    navSections: NavSection[];
    selectedItem: number;
    setSelectedItem: Dispatch<SetStateAction<number>>;
}

export default function AppSidebar(props: AppSidebarProps) {
    return (
        <Sidebar collapsible="icon" variant="inset">
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton
                            asChild
                            className="data-[slot=sidebar-menu-button]:!p-1.5"
                        >
                            <a href="#">
                                <img src={logo} />
                                <span className="text-base font-semibold">OPENMC PANEL</span>
                            </a>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarGroup>
                <SidebarGroupContent className="flex flex-col gap-2">
                    <SidebarMenu>
                        {props.navSections.map((section) => (
                            <SidebarMenuItem key={section.id}>
                                <SidebarMenuButton
                                    isActive={section.id === props.selectedItem}
                                    asChild
                                    onClick={() => {
                                        window.location.hash = section.name.replace(" ", "-")
                                        props.setSelectedItem(section.id);
                                    }}
                                >
                                    <span>
                                        {section.icon}
                                        {section.name}
                                    </span>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        ))}
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </Sidebar>
    );
}
