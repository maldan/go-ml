import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export interface IMethodInput {
  name: string;
  kind: string;
  type: string;
  fieldList: IMethodInput[];
}

export interface IMethod {
  uid: string;
  controller: string;
  fullPath: string;
  httpMethod: string;
  name: string;
  url: string;
  inputMethod: string;
  input: IMethodInput;
}

export type MethodState = {
  items: IMethod[];
  methodsMap: Record<string, any>;
  params: Record<string, any>;
  formData: Record<string, any>;
  response: Record<string, any>;
  responseInfo: Record<string, { time: number; status: number }>;
};

function fillArgs(args: IMethodInput, out: Record<string, any>) {
  const fl = args?.fieldList || [];
  for (let i = 0; i < fl.length; i++) {
    if (fl[i].name == "accessToken") continue;
    if (fl[i].kind === "string") out[fl[i].name] = "";
    if (fl[i].kind === "bool") out[fl[i].name] = false;
    if (fl[i].kind === "slice") out[fl[i].name] = [];
    if (fl[i].kind === "map") out[fl[i].name] = {};
    if (fl[i].kind === "struct") {
      out[fl[i].name] = {};
      fillArgs(fl[i], out[fl[i].name]);
    }
    if (fl[i].kind.match(/^int|^float/)) out[fl[i].name] = 0;
  }
  return out;
}

export const useMethodStore = defineStore({
  id: "method",
  state: () =>
    ({
      items: [],
      params: {},
      formData: {},
      response: {},
      responseInfo: {},
      methodsMap: {},
    } as MethodState),
  actions: {
    async getList() {
      let items = (await Axios.get(`${HOST}/debug/api/methodList`))
        .data as IMethod[];
      for (let i = 0; i < items.length; i++) {
        this.params[items[i].uid] = JSON.stringify(
          fillArgs(items[i].input, {}),
          null,
          4
        );
        this.formData[items[i].uid] = {};
        this.responseInfo[items[i].uid] = { status: 0, time: 0 };
        this.methodsMap[items[i].uid] = items[i];
      }

      items = items.sort((a: IMethod, b: IMethod) => {
        return (
          a.controller.localeCompare(b.controller) ||
          a.name.localeCompare(b.name)
        );
      });

      this.items = items;
    },
    async run(uid: string, method: string, url: string, args: any) {
      Axios.defaults.headers.common["Authorization"] =
        localStorage.getItem("debug__accessToken") || "";

      try {
        let time = new Date().getTime();
        let response = null;
        if (method === "GET") {
          response = await Axios.get(url, {
            params: args,
          });
        }
        if (method === "DELETE") {
          response = await Axios.delete(url, {
            params: args,
          });
        }
        if (method === "POST") {
          response = await Axios.post(url, args);
        }
        if (method === "PUT") {
          response = await Axios.put(url, args);
        }
        if (method === "PATCH") {
          response = await Axios.patch(url, args);
        }
        if (response) {
          this.response[uid] = response.data;
          this.responseInfo[uid].status = response.status;
          this.responseInfo[uid].time = new Date().getTime() - time;
        }
      } catch (e: any) {
        this.response[uid] = e.response?.data || {};
        this.responseInfo[uid].status = e.response.status;
      }
    },
  },
});
