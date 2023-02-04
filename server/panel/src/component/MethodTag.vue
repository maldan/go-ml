<template>
  <div :class="$style.main">
    <!-- Get -->
    <el-tag v-if="props.tag === 'GET'" effect="dark">{{ tag }}</el-tag>

    <!-- Post -->
    <el-tag
      v-else-if="props.tag === 'POST' || props.tag === 'PUT'"
      effect="dark"
      type="warning"
      >{{ tag }}</el-tag
    >

    <!-- Post -->
    <el-tag v-else-if="props.tag === 'PATCH'" effect="dark" type="success">{{
      tag
    }}</el-tag>

    <!-- Delete -->
    <el-tag v-else-if="props.tag === 'DELETE'" effect="dark" type="danger">{{
      tag
    }}</el-tag>

    <!-- Other -->
    <el-tag v-else effect="dark" type="info">{{ tag }}</el-tag>
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
}
</style>
