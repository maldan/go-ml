<template>
  <div :class="$style.main">1</div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useDbStore } from "@/store/db";

// Stores
const dbStore = useDbStore();

// Vars

// Hooks
onMounted(async () => {});

// Methods
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);
}
</style>
