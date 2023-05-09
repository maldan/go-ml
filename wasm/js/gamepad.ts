class GoGamepad {
  static poll(id: number) {
    const gamepads = navigator.getGamepads();
    const gp = gamepads[id];
    if (!gp) {
      return;
    }

    // Set buttons
    let ptr = window.go.pointer.gamepadKeyState_0;
    for (let i = 0; i < gp.buttons.length; i++) {
      window.go.memory.writeI8(ptr + i, gp.buttons[i].pressed ? 1 : 0);
    }
    // console.log(gp.buttons.map((x) => x.pressed));

    // Set axes
    ptr = window.go.pointer.gamepadAxisState_0;
    for (let i = 0; i < gp.axes.length; i++) {
      window.go.memory.writeF32(ptr, gp.axes[i]);
      ptr += 4;
    }
  }
}
