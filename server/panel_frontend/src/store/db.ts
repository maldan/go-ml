import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import type { IMethod } from "@/store/method";

export interface DB {
  path: string;
  type: string;
}

export interface DBState {
  list: string[];
  recordList: any[];
  selectedTable: string;
  where: string;
}

export const useDBStore = defineStore({
  id: "db",
  state: () =>
    ({
      list: [],
      recordList: [],
      selectedTable: "",
      where: "",
    } as DBState),
  actions: {
    async getList() {
      this.list = (await Axios.get(`${HOST}/debug/api/db/list`)).data;
    },
    async getSearch() {
      try {
        this.recordList = (
          await Axios.get(
            `${HOST}/debug/api/db/search?table=${
              this.selectedTable
            }&where=${btoa(this.where)}`
          )
        ).data;
      } catch {
        this.recordList = [];
      }
    },
    async selectTable(table: string) {
      this.selectedTable = table;
      await this.getSearch();
    },
  },
});
