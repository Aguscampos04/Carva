// US-11 BDD-01/04/05/06: exercise the UI controller without external libraries.
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const path = require("node:path");
class Element {
  constructor(tag = "div") { this.tagName = tag; this.children = []; this.listeners = {}; this.textContent = ""; this.value = ""; }
  append(...nodes) { this.children.push(...nodes); }
  replaceChildren(...nodes) { this.children = nodes; }
  addEventListener(name, callback) { this.listeners[name] = callback; }
}
function setup(requestJSON) {
  const elements = Object.fromEntries(["form", "projectInput", "select", "message", "details", "goal"].map(name => [name, new Element()]));
  const context = vm.createContext({ document: { createElement: tag => new Element(tag) }, encodeURIComponent });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/sprint-history.js"), "utf8"), context);
  return { ...elements, ui: context.createSprintHistory({ ...elements, requestJSON }) };
}
const event = { preventDefault() {} };
const closedSprint = { id: 7, projectID: "p", goal: "Original goal", closed: true };
const snapshot = { backlog_item_id: "story", title: "<b>Original</b>", description: "Original description", priority: "Alta", planned_status: "Pendiente", status: "Pendiente", acceptance_criteria_verified: false, provisional_story_points: 3, story_points: null, acceptance_criteria: ["Original criterion"] };
function allText(element) { return [element.textContent, ...element.children.map(allText)].join(" "); }

test("BDD-01: select history renders original fields and results as safe text using GET only", async () => {
  const calls = [];
  const e = setup(async (...args) => { calls.push(args); return args[0].endsWith("/closed") ? [closedSprint] : args[0].endsWith("/backlog-items") ? [snapshot, { ...snapshot, title: "Completed", acceptance_criteria_verified: true, status: "Completada" }] : closedSprint; });
  await e.ui.showProject("p");
  assert.equal(e.select.children[1].textContent, "Sprint 7 — Original goal");
  e.select.value = "7";
  await e.select.listeners.change();
  assert.match(e.goal.textContent, /Original goal/);
  assert.match(allText(e.details), /<b>Original<\/b>/);
  assert.match(allText(e.details), /Original criterion/);
  assert.match(allText(e.details), /Incompleta/);
  assert.match(allText(e.details), /Completada/);
  assert.match(allText(e.details), /Story Points provisionales: 3/);
  assert.equal(calls.every(args => args.length === 1), true);
});

test("BDD-04/05: empty list, empty Sprint and retry after errors", async () => {
  let fail = false; let list = [];
  const e = setup(async url => { if (fail) throw new Error("No encontrado"); return url.endsWith("/closed") ? list : url.endsWith("/backlog-items") ? [] : closedSprint; });
  await e.ui.showProject("p");
  assert.match(e.message.textContent, /No hay Sprints cerrados/);
  fail = true;
  await e.form.listeners.submit(event);
  assert.match(e.message.textContent, /No encontrado/);
  assert.equal(e.details.children.length, 0);
  fail = false; list = [closedSprint];
  await e.form.listeners.submit(event);
  e.select.value = "7";
  await e.select.listeners.change();
  assert.match(allText(e.details), /Sin historias planificadas/);
});

test("BDD-06: late list and detail responses cannot overwrite another project", async () => {
  let finishList, finishDetail;
  const e = setup(url => url.includes("/old/") ? new Promise(resolve => { finishList = resolve; }) : url.endsWith("/closed") ? Promise.resolve([closedSprint]) : url.endsWith("/backlog-items") ? Promise.resolve([snapshot]) : new Promise(resolve => { finishDetail = resolve; }));
  const old = e.ui.showProject("old");
  await e.ui.showProject("p");
  finishList([{ id: 99, goal: "Stale list" }]); await old;
  assert.equal(e.select.children.length, 2);
  e.select.value = "7";
  const detail = e.select.listeners.change();
  await e.ui.showProject("new");
  finishDetail(closedSprint); await detail;
  assert.equal(e.goal.textContent, "");
  assert.equal(e.details.children.length, 0);
});

test("BDD-06: clearing a selection invalidates a pending detail; detail errors stay visible", async () => {
  let finish; let fail = false;
  const e = setup(url => url.endsWith("/closed") ? Promise.resolve([closedSprint]) : fail ? Promise.reject(new Error("Storage error")) : url.endsWith("/backlog-items") ? Promise.resolve([]) : new Promise(resolve => { finish = resolve; }));
  await e.ui.showProject("p"); e.select.value = "7";
  const loading = e.select.listeners.change();
  e.select.value = ""; await e.select.listeners.change();
  finish(closedSprint); await loading;
  assert.equal(e.details.children.length, 0);
  fail = true; e.select.value = "7"; await e.select.listeners.change();
  assert.match(e.message.textContent, /Storage error/);
});
