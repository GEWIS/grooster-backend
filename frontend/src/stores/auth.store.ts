import { jwtDecode, type JwtPayload } from 'jwt-decode';
import { defineStore } from 'pinia';
import { useOrganStore } from '@/stores/organ.store';

export enum Role {
    Owner = 'Owner',
    Admin = 'Admin',
    Member = 'Member',
}

interface CustomJwtPayload extends JwtPayload {
    organs: Record<string, string>;
}

export const useAuthStore = defineStore('auth', {
    state: () => ({
        token: localStorage.getItem('access_token') || null,
    }),
    getters: {
        userOrgans(state): Record<string, string> {
            if (!state.token) return {};
            try {
                const decoded = jwtDecode<CustomJwtPayload>(state.token);
                return decoded.organs || {};
            } catch {
                return {};
            }
        },
        currentRole(): Role | null {
            const organStore = useOrganStore();
            const activeId = String(organStore.organ);

            const roleName = this.userOrgans[activeId];

            if (!roleName) return null;

            const normalizedRole = roleName.charAt(0).toUpperCase() + roleName.slice(1);
            return normalizedRole as Role;
        },
        can() {
            return (requiredRoles: Role[]): boolean => {
                const role = this.currentRole;
                return role ? requiredRoles.includes(role) : false;
            };
        },
    },
});
