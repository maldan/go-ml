import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import type { IMethod } from "@/store/method";

export interface Router {
  path: string;
  type: string;
}

export interface RouterState {
  list: Router[];
  selectedRouter: Router;
  selectedController: string;
  controllerList: string[];
  methodList: any[];
}

export const useRouterStore = defineStore({
  id: "router",
  state: () =>
    ({
      list: [],
      selectedRouter: { path: "", type: "" },
      selectedController: "",
      controllerList: [],
      methodList: [],
    } as RouterState),
  actions: {
    async getList() {
      this.list = (await Axios.get(`${HOST}/debug/api/router/list`)).data;
    },
    async getControllerList() {
      this.controllerList = (
        await Axios.get(
          `${HOST}/debug/api/router/controllerList?path=${this.selectedRouter.path}`
        )
      ).data;
    },
    async getMethodList() {
      this.methodList = (
        await Axios.get(
          `${HOST}/debug/api/router/methodList?path=${this.selectedRouter.path}&controller=${this.selectedController}`
        )
      ).data;
      console.log(this.methodList);
    },
    async selectRouter(x: Router) {
      this.selectedRouter = x;
      await this.getControllerList();
    },
    async selectController(x: string) {
      this.selectedController = x;
      await this.getMethodList();
    },
  },
});
