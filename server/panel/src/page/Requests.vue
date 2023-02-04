<template>
  <div :class="$style.main">
    <!--    <div :class="$style.header">
      <el-input
        placeholder="Offset..."
        v-model="offset"
        style="margin-right: 5px"
        @change="refresh"
      />
      <el-input placeholder="Limit..." v-model="limit" @change="refresh" />
    </div>-->

    <!-- Pagination -->
    <el-pagination
      background
      layout="prev, pager, next"
      :total="requestStore.search.total"
      :page-size="requestStore.limit"
      style="margin-bottom: 10px; width: 100%"
      @current-change="changePage"
    />

    <el-table
      :data="requestStore.search.result"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
      :cell-style="{ verticalAlign: 'top' }"
    >
      <!-- Method tag -->
      <el-table-column label="Method" width="100">
        <template #default="scope">
          <MethodTag :tag="scope.row.httpMethod" />
        </template>
      </el-table-column>

      <!-- Url -->
      <el-table-column label="Url">
        <template #header>
          <el-input
            v-model="requestStore.filter['url']"
            @change="refresh"
            size="small"
            placeholder="Filter by url..."
          />
        </template>
        <template #default="scope">
          {{ scope.row.url }}
        </template>
      </el-table-column>

      <el-table-column label="Args">
        <template #default="scope">
          <pre
            v-if="toggleArgs[scope.row.id]"
            v-html="formatHighlight(scope.row.args || {}, customColorOptions)"
          ></pre>
        </template>
      </el-table-column>

      <el-table-column label="Response">
        <template #default="scope">
          <pre
            v-if="toggleArgs[scope.row.id]"
            v-html="
              formatHighlight(scope.row.response || {}, customColorOptions)
            "
          ></pre>
        </template>
      </el-table-column>

      <el-table-column label="Error">
        <template #default="scope">
          <pre
            v-if="toggleArgs[scope.row.id]"
            v-html="scope.row.error ? formatHighlight(scope.row.error) : '-'"
          ></pre>
        </template>
      </el-table-column>

      <!-- Remote addr -->
      <el-table-column label="Remote IP" width="150">
        <template #default="scope"> {{ scope.row.remoteAddr }} </template>
      </el-table-column>

      <!-- Status -->
      <el-table-column label="Status" width="80">
        <template #default="scope">
          <el-tag v-if="scope.row.error?.code" type="danger">{{
            scope.row.error?.code
          }}</el-tag>
          <el-tag v-else type="success">{{ 200 }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="Created" width="180">
        <template #default="scope">
          {{ dayjs(scope.row.created).format("MMM DD HH:mm:ss (SSS)") }}
        </template>
      </el-table-column>

      <!-- Expand -->
      <el-table-column label="Expand" width="100">
        <template #default="scope">
          <el-button
            @click="toggleArgs[scope.row.id] = !toggleArgs[scope.row.id]"
            size="small"
            >Expand</el-button
          >
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useRequestStore } from "@/store/request";
import MethodTag from "@/component/MethodTag.vue";
import dayjs from "dayjs";
import formatHighlight from "json-format-highlight";

// Stores
const requestStore = useRequestStore();

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
const toggleArgs = ref({});

// Hooks
onMounted(async () => {
  tableHeight.value = window.innerHeight - 120;
  await refresh();
});

// Methods
async function refresh() {
  await requestStore.getSearch();
  console.log(requestStore.search.result);
}

async function changePage(page: number) {
  requestStore.offset = (page - 1) * requestStore.limit;
  requestStore.search.result = [];
  await requestStore.getSearch();
}
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  //height: calc(100% - 80px);

  .header {
    display: flex;
    margin-bottom: 10px;
  }

  pre {
    background: #101010;
    padding: 5px;
    box-sizing: border-box;
    word-break: break-all;
    font-size: 14px;
    margin: 0;
    max-height: 300px;
    overflow-y: auto;
    line-height: 16px;
  }

  .request {
    display: flex;
    margin-bottom: 10px;

    > div {
      flex: 1;
    }
  }
}
</style>
