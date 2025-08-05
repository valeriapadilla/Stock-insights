import { createRouter, createWebHistory } from 'vue-router'
import TestView from '../views/TestView.vue'

const routes = [
  {
    path: '/',
    name: 'Test',
    component: TestView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 