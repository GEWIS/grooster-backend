// Answer values come from the backend per roster (`Roster.values`), e.g. J / X / N / L / S / B.
// The letter itself is the visible label, so colour is never the only signal.

export const NOT_SET_LABEL = 'Not set';

const PREFERENCE_TONES: Record<string, string> = {
    J: 'bg-pref-prefer text-pref-prefer-fg',
    X: 'bg-pref-avoid text-pref-avoid-fg',
    N: 'bg-pref-no text-pref-no-fg',
    L: 'bg-pref-l text-pref-l-fg',
    S: 'bg-pref-rainbow motion-safe:animate-pref-rainbow text-pref-special-fg',
    B: 'bg-pref-beer text-pref-special-fg',
};

const NOT_SET_TONE = 'bg-pref-none text-pref-none-fg';
const UNKNOWN_TONE = 'bg-pref-ok text-pref-ok-fg';

export function isPreferenceSet(value: string | null | undefined): value is string {
    return !!value;
}

export function preferenceLabel(value: string | null | undefined): string {
    return isPreferenceSet(value) ? value : NOT_SET_LABEL;
}

export function preferenceToneClass(value: string | null | undefined): string {
    if (!isPreferenceSet(value)) return NOT_SET_TONE;
    return PREFERENCE_TONES[value] ?? UNKNOWN_TONE;
}
