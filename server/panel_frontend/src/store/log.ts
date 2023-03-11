import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export type LogState = {
  list: { kind: string; created: string }[];
};

export const useLogStore = defineStore({
  id: "log",
  state: () =>
    ({
      list: [],
    } as LogState),
  actions: {
    async getList() {
      this.list.length = 0;

      const lines = (await Axios.get(`${HOST}/debug/api/log/list`)).data;

      for (let i = 0; i < lines.length; i++) {
        try {
          this.list.push(JSON.parse(lines[i]));
        } catch {}
      }
    },
  },
});
