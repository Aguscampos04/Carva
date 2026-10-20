// US-02 AC-05/06/12; BDD-08/09/20. Node built-in runner, no browser framework.
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const path = require("node:path");
function setup(requestJSON, confirmFinalization = () => true) {
  const elements = Object.fromEntries(["form", "nameInput", "startInput", "endOutput", "finalizeButton", "saveButton", "message"].map(key => [key, {
    value: "", textContent: "", disabled: false, hidden: false, listeners: {},
    addEventListener(event, callback) { this.listeners[event] = callback; },
  }]));
  const context = vm.createContext({ encodeURIComponent });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/projects.js"), "utf8"), context);
  const saved = [];
  const ui = context.createProjectEditor({ ...elements, requestJSON, confirmFinalization, onSaved(project) { saved.push(project); ui.showProject(project); } });
  ui.showProject({ id: "demo", name: "Carva Demo", status: "Planificado", start_date: null, end_date: null });
  return { ...elements, ui, saved };
}
const event = { preventDefault() {} };
test("BDD-08: confirmation before write; cancel preserves draft and sends no request", async () => {
  const calls = []; let prompt;
  const e = setup(async (...args) => calls.push(args), text => { prompt = text; return false; });
  e.nameInput.value = "Unsaved";
  await e.finalizeButton.listeners.click();
  assert.equal(calls.length, 0);
  assert.equal(e.nameInput.value, "Unsaved");
  assert.match(prompt, /Carva Demo/); assert.match(prompt, /datos guardados/);
});
test("BDD-07/08/09: finalize sends only confirmation; finalized still editable", async () => {
  const calls = [];
  const e = setup(async (url, options) => {
    calls.push({ url, options });
    return { id: "demo", name: "Carva Demo", status: "Finalizado", start_date: null, end_date: "2026-10-20" };
  });
  e.nameInput.value = "Unsaved"; await e.finalizeButton.listeners.click();
  assert.equal(calls[0].url, "/api/projects/demo/finalize");
  assert.equal(calls[0].options.body, '{"confirmed":true}');
  assert.equal(e.finalizeButton.hidden, true); assert.equal(e.nameInput.disabled, false);
  assert.equal(e.startInput.disabled, false); assert.equal(e.saveButton.disabled, false);
  assert.equal(e.endOutput.textContent, "2026-10-20");
  e.nameInput.value = "Edited"; e.startInput.value = "";
  await e.form.listeners.submit(event);
  assert.equal(calls[1].options.method, "PATCH");
  assert.equal(calls[1].options.body, '{"name":"Edited","start_date":null}');
});
test("errors preserve draft, no phantom success, duplicate submit blocked", async () => {
  let reject; let count = 0;
  const e = setup(() => { ++count; return new Promise((_, r) => { reject = r; }); });
  e.nameInput.value = "Draft";
  const first = e.form.listeners.submit(event);
  await e.form.listeners.submit(event);
  assert.equal(count, 1); assert.equal(e.saveButton.disabled, true);
  reject(new Error("Falló; consultar de nuevo antes de reintentar")); await first;
  assert.equal(e.saved.length, 0); assert.equal(e.nameInput.value, "Draft");
  assert.match(e.message.textContent, /Falló/); assert.equal(e.saveButton.disabled, false);
});
test("late mutations from another project cannot replace selection or reset draft", async () => {
  let resolve;
  const e = setup(() => new Promise(r => { resolve = r; }));
  const pending = e.finalizeButton.listeners.click();
  e.ui.showProject({ id: "new", name: "New", status: "Planificado", start_date: null, end_date: null });
  e.nameInput.value = "New draft";
  resolve({ id: "demo", status: "Finalizado", end_date: "2026-10-20" }); await pending;
  assert.equal(e.saved.length, 0); assert.equal(e.nameInput.value, "New draft");
  assert.equal(e.finalizeButton.hidden, false);
});
