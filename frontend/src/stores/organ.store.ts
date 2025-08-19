import { defineStore } from 'pinia';

export const useOrganStore = defineStore('organ', {
   state: () => ({
       organ: undefined as number|undefined,
   }),
    getters: {
       getOrgan: (state) => state.organ,
    },
    actions: {
       setOrgan(id: number) {
           this.organ = id;
       },
        clearOrgan() {
           this.organ = undefined;
        }
    }
});