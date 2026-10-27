// US-18: the backend owns defect validation; this widget scopes registration to one story.
function createDefects({ requestJSON, sprintOptions }) {
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
    let availableSprints = [];
    function sprintField(form, name, title, value = null) {
      const input = addField(form, name, title, sprintOptions ? 'select' : 'input');
      if (!sprintOptions) { input.type = 'number'; input.min = '1'; input.step = '1'; input.value = value == null ? '' : String(value); }
      else fillSprintOptions(input, value);
      return input;
    }
    function fillSprintOptions(input, value) {
      input.replaceChildren();
      const none = element('option', 'Sin Sprint'); none.value = ''; input.append(none);
      for (const sprint of availableSprints) {
        const option = element('option', `Sprint ${sprint.id} — ${sprint.goal}${sprint.closed ? ' (cerrado)' : ''}`);
        option.value = String(sprint.id); input.append(option);
      }
      if (value != null && !availableSprints.some(s => String(s.id) === String(value))) {
        const existing = element('option', `Sprint ${value} — asociación actual`); existing.value = String(value); input.append(existing);
      }
      input.value = value == null ? '' : String(value);
    }
    const form = element("form");
    const description = addField(form, "description", "Descripción", "textarea"); description.required = true;
    const severity = addField(form, "severity", "Severidad", "select"); severity.required = true; fillSelect(severity, severityValues);
    const status = addField(form, "status", "Estado", "select"); status.required = true; fillSelect(status, statusValues);
    const detection = sprintField(form, "detection_sprint_id", "Sprint de detección (opcional)");
    const resolution = sprintField(form, "resolution_sprint_id", "Sprint de resolución (opcional)");
    const submit = element("button", "Registrar defecto"); submit.type = "submit"; form.append(submit);
    const list = element("ul");
    let records = [];
    let editing = false;

    function editDefect(defect, row) {
      if (!current() || editing || submit.disabled) return;
      editing = true; open.disabled = true; submit.disabled = true;
      const editor = element("form");
      const state = addField(editor, "status", "Estado", "select");
      fillSelect(state, statusValues); state.value = defect.status;
      const detected = sprintField(editor, "detection_sprint_id", "Sprint de detección (opcional)", defect.detection_sprint_id);
      const resolved = sprintField(editor, "resolution_sprint_id", "Sprint de resolución (opcional)", defect.resolution_sprint_id);
      const save = element("button", "Guardar cambios"); save.type = "submit";
      const cancel = element("button", "Cancelar"); cancel.type = "button";
      editor.append(save, cancel); row.append(editor);
      const finish = () => {
        editing = false; open.disabled = false; submit.disabled = false; render(records);
      };
      cancel.addEventListener("click", () => { if (current() && !save.disabled) { finish(); message.textContent = "Edición cancelada."; } });
      editor.addEventListener("submit", async event => {
        event.preventDefault();
        if (!current() || save.disabled) return;
        const body = {};
        if (state.value !== defect.status) body.status = state.value;
        for (const [name, field] of [["detection_sprint_id", detected], ["resolution_sprint_id", resolved]]) {
          const value = field.value === "" ? null : Number(field.value);
          if (value !== (defect[name] ?? null)) body[name] = value;
        }
        if (Object.keys(body).length === 0) { message.textContent = "No hay cambios para guardar."; return; }
        save.disabled = true; cancel.disabled = true; message.textContent = "Actualizando defecto…";
        try {
          const updated = await requestJSON(`${base}/${encodeURIComponent(defect.id)}`, {
            method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body),
          });
          if (!current()) return;
          records = records.map(item => item.id === updated.id ? updated : item);
          finish(); message.textContent = "Defecto actualizado.";
        } catch (error) {
          if (current()) message.textContent = error.message;
        } finally {
          if (current()) { save.disabled = false; cancel.disabled = false; }
        }
      });
    }

    function render(items) {
      records = items;
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
        const edit = element("button", "Editar defecto"); edit.type = "button";
        edit.addEventListener("click", () => editDefect(defect, row));
        row.append(details, edit); list.append(row);
      }
      if (items.length === 0) list.append(element("li", "No hay defectos registrados para esta historia."));
    }

    async function load() {
      if (!current() || editing || open.disabled) return;
      open.disabled = true; message.textContent = "Consultando defectos…";
      try {
        const [items, sprintList] = await Promise.all([requestJSON(base), sprintOptions ? sprintOptions() : Promise.resolve([])]);
        if (!current()) return;
        availableSprints = sprintList;
        if (sprintOptions) { fillSprintOptions(detection, detection.value || null); fillSprintOptions(resolution, resolution.value || null); }
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
