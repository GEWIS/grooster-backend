import { defineStore } from 'pinia';

export type ThemeMode = 'system' | 'light' | 'dark';

const systemQuery = window.matchMedia('(prefers-color-scheme: dark)');

export const useThemeStore = defineStore('theme', {
    state: () => ({
        mode: 'system' as ThemeMode,
        systemDark: systemQuery.matches,
    }),
    persist: { pick: ['mode'] },
    getters: {
        isDark: (state) => (state.mode === 'system' ? state.systemDark : state.mode === 'dark'),
    },
    actions: {
        init() {
            systemQuery.addEventListener('change', (e) => {
                this.systemDark = e.matches;
                this.apply();
            });
            this.apply();
        },
        setMode(mode: ThemeMode) {
            this.mode = mode;
            this.apply();
        },
        cycleMode() {
            const order: ThemeMode[] = ['system', 'light', 'dark'];
            this.setMode(order[(order.indexOf(this.mode) + 1) % order.length]);
        },
        apply() {
            document.documentElement.classList.toggle('app-dark', this.isDark);
        },
    },
});
