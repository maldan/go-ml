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
  fileList: any[];
  typeMap: any;
  responseData: any;
  methodShowDetailed: any;
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
      fileList: [],
      typeMap: {},
      responseData: {},
      methodShowDetailed: {},
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
    async getTypeList() {
      const list = (
        await Axios.get(
          `${HOST}/debug/api/router/typeList?path=${this.selectedRouter.path}&controller=${this.selectedController}`
        )
      ).data;
      for (let i = 0; i < list.length; i++) {
        this.typeMap[list[i].name] = list[i];
      }
    },
    async getFileList() {
      this.fileList = (
        await Axios.get(
          `${HOST}/debug/api/router/fileList?path=${this.selectedRouter.path}`
        )
      ).data;
    },
    async selectRouter(x: Router) {
      this.selectedRouter = x;
      if (x.type == "API") {
        await this.getControllerList();
      } else {
        await this.getFileList();
      }
    },
    async selectController(x: string) {
      this.selectedController = x;
      await this.getTypeList();
      await this.getMethodList();
    },
    async executeMethod(
      id: string,
      httpMethod: string,
      url: string,
      args: any
    ) {
      Axios.defaults.headers.common["Authorization"] =
        localStorage.getItem("debug__accessToken") || "";

      url = `${HOST}${url}`;
      console.log(url);

      try {
        let time = new Date().getTime();
        let response = null;
        if (httpMethod === "GET") {
          response = await Axios.get(url, {
            params: args,
          });
        }
        if (httpMethod === "DELETE") {
          response = await Axios.delete(url, {
            params: args,
          });
        }
        if (httpMethod === "POST") {
          response = await Axios.post(url, args);
        }
        if (httpMethod === "PUT") {
          response = await Axios.put(url, args);
        }
        if (httpMethod === "PATCH") {
          response = await Axios.patch(url, args);
        }
        if (response) {
          this.responseData[id] = {
            headers: response.headers,
            status: response.status,
            body: response.data,
            time: new Date().getTime() - time,
          };
          //this.responseInfo[uid].status = response.status;
          //this.responseInfo[uid].time = new Date().getTime() - time;
        }
      } catch (e: any) {
        console.log(e);
        //this.response[uid] = e.response?.data || {};
        //this.responseInfo[uid].status = e.response.status;
      }
    },
  },
});
