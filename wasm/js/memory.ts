if (!window.go) window.go = {};

window.go.memory = {
  getMemory() {
    const mm =
      window.go.instance.exports.mem?.buffer ??
      window.go.instance.exports.memory?.buffer;
    return new DataView(mm);
  },
  writeI8(ptr: number, val: number) {
    const m = this.getMemory();
    m.setUint8(ptr, val);
  },
  writeI32(ptr: number, val: number) {
    const m = this.getMemory();
    m.setUint32(ptr, val, true);
  },
  readF32(ptr: number): number {
    const m = this.getMemory();
    return m.getFloat32(ptr, true);
  },
};
