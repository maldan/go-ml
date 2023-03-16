<template>
  <div :class="$style.main">
    <!-- Header list -->
    <div :class="$style.list">
      <el-button
        @click="routerStore.selectRouter(x)"
        v-for="x in routerStore.list"
        :key="x.path"
        :type="routerStore.selectedRouter.path === x.path ? 'primary' : ''"
        effect="plain"
      >
        {{ x.path }}
        <el-tag
          size="small"
          type="success"
          effect="light"
          style="margin-left: 10px"
          >{{ x.type }}</el-tag
        >
      </el-button>
    </div>

    <!-- Controller list -->
    <div :class="$style.list">
      <el-button
        @click="routerStore.selectController(x)"
        v-for="x in routerStore.controllerList"
        :key="x"
        :type="routerStore.selectedController === x ? 'primary' : ''"
        effect="plain"
      >
        {{ x }}
      </el-button>
    </div>

    <!-- Method list -->
    <div>
      <div :class="$style.method" v-for="x in routerStore.methodList" :key="x">
        <MethodTag :class="$style.httpMethod" :tag="x.httpMethod" />
        <div :class="$style.url">{{ x.url }}</div>
        <el-button
          v-for="y in x.args"
          :key="y"
          plain
          type="success"
          size="small"
          style="font-size: 14px"
          >{{ y }}</el-button
        >

        <el-button
          v-for="y in x.return"
          :key="y"
          plain
          type="warning"
          size="small"
          style="font-size: 14px"
          >{{ y }}</el-button
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import { useRouterStore } from "@/store/router";
import MethodTag from "@/component/MethodTag.vue";

// Stores
const routerStore = useRouterStore();

// Vars

// Hooks
onMounted(async () => {
  await routerStore.getList();
});

// Methods
</script>

<style module lang="scss">
.main {
  padding: 10px;
  display: flex;
  flex-direction: column;

  .list {
    display: flex;
    margin-bottom: 10px;
  }

  .method {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    border: 1px solid rgba(255, 255, 255, 0.25);
    padding: 10px;

    .httpMethod {
      width: 80px;
    }

    .url {
      width: 240px;
      font-size: 14px;
    }
  }
}
</style>
