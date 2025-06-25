import type { ReactNode } from "react";

export interface NavSection {
    id: number,
    name: string,
    icon: ReactNode
    sideBar: boolean
}