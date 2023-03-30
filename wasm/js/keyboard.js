document.addEventListener("keydown", (e) => {
  console.log(e.keyCode);
  __golangWebglKeyboardBind_KeyDown(e.keyCode);
});

document.addEventListener("keyup", (e) => {
  __golangWebglKeyboardBind_KeyUp(e.keyCode);
});
