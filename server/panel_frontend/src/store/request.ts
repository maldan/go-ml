import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export interface IRequest {
  id: string;
  httpMethod: string;
  inputHeader: string;
  inputBody: string;
  outputHeader: string;
  outputBody: string;
  url: string;
  created: string;
}

export type RequestState = {
  list: IRequest[];
  filter: Record<string, string>;
  /*search: {
    total: number;
    count: number;
    page: number;
    result: any[];
  };
  offset: number;
  limit: number;
  error: string;
  filter: Record<string, string>;*/
};

export const useRequestStore = defineStore({
  id: "request",
  state: () =>
    ({
      list: [],
      filter: {
        httpMethod: "",
        created: "",
      },
      /*search: {
        result: [],
        total: 0,
        page: 0,
        count: 0,
      },
      filter: {},
      offset: 0,
      limit: 20,
      error: "",*/
    } as RequestState),
  actions: {
    async getList() {
      this.list.length = 0;

      const r = (
        await Axios.get(
          `${HOST}/debug/api/request/list?httpMethod=${
            this.filter["httpMethod"]
          }&created=${encodeURIComponent(
            this.filter["created"]
          )}&timezoneOffset=${new Date().getTimezoneOffset()}`
        )
      ).data;

      this.list = r.result;

      /*this.error = "";
      this.search.result = [];

      try {
        this.search = (
          await Axios.get(
            `${HOST}/debug/api/requestList?offset=${this.offset}&limit=${this.limit}`
          )
        ).data;
      } catch (e: any) {
        this.error = e.response.data.description;
        this.search.total = 0;
      }

      console.log(this.search);*/
    },
  },
});
