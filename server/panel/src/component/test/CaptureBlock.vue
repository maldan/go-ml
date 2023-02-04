<template>
  <div :class="$style.main">
    <div v-for="(v, k) in map" :key="k">{{ k }} -> {{ v }}</div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";

// Props
const props = defineProps<{
  map: Record<string, string>;
}>();

// Vars
const jsonBody = ref("");

// Hooks
onMounted(async () => {});
</script>

<style module lang="scss">
.main {
  background: #1b1b1b;
  padding: 5px;
  width: 140px;
  display: flex;
  align-items: center;
  justify-content: center;

  > div {
    margin-bottom: 10px;
  }
}
</style>
