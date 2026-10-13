import { defineStore } from 'pinia';
import { User } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';

interface UserState {
    user: User | undefined;
}

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        user: undefined,
    }),
    getters: {
        getUser: (state) => state.user,
    },
    actions: {
        async fetchUser(gewisId: number) {
            if (this.user) return this.user;

            try {
                const response = await ApiService.user.userGet(undefined, gewisId);
                this.user = response.data[0];

                return this.user;
            } catch (e) {
                console.error('Failed to fetch user:', e);
                return null;
            }
        },
    },
});
