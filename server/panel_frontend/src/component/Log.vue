<template>
  <div :class="$style.log">
    <div :class="[$style.kind, $style[kind]]">{{ kind }}</div>
    <div :class="$style.time">
      {{ dayjs(time).format("YYYY-MM-DD") }}
    </div>
    <div :class="$style.time">
      {{ dayjs(time).format("HH:mm:ss.SSS") }}
    </div>
    <div :class="$style.message">{{ message }}</div>
  </div>
</template>

<script setup lang="ts">
import dayjs from "dayjs";

defineProps<{
  kind: string;
  message: string;
  time: string;
}>();
</script>

<style lang="scss" module>
.log {
  display: flex;
  font-size: 16px;
  align-items: center;
  font-family: Consolas, sans-serif;

  .kind {
    padding: 2px 15px 2px 5px;
    background: #a31c1c;
    color: #fefefe;
    // border-radius: 4px;
    font-weight: bold;
    text-transform: uppercase;
    font-size: 12px;
    margin-right: 15px;
    position: relative;
    width: 32px;

    &::after {
      content: "";
      display: block;
      position: absolute;
      width: 12px;
      height: 12px;
      transform: rotate(45deg);
      right: -7px;
      top: 3px;
      background: #2b2b2b;
    }

    &.error {
      background: #a31c1c;

      &::after {
        background: #a31c1c;
      }
    }

    &.info {
      background: #198bb4;

      &::after {
        background: #198bb4;
      }
    }
  }

  .time {
    color: #51df35;
    margin-right: 5px;
    background: #0a5e05;
    font-size: 12px;
    padding: 2px 5px;
  }

  .message {
    color: #999999;
  }
}
</style>
