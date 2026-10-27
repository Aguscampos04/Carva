const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const path = require("node:path");

class Element {
  constructor() { this.children = []; this.listeners = {}; this.textContent = ""; this.value = ""; this.disabled = false; this.hidden = false; }
  append(...nodes) { this.children.push(...nodes); }
  replaceChildren(...nodes) { this.children = nodes; }
  addEventListener(name, callback) { this.listeners[name] = callback; }
}

function setup(requestJSON) {
  const elements = Object.fromEntries(["section", "message", "projectMetrics", "sprintSelect", "sprintMetrics"].map(name => [name, new Element()]));
  const context = vm.createContext({ document: { createElement: () => new Element() }, encodeURIComponent, Intl });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/metrics.js"), "utf8"), context);
  return { ...elements, ui: context.createMetricsSection({ ...elements, requestJSON }) };
}

const metrics = { planned_story_points: 10, completed_story_points: 4, team_velocity: null, completion_percentage: 40 };

test("BDD-US-20-01: project and selected Sprint metrics are shown", async () => {
  const calls = [];
  const e = setup(async url => {
    calls.push(url);
    return url.includes("/sprints/") ? metrics : { metrics, sprints: [{ id: 3, goal: "Goal", closed: true }] };
  });
  await e.ui.showProject("project/id");
  assert.equal(e.section.hidden, false);
  assert.equal(e.projectMetrics.children[0].textContent, "Story Points planificados");
  assert.match(e.projectMetrics.children[5].textContent, /No disponible/);
  assert.equal(e.sprintSelect.children[1].textContent, "Sprint 3 — Goal (cerrado)");
  e.sprintSelect.value = "3";
  await e.sprintSelect.listeners.change();
  assert.match(calls[0], /project%2Fid/);
  assert.match(calls[1], /sprints\/3\/metrics\/story-points/);
  assert.match(e.sprintMetrics.children[7].textContent, /40,00 %/);
});

test("BDD-US-20-05: no Sprint history displays project data and resets selection", async () => {
  const e = setup(async () => ({ metrics: { planned_story_points: null, completed_story_points: null, team_velocity: null, completion_percentage: null }, sprints: [] }));
  await e.ui.showProject("project");
  assert.equal(e.sprintSelect.disabled, true);
  assert.equal(e.sprintSelect.children[0].textContent, "No hay Sprints");
  assert.equal(e.projectMetrics.children[1].textContent, "No disponible");
  await e.ui.showProject("");
  assert.equal(e.section.hidden, true);
  assert.equal(e.projectMetrics.children.length, 0);
});

test("metric request errors stay visible and can be retried", async () => {
  let fail = true;
  const e = setup(async () => { if (fail) throw new Error("Sin conexión"); return { metrics, sprints: [] }; });
  await e.ui.showProject("project");
  assert.equal(e.message.textContent, "Sin conexión");
  fail = false;
  await e.ui.showProject("project");
  assert.match(e.message.textContent, /Métricas calculadas/);
});
