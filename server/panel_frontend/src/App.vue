<template>
  <div :class="$style.mainApp">
    <el-menu
      :default-active="$route.path"
      class="el-menu-demo"
      mode="horizontal"
      @select="handleSelect"
    >
      <!--      <el-menu-item index="/logs" style="color: #fe6e3d">
        <el-icon><Promotion /></el-icon>Logs
      </el-menu-item>
      <el-menu-item index="/requests" style="color: #fe6e3d">
        <el-icon><Promotion /></el-icon>Requests
      </el-menu-item>
      <el-menu-item index="/methods" style="color: #c1fe48">
        <el-icon><EditPen /></el-icon>Methods
      </el-menu-item>
      <el-menu-item index="/tests" style="color: #fea048">
        <el-icon><WarnTriangleFilled /></el-icon>Tests
      </el-menu-item>
      <el-menu-item index="/db" style="color: #ec48fe">
        <el-icon><Document /></el-icon>DB
      </el-menu-item>
      <el-menu-item index="/control" style="color: #fed448">
        <el-icon><Operation /></el-icon>Control
      </el-menu-item>
      <el-menu-item index="/settings">
        <el-icon><Tools /></el-icon>Settings
      </el-menu-item>-->
      <el-menu-item :index="x.url" v-for="x in links" :key="x.url">
        <el-icon><component :is="x.icon" /></el-icon>{{ x.title }}
      </el-menu-item>
    </el-menu>

    <RouterView />
  </div>
</template>

<script lang="ts" setup>
import { RouterView, useRouter } from "vue-router";
import { onMounted, ref } from "vue";
import { useMainStore } from "@/store/main";

const router = useRouter();
const mainStore = useMainStore();
const links = ref<any[]>([
  /*
    { url: "/logs", title: "Logs", icon: "Promotion" },
    { url: "/requests", title: "Requests", icon: "Promotion" },
    { url: "/settings", title: "Settings", icon: "Tools" },
  */
]);

const handleSelect = (key: string, keyPath: string[]) => {
  router.push(key);
};

onMounted(async () => {
  await mainStore.getSetting();

  if (mainStore.hasLogTab) {
    links.value.push({ url: "/logs", title: "Logs", icon: "Promotion" });
  }

  links.value.push({ url: "/requests", title: "Requests", icon: "Promotion" });
  links.value.push({ url: "/router", title: "Router", icon: "Promotion" });
  links.value.push({ url: "/db", title: "DB", icon: "Promotion" });
  links.value.push({ url: "/forms", title: "Forms", icon: "Promotion" });
});
</script>

<style lang="scss" module>
.mainApp {
  height: 100%;
}
</style>
