import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/orders',
      name: 'orders',
      component: () => import('@/views/OrdersView.vue'),
    },
    {
      path: '/orders/:orderId',
      name: 'order-detail',
      component: () => import('@/views/OrderDetailView.vue'),
    },
    {
      path: '/carriers',
      name: 'carriers',
      component: () => import('@/views/CarriersView.vue'),
    },
  ],
});

export default router;
