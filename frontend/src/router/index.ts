import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router';
import Home from "@/views/Home.vue";
import FileBrowser from "@/views/FileBrowser.vue";
import Config from '@/views/Config.vue';

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        redirect: '/home'
    },
    {
        path: '/home',
        component: Home,
        redirect: '/home/file-browser',
        children: [
            {
                path: '/home/config',
                component: Config,
            },
            {
                path: '/home/file-browser',
                component: FileBrowser,
            },
        ]
    },

]

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes
})

export default router
