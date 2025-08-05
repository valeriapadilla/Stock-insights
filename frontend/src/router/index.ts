import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: DashboardView
  },
  {
    path: '/test',
    name: 'Test',
    component: () => import('../views/TestView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 