<template>
  <div :class="$style.main">
    <el-table :data="fileList" border style="width: 100%">
      <!-- Icon -->
      <el-table-column label="icon" width="64">
        <template #default="scope">
          <el-icon v-if="scope.row.isDir" :size="32"><Folder /></el-icon>
          <el-icon v-if="!scope.row.isDir" :size="32"><Document /></el-icon
        ></template>
      </el-table-column>

      <el-table-column prop="name" label="Name" />

      <!-- Size -->
      <el-table-column label="Size">
        <template #default="scope">
          {{ prettyBytes(scope.row.size) }}
        </template>
      </el-table-column>

      <!-- Created -->
      <el-table-column label="Created" width="220">
        <template #default="scope">
          <Created :date="scope.row.created" />
        </template>
      </el-table-column>
    </el-table>
    <!--    <div v-for="x in fileList" :key="x.name">{{ x.name }}</div>-->
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import dayjs from "dayjs";
import Created from "@/component/Created.vue";
import prettyBytes from "pretty-bytes";

const props = defineProps<{
  fileList: { name: string; fullPath: string }[];
}>();

// Hooks
onMounted(async () => {});
</script>

<style module lang="scss">
.main {
  user-select: none;
}
</style>
