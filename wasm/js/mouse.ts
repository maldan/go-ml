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

      // @ts-ignore
      window.go.setMousePosition(px, -py);
    });

    document.addEventListener("mousedown", (e: MouseEvent) => {
      // @ts-ignore
      window.go.setMouseDown(e.button, true);
    });

    document.addEventListener("mouseup", (e: MouseEvent) => {
      // @ts-ignore
      window.go.setMouseDown(e.button, false);
    });

    document.addEventListener("click", (e: MouseEvent) => {
      // @ts-ignore
      window.go.setMouseClick(e.button, true);
    });
  }
}
