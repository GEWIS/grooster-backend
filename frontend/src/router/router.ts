import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHistory, NavigationGuardNext, RouteLocationNormalizedGeneric } from 'vue-router';
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
        beforeEnter: handleOrganAccess,
    },
    {
        path: '/templates/:id',
        name: 'templates',
        component: RosterTemplateView,
        beforeEnter: handleOrganAccess,
    },
    { path: '/callback', component: OrganView },
    {
        path: '/organ/:id/profile',
        name: 'profile',
        component: OrganMemberSettings,
        beforeEnter: handleOrganAccess,
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
        return true;
    } else {
        loginRedirect();
        return false;
    }
});

/**
 * Reusable guard to ensure the organ store is set and user has access
 */
function handleOrganAccess(
    to: RouteLocationNormalizedGeneric,
    from: RouteLocationNormalizedGeneric,
    next: NavigationGuardNext,
) {
    const organStore = useOrganStore();
    const userStore = useUserStore();

    const id = parseInt(to.params.id as string);

    if (isNaN(id)) {
        return next('/');
    }

    const user = userStore.getUser;

    if (!user) {
        return next('/');
    }

    const organ = user.organs.find((org) => org.id === id);

    if (!organ) {
        return next('/');
    }

    organStore.setOrgan(id, organ.name);
    next();
}

export default router;
