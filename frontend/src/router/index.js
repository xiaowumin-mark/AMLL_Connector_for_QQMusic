import { createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/home.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta:{
        show_nav: true,
        show_top_bar: true
      }
    },
    {
      path: '/init',
      name: 'init',
      component: () => import('../views/init.vue'),
      meta:{
        show_nav: false,
        show_top_bar: false
      }
    },
    {
      path: '/setting',
      name: 'setting',
      component: () => import('../views/setting.vue'),
      meta:{
        show_nav: true,
        show_top_bar: true
      }
    },
    {
      path: '/lyrics',
      name: 'lyrics',
      component: () => import('../views/lyrics.vue'),
      meta:{
        show_nav: false,
        show_top_bar: false
      }
    },
    {
      path: '/smtcs',
      name: 'smtcs',
      component: () => import('../views/smtc-manager.vue'),
      meta:{
        show_nav: true,
        show_top_bar: true
      }
    },
  ]
})

export default router
