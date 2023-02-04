import { createApp } from "vue";
import ElementPlus from "element-plus";
import { createPinia } from "pinia";
import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";
import "./main.scss";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

import App from "./App.vue";
import router from "./router";
import dayjs from "dayjs";

const app = createApp(App);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// Relative
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

app.use(createPinia());
app.use(router);
app.use(ElementPlus);

app.mount("#app");
