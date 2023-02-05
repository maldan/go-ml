import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import { escape } from "lodash";

export interface ControlState {
  commandList: { folder: string; name: string }[];
  status: Record<string, boolean>;
}

export const useControlStore = defineStore({
  id: "control",
  state: () =>
    ({
      commandList: [],
      status: {},
    } as ControlState),
  actions: {
    async getCommandList() {
      this.commandList = (await Axios.get(`${HOST}/debug/control/list`)).data;
    },
    async execute(folder: string, name: string) {
      this.status[folder + "/" + name] = true;
      await Axios.post(
        `${HOST}/debug/control/execute?folder=${folder}&name=${name}`
      );
      this.status[folder + "/" + name] = false;
    },
  },
});
