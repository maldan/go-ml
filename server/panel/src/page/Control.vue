<template>
  <div :class="$style.main" style="padding: 10px">
    <!-- Left menu -->
    <el-tabs tab-position="left" class="demo-tabs">
      <el-tab-pane v-for="v in tabs" :label="v">
        <!-- Command list -->
        <div
          v-for="cmd in controlStore.commandList.filter((x) => x.folder === v)"
          :key="cmd"
        >
          <el-button
            type="primary"
            :loading="controlStore.status[cmd.folder + '/' + cmd.name]"
            @click="controlStore.execute(cmd.folder, cmd.name)"
            >{{ cmd.name }}</el-button
          >
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useControlStore } from "@/store/control";

// Stores
const controlStore = useControlStore();

// Vars
const tabs = computed(() => {
  return [...new Set(controlStore.commandList.map((x) => x.folder))];
});

// Hooks
onMounted(async () => {
  await controlStore.getCommandList();
});

// Methods
</script>

<style module lang="scss">
.main {
  display: flex;
  flex-direction: column;
  font-size: 14px;
}
</style>
