<template>
  <div :class="$style.main">
    <!-- Label -->
    <div v-if="info.type === 'bitmask'" style="margin-bottom: 5px">
      {{ info.name }}
    </div>
    <div v-else style="margin-bottom: 5px">{{ info.label || info.name }}</div>

    <el-input
      v-if="info.type === 'string'"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      :disabled="!info.isEdit"
    />
    <el-input-number
      v-if="info.type === 'int'"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      style="width: 100%"
      :disabled="!info.isEdit"
    />
    <el-checkbox
      v-if="info.type === 'bool'"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      :disabled="!info.isEdit"
    />

    <!-- Bitmask -->
    <div v-if="info.type === 'bitmask'">
      <el-checkbox
        :label="x"
        v-for="(x, i) in info.label.split(',')"
        :placeholder="x"
        @change="
          tempRow[info.name] = changeBitMask($event, tempRow[info.name], i)
        "
        :disabled="!info.isEdit"
        v-model="tempRow[info.name + '_mask_' + i]"
      />
      <!--      <div>
        Bitmask: {{ tempRow[info.name]?.toString(2) }}
        {{ tempRow[info.name] }}
      </div>-->
    </div>

    <!-- Default non edit field -->
    <!-- <el-input
       v-if="!info.isEdit"
       :disabled="!info.isEdit"
       :placeholder="info.name"
       v-model="tempRow[info.name]"
     /> -->
  </div>
</template>

<script setup lang="ts">
import dayjs from "dayjs";
import { ref } from "vue";

const props = defineProps<{
  info: { label: string; isEdit: boolean; type: string; name: string };
  tempRow: any;
}>();

// const tempRow = ref({});

function changeBitMask(isSet: boolean, current: number, pos: number): number {
  if (isSet) {
    return current | (1 << pos);
  }
  return current & ~(1 << pos);
}
</script>

<style lang="scss" module>
.main {
  // border: 2px solid #fe0000;
  // border: 1px dashed rgba(255, 255, 255, 0.3);
  border-bottom: 1px dashed rgba(255, 255, 255, 0.15);
  padding-bottom: 10px;
  // border-radius: 4px;
}
</style>
