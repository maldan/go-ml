class Preloader {
  static fileMap: Record<string, any> = {};
  static state = {
    totalFiles: 0,
    loadedFiles: 0,
    files: {},
  } as any;
  static isInit = false;

  static load(url: string) {
    return new Promise<any>((resolve, reject) => {
      let blob;
      let xmlHTTP = new XMLHttpRequest();
      xmlHTTP.open("GET", url, true);
      xmlHTTP.responseType = "arraybuffer";
      xmlHTTP.onload = async function (e) {
        blob = new Blob([this.response]);

        if (url.includes(".js")) {
          Preloader.fileMap[url] = await blob.text();
          resolve(true);
          //resolve(await blob.text());
        } else {
          Preloader.fileMap[url] = await blob.arrayBuffer();
          resolve(true);
        }

        Preloader.state.loadedFiles += 1;
      };
      xmlHTTP.onprogress = function (pr) {
        // console.log(pr);
        if (!Preloader.state.files[url]) {
          const contentLength = Number.parseInt(
            xmlHTTP.getResponseHeader("Content-Length") || "0"
          );
          console.log(contentLength);
          Preloader.state.files[url] = {
            total: pr.total || contentLength,
            loaded: pr.loaded,
          };
        } else {
          Preloader.state.files[url].loaded = pr.loaded;
        }

        /*const percentage = (pr.loaded / pr.total) * 100;

        const out = document.getElementById("preloader_" + url);
        if (out)
          out.innerHTML = `
            <div style="width: ${percentage}%; background: #c32727; height: 100%;"></div>
                
            <div style="width: 100%; position: absolute; left: 0; top: 0; height: 100%; color: #fefefe; display: flex; align-items: center; justify-content: center;">
            ${(pr.loaded / 1048576).toFixed(2)} MB / 
            ${(pr.total / 1048576).toFixed(2)} MB
           </div>
        `;*/
      };
      xmlHTTP.onloadend = function (e) {
        // onEnd();
      };
      xmlHTTP.send();
    });
  }

  static init() {
    if (Preloader.isInit) return;
    Preloader.isInit = true;

    const head = document.createElement("style");
    head.innerHTML = `
      #preloader {
        width: 320px;
        position: absolute;
        left: calc(50% - 320px/2);
        top: 50%;
        border: 2px solid rgba(255, 255, 255, 0.25);
        padding: 5px;
        font-size: 12px;
        height: 16px;
      }
    `;
    document.head.appendChild(head);

    const div = document.createElement("div");
    document.body.appendChild(div);

    div.innerHTML = `<div id="preloader"></div>`;
  }

  static async loadList(list: string[]) {
    Preloader.init();

    const intervalId = setInterval(() => {
      const out = document.getElementById("preloader");
      if (!out) return;
      let totalSize = 0;
      let loadedSize = 0;
      let loaded = 0;
      for (let key in Preloader.state.files) {
        loadedSize += Preloader.state.files[key].loaded;
        totalSize += Preloader.state.files[key].total;
        if (loadedSize >= totalSize) loaded += 1;
      }
      let percentage = loadedSize / totalSize;
      if (percentage > 1.0) {
        percentage = 1.0;
      }

      out.innerHTML = `
          <div style="width: ${
            percentage * 100
          }%; background: #c32727; height: 100%;"></div>
          
          <div style="width: 100%; position: absolute; left: 0; top: 0; height: 100%; color: #fefefe; display: flex; align-items: center; justify-content: center;">
          ${(loadedSize / 1048576).toFixed(2)} MB / 
          ${(totalSize / 1048576).toFixed(2)} MB
          ${loaded} / ${list.length}
         </div>
      `;
    }, 16);

    let list2 = [];
    for (let i = 0; i < list.length; i++) {
      list2.push(this.load(list[i]));
      // this.fileMap[list[i]] = await this.load(list[i]);
    }
    await Promise.all(list2);

    clearInterval(intervalId);
  }

  static hide() {
    const out = document.getElementById("preloader");
    if (out) out.style.display = `none`;
  }

  static injectJs() {
    for (const key in this.fileMap) {
      if (key.match(/\.js$/)) {
        const script = document.createElement("script");
        script.type = "text/javascript";
        document.body.prepend(script);
        script.text = this.fileMap[key];
      }
    }
  }
}
