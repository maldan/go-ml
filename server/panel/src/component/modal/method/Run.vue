<template>
  <div :class="$style.container">
    <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 10px">
      <Tree
        :class="$style.field"
        :data="modalStore.data.input"
        style="height: 400px; overflow-y: auto"
      />
      <pre :class="$style.field">{{
        JSON.stringify(modalStore.data.input, null, 4)
      }}</pre>
      <pre :class="$style.field">{{ JSON.stringify(response, null, 4) }}</pre>
    </div>

    <Row size="1fr 1fr" style="margin-top: 8px; gap: 10px">
      <Button @click="run" color="gray" text="Run" />
      <Button @click="modalStore.cancel()" color="gray" text="Close" />
    </Row>
  </div>
</template>

<script setup lang="ts">
import { Button, Row, Tree } from "@/gam-lib-ui/vue/component/ui";
import { useModalStore } from "@/gam-lib-ui/vue/store/modal";
import { ref } from "vue";
import Axios from "axios";

const modalStore = useModalStore();
const response = ref({});

const run = async () => {
  if (modalStore.data.method === "GET") {
    response.value = (
      await Axios.get(modalStore.data.url, {
        params: modalStore.data.input,
      })
    ).data;
  }
  if (modalStore.data.method === "DELETE") {
    response.value = (
      await Axios.delete(modalStore.data.url, {
        params: modalStore.data.input,
      })
    ).data;
  }
  if (modalStore.data.method === "POST") {
    response.value = (
      await Axios.post(modalStore.data.url, modalStore.data.input)
    ).data;
  }
  if (modalStore.data.method === "PUT") {
    response.value = (
      await Axios.put(modalStore.data.url, modalStore.data.input)
    ).data;
  }
  if (modalStore.data.method === "PATCH") {
    response.value = (
      await Axios.patch(modalStore.data.url, modalStore.data.input)
    ).data;
  }
};
</script>

<style module lang="scss">
@import "@/gam-lib-ui/vue/vars";

.container {
  .field {
    border: 1px solid $color-white-020;
    padding: 10px;
    margin: 0;
  }
}
</style>
