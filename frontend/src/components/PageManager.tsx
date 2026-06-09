import { useMemo } from "react";
import { checkPermissions, Permissions } from "@/lib/utils";
import ConsolePage from "./pages/ConsolePage";
import StatusPage from "./pages/StatusPage";
import UserManagementPage from "./pages/UserManagementPage";

interface PageHandlerProps {
    page: number;
    permissions: bigint;
}

export default function PageHandler(props: PageHandlerProps) {
    const canViewStatus = useMemo(() =>
        checkPermissions(props.permissions, Permissions.VIEW_SERVER_STATUS),
        [props.permissions]
    );

    const canReadLogs = useMemo(() =>
        checkPermissions(props.permissions, Permissions.READ_SERVER_LOGS),
        [props.permissions]
    );

    const canViewUsers = useMemo(() =>
        checkPermissions(props.permissions, Permissions.VIEW_USERS),
        [props.permissions]
    )

    return (
        <div className="flex-1 min-h-0 p-2">
            {canViewStatus && <StatusPage permissions={props.permissions} visible={props.page === 0} />}
            {canReadLogs && <ConsolePage permissions={props.permissions} visible={props.page === 1} isPage={true} />}
            {canViewUsers && <UserManagementPage permissions={props.permissions} visible={props.page === 100 } />}
        </div>
    );
}
