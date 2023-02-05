<template>
  <div :class="$style.main">
    <el-button @click="tag = x" v-for="x in tags" :type="x === tag ? 'primary' : ''" :key="x">{{ x }}</el-button>

    <div :class="$style.caseList" style="margin-top: 10px" v-if="tag">
      <div v-for="testCase in testCaseList" :key="testCase.name" :class="$style.case">
        <el-button @click="testStore.runCase(testCase.tag, testCase.name)" style="margin-bottom: 10px">Run</el-button>
        <div :class="$style.list">
          <component
            v-for="x in testCase.blockList"
            :key="x.id"
            :id="x.id"
            :args="x.args"
            :http-method="x.httpMethod"
            :url="x.url"
            :input="x.input"
            :output="x.output"
            :constraints="x.constraints"
            :is="requestBlock"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from 'vue';
import type { IMethod } from '@/store/method';
import { useMethodStore } from '@/store/method';
import RequestBlock from '@/component/test/RequestBlock.vue';
import CaptureBlock from '@/component/test/CaptureBlock.vue';
import { useTestStore } from '@/store/test';

// Stores
const methodStore = useMethodStore();
const testStore = useTestStore();

// Vars
const requestBlock = RequestBlock;
const tags = computed(() => {
  return [...new Set(testStore.list.map((x) => x.tag))];
});
const testCaseList = computed(() => {
  return testStore.list.filter((x) => x.tag === tag.value);
});
const tag = ref('');

// Hooks
onMounted(async () => {
  await testStore.getList();
});

// Methods
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);

  .caseList {
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    height: calc(100% - 40px);

    .case {
      background: #2b2b2b;
      padding: 10px;
      margin-bottom: 10px;
      display: flex;
      flex-direction: column;

      .list {
        display: flex;

        > div {
          margin-right: 20px;
        }
      }
    }
  }
}
</style>
