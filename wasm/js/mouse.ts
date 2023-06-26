class GoMouse {
  static init() {
    document.addEventListener("mousemove", (e: MouseEvent) => {
      if (!e.view) return;

      const width = e.view.innerWidth;
      const height = e.view.innerHeight;

      let px = (e.pageX / width) * 2 - 1;
      let py = (e.pageY / height) * 2 - 1;

      // @ts-ignore
      const canvas = window.go.canvas;

      if (canvas) {
        const fx = canvas.getBoundingClientRect().width / width;
        const fy = canvas.getBoundingClientRect().height / height;

        px /= fx;
        py /= fy;
      }

      if (window.go.memoryOperation) {
        window.go.memoryOperation.push({
          offset: window.go.pointer.mousePosition,
          value: px,
          type: "float32",
        });
        window.go.memoryOperation.push({
          offset: window.go.pointer.mousePosition + 4,
          value: -py,
          type: "float32",
        });
      }
    });

    document.addEventListener("mousedown", (e: MouseEvent) => {
      if (e.button > 2) return;

      const mp = window.go.pointer.mouseDown;
      window.go.memoryOperation.push({
        offset: mp + e.button,
        value: 1,
        type: "uint8",
      });
    });

    document.addEventListener("mouseup", (e: MouseEvent) => {
      if (e.button > 2) return;

      const mp = window.go.pointer.mouseDown;
      window.go.memoryOperation.push({
        offset: mp + e.button,
        value: 0,
        type: "uint8",
      });
    });

    document.addEventListener("click", (e: MouseEvent) => {
      if (e.button > 2) return;

      const mp = window.go.pointer.mouseClick;
      window.go.memoryOperation.push({
        offset: mp + e.button,
        value: 1,
        type: "uint8",
      });
    });
  }
}
