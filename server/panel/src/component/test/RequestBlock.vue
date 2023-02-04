<template>
  <div :class="[$style.main, $style['status_' + blockStatus]]">
    <div :class="$style.header">
      <MethodTag :tag="props.httpMethod" />
      <div :class="$style.url">{{ props.url }}</div>
    </div>
    <!--    <el-input type="textarea" placeholder="Name..." v-model="jsonBody" />-->
    <pre
      style="word-break: break-all; width: 100%; font-size: 12px"
      v-html="formatHighlight(props.args, customColorOptions)"
    ></pre>

    <div
      :class="[$style.constraint, $style[constraintStatus[x]]]"
      v-for="x in constraints"
      :key="x"
    >
      <div>{{ x }}</div>
      <div :class="$style.status">⦿</div>
    </div>

    <div :class="$style.input" v-for="x in input" :key="x">{{ x }}</div>
    <div :class="$style.output" v-for="x in output" :key="x">{{ x }}</div>

    <pre
      style="
        word-break: break-all;
        width: 100%;
        font-size: 12px;
        max-height: 400px;
        overflow-y: auto;
        overflow-x: hidden;
      "
      v-html="formatHighlight(responseDetails.body, customColorOptions)"
    ></pre>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useTestStore } from "@/store/test";
import Axios from "axios";
import formatHighlight from "json-format-highlight";
import MethodTag from "@/component/MethodTag.vue";
import { HOST } from "@/const";

// Props
const props = defineProps<{
  id: string;
  httpMethod: string;
  url: string;
  input?: string[];
  output?: string[];
  constraints: string[];
  args?: Record<string, any>;
}>();

// Store
const testStore = useTestStore();

// Vars
const jsonBody = ref(JSON.stringify(props.args, null, 4));
const blockStatus = ref("idle");
//const responseData = ref<Record<string, any>>({});
//const statusCode = ref(0);
const responseDetails = ref({ status: 0, body: {} });
const constraintStatus = ref<Record<string, string>>({});
const customColorOptions = ref({
  keyColor: "#af6ed1",
  numberColor: "#77b0fc",
  stringColor: "#57ab51",
  trueColor: "#ff8080",
  falseColor: "#ff8080",
  nullColor: "#e54b4b",
});
const musor = {} as any;

// Hooks
onMounted(async () => {
  testStore.setRef(props.id, {
    execute,
    musor,
  });
});

async function execute() {
  Axios.defaults.headers.common["Authorization"] =
    localStorage.getItem("debug__accessToken") || "";

  try {
    let resp = null;
    if (props.httpMethod === "GET") {
      resp = await Axios.get(`${HOST}${props.url}`, {
        params: props.args,
      });
    }
    if (props.httpMethod === "DELETE") {
      resp = await Axios.delete(`${HOST}${props.url}`, {
        params: props.args,
      });
    }

    if (props.httpMethod === "POST") {
      resp = await Axios.post(`${HOST}${props.url}`, props.args);
    }
    if (props.httpMethod === "PATCH") {
      resp = await Axios.patch(`${HOST}${props.url}`, props.args);
    }
    if (props.httpMethod === "PUT") {
      resp = await Axios.patch(`${HOST}${props.url}`, props.args);
    }

    if (!resp) {
      return;
    }

    responseDetails.value.body = resp.data;
    responseDetails.value.status = resp.status;
    testStore.setResponse(props.id, resp.data);
  } catch (e: any) {
    responseDetails.value.body = e.response.data;
    responseDetails.value.status = e.response.status;
  }

  return checkConstraints();
}

function checkConstraints() {
  for (let i = 0; i < props.constraints.length; i++) {
    const constraint = props.constraints[i];

    const response = responseDetails.value.body;
    const statusCode = responseDetails.value.status;
    const args = props.args;

    constraintStatus.value[constraint] = eval(constraint) ? "ok" : "error";
    if (!constraintStatus.value[constraint]) return false;

    musor.r1 = response;
    musor.r2 = statusCode;
    musor.r3 = args;
  }

  return true;
}
</script>

<style module lang="scss">
.main {
  background: #1b1b1b;
  padding: 5px;
  // width: 200px;
  max-width: 240px;
  position: relative;
  border-radius: 4px;
  box-sizing: border-box;

  .header {
    display: flex;
    align-items: center;

    .url {
      margin-left: 10px;
    }
  }

  pre {
    background: #121212;
    padding: 5px;
    box-sizing: border-box;
  }

  &.status_idle {
    border: 1px solid #3b3b3b;
  }

  &.status_ok {
    border: 1px solid #55c729;
  }

  &.status_error {
    border: 1px solid #da3232;
  }

  > div {
    margin-bottom: 10px;
  }

  .constraint {
    padding: 2px 5px;
    background: #212121;
    border: 1px solid #4d4d4d;
    border-radius: 4px;
    display: flex;

    .status {
      margin-left: auto;
      font-size: 16px;
    }

    &.ok {
      color: #55c729;
    }

    &.error {
      color: #da3232;
    }
  }

  .input,
  .output {
    position: relative;
    padding-left: 10px;

    &::before {
      content: "";
      display: block;
      width: 8px;
      height: 8px;
      border-radius: 6px;
      background: #309f0b;
      position: absolute;
      left: -10px;
      top: 5px;
      border: 2px solid #fefefe66;
    }
  }

  .output {
    text-align: right;
    padding-right: 10px;

    &::before {
      left: unset;
      right: -10px;
      background: #9f210b;
    }
  }
}
</style>
