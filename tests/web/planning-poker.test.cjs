// US-12: member identities are selected per project and open rounds expose only the selected member's vote.
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

class Element {
  constructor(tagName = "div") {
    this.tagName = tagName;
    this.children = [];
    this.listeners = {};
    this.textContent = "";
    this.value = "";
    this.hidden = false;
    this.disabled = false;
    this.validity = { valid: true };
  }
  append(...children) { this.children.push(...children); }
  replaceChildren(...children) { this.children = [...children]; }
  addEventListener(event, callback) { this.listeners[event] = callback; }
  setAttribute() {}
}

test("BDD-US-12-02: selector loads only the chosen participant's vote and saves corrections", async () => {
  const calls = [];
  const members = [{ id: "member-1", name: "Ana" }, { id: "member-2", name: "Bea" }];
  const context = vm.createContext({
    document: { createElement: tag => new Element(tag) },
    encodeURIComponent,
    URL,
  });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/planning-poker.js"), "utf8"), context);
  const widget = context.createPlanningPoker({
    async requestJSON(url, options = {}) {
      calls.push({ url, options });
      if (url.endsWith("/rounds")) return [{ id: "round-1", status: "open", participants: members.map(item => item.id) }];
      if (url.endsWith("/members")) return members;
      if (options.method === "PUT") return { saved: true };
      const selected = new URL(`http://carva${url}`).searchParams.get("member_id");
      return [{ member_id: selected, story_points: selected === "member-1" ? 3 : 8 }];
    },
  });

  const section = widget.createWidget("project-1", { id: "story-1" });
  await new Promise(resolve => setImmediate(resolve));
  const [, message, , form] = section.children;
  const finalize = section.children[4];
  const [memberLabel, pointsLabel, submit] = form.children;
  const select = memberLabel.children[0];
  const points = pointsLabel.children[0];
  assert.equal(select.children.length, 2);
  assert.equal(select.value, "member-1");
  assert.equal(points.value, "3");
  assert.equal(finalize.hidden, false);
  assert.match(message.textContent, /oculto/);

  select.value = "member-2";
  await select.listeners.change();
  assert.equal(points.value, "8");
  assert.match(calls.at(-1).url, /member_id=member-2$/);

  points.value = "13";
  await form.listeners.submit({ preventDefault() {} });
  const save = calls.find(call => call.options.method === "PUT");
  assert.match(save.url, /votes\/member-2$/);
  assert.deepEqual(JSON.parse(save.options.body), { story_points: 13 });
});

test("BDD-US-13-02/04: finalizing reveals every vote and displays differences", async () => {
  const members = [{ id: "member-1", name: "Ana" }, { id: "member-2", name: "Bea" }];
  const votes = [
    { member_id: "member-1", story_points: 3 },
    { member_id: "member-2", story_points: 8 },
  ];
  let finalized = false;
  const context = vm.createContext({
    document: { createElement: tag => new Element(tag) },
    encodeURIComponent,
    URL,
  });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/planning-poker.js"), "utf8"), context);
  const widget = context.createPlanningPoker({
    async requestJSON(url, options = {}) {
      if (url.endsWith("/rounds")) {
        return [{ id: "round-1", number: 1, status: finalized ? "finalized" : "open", participants: members.map(item => item.id) }];
      }
      if (url.endsWith("/members")) return members;
      if (url.includes("/finalize?")) {
        assert.equal(options.method, "POST");
        finalized = true;
        return { round: { status: "finalized" }, votes, has_differences: true };
      }
      if (url.includes("/result?")) return { round: { status: "finalized" }, votes, has_differences: true };
      if (url.includes("/votes?")) return [];
      throw new Error(`Unexpected request ${url}`);
    },
  });

  const section = widget.createWidget("project-1", { id: "story-1" });
  await new Promise(resolve => setImmediate(resolve));
  const finalize = section.children[4];
  const result = section.children[5];
  assert.equal(finalize.hidden, false);
  assert.equal(result.hidden, true);

  await finalize.listeners.click();
  assert.equal(finalize.hidden, true);
  assert.equal(result.hidden, false);
  assert.equal(result.children[0].textContent, "Ronda 1");
  assert.match(result.children[1].textContent, /Hay diferencias/);
  assert.deepEqual(result.children[2].children.map(item => item.textContent), [
    "Ana: 3 Story Points",
    "Bea: 8 Story Points",
  ]);
});

test("BDD-US-14-01/02: previous results stay visible when a new round starts", async () => {
  const members = [{ id: "member-1", name: "Ana" }, { id: "member-2", name: "Bea" }];
  const previousResult = {
    round: { id: "round-1", number: 1, status: "finalized", participants: members.map(item => item.id) },
    votes: [
      { member_id: "member-1", story_points: 3 },
      { member_id: "member-2", story_points: 8 },
    ],
    has_differences: true,
  };
  let started = false;
  const context = vm.createContext({
    document: { createElement: tag => new Element(tag) },
    encodeURIComponent,
    URL,
  });
  vm.runInContext(fs.readFileSync(path.join(__dirname, "../../web/static/js/planning-poker.js"), "utf8"), context);
  const widget = context.createPlanningPoker({
    async requestJSON(url, options = {}) {
      if (url.endsWith("/rounds") && options.method === "POST") {
        started = true;
        return { id: "round-2", number: 2, status: "open", participants: members.map(item => item.id) };
      }
      if (url.endsWith("/rounds")) {
        return started
          ? [previousResult.round, { id: "round-2", number: 2, status: "open", participants: members.map(item => item.id) }]
          : [{ ...previousResult.round, participants: members.map(item => item.id) }];
      }
      if (url.endsWith("/members")) return members;
      if (url.includes("/round-1/result?")) return previousResult;
      if (url.includes("/round-2/votes?")) return [];
      throw new Error(`Unexpected request ${url}`);
    },
  });

  const section = widget.createWidget("project-1", { id: "story-1" });
  await new Promise(resolve => setImmediate(resolve));
  const [, message, start, form, , result] = section.children;
  assert.equal(start.hidden, false);
  assert.equal(start.textContent, "Iniciar nueva ronda");
  assert.equal(result.hidden, false);
  assert.equal(result.children[0].textContent, "Ronda 1");
  assert.match(result.children[1].textContent, /Hay diferencias/);

  await start.listeners.click();
  assert.equal(start.hidden, true);
  assert.equal(form.hidden, false);
  assert.match(message.textContent, /ocultos/);
  assert.equal(result.hidden, false);
  const previousRoundCard = result.children[3].children[0];
  assert.equal(previousRoundCard.children[0].textContent, "Ronda 1");
  assert.deepEqual(previousRoundCard.children[2].children.map(item => item.textContent), [
    "Ana: 3 Story Points",
    "Bea: 8 Story Points",
  ]);
});
