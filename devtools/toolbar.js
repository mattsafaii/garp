// garp dev toolbar — served at /_garp/toolbar.js by `garp dev` only.
// Never written to site/; nothing here ships with the deliverable.
(function () {
  "use strict";

  // Live reload — shipped behavior, must always run regardless of the UI.
  new EventSource("/_garp/reload").onmessage = function () {
    location.reload();
  };

  // All UI lives in a shadow root so page CSS and the toolbar can't touch
  // each other. The host is a custom tag name to dodge page element selectors.
  var host = document.createElement("garp-devtools");
  host.style.cssText = "position:fixed;z-index:2147483647";
  var root = host.attachShadow({ mode: "open" });

  var style = document.createElement("style");
  style.textContent = [
    ":host { all: initial; }",
    "* { box-sizing: border-box; }",
    ".btn {",
    "  position: fixed; right: 16px; bottom: 16px; width: 36px; height: 36px;",
    "  border: 0; border-radius: 50%; cursor: pointer;",
    "  background: rgba(20, 20, 24, 0.65); color: rgba(255, 255, 255, 0.85);",
    "  font: 600 16px/36px ui-sans-serif, system-ui, sans-serif; text-align: center;",
    "  padding: 0;",
    "}",
    ".btn:hover { background: rgba(20, 20, 24, 0.9); }",
    ".panel {",
    "  position: fixed; right: 16px; bottom: 60px; width: 380px; max-width: calc(100vw - 32px);",
    "  max-height: 70vh; overflow-y: auto;",
    "  background: #16161a; color: #d8d8de;",
    "  font: 12.5px/1.5 ui-sans-serif, system-ui, sans-serif;",
    "  border: 1px solid #2c2c33; border-radius: 8px;",
    "  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.45);",
    "  padding: 12px 14px;",
    "}",
    ".hdr { margin-bottom: 10px; }",
    ".hdr .url { font-weight: 600; color: #fff; }",
    ".hdr .src { color: #85858f; font-size: 11.5px; }",
    ".row { margin: 4px 0; }",
    ".row .label { color: #85858f; margin-right: 6px; }",
    "table { width: 100%; border-collapse: collapse; margin-top: 10px; }",
    "td { padding: 4px 6px 4px 0; border-top: 1px solid #26262c; vertical-align: top; }",
    ".key { font: 700 12px ui-monospace, monospace; color: #fff; white-space: nowrap; }",
    ".badge {",
    "  display: inline-block; padding: 1px 6px; border-radius: 8px;",
    "  font-size: 10.5px; white-space: nowrap;",
    "}",
    ".badge.frontmatter { background: #143a24; color: #6fd695; }",
    ".badge.dirdata { background: #3a2e12; color: #e0b45c; }",
    ".badge.globaldata { background: #142c40; color: #6cb2e8; }",
    ".badge.config { background: #2e1c3e; color: #c390e8; }",
    ".badge.other { background: #2c2c33; color: #a0a0aa; }",
    ".val { font: 12px ui-monospace, monospace; color: #c5c5cf; word-break: break-word; }",
    ".val pre { margin: 4px 0 0; white-space: pre-wrap; word-break: break-word; }",
    "details > summary { cursor: pointer; color: #c5c5cf; }",
    ".err { color: #e88; }",
    ".muted { color: #85858f; }",
  ].join("\n");
  root.appendChild(style);

  var btn = document.createElement("button");
  btn.className = "btn";
  btn.type = "button";
  btn.title = "garp — page data";
  btn.tabIndex = -1;
  btn.textContent = "g";
  // Don't let the toggle steal focus from whatever the page has focused.
  btn.addEventListener("mousedown", function (e) { e.preventDefault(); });
  root.appendChild(btn);

  var panel = document.createElement("div");
  panel.className = "panel";
  panel.hidden = true;
  root.appendChild(panel);

  document.body.appendChild(host);

  function el(tag, className, text) {
    var node = document.createElement(tag);
    if (className) node.className = className;
    if (text !== undefined) node.textContent = text;
    return node;
  }

  // Which cascade layer a data key came from, for badge coloring.
  function sourceKind(source) {
    if (source === "frontmatter") return "frontmatter";
    if (source === "config.yaml") return "config";
    if (/(^|\/)_data\.ya?ml$/.test(source)) return "dirdata";
    if (source.indexOf("data/") === 0) return "globaldata";
    return "other";
  }

  // Pretty value, one line when short; <details> when it would be noisy.
  function valueNode(value) {
    var pretty = JSON.stringify(value, null, 1);
    if (pretty === undefined) pretty = "undefined";
    var oneLine = pretty.replace(/\n\s*/g, " ");
    var wrap = el("div", "val");
    if (pretty.length <= 120) {
      wrap.textContent = oneLine;
      return wrap;
    }
    var details = document.createElement("details");
    details.appendChild(el("summary", null, oneLine.slice(0, 100) + "…"));
    details.appendChild(el("pre", null, pretty));
    wrap.appendChild(details);
    return wrap;
  }

  function renderError(message) {
    panel.textContent = "";
    panel.appendChild(el("div", "err", message));
  }

  function render(info) {
    panel.textContent = "";

    var hdr = el("div", "hdr");
    hdr.appendChild(el("div", "url", info.url));
    hdr.appendChild(el("div", "src", info.source));
    panel.appendChild(hdr);

    var layouts = el("div", "row");
    layouts.appendChild(el("span", "label", "layouts"));
    var chain = info.layout_chain || [];
    layouts.appendChild(
      chain.length ? el("span", null, chain.join(" → ")) : el("span", "muted", "no layout")
    );
    panel.appendChild(layouts);

    var colls = el("div", "row");
    colls.appendChild(el("span", "label", "collections"));
    var names = info.collections || [];
    colls.appendChild(
      names.length ? el("span", null, names.join(", ")) : el("span", "muted", "none")
    );
    panel.appendChild(colls);

    var table = document.createElement("table");
    (info.data || []).forEach(function (entry) {
      var tr = document.createElement("tr");
      var keyCell = document.createElement("td");
      keyCell.appendChild(el("span", "key", entry.key));
      tr.appendChild(keyCell);
      var srcCell = document.createElement("td");
      srcCell.appendChild(el("span", "badge " + sourceKind(entry.source), entry.source));
      tr.appendChild(srcCell);
      var valCell = document.createElement("td");
      valCell.appendChild(valueNode(entry.value));
      tr.appendChild(valCell);
      table.appendChild(tr);
    });
    panel.appendChild(table);
  }

  // Fetched on every open, never on load — zero network cost until used.
  function load() {
    panel.textContent = "";
    panel.appendChild(el("div", "muted", "loading…"));
    var path = location.pathname;
    fetch("/_garp/page?path=" + encodeURIComponent(path))
      .then(function (res) {
        if (!res.ok) {
          renderError("no page data for " + path + " — is this a content page?");
          return;
        }
        return res.json().then(render);
      })
      .catch(function () {
        renderError("couldn't reach garp dev at /_garp/page");
      });
  }

  function toggle(open) {
    var show = open !== undefined ? open : panel.hidden;
    panel.hidden = !show;
    if (show) load();
  }

  btn.addEventListener("click", function () { toggle(); });
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape" && !panel.hidden) toggle(false);
  });
})();
