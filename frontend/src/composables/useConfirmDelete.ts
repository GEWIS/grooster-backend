import { useConfirm } from 'openvue/useconfirm';

interface ConfirmDeleteOptions {
    /** e.g. "Delete roster?" */
    header: string;
    message: string;
    /** Label of the destructive button, e.g. "Delete roster" */
    acceptLabel: string;
    onAccept: () => void | Promise<void>;
}

export function useConfirmDelete() {
    const confirm = useConfirm();

    return (options: ConfirmDeleteOptions) =>
        confirm.require({
            header: options.header,
            message: options.message,
            acceptProps: {
                label: options.acceptLabel,
                severity: 'danger',
                // Aura's red-500 with white text is below 4.5:1; red-700 passes
                class: 'border-red-700 bg-red-700 enabled:hover:border-red-800 enabled:hover:bg-red-800',
            },
            rejectProps: { label: 'Cancel', severity: 'secondary', outlined: true },
            accept: () => void options.onAccept(),
        });
}
