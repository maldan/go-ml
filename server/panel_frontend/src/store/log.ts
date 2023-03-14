import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import { escape } from "lodash";

export type LogState = {
  page: number;
  total: number;
  list: { kind: string; body: string; created: string; isJsonBody: boolean }[];
  filter: Record<string, string>;
};

export const useLogStore = defineStore({
  id: "log",
  state: () =>
    ({
      page: 1,
      total: 0,
      list: [],
      filter: {
        kind: "",
        created: "",
      },
    } as LogState),
  actions: {
    async getList() {
      this.list.length = 0;

      const r = (
        await Axios.get(
          `${HOST}/debug/api/log/list?page=${this.page}&kind=${
            this.filter["kind"]
          }&created=${encodeURIComponent(
            this.filter["created"]
          )}&timezoneOffset=${new Date().getTimezoneOffset()}`
        )
      ).data;

      for (let i = 0; i < r.result.length; i++) {
        try {
          const newBody = JSON.parse(r.result[i].body);
          if (typeof newBody === "object") {
            r.result[i].body = newBody;
          }
          r.result[i].isJsonBody = true;
        } catch {}
      }

      this.list = r.result;
    },
  },
});
