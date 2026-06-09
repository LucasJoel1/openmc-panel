import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export const Permissions = {
	VIEW_SERVER_STATUS: 0x1n,
	START_SERVER: 0x2n,
	STOP_SERVER: 0x4n,
	RESTART_SERVER: 0x8n,
	VIEW_PLAYERS: 0x10n,
	VIEW_SERVER_PROPERTIES: 0x20n,
	READ_SERVER_LOGS: 0x40n,
	SEND_COMMANDS: 0x80n,
	VIEW_LOGS_HISTORY: 0x100n,
	DELETE_LOG: 0x200n,
	CREATE_USER: 0x400n,
	DELETE_USER: 0x800n,
	MODIFY_USER: 0x1000n,
	VIEW_USERS: 0x2000n,
} as const;

export function checkPermissions(permissions: bigint, permissionToCheck: bigint): boolean {
	return (permissions & permissionToCheck) === permissionToCheck
}

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function formatTime(seconds: number) {
	const units = [
		{ label: "yr", secs: 31536000 },
		{ label: "mo", secs: 2628000 },
		{ label: "wk", secs: 604800 },
		{ label: "d", secs: 86400 },
		{ label: "hr", secs: 3600 },
		{ label: "min", secs: 60 },
		{ label: "sec", secs: 1 },
	];

	const result: string[] = [];
	let vals = 0

	for (const unit of units) {
		const value = Math.floor(seconds / unit.secs);
		if (value > 0) {
			result.push(`${value} ${unit.label}${value !== 1 ? "s" : ""}`);
			seconds %= unit.secs;
			vals++;
		}

		// only have 2 units displaying at a time
		if (vals > 1) {
			break;
		}
	}

	return result.length > 0 ? result.join(", ") : "0 seconds";
}

