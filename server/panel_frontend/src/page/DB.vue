<template>
  <div :class="$style.main">
    <div :class="$style.header">
      <el-button
        @click="changeTab(x)"
        v-for="x in dbStore.tableList"
        :type="x === dbStore.table ? 'primary' : ''"
        :key="x"
        >{{ x }}</el-button
      >
    </div>

    <!-- Pagination -->
    <el-pagination
      background
      layout="prev, pager, next"
      :total="dbStore.search.total"
      :page-size="dbStore.limit"
      style="margin-bottom: 10px; width: 100%"
      @current-change="changePage"
    />

    <!-- Data -->
    <el-table
      :data="dbStore.search.result"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
      :empty-text="dbStore.error"
    >
      <!-- Custom -->
      <el-table-column
        v-for="v in dbStore.settings.fieldList.filter((x) => !x.isHide)"
        :label="v.name"
        :width="v.width"
      >
        <!-- Filter -->
        <template v-if="v.hasFilter" #header>
          <el-input
            v-model="dbStore.filter[v.name]"
            size="small"
            @change="refresh"
            :placeholder="v.name"
          />
        </template>

        <!-- Body -->
        <template #default="scope">
          <div v-if="v.type === 'bool'">
            <el-checkbox size="large" :model-value="scope.row[v.name]" />
          </div>
          <div v-else>{{ scope.row[v.name] }}</div>
        </template>
      </el-table-column>

      <!-- Edit -->
      <el-table-column
        v-if="dbStore.settings.isEditable"
        label="Edit"
        width="70"
      >
        <template #default="scope">
          <el-button @click="editById(scope.row.id)" style="width: 100%"
            ><el-icon><EditPen /></el-icon
          ></el-button>
        </template>
      </el-table-column>

      <!-- Delete -->
      <el-table-column
        v-if="dbStore.settings.isDeletable"
        label="Delete"
        width="70"
      >
        <template #default="scope">
          <el-button @click="deleteById(scope.row.id)" style="width: 100%"
            ><el-icon><DeleteFilled /></el-icon
          ></el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-button
      v-if="dbStore.settings.isCreatable"
      @click="enableCreateMode()"
      style="margin-top: 10px"
      type="success"
      >Create new</el-button
    >

    <!-- Modal Edit -->
    <el-dialog v-model="isEditMode" title="Edit" width="40%" draggable>
      <!-- Content -->
      <div style="max-height: 600px; overflow-y: scroll">
        <div
          v-for="v in dbStore.settings.fieldList"
          :key="v.name"
          style="margin-bottom: 10px"
        >
          <ContentEditor :info="v" :temp-row="tempRow" />
        </div>
      </div>

      <template #footer>
        <el-button type="primary" @click="update"> Save </el-button>
      </template>
    </el-dialog>

    <!-- Modal Create -->
    <el-dialog v-model="isCreateMode" title="Create" width="40%" draggable>
      <div style="max-height: 600px; overflow-y: scroll">
        <div
          v-for="v in dbStore.settings.fieldList"
          :key="v.name"
          style="margin-bottom: 10px"
        >
          <ContentEditor :info="v" :temp-row="tempRow" />
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="create"> Save </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useDbStore } from "@/store/db";
import ContentEditor from "@/component/db/ContentEditor.vue";

// Stores
const dbStore = useDbStore();

// Vars
const offset = ref("0");
const tableHeight = ref(400);
const editId = ref(0);
const isCreateMode = ref(false);
const isEditMode = ref(false);
const tempRow = ref({});

// Hooks
onMounted(async () => {
  tableHeight.value = window.innerHeight - 205;

  await dbStore.getTableList();
  dbStore.table = dbStore.tableList[0];
  await dbStore.getSettings();
  await dbStore.getSearch();
});

// Methods
async function refresh() {
  await dbStore.getSearch();
}

async function deleteById(id: string) {
  if (confirm("Sure?")) {
    await dbStore.deleteById(id);
    await dbStore.getSearch();
  }
}

async function editById(id: number) {
  editId.value = id;
  tempRow.value = {};
  tempRow.value = await dbStore.getById(id);

  for (let i = 0; i < dbStore.settings.fieldList.length; i++) {
    const field = dbStore.settings.fieldList[i];
    if (field.type === "bitmask") {
      for (let j = 0; j < 64; j++) {
        // @ts-ignore
        tempRow.value[field.name + "_mask_" + j] =
          (tempRow.value[field.name] & (1 << j)) === 1 << j;
      }
    }
  }

  isEditMode.value = true;
}

async function enableCreateMode() {
  tempRow.value = {};

  for (let i = 0; i < dbStore.settings.fieldList.length; i++) {
    const field = dbStore.settings.fieldList[i];
    if (field.type === "bitmask") {
      for (let j = 0; j < 64; j++) {
        // @ts-ignore
        tempRow.value[field.name + "_mask_" + j] =
          (tempRow.value[field.name] & (1 << j)) === 1 << j;
      }
    }
  }

  isCreateMode.value = true;
}

async function changeTab(tab: string) {
  dbStore.table = tab;
  dbStore.search.result = [];
  await dbStore.getSettings();
  await dbStore.getSearch();
}

async function changePage(page: number) {
  dbStore.offset = (page - 1) * dbStore.limit;
  dbStore.search.result = [];
  await dbStore.getSearch();
}

async function update() {
  await dbStore.update(editId.value, tempRow.value);
  await dbStore.getSearch();
  isEditMode.value = false;
}

async function create() {
  await dbStore.create(tempRow.value);
  await dbStore.getSearch();
  isCreateMode.value = false;
}
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);

  .header {
    margin-bottom: 10px;
  }
}
</style>
