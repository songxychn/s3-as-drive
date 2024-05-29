import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router';
import Home from "@/views/Home.vue";
import FileBrowser from "@/views/FileBrowser.vue";
import Config from '@/views/Config.vue';
import About from "../views/About.vue";
import SyncDir from "../views/SyncDir.vue";

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
            {
                path: '/home/about',
                component: About,
            },
            {
                path: '/home/sync-dir',
                component: SyncDir,
            },
        ]
    },

]

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes
})

export default router
