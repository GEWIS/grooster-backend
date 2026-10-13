// Shared pass-through (pt) styling for OpenVue components that appear in several places.

/** NewRoster.html dialog: left-aligned title, small icon close button, footer bar. */
export const dialogPt = {
    root: { class: 'w-[520px] max-w-[calc(100vw-2rem)] overflow-hidden rounded-[14px] border-0 shadow-xl' },
    header: { class: 'items-start gap-3 px-6 pb-1 pt-[22px]' },
    title: { class: 'font-display text-[22px] font-semibold tracking-[-0.01em]' },
    pcCloseButton: { root: { class: 'size-9 text-surface-700 dark:text-surface-200' } },
    content: { class: 'flex flex-col gap-[18px] px-6 pb-[22px] pt-[18px]' },
    footer: {
        class: 'flex justify-end gap-2.5 border-t border-surface bg-surface-50 px-6 py-4 dark:bg-surface-800',
    },
};

/** ConfirmDialog with the same look; narrower and with the message as body text. */
export const confirmDialogPt = {
    ...dialogPt,
    root: { class: 'w-[440px] max-w-[calc(100vw-2rem)] overflow-hidden rounded-[14px] border-0 shadow-xl' },
    content: { class: 'flex items-start px-6 pb-[22px] pt-3' },
    icon: { class: 'hidden' },
    message: { class: 'm-0 text-left text-[15px] text-color' },
    pcAcceptButton: { root: { class: 'h-10 text-sm' } },
    pcRejectButton: { root: { class: 'h-10 text-sm' } },
};

/** Disabled primary buttons stay readable (HANDOFF §7) instead of Aura's 60% opacity. */
export const readableDisabled =
    'disabled:border-surface-300 disabled:bg-surface-200 disabled:text-surface-700 disabled:opacity-100 dark:disabled:border-surface-600 dark:disabled:bg-surface-700 dark:disabled:text-surface-100';

/** Form-field Selects use the same text size as inputs: 16px on phones (no iOS zoom), 14px from md. */
export const fieldSelectPt = { label: { class: 'text-base md:text-sm' } };
