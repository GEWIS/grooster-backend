import { createRouter, createWebHistory, NavigationGuardNext, RouteLocationNormalizedGeneric } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import RosterView from '@/views/RosterView.vue';
import OrganView from '@/views/OrganView.vue';
import { getGEWISId, isAuthenticated, loginRedirect } from '@/helpers/TokenHelper';
import ApiService from '@/services/ApiService';
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

  if (to.path === '/home') {
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

async function canAccessOrgan(id: number): Promise<boolean> {
  try {
    const gewisId = getGEWISId();
    if (!gewisId) return false;

    const response = await ApiService.user.userGet(undefined, gewisId);
    const user = response.data?.[0];

    if (!user?.organs?.length) return false;

    const hasAccess = user.organs.some((organ) => organ.id === id);
    if (!hasAccess) return false;

    await ApiService.roster.getRoster(id);
    return true;
  } catch (error) {
    console.warn(`Failed to check organ access for ID ${id}:`, error);
    return false;
  }
}

/**
 * Reusable guard to ensure the organ store is set and user has access
 */
async function handleOrganAccess(
  to: RouteLocationNormalizedGeneric,
  from: RouteLocationNormalizedGeneric,
  next: NavigationGuardNext,
) {
  const store = useOrganStore();
  const id = parseInt(to.params.id as string);

  if (isNaN(id)) {
    return next('/');
  }

  store.setOrgan(id);

  if (await canAccessOrgan(id)) {
    next();
  } else {
    next('/');
  }
}

export default router;
