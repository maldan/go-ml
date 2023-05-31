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
  writeBytes(ptr: number, val: Uint8Array) {
    const m = this.getMemory();
    for (let i = 0; i < val.length; i++) {
      m.setUint8(ptr + i, val[i]);
    }
  },
  writeI32(ptr: number, val: number) {
    const m = this.getMemory();
    m.setUint32(ptr, val, true);
  },
  writeF32(ptr: number, val: number) {
    const m = this.getMemory();
    m.setFloat32(ptr, val, true);
  },
  readF32(ptr: number): number {
    const m = this.getMemory();
    return m.getFloat32(ptr, true);
  },
  readI8(ptr: number): number {
    const m = this.getMemory();
    return m.getInt8(ptr);
  },
  readU32(ptr: number): number {
    const m = this.getMemory();
    return m.getUint32(ptr, true);
  },
  /*  readU64(ptr: number): number {
    const m = this.getMemory();
    return m.getUint64(ptr, true);
  },*/
  readString(ptr: number, len: number): string {
    const m = this.getMemory();
    const x = [];
    for (let i = 0; i < len; i++) {
      x.push(m.getUint8(ptr + i));
    }
    return new TextDecoder().decode(new Uint8Array(x));
    // return m.getUint32(ptr, true);
  },
  readSlice(ptr: number, len: number): Uint8Array {
    const m = this.getMemory();
    const x = [];
    for (let i = 0; i < len; i++) {
      x.push(m.getUint8(ptr + i));
    }
    return new Uint8Array(x);
  },
};
