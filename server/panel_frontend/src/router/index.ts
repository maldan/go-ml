import { createRouter, createWebHashHistory } from "vue-router";
import Methods from "../page/Methods.vue";
import Settings from "../page/Settings.vue";
import Tests from "../page/Tests.vue";
import Requests from "../page/Requests.vue";
import DB from "../page/DB.vue";
import Control from "../page/Control.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: "Methods",
      component: Methods,
    },
    {
      path: "/methods",
      name: "Methods",
      component: Methods,
    },
    {
      path: "/settings",
      name: "Settings",
      component: Settings,
    },
    {
      path: "/tests",
      name: "Tests",
      component: Tests,
    },
    {
      path: "/requests",
      name: "Requests",
      component: Requests,
    },
    {
      path: "/db",
      name: "DB",
      component: DB,
    },
    {
      path: "/control",
      name: "Control",
      component: Control,
    },
  ],
});

export default router;
