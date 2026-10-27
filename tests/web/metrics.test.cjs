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
  const elements = Object.fromEntries(["section", "message", "projectMetrics", "sprintSelect", "sprintMetrics", "projectHours", "sprintHours", "projectDefects", "sprintDefects"].map(name => [name, new Element()]));
  const context = vm.createContext({ document: { createElement: () => new Element() }, encodeURIComponent, Intl });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/metrics.js"), "utf8"), context);
  return { ...elements, ui: context.createMetricsSection({ ...elements, requestJSON }) };
}

const metrics = { planned_story_points: 10, completed_story_points: 4, team_velocity: null, completion_percentage: 40 };
const hours = { estimated_hours: 5.75, actual_hours: 6.25, deviation: 0.5 };

test("BDD-US-20-01: project and selected Sprint metrics are shown", async () => {
  const calls = [];
  const e = setup(async url => {
    calls.push(url);
    if (url.endsWith("/metrics/hours")) return url.includes("/sprints/") ? hours : { metrics: hours, sprints: [{ id: 3, goal: "Goal", closed: true }] };
    return url.includes("/sprints/") ? metrics : { metrics, sprints: [{ id: 3, goal: "Goal", closed: true }] };
  });
  await e.ui.showProject("project/id");
  assert.equal(e.section.hidden, false);
  assert.equal(e.projectMetrics.children[0].textContent, "Story Points planificados");
  assert.match(e.projectMetrics.children[5].textContent, /No disponible/);
  assert.equal(e.projectHours.children[0].textContent, "Horas estimadas");
  assert.match(e.projectHours.children[5].textContent, /0,50/);
  assert.equal(e.sprintSelect.children[1].textContent, "Sprint 3 — Goal (cerrado)");
  e.sprintSelect.value = "3";
  await e.sprintSelect.listeners.change();
  assert.ok(calls.some(url => url.includes("project%2Fid") && url.endsWith("/metrics/hours")));
  assert.ok(calls.some(url => url.includes("sprints/3/metrics/story-points")));
  assert.ok(calls.some(url => url.includes("sprints/3/metrics/hours")));
  assert.match(e.sprintMetrics.children[7].textContent, /40,00 %/);
  assert.match(e.sprintHours.children[5].textContent, /0,50/);
});

test("BDD-US-20-05: no Sprint history displays project data and resets selection", async () => {
  const e = setup(async url => url.endsWith("/metrics/hours")
    ? { metrics: { estimated_hours: null, actual_hours: null, deviation: null }, sprints: [] }
    : { metrics: { planned_story_points: null, completed_story_points: null, team_velocity: null, completion_percentage: null }, sprints: [] });
  await e.ui.showProject("project");
  assert.equal(e.sprintSelect.disabled, true);
  assert.equal(e.sprintSelect.children[0].textContent, "No hay Sprints");
  assert.equal(e.projectMetrics.children[1].textContent, "No disponible");
  assert.equal(e.projectHours.children[1].textContent, "No disponible");
  await e.ui.showProject("");
  assert.equal(e.section.hidden, true);
  assert.equal(e.projectMetrics.children.length, 0);
  assert.equal(e.projectHours.children.length, 0);
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

test("BDD-US-22-01: project and selected Sprint show defect counts independently", async () => {
  const calls = [];
  const defectMetrics = { detected_defects: 3, resolved_defects: null };
  const e = setup(async url => {
    calls.push(url);
    if (url.endsWith("/metrics/defects")) {
      return url.includes("/sprints/") ? defectMetrics : { metrics: defectMetrics, sprints: [{ id: 3, goal: "Goal", closed: true }] };
    }
    if (url.endsWith("/metrics/hours")) return { metrics: hours, sprints: [{ id: 3, goal: "Goal", closed: true }] };
    return { metrics, sprints: [{ id: 3, goal: "Goal", closed: true }] };
  });

  await e.ui.showProject("project");
  assert.equal(e.projectDefects.children[0].textContent, "Defectos detectados");
  assert.equal(e.projectDefects.children[1].textContent, "3");
  assert.equal(e.projectDefects.children[3].textContent, "No disponible");
  assert.ok(calls.some(url => url === "/api/projects/project/metrics/defects"));

  e.sprintSelect.value = "3";
  await e.sprintSelect.listeners.change();
  assert.equal(e.sprintDefects.children[1].textContent, "3");
  assert.equal(e.sprintDefects.children[3].textContent, "No disponible");
  assert.ok(calls.some(url => url === "/api/projects/project/sprints/3/metrics/defects"));
});
