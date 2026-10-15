import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHistory } from 'vue-router';
import RosterView from '@/views/RosterView.vue';
import OrganView from '@/views/OrganView.vue';
import { getGEWISId, isAuthenticated, loginRedirect } from '@/helpers/TokenHelper';
import RosterTemplateView from '@/views/RosterTemplateView.vue';
import { useOrganStore } from '@/stores/organ.store';
import { useRosterStore } from '@/stores/roster.store';
import { useUserStore } from '@/stores/user.store';
import OrganMemberSettings from '@/views/OrganMemberSettings.vue';

const routes: RouteRecordRaw[] = [
    { path: '/', component: OrganView },
    {
        path: '/rosters/:id',
        name: 'rosters',
        component: RosterView,
        meta: { requireOrgan: true },
    },
    {
        path: '/templates/:id',
        name: 'templates',
        component: RosterTemplateView,
        meta: { requireOrgan: true },
    },
    { path: '/callback', component: OrganView },
    {
        path: '/organ/:id/profile',
        name: 'profile',
        component: OrganMemberSettings,
        meta: { requireOrgan: true },
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach(async (to, from) => {
    const rosterStore = useRosterStore();
    const userStore = useUserStore();
    const organStore = useOrganStore();

    if (from.name && to.params.id !== from.params.id) {
        rosterStore.clearRosters();
    }

    if (to.path === '/home' || to.path === '/') {
        organStore.clearOrgan();
    }

    if (to.path === '/callback') {
        const token = to.query.token;
        if (typeof token === 'string') {
            localStorage.setItem('access_token', token);
            return '/';
        }
    }

    if (isAuthenticated()) {
        const gewisId = getGEWISId();
        if (gewisId) {
            await userStore.fetchUser(gewisId);
        }

        // Runs on every navigation (including param-only changes like /rosters/1 -> /rosters/2),
        // which per-route beforeEnter guards do not.
        if (to.meta.requireOrgan) {
            const id = parseInt(to.params.id as string);
            const user = userStore.getUser;
            const organ = user?.organs.find((org) => org.id === id);

            if (isNaN(id) || !organ) {
                return '/';
            }

            if (organStore.getOrgan !== id) {
                rosterStore.clearRosters();
            }
            organStore.setOrgan(id, organ.name);
        }

        return true;
    } else {
        loginRedirect();
        return false;
    }
});

export default router;
