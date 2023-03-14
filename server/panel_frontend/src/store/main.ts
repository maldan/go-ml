import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export type MainState = {
  hasLogTab: boolean;
};

export const useMainStore = defineStore({
  id: "main",
  state: () =>
    ({
      hasLogTab: false,
    } as MainState),
  actions: {
    async getSetting() {
      const response = (await Axios.get(`${HOST}/debug/api/panel/setting`))
        .data;
      this.hasLogTab = response.hasLogTab;
    },
  },
});
