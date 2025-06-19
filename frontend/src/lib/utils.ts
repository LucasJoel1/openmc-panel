import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

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