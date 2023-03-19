<template>
  <div :class="$style.main">
    <!-- Header list -->
    <div :class="$style.list">
      <el-button
        @click="routerStore.selectRouter(x)"
        v-for="x in routerStore.list"
        :key="x.path"
        :type="routerStore.selectedRouter.path === x.path ? 'primary' : ''"
        effect="plain"
      >
        {{ x.path }}
        <el-tag
          size="small"
          type="success"
          effect="light"
          style="margin-left: 10px"
          >{{ x.type }}</el-tag
        >
      </el-button>
    </div>

    <!-- Controller list -->
    <div :class="$style.list">
      <el-button
        @click="routerStore.selectController(x)"
        v-for="x in routerStore.controllerList"
        :key="x"
        :type="routerStore.selectedController === x ? 'primary' : ''"
        effect="plain"
      >
        {{ x }}
      </el-button>
    </div>

    <!-- Method list -->
    <div :class="$style.methodList">
      <div
        :class="[$style.method, $style[x.httpMethod]]"
        v-for="x in routerStore.methodList"
        :key="x"
      >
        <div :class="$style.header">
          <div :class="$style.httpMethod">{{ x.httpMethod }}</div>
          <div :class="$style.url">{{ x.url }}</div>

          <!-- Args -->
          <el-popover
            v-for="y in x.args"
            :key="y"
            placement="bottom"
            trigger="click"
          >
            <template #reference>
              <el-button
                plain
                type="success"
                size="small"
                style="font-size: 14px"
                >{{ y }}</el-button
              >
            </template>
            <TypeInfo :type="routerStore.typeMap[y]" />
          </el-popover>

          <!--          <el-button
            v-for="y in x.args"
            :key="y"
            plain
            type="success"
            size="small"
            style="font-size: 14px"
            >{{ y }}</el-button
          >-->

          <el-button
            v-for="y in x.return"
            :key="y"
            plain
            type="warning"
            size="small"
            style="font-size: 14px"
            >{{ y }}</el-button
          >

          <el-button
            @click="
              routerStore.methodShowDetailed[x.id] =
                !routerStore.methodShowDetailed[x.id]
            "
            style="margin-left: auto"
            >O</el-button
          >
        </div>

        <div v-if="routerStore.methodShowDetailed[x.id]" :class="$style.body">
          <pre
            v-if="
              routerStore.responseData[x.id]?.headers?.['content-type'] ===
              'application/json'
            "
            v-html="
              formatHighlight(
                routerStore.responseData[x.id]?.body,
                customColorOptions
              )
            "
          ></pre>
          <el-button @click="execute(x)">Execute</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import { useRouterStore } from "@/store/router";
import formatHighlight from "json-format-highlight";
import TypeInfo from "@/component/router/TypeInfo.vue";

// Stores
const routerStore = useRouterStore();

// Vars
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
  await routerStore.getList();
});

// Methods
async function execute(x: any) {
  await routerStore.executeMethod(x.id, x.httpMethod, x.url, {});
}
</script>

<style module lang="scss">
$httpGet: rgba(52, 186, 241, 0.85);
$httpPost: rgba(109, 241, 52, 0.85);

.main {
  padding: 10px;
  display: flex;
  flex-direction: column;
  height: calc(100% - 80px);

  .list {
    display: flex;
    margin-bottom: 10px;
  }

  .methodList {
    overflow-y: auto;
    height: calc(100% - 85px);

    .method {
      border: 1px solid rgba(255, 255, 255, 0.4);
      padding: 10px;
      margin-bottom: 10px;

      .header {
        display: flex;
        align-items: center;

        .httpMethod {
          width: max-content;
          margin-right: 10px;
          border: 1px solid rgba(255, 255, 255, 0.4);
          padding: 2px 8px;
          font-size: 14px;
          border-radius: 4px;
        }

        .url {
          width: 240px;
          font-size: 16px;
          background: rgba(255, 255, 255, 0.1);
          padding: 2px 8px;
          border-radius: 4px;
          margin-right: 10px;
          border: 1px solid rgba(255, 255, 255, 0.2);
        }
      }

      .body {
        margin-top: 10px;

        pre {
          word-break: break-all;
          white-space: pre-wrap;
          padding: 10px;
          border: 1px solid rgba(255, 255, 255, 0.4);
          border-radius: 4px;
        }
      }

      &.GET {
        border: 1px solid $httpGet;

        .httpMethod {
          border: 1px solid $httpGet;
          color: $httpGet;
        }
      }

      &.POST {
        border: 1px solid $httpPost;

        .httpMethod {
          border: 1px solid $httpPost;
          color: $httpPost;
        }
      }
    }
  }
}
</style>
