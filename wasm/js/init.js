class WasmLoader {
  static async load(url) {
    const go = new Go();
    const b = await fetch(url);
    const bytes = await b.blob();

    return {
      go,
      module: await WebAssembly.instantiate(
        await bytes.arrayBuffer(),
        go.importObject
      ),
    };
  }

  static getMemory(module) {
    if (module.instance.exports.mem) {
      return module.instance.exports.mem.buffer;
    }
    return module.instance.exports.memory.buffer;
  }

  static sliceMemoryF64(module, start, length) {
    const memory = this.getMemory(module);
    const f32 = new Float32Array(memory);
    return f32.subarray(start / 4, start / 4 + length);
  }
}
window.WasmLoader = WasmLoader;
