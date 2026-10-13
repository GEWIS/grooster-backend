// Roster dates come from the backend as "YYYY-MM-DDT00:00:00Z". Read the date part as a local date
// so the day never shifts with the time zone.

export function parseRosterDate(value: string | undefined): Date | null {
    const match = value?.match(/^(\d{4})-(\d{2})-(\d{2})/);
    if (!match) return null;
    return new Date(Number(match[1]), Number(match[2]) - 1, Number(match[3]));
}

export function toRosterDate(date: Date): string {
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${date.getFullYear()}-${month}-${day}T00:00:00Z`;
}

/** "Tue 29 Sep 2026" */
export function formatRosterDate(date: Date): string {
    const weekday = date.toLocaleDateString('en-GB', { weekday: 'short' });
    return `${weekday} ${date.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })}`;
}

/** ISO 8601 week number */
export function isoWeek(date: Date): number {
    const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
    const day = d.getUTCDay() || 7;
    d.setUTCDate(d.getUTCDate() + 4 - day);
    const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
    return Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7);
}
