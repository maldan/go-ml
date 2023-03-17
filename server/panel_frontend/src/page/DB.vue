<template>
  <div :class="$style.main">
    <!-- Header list -->
    <div :class="$style.list">
      <el-button
        @click="dbStore.selectTable(x)"
        v-for="x in dbStore.list"
        :key="x"
        :type="dbStore.selectedTable === x ? 'primary' : ''"
        effect="plain"
      >
        {{ x }}
      </el-button>
    </div>

    <!-- Body -->
    <div :class="$style.recordList">
      <div :class="$style.record" v-for="x in dbStore.recordList" :key="x">
        <div :class="$style.field" v-for="(v, k) in x">
          <div :class="$style.name">{{ k }}</div>
          <div :class="$style.string" v-if="typeof v === 'string'">
            "{{ v }}"
          </div>
          <div :class="$style.number" v-if="typeof v === 'number'">
            {{ v }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useDBStore } from "@/store/db";

// Stores
const dbStore = useDBStore();

// Vars

// Hooks
onMounted(async () => {
  await dbStore.getList();
});

// Methods
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);

  .list {
    margin-bottom: 10px;
  }

  .recordList {
    overflow-y: auto;
    height: calc(100% - 40px);

    .record {
      display: flex;
      border: 1px solid rgba(255, 255, 255, 0.25);
      margin-bottom: 10px;
      padding: 5px;
      border-radius: 4px;

      .field {
        padding: 2px 5px;
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 4px;
        margin-right: 10px;
        display: flex;
        font-size: 14px;

        .name {
          // font-weight: bold;
          color: #ffc46a;
          margin-right: 10px;
        }

        .string {
          color: #3ad81c;
        }
        .number {
          color: #63a3f2;
        }
      }
    }
  }
}
</style>
