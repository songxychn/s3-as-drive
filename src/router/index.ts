import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router';
import Home from "@/views/Home.vue";
import Settings from "../views/settings/Settings.vue";
import FileBrowser from "../views/FileBrowser.vue";

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        redirect: '/home'
    },
    {
        path: '/home',
        component: Home,
        children: [
            {
                path: '/home/settings',
                component: Settings,
            },
            {
                path: '/home/file-browser',
                component: FileBrowser,
            },
        ]
    },

]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

export default router
