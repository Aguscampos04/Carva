// US-18: the backend owns defect validation; this widget scopes registration to one story.
function createDefects({ requestJSON }) {
  let generation = 0;
  let fieldNumber = 0;
  const statusValues = ["Abierto", "En progreso", "Resuelto", "Cerrado"];
  const severityValues = ["Baja", "Media", "Alta", "Crítica"];
  const element = (tag, text) => { const node = document.createElement(tag); if (text !== undefined) node.textContent = text; return node; };

  function createWidget(project, story) {
    const version = generation;
    const current = () => version === generation;
    const base = `/api/projects/${encodeURIComponent(project.id)}/backlog/${encodeURIComponent(story.id)}/defects`;
    const host = element("section"); host.className = "defect-widget";
    const open = element("button", "Registrar y consultar defectos"); open.type = "button";
    const message = element("p"); message.setAttribute("aria-live", "polite");
    const content = element("div");
    host.append(open, message, content);

    function addField(form, name, title, tag = "input") {
      const label = element("label", title), input = element(tag);
      input.name = name; input.id = `defect-${++fieldNumber}`; label.htmlFor = input.id;
      form.append(label, input);
      return input;
    }
    function fillSelect(select, values) {
      for (const value of values) {
        const option = element("option", value); option.value = value; select.append(option);
      }
    }
    const form = element("form");
    const description = addField(form, "description", "Descripción", "textarea"); description.required = true;
    const severity = addField(form, "severity", "Severidad", "select"); severity.required = true; fillSelect(severity, severityValues);
    const status = addField(form, "status", "Estado", "select"); status.required = true; fillSelect(status, statusValues);
    const detection = addField(form, "detection_sprint_id", "Sprint de detección (ID, opcional)");
    detection.type = "number"; detection.min = "1"; detection.step = "1";
    const resolution = addField(form, "resolution_sprint_id", "Sprint de resolución (ID, opcional)");
    resolution.type = "number"; resolution.min = "1"; resolution.step = "1";
    const submit = element("button", "Registrar defecto"); submit.type = "submit"; form.append(submit);
    const list = element("ul");

    function render(items) {
      list.replaceChildren();
      for (const defect of items) {
        const row = element("li");
        row.append(element("p", defect.description));
        const details = element("dl");
        for (const [label, value] of [
          ["Severidad", defect.severity], ["Estado", defect.status],
          ["Sprint de detección", defect.detection_sprint_id ?? "Sin Sprint"],
          ["Sprint de resolución", defect.resolution_sprint_id ?? "Sin Sprint"],
        ]) {
          details.append(element("dt", label), element("dd", value));
        }
        row.append(details); list.append(row);
      }
      if (items.length === 0) list.append(element("li", "No hay defectos registrados para esta historia."));
    }

    async function load() {
      if (!current()) return;
      open.disabled = true; message.textContent = "Consultando defectos…";
      try {
        const items = await requestJSON(base);
        if (!current()) return;
        render(items);
        content.replaceChildren(form, list);
        message.textContent = items.length ? `${items.length} defectos registrados.` : "Todavía no hay defectos para esta historia.";
      } catch (error) {
        if (current()) message.textContent = error.message;
      } finally {
        if (current()) open.disabled = false;
      }
    }
    open.addEventListener("click", load);
    form.addEventListener("submit", async event => {
      event.preventDefault();
      if (!current() || submit.disabled) return;
      submit.disabled = true; message.textContent = "Registrando defecto…";
      const sprintID = field => field.value === "" ? null : Number(field.value);
      const body = {
        description: description.value, severity: severity.value, status: status.value,
        detection_sprint_id: sprintID(detection), resolution_sprint_id: sprintID(resolution),
      };
      try {
        await requestJSON(base, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
        if (!current()) return;
        form.reset();
        await load();
      } catch (error) {
        if (current()) message.textContent = error.message;
      } finally {
        if (current()) submit.disabled = false;
      }
    });
    return host;
  }

  return { createWidget, invalidate() { generation++; } };
}
