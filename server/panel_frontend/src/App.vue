<template>
  <div :class="$style.mainApp">
    <el-menu
      :default-active="$route.path"
      class="el-menu-demo"
      mode="horizontal"
      @select="handleSelect"
    >
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
      </el-menu-item>
    </el-menu>

    <RouterView />
  </div>
</template>

<script lang="ts" setup>
import { RouterView, useRouter } from "vue-router";

const router = useRouter();

const handleSelect = (key: string, keyPath: string[]) => {
  router.push(key);
};
</script>

<style lang="scss" module>
.mainApp {
  height: 100%;
}
</style>
