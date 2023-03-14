<template>
  <div :class="$style.main">
    <!-- Pagination -->
    <el-pagination
      background
      layout="prev, pager, next"
      :total="logStore.total"
      :page-size="1"
      style="margin-bottom: 10px; width: 100%"
      @current-change="changePage"
    />

    <!-- Data -->
    <el-table
      :data="logStore.list"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
    >
      <el-table-column prop="kind" label="Kind" width="110" />

      <!-- Body -->
      <el-table-column label="Body">
        <template #default="scope">
          <pre
            v-if="scope.row.isJsonBody"
            v-html="formatHighlight(scope.row.body, customColorOptions)"
          ></pre>
          <div v-else>{{ scope.row.body }}</div>
        </template>
      </el-table-column>

      <!-- File -->
      <el-table-column label="File" width="250">
        <template #default="scope">
          <div v-if="!scope.row.file">-</div>
          <div v-else>{{ scope.row.file }}:{{ scope.row.line }}</div>
        </template>
      </el-table-column>

      <!-- Created -->
      <el-table-column label="Created" width="220">
        <template #default="scope">
          <div :class="$style.created">
            <el-tag type="danger" effect="plain">
              {{ dayjs(scope.row.created).format("YYYY-MM-DD") }}
            </el-tag>
            <el-tag type="success" effect="plain">
              {{ dayjs(scope.row.created).format("HH:mm:ss") }}
            </el-tag>
            <el-tag type="warning" effect="plain">
              {{ dayjs(scope.row.created).format("SSS") }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!--    <div :class="$style.item" v-for="x in logStore.list" :key="x.created">
      <div style="flex: none; width: 128px">{{ x.kind }}</div>

      &lt;!&ndash; Raw &ndash;&gt;
      <div>
        <div v-if="x.kind === 'raw'">-</div>
        <div v-else>{{ x.file }}:{{ x.line }}</div>
      </div>

      <div>{{ x.body }}</div>

      <div v-if="x.kind === 'raw'" :class="$style.created">-</div>

      <div v-if="x.kind !== 'raw'" :class="$style.created">
        <div :class="$style.date">
          {{ dayjs(x.created).format("YYYY-MM-DD") }}
        </div>
        <div :class="$style.time">
          {{ dayjs(x.created).format("HH:mm:ss") }}
        </div>
        <div :class="$style.ms">{{ dayjs(x.created).format("SSS") }}</div>
      </div>
    </div>-->
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useRequestStore } from "@/store/request";
import { useMainStore } from "@/store/main";
import { useLogStore } from "@/store/log";
import dayjs from "dayjs";
import formatHighlight from "json-format-highlight";

// Stores
const logStore = useLogStore();

// Vars
const tableHeight = ref(400);
const customColorOptions = ref({
  keyColor: "#af6ed1",
  numberColor: "#77b0fc",
  stringColor: "#57ab51",
  trueColor: "#ff8080",
  falseColor: "#ff8080",
  nullColor: "#e54b4b",
});

// Hooks
onMounted(async () => {
  await refresh();
  tableHeight.value = window.innerHeight - 120;
});

// Methods
async function changePage(page: number) {
  logStore.page = page;
  await refresh();
}
async function refresh() {
  await logStore.getList();
  console.log(logStore.list);
}
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  // overflow-y: scroll;
  height: calc(100% - 85px);

  .item {
    display: flex;
    border-bottom: 1px solid rgba(255, 255, 255, 0.5);
    padding: 10px;

    > div {
      flex: 1;
    }
  }

  pre {
    padding: 5px 10px;
    margin: 0;
    background: rgba(0, 0, 0, 0.5);
    border-radius: 4px;
  }

  .created {
    display: grid;
    grid-template-columns: 80px 65px 1fr;
    gap: 7px;

    .date,
    .time,
    .ms {
      background: rgba(255, 0, 0, 0.1);
      border: 1px solid rgba(255, 0, 0, 0.4);
      color: rgba(255, 0, 0, 0.8);
      padding: 1px 3px;
      border-radius: 4px;
      font-size: 12px;
      text-align: center;
    }

    .time {
      background: rgba(110, 178, 41, 0.1);
      border: 1px solid rgba(110, 178, 41, 0.4);
      color: rgba(110, 178, 41, 0.8);
    }

    .ms {
      background: rgba(208, 0, 255, 0.1);
      border: 1px solid rgba(208, 0, 255, 0.4);
      color: rgba(208, 0, 255, 0.8);
    }
  }
}
</style>
