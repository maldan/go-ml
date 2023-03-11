<template>
  <div :class="$style.main">
    <div :class="$style.item" v-for="x in logStore.list" :key="x.created">
      <div style="flex: none; width: 128px">{{ x.kind }}</div>
      <div>{{ x.file }}:{{ x.line }}</div>
      <div>{{ x.body }}</div>
      <div :class="$style.created">
        <div :class="$style.date">
          {{ dayjs(x.created).format("YYYY-MM-DD") }}
        </div>
        <div :class="$style.time">
          {{ dayjs(x.created).format("HH:mm:ss") }}
        </div>
        <div :class="$style.ms">{{ dayjs(x.created).format("SSS") }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useRequestStore } from "@/store/request";
import { useMainStore } from "@/store/main";
import { useLogStore } from "@/store/log";
import dayjs from "dayjs";

// Stores
const logStore = useLogStore();

// Vars

// Hooks
onMounted(async () => {
  await logStore.getList();
});

// Methods
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  overflow-y: scroll;
  height: calc(100% - 85px);

  .item {
    display: flex;
    border-bottom: 1px solid rgba(255, 255, 255, 0.5);
    padding: 10px;

    > div {
      flex: 1;
    }

    .created {
      display: flex;
      align-items: flex-start;
      justify-content: flex-end;

      .date,
      .time,
      .ms {
        background: rgba(255, 0, 0, 0.5);
        padding: 3px 7px;
        border-radius: 4px;
        font-size: 12px;
        margin-right: 10px;
      }

      .time {
        background: rgba(0, 255, 17, 0.5);
      }

      .ms {
        background: rgba(123, 0, 255, 0.5);
      }
    }
  }
}
</style>
