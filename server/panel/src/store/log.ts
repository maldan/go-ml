import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export type LogState = {
  items: { kind: string; message: string; time: string }[];
};

export const useLogStore = defineStore({
  id: "log",
  state: () =>
    ({
      items: [],
    } as LogState),
  actions: {
    async getList(date: string) {
      this.items = (
        await Axios.get(`${HOST}/debug/log/search?date=${date}`)
      ).data.response;
    },
  },
});
