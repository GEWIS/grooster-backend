import {createRouter, createWebHistory, RouteRecordRaw} from "vue-router";
import RosterView from "@/views/RosterView.vue";
import OrganView from "@/views/OrganView.vue";
import {getGEWISId, isAuthenticated, loginRedirect} from "@/helpers/TokenHelper";
import ApiService from "@/services/ApiService";
import RosterTemplateView from "@/views/RosterTemplateView.vue";
import {useOrganStore} from "@/stores/organ.store";

const routes: RouteRecordRaw[] = [
    { path: '/', component: OrganView },
    {
        path: '/rosters/:id',
        name: 'rosters',
        component: RosterView,
        beforeEnter: async (to, from, next) => {
            const store = useOrganStore();
            const id = parseInt(to.params.id as string);

            store.setOrgan(id);

            if (await canAccessOrgan(id)) {
                next();
            } else {
                next('/');
            }
        }
    },
    {
        path: '/templates/:id',
        name: 'templates',
        component: RosterTemplateView,
        beforeEnter: async (to, from, next) => {
            const store = useOrganStore();
            const id = parseInt(to.params.id as string);

            store.setOrgan(id);

            if (await canAccessOrgan(id)) {
                next();
            } else {
                next('/');
            }
        }
    },
    { path: '/callback', component: OrganView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    if (to.path === '/callback') {
        const token = to.query.token;
        if (typeof token === 'string') {
            localStorage.setItem("access_token", token);
            // Optionally redirect to home or dashboard after storing token
            next('/');
            return;
        }
    }

    const store = useOrganStore()

    if (to.path === '/home') {
        store.clearOrgan();
    }

   if (isAuthenticated()) {
       next();
   } else {
       loginRedirect();
   }
});

async function canAccessOrgan(id: number): Promise<boolean> {
    try {
        const gewisId = getGEWISId();
        const response = await ApiService.user.userGet(undefined, gewisId);
        const user = response.data[0];
        const organs = user.organs;

        if (organs.some((organ) => organ.id === id)) {
            await ApiService.roster.getRoster(id);
            return true;
        } else {
            return false;
        }
    } catch (error) {
        return false;
    }
}

export default router;