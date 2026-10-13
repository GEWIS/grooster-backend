import ApiService from '@/services/ApiService';

/** Downloads the backend-rendered PNG of a locked roster. */
export async function exportRoster(rosterId: number) {
    const response = await ApiService.exportApi.exportRosterIdGet(rosterId, { responseType: 'blob' });
    const blob = response.data instanceof Blob ? response.data : new Blob([response.data], { type: 'image/png' });

    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `roster-${rosterId}.png`);
    document.body.appendChild(link);
    link.click();

    window.URL.revokeObjectURL(url);
    link.remove();
}
