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
  page: number;
  selectedTable: string;
  where: string;
  search: {
    total: number;
    count: number;
    page: number;
    perPage: number;
    result: any[];
  };
}

export const useDBStore = defineStore({
  id: "db",
  state: () =>
    ({
      list: [],
      page: 0,
      search: {
        total: 0,
        count: 0,
        page: 0,
        perPage: 0,
        result: [],
      },
      selectedTable: "",
      where: "",
    } as DBState),
  actions: {
    async getList() {
      this.list = (await Axios.get(`${HOST}/debug/api/db/list`)).data;
    },
    async getSearch() {
      try {
        this.search = (
          await Axios.get(
            `${HOST}/debug/api/db/search?table=${
              this.selectedTable
            }&where=${btoa(this.where)}&page=${this.page}`
          )
        ).data;
      } catch {
        this.search.result = [];
      }
    },
    async selectTable(table: string) {
      this.selectedTable = table;
      await this.getSearch();
    },
  },
});
