var gas = {};

globalThis.fs.open = async function (path, flags, mode, callback) {
  const p = await fetch(path);
  const b = await p.blob();
  const body = await b.arrayBuffer();
  // console.log(callback);
  // console.log(path, flags, mode);
  // console.log(path);
  // console.log(callback);
  gas[10] = new Uint8Array(body);
  callback(null, 10);
};

globalThis.fs.read = async function (
  fd,
  buffer,
  offset,
  length,
  position,
  callback
) {
  for (let i = 0; i < gas[fd].length; i++) {
    buffer[i] = gas[fd][i];
  }
  //buffer[0] = 1;
  //console.log(fd, buffer, offset, length, position);
  callback(null, buffer.length);
};

globalThis.fs.fstat = async function (fd, callback) {
  console.log("fstat", fd);
  callback(null, {
    isDirectory: function () {
      return false;
    },
    dev: 0,
    ino: 0,
    mode: 0,
    nlink: 0,
    uid: 0,
    gid: 0,
    rdev: 0,
    size: gas[fd].length,
    blksize: 0,
    blocks: 0,
    atimeMs: 0,
    mtimeMs: 0,
    ctimeMs: 0,

    /* GetSize() {
           return 10;
       },
       IsDir() {
           return false;
       }*/
  });
};

globalThis.fs.stat = async function (path, callback) {
  console.log(path);
};

globalThis.fs.close = async function (fd, callback) {
  console.log(fd);
};

globalThis.process.cwd = function () {
  return "/";
};
