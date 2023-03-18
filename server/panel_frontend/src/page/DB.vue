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
      <el-input
        placeholder="Filter..."
        v-model="dbStore.where"
        style="margin-top: 10px"
        @change="dbStore.getSearch"
      />

      <!-- Pagination -->
      <div style="display: flex; align-items: center; margin-top: 10px">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="dbStore.search.total"
          :page-size="dbStore.search.perPage"
          @current-change="changePage"
        />
        <div style="margin-left: 10px">{{ dbStore.search.total }}</div>
      </div>
    </div>

    <!-- Body -->
    <div :class="$style.recordList">
      <div :class="$style.record" v-for="x in dbStore.search.result" :key="x">
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
async function changePage(page: number) {
  dbStore.page = page - 1;
  dbStore.search.result = [];
  await dbStore.getSearch();
}
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
    height: calc(100% - 125px);

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
        flex-direction: column;
        font-size: 14px;

        .name {
          // font-weight: bold;
          color: #ffc46a;
          margin-bottom: 5px;
        }

        .string {
          color: #3ad81c;
          word-break: break-all;
        }
        .number {
          color: #63a3f2;
        }
      }
    }
  }
}
</style>
