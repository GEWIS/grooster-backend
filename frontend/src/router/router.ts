import {createRouter, createWebHistory, RouteRecordRaw} from "vue-router";
import RosterView from "@/views/RosterView.vue";
import OrganView from "@/views/OrganView.vue";
import {isAuthenticated, loginRedirect} from "@/helpers/TokenHelper";
import ApiService from "@/services/ApiService";

const routes: RouteRecordRaw[] = [
    { path: '/', component: OrganView },
    {
        path: '/roster/:id',
        name: 'roster',
        component: RosterView,
        beforeEnter: async (to, from, next) => {
            try {
                const id = to.params.id as string;
                await ApiService.roster.getRoster(parseInt(id));
                next();
            } catch (error) {
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

   if (isAuthenticated()) {
       next();
   } else {
       loginRedirect();
   }
});

export default router;