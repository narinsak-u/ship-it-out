import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/HomeView.vue"),
    },
    {
      path: "/orders",
      name: "orders",
      component: () => import("@/views/OrdersView.vue"),
    },
    {
      path: "/orders/create",
      name: "order-create",
      meta: { staffOnly: true },
      component: () => import("@/views/OrderFormView.vue"),
    },
    {
      path: "/orders/:orderId/edit",
      meta: { staffOnly: true },
      name: "order-edit",
      component: () => import("@/views/OrderFormView.vue"),
    },
    {
      path: "/orders/:orderId",
      name: "order-detail",
      component: () => import("@/views/OrderDetailView.vue"),
    },
    {
      path: "/carriers",
      name: "carriers",
      component: () => import("@/views/CarriersView.vue"),
    },
  ],
});
router.beforeEach(async (to) => {
  if (!to.meta.staffOnly) return true;
  const auth = useAuthStore();
  if (auth.loading) await auth.init();
  if (!auth.isStaff) return { name: "orders" };
});

export default router;
