export function getMonthStr(numeric: string): string {
    switch (numeric) {
        case "10":
            return "Oct";
        case "11":
            return "Nov";
        case "12":
            return "Dec";
        default:
            return "";
    }
}

export function getDayStr(numeric: string): string {
    if (numeric.startsWith("0")) {
        return numeric.substring(1);
    }
    return numeric;
}
