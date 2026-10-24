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
  const [memberLabel, pointsLabel, submit] = form.children;
  const select = memberLabel.children[0];
  const points = pointsLabel.children[0];
  assert.equal(select.children.length, 2);
  assert.equal(select.value, "member-1");
  assert.equal(points.value, "3");
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
