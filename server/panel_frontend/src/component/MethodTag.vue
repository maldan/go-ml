<template>
  <div :class="$style.main">
    <!-- Get -->
    <el-tag v-if="props.tag === 'GET'" effect="plain" :class="$style.get">{{
      tag
    }}</el-tag>

    <!-- Post -->
    <el-tag
      v-else-if="props.tag === 'POST' || props.tag === 'PUT'"
      effect="plain"
      type="success"
      :class="$style.post"
      >{{ tag }}</el-tag
    >

    <!-- Post -->
    <el-tag
      v-else-if="props.tag === 'PATCH'"
      effect="plain"
      type="success"
      style="font-size: 14px"
      >{{ tag }}</el-tag
    >

    <!-- Delete -->
    <el-tag
      v-else-if="props.tag === 'DELETE'"
      effect="plain"
      type="danger"
      style="font-size: 14px"
      >{{ tag }}</el-tag
    >

    <!-- Other -->
    <el-tag v-else effect="plain" type="info" style="font-size: 14px">{{
      tag
    }}</el-tag>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";

const props = defineProps<{
  tag: string;
}>();

// Hooks
onMounted(async () => {});
</script>

<style module lang="scss">
.main {
  user-select: none;

  .get,
  .post {
    font-size: 14px;
    color: #41c3ff;
    border-color: #41c3ff;
  }

  .post {
    color: #80df22;
    border-color: #80df22;
  }
}
</style>
