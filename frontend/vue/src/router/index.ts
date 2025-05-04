import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthView from "@/views/AuthView.vue";
import CalculateView from "@/views/CalculateView.vue";
import FindVue from "@/views/FindVue.vue";
import ExpressionsView from "@/views/ExpressionsView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: 'Home' }
    },
    {
      path: '/auth',
      name: 'auth',
      component: AuthView,
      meta: { title: 'Login' }
    },
    {
      path: '/expressions',
      name: 'expressions',
      component: ExpressionsView,
      meta: { title: 'Expressions' }
    },
    {
      path: '/calculate',
      name: 'calculate',
      component: CalculateView,
      meta: { title: 'Calculate' }
    },
    {
      path: '/find',
      name: 'find',
      component: FindVue,
      meta: { title: 'Find' }
    }
  ],
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
  next()
})

export default router
