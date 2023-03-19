<template>
  <div :class="$style.main">
    <div v-for="x in type.fieldList" :key="x.name">
      {{ x.name }} {{ x.kind }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";

const props = defineProps<{
  type: { name: string; type: string; fieldList: any };
}>();

// Hooks
onMounted(async () => {});
</script>

<style module lang="scss">
.main {
  user-select: none;
}
</style>
