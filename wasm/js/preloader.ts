class Preloader {
  static fileMap: Record<string, any> = {};

  static load(url: string) {
    return new Promise<any>((resolve, reject) => {
      let blob;
      let xmlHTTP = new XMLHttpRequest();
      xmlHTTP.open("GET", url, true);
      xmlHTTP.responseType = "arraybuffer";
      xmlHTTP.onload = async function (e) {
        blob = new Blob([this.response]);

        if (url.includes(".js")) {
          resolve(await blob.text());
        } else {
          resolve(await blob.arrayBuffer());
        }
      };
      xmlHTTP.onprogress = function (pr) {
        const percentage = (pr.loaded / pr.total) * 100;

        const out = document.getElementById("preloader_" + url);
        if (out)
          out.innerHTML = `
            <div style="width: ${percentage}%; background: #c32727; height: 100%;"></div>
                
            <div style="width: 100%; position: absolute; left: 0; top: 0; height: 100%; color: #fefefe; display: flex; align-items: center; justify-content: center;">
            ${(pr.loaded / 1048576).toFixed(2)} MB / 
            ${(pr.total / 1048576).toFixed(2)} MB
           </div>
        `;
      };
      xmlHTTP.onloadend = function (e) {
        // onEnd();
      };
      xmlHTTP.send();
    });
  }

  static init(list: string[]) {
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
      }
    `;
    document.head.appendChild(head);

    const div = document.createElement("div");
    document.body.appendChild(div);

    div.innerHTML = `
      <div id="preloader">
        ${list
          .map(
            (x) =>
              `<div style="width: 100%; height: 16px; position: relative; margin-bottom: 5px;" id="preloader_${x}"></div>`
          )
          .join("")} 
      </div>
    `;
  }

  static async loadList(list: string[]) {
    Preloader.init(list);

    for (let i = 0; i < list.length; i++) {
      this.fileMap[list[i]] = await this.load(list[i]);
    }

    const out = document.getElementById("preloader");
    if (out) out.style.display = `none`;
  }

  static injectJs() {
    for (const key in this.fileMap) {
      if (key.includes(".js")) {
        const script = document.createElement("script");
        script.type = "text/javascript";
        document.body.prepend(script);
        script.text = this.fileMap[key];
      }
    }
  }
}
