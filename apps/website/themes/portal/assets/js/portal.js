(function () {
  "use strict";

  function pick(list) {
    return list[Math.floor(Math.random() * list.length)];
  }

  // tabs (blog / shrines / changelog). without js every pane just shows stacked.
  document.querySelectorAll("[data-tabs]").forEach(function (root) {
    var tabs = Array.prototype.slice.call(root.querySelectorAll("[role=tab]"));
    var panes = tabs.map(function (t) {
      return document.getElementById(t.getAttribute("aria-controls"));
    });

    function show(index) {
      tabs.forEach(function (t, i) {
        var on = i === index;
        t.setAttribute("aria-selected", on ? "true" : "false");
        t.tabIndex = on ? 0 : -1;
        if (panes[i]) panes[i].hidden = !on;
      });
    }

    tabs.forEach(function (t, i) {
      t.addEventListener("click", function () { show(i); });
      t.addEventListener("keydown", function (e) {
        var next = null;
        if (e.key === "ArrowRight") next = (i + 1) % tabs.length;
        if (e.key === "ArrowLeft") next = (i - 1 + tabs.length) % tabs.length;
        if (next !== null) { show(next); tabs[next].focus(); }
      });
    });
    show(0);
  });

  // random propaganda poster + random header gif on every load
  function randomSrc(id, attr) {
    var el = document.getElementById(id);
    if (!el) return;
    try {
      var list = JSON.parse(el.getAttribute(attr));
      if (list && list.length) el.src = pick(list);
    } catch (e) { /* keep the build-time default */ }
  }
  randomSrc("poster", "data-posters");
  randomSrc("mascot-img", "data-gifs");

  // site search: fetches /index.json once, matches every term against the whole entry, shows a modal
  var input = document.getElementById("q");
  var box = document.getElementById("results");
  var modal = document.getElementById("search-modal");
  var closeBtn = document.getElementById("search-close");
  var form = document.getElementById("search-form");
  if (!input || !box || !modal || !form) return;

  var index = null;
  var loading = null;
  var lastFocus = null;

  function load() {
    if (index) return Promise.resolve(index);
    if (!loading) {
      loading = fetch("/index.json")
        .then(function (r) { return r.json(); })
        .then(function (data) { index = data || []; return index; })
        .catch(function () { index = []; return index; });
    }
    return loading;
  }

  function path(url) {
    try { return new URL(url, location.origin).pathname; } catch (e) { return url; }
  }

  function norm(text) {
    return String(text).normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
  }

  function snippet(text) {
    var parts = String(text || "").trim().split(/(?<=[.!?])\s+/);
    return parts.slice(0, 2).join(" ") + (parts.length > 2 ? "..." : "");
  }

  function openModal() {
    if (!modal.hidden) return;
    lastFocus = document.activeElement;
    modal.hidden = false;
    document.body.classList.add("modal-open");
    closeBtn.focus();
  }

  function closeModal() {
    if (modal.hidden) return;
    modal.hidden = true;
    document.body.classList.remove("modal-open");
    if (lastFocus && lastFocus.focus) lastFocus.focus();
  }

  function text(tag, cls, value) {
    var el = document.createElement(tag);
    if (cls) el.className = cls;
    el.textContent = value;
    return el;
  }

  function render(items, term) {
    box.textContent = "";
    if (!items.length) {
      box.appendChild(text("li", "none", "nothing found for \"" + term + "\""));
    }
    items.forEach(function (item) {
      var li = document.createElement("li");
      var h = document.createElement("h3");
      var a = document.createElement("a");
      a.href = path(item.permalink);
      a.textContent = item.title;
      h.appendChild(a);
      li.appendChild(h);
      var s = snippet(item.summary);
      if (s) li.appendChild(text("p", "", s));
      if (item.tags && item.tags.length) li.appendChild(text("p", "tags", "tagged: " + [].concat(item.tags).join(", ")));
      box.appendChild(li);
    });
    openModal();
  }

  function run() {
    var raw = input.value.trim();
    if (!raw) { closeModal(); return; }
    var terms = norm(raw).split(/\s+/);
    load().then(function (data) {
      render(data.filter(function (p) {
        var hay = norm(JSON.stringify(p));
        return terms.every(function (t) { return hay.indexOf(t) !== -1; });
      }), raw);
    });
  }

  input.addEventListener("focus", load);
  form.addEventListener("submit", function (e) { e.preventDefault(); run(); });
  closeBtn.addEventListener("click", closeModal);
  modal.addEventListener("click", function (e) { if (e.target === modal) closeModal(); });
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") closeModal();
  });
})();
