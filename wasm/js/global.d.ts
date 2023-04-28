export {};

declare global {
  interface Window {
    WasmLoader: {
      load(url: string): Promise<{ go: any; module: any }>;
      getMemory(module: any): ArrayBuffer;
      sliceMemoryF64(module: any, start: number, length: number): Float32Array;
    };
    go:
      | any
      | {
          instance: WebAssembly.Instance;
          pointer: Record<string, number>;
          memory: {
            getMemory(): DataView;
          };
          memoryView: DataView;
          memoryOperation: {
            offset: number;
            value: number;
            type: string;
          }[];
        };
  }
}
