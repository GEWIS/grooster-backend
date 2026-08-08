import { defineStore } from 'pinia';

export const useOrganStore = defineStore('organ', {
    state: () => ({
        organ: undefined as number | undefined,
        organName: undefined as string | undefined,
    }),
    getters: {
        getOrgan: (state) => state.organ,
    },
    actions: {
        setOrgan(id: number, name: string) {
            this.organ = id;
            this.organName = name;
        },
        clearOrgan() {
            this.organ = undefined;
            this.organName = undefined;
        },
    },
});
