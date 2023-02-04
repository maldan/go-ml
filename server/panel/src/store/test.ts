import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import type { IMethod } from "@/store/method";

export interface ITest {
  id: string;
  api: {
    execute: () => boolean;
  };
  httpMethod: string;
  url: string;
  input: string[];
  output: string[];
  constraints: string[];
  args: Record<string, any>;
}

export interface ITestCase {
  name: string;
  tag: string;
  blockList: ITest[];
  connectionList: {
    fromId: string;
    toId: string;
    fromField: string;
    toField: string;
  }[];
}

export type TestState = {
  list: ITestCase[];
  response: Record<string, any>;
};

export const useTestStore = defineStore({
  id: "test",
  state: () =>
    ({
      list: [],
      response: {},
    } as TestState),
  actions: {
    async getList() {
      this.list = (await Axios.get(`${HOST}/debug/api/testList`)).data;
      for (let i = 0; i < this.list.length; i++) {
        if (!this.list[i].blockList) this.list[i].blockList = [];
        if (!this.list[i].connectionList) this.list[i].connectionList = [];
      }
    },
    setRef(blockId: string, api: any) {
      for (let i = 0; i < this.list.length; i++) {
        const block = this.list[i].blockList.find((x) => x.id == blockId);
        if (block) {
          block.api = api;
        }
      }
    },
    setResponse(blockId: string, response: any) {
      this.response[blockId] = response;
    },
    getOutput(blockInId: string) {
      for (let i = 0; i < this.list.length; i++) {
        const conn = this.list[i].connectionList.find(
          (x) => x.toId === blockInId
        );
        if (conn) {
          return eval(`this.response.${conn.fromId}.${conn.fromField}`);
        }
      }
    },
    /*async run() {
      for (let i = 0; i < this.blocks.length; i++) {
        const block = this.blocks[i];
        if (block.input.length) {
          block.input.forEach((x) => {
            console.log(this.getOutput(block.id));
            block.args[x] = this.getOutput(block.id);
          });
        }
        const status = await block.api.execute();
        if (!status) break;
      }
    },*/
    async runCase(tag: string, name: string) {
      const testCase = this.list.find((x) => x.tag === tag && x.name === name);
      if (!testCase) return;

      for (let i = 0; i < testCase.blockList.length; i++) {
        const block = testCase.blockList[i];
        if (block.input.length) {
          block.input.forEach((x) => {
            block.args[x] = this.getOutput(block.id);
          });
        }
        const status = await block.api.execute();
        if (!status) break;
      }
    },
  },
});
