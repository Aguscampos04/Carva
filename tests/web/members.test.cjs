// US-03 BDD 04/07/12/13/15/16. Built-in Node runner; no added dependencies.
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

class Element {
  constructor(tagName = "div") { this.tagName = tagName; this.children = []; this.listeners = {}; this.textContent = ""; this.value = ""; this.hidden = true; this.disabled = false; }
  append(...children) { this.children.push(...children); }
  replaceChildren(...children) { this.children = [...children]; }
  addEventListener(event, callback) { this.listeners[event] = callback; }
  reset() { this.resets = (this.resets || 0) + 1; }
}

function setup(requestJSON) {
  const elements = Object.fromEntries(["section", "form", "nameInput", "message", "list", "refresh", "submit"].map(key => [key, new Element()]));
  const context = vm.createContext({ document: { createElement: tag => new Element(tag) }, encodeURIComponent });
  const source = fs.readFileSync(path.join(__dirname, "../../web/static/js/members.js"), "utf8");
  vm.runInContext(source, context);
  return { ...elements, ui: context.createMembersSection({ ...elements, requestJSON }) };
}
const submitEvent = { preventDefault() {} };

test("BDD-07/04/15: name-only form submits once and renders homonyms as safe text with secondary full UUID", async () => {
  const records = []; const calls = [];
  const e = setup(async (url, options = {}) => {
    calls.push({ url, options });
    if (options.method === "POST") { records.push({ id: records.length === 0 ? "b1d168da-325a-40d1-90b7-b5a1ba006389" : "a20a4fad-e66e-4889-bcaa-65af6a39bdec", name: JSON.parse(options.body).name }); return records.at(-1); }
    return [...records];
  });
  await e.ui.showProject("demo");
  assert.equal(e.section.hidden, false);
  assert.match(e.message.textContent, /Todavía no hay miembros/);
  e.nameInput.value = "<b>Ana</b>";
  await Promise.all([e.form.listeners.submit(submitEvent), e.form.listeners.submit(submitEvent)]);
  assert.equal(calls.filter(c => c.options.method === "POST").length, 1);
  assert.equal(calls.at(-2).options.body, '{"name":"<b>Ana</b>"}');
  assert.equal(e.list.children[0].children[0].textContent, "<b>Ana</b>");
  assert.equal(e.list.children[0].children[1].tagName, "small");
  assert.equal(e.list.children[0].children[1].textContent, "UUID: b1d168da-325a-40d1-90b7-b5a1ba006389");
  assert.equal(e.form.resets, 2); // project context reset, then successful submission reset
  assert.equal(e.submit.disabled, false);
  e.nameInput.value = "<b>Ana</b>";
  await e.form.listeners.submit(submitEvent);
  assert.equal(e.list.children.length, 2);
  assert.equal(e.list.children[0].children[0].textContent, e.list.children[1].children[0].textContent);
  assert.notEqual(e.list.children[0].children[1].textContent, e.list.children[1].children[1].textContent);
  await e.ui.showProject("");
  assert.equal(e.section.hidden, true);
  assert.equal(e.list.children.length, 0);
});

test("BDD-16: a late list from another project cannot replace the current list", async () => {
  let resolveOld;
  const e = setup(url => url.includes("/old/") ? new Promise(resolve => { resolveOld = resolve; }) : Promise.resolve([{ name: "Actual", id: "current" }]));
  const oldLoad = e.ui.showProject("old");
  await e.ui.showProject("new");
  resolveOld([{ name: "Anterior", id: "old" }]);
  await oldLoad;
  assert.equal(e.list.children[0].children[0].textContent, "Actual");
});

test("BDD-16: a late registration cannot reset or populate another project's form", async () => {
  let resolvePost;
  const e = setup((url, options = {}) => options.method === "POST" ? new Promise(resolve => { resolvePost = resolve; }) : Promise.resolve([]));
  await e.ui.showProject("old"); e.nameInput.value = "Ana";
  const pending = e.form.listeners.submit(submitEvent);
  await e.ui.showProject("new");
  const resets = e.form.resets;
  e.nameInput.value = "Bruno";
  resolvePost({ id: "old", name: "Ana" }); await pending;
  assert.equal(e.nameInput.value, "Bruno");
  assert.equal(e.form.resets, resets);
  assert.equal(e.list.children.length, 0);
});

test("BDD-12/13: storage errors are visible without creating a phantom member or a successful empty list", async () => {
  let fail = false;
  const e = setup(async () => { if (fail) throw new Error("No se pudo consultar"); return []; });
  await e.ui.showProject("demo"); fail = true;
  await e.refresh.listeners.click();
  assert.equal(e.message.textContent, "No se pudo consultar");
  e.nameInput.value = "Ana";
  await e.form.listeners.submit(submitEvent);
  assert.equal(e.message.textContent, "No se pudo consultar");
  assert.equal(e.list.children.length, 0);
  assert.equal(e.form.resets, 1); // only the project change reset the form
  assert.equal(e.submit.disabled, false);
});

test("successful POST followed by failed refresh reports that the member was already registered", async () => {
  let posted = false;
  const e = setup(async (_, options = {}) => {
    if (options.method === "POST") { posted = true; return {}; }
    if (posted) throw new Error("Falló la consulta");
    return [];
  });
  await e.ui.showProject("demo"); e.nameInput.value = "Ana";
  await e.form.listeners.submit(submitEvent);
  assert.match(e.message.textContent, /Miembro registrado.*Falló la consulta/);
});

test("same-project overlapping refreshes keep the newest response", async () => {
  const pending = [];
  const e = setup(() => new Promise(resolve => pending.push(resolve)));
  const first = e.ui.showProject("demo");
  const second = e.refresh.listeners.click();
  pending[1]([{ name: "Nuevo", id: "2" }]); await second;
  pending[0]([{ name: "Viejo", id: "1" }]); await first;
  assert.equal(e.list.children[0].children[0].textContent, "Nuevo");
});
