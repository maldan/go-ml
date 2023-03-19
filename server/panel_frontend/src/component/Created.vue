<template>
  <div :class="$style.main">
    <el-tag type="danger" effect="plain">
      {{ dayjs(date).format("DD MMM YYYY") }}
    </el-tag>
    <el-tag type="success" effect="plain">
      {{ dayjs(date).format("HH:mm:ss") }}
    </el-tag>
    <el-tag type="warning" effect="plain">
      {{ dayjs(date).format("SSS") }}
    </el-tag>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import dayjs from "dayjs";

const props = defineProps<{
  date: string;
}>();

// Hooks
onMounted(async () => {});
</script>

<style module lang="scss">
.main {
  display: grid;
  grid-template-columns: 80px 65px 1fr;
  gap: 7px;
}
</style>
