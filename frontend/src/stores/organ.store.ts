import { defineStore } from 'pinia';

interface OrganState {
    organ: number | undefined;
    organName: string | undefined;
}

export const useOrganStore = defineStore('organ', {
    state: (): OrganState => ({
        organ: undefined,
        organName: undefined,
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
