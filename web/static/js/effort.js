// Tokenize whole JSON strings so text inside activities can never become a field.
// Keep only hour numbers as decimal text, before JavaScript's binary conversion.
function parseHoursJSON(text) {
  return JSON.parse(text.replace(/"(?:estimated_hours|actual_hours|hours)"\s*:\s*(\d+(?:\.\d+)?)|"(?:\\.|[^"\\])*"/g,
    (token, decimal) => decimal === undefined ? token : token.slice(0, token.length - decimal.length) + JSON.stringify(decimal)));
}

// Preserve decimal text until JSON serialization; Go owns all business validations.
function exactHoursJSON(value, hourFields) {
  return `{${Object.entries(value).map(([key, data]) => {
    if (!hourFields.includes(key)) return `${JSON.stringify(key)}:${JSON.stringify(data)}`;
    const text = String(data).replace(",", ".");
    if (!/^(0|[1-9][0-9]*)(\.[0-9]{1,2})?$/.test(text)) throw new Error("Las horas requieren un número con hasta dos decimales.");
    return `${JSON.stringify(key)}:${text}`;
  }).join(",")}}`;
}

function createEffort({ requestJSON, confirmDelete, onChanged, entryOnly = false, confirmDiscard = () => true }) {
  let epoch = 0;
  let sequence = 0;
  const element = (tag, text) => { const node = document.createElement(tag); if (text !== undefined) node.textContent = text; return node; };
  function createWidget(project, story, preferredSprintID = null) {
    const version = epoch;
    const host = element("section");
    const open = element("button", entryOnly ? "Cargar horas trabajadas" : "Registrar y consultar esfuerzo"); open.type = "button";
    const message = element("p"); message.setAttribute("aria-live", "polite");
    const content = element("div"); host.append(open, message, content);
    const base = `/api/projects/${encodeURIComponent(project.id)}`;
    const storyPath = `${base}/backlog/${encodeURIComponent(story.id)}`;
    let loadID = 0, busy = false;
    let currentProject, members, stories, context;
    const current = () => epoch === version;
    function field(form, name, title, tag = "input") {
      const label = element("label", title), input = element(tag);
      input.name = name; input.id = `effort-${++sequence}`; label.htmlFor = input.id;
      form.append(label, input); return input;
    }
    function option(select, value, title) { const node = element("option", title); node.value = String(value); select.append(node); }
    function renderForm(entry = null) {
      // The form always lives in its own container, separate from the active list.
      const form = element("form");
      const f = {};
      f.member_id = field(form, "member_id", entryOnly ? "Miembro" : "Miembro (nombre y UUID)", "select");
      option(f.member_id, "", "Seleccionar miembro");
      for (const member of members) option(f.member_id, member.id, entryOnly && members.filter(m => m.name === member.name).length === 1 ? member.name : `${member.name} — ${member.id}`);
      f.member_id.value = entry?.member_id ?? ""; f.member_id.required = true;
      f.backlog_item_id = field(form, "backlog_item_id", "User Story", "select");
      for (const item of stories) option(f.backlog_item_id, item.id, `${item.title} — ${item.id}`);
      f.backlog_item_id.value = entry?.backlog_item_id ?? story.id; f.backlog_item_id.disabled = !entry;
      if (entryOnly) {
        f.backlog_item_id.hidden = true;
        form.children[form.children.length - 2].hidden = true;
      }
      f.work_date = field(form, "work_date", "Fecha de trabajo"); f.work_date.type = "date"; f.work_date.required = true;
      f.work_date.value = entry?.work_date ?? context.today;
      f.activity = field(form, "activity", "Actividad", "textarea"); f.activity.value = entry?.activity ?? ""; f.activity.required = true;
      f.hours = field(form, "hours", "Horas trabajadas"); f.hours.type = "number"; f.hours.min = "0.01"; f.hours.step = "0.01"; f.hours.required = true; f.hours.value = String(entry?.hours ?? "");
      f.sprint_id = field(form, "sprint_id", "Sprint (opcional)", "select");
      const save = element("button", entryOnly ? "Guardar horas" : "Guardar esfuerzo"); save.type = "submit";
      const cancel = element("button", entryOnly ? "Cerrar" : "Cancelar edición"); cancel.type = "button";
      let sprintRequest = 0, waiting = false;
      const blocked = currentProject.status === "Finalizado" && currentProject.end_date == null;
      const updateSave = () => { save.disabled = busy || waiting || blocked || members.length === 0 || f.sprint_id.value === "__resolve__"; cancel.disabled = busy; };
      function sprints(options, previous) {
        f.sprint_id.replaceChildren(); option(f.sprint_id, "", "Sin Sprint");
        for (const sprint of options) option(f.sprint_id, sprint.id, `Sprint ${sprint.id}${sprint.closed ? " (cerrado)" : ""}`);
        const valid = previous == null || options.some(s => String(s.id) === String(previous));
        if (!valid) option(f.sprint_id, "__resolve__", "Seleccionar otro Sprint o quitar asociación");
        f.sprint_id.value = valid ? String(previous ?? "") : "__resolve__"; updateSave();
      }
      const initialSprint = entry ? entry.sprint_id : (context.sprints.some(s => String(s.id) === String(preferredSprintID)) ? preferredSprintID : null);
      sprints(context.sprints, initialSprint);
      if (entryOnly && initialSprint != null) {
        f.sprint_id.hidden = true;
        // The label is paired with this field, while the submitted Sprint stays explicit.
        for (const child of form.children) if (child.htmlFor === f.sprint_id.id) child.hidden = true;
      }
      f.sprint_id.addEventListener("change", updateSave);
      f.backlog_item_id.addEventListener("change", async () => {
        const request = ++sprintRequest;
        const previous = entry?.sprint_id ?? null; waiting = true; updateSave();
        try {
          const selected = await requestJSON(`${base}/backlog/${encodeURIComponent(f.backlog_item_id.value)}/effort-context`);
          if (!current() || request !== sprintRequest) return;
          waiting = false; sprints(selected.sprints, previous);
        } catch (error) { if (current() && request === sprintRequest) { message.textContent = error.message; } }
      });
      const formValues = () => JSON.stringify(Object.values(f).map(input => String(input.value ?? "")));
      const initialValues = formValues();
      cancel.addEventListener("click", () => {
        if (!entryOnly) { renderForm(); return; }
        if (!current() || busy) return;
        if (formValues() !== initialValues && !confirmDiscard("Tenés horas sin guardar. ¿Descartar los cambios y cerrar?")) return;
        ++loadID; content.replaceChildren(); open.hidden = false; message.textContent = ""; open.focus();
      });
      form.addEventListener("submit", async event => {
        event.preventDefault(); if (!current() || save.disabled || busy) return;
        busy = true; updateSave(); let saved = false;
        try {
          const body = { member_id: f.member_id.value, work_date: f.work_date.value, activity: f.activity.value, hours: f.hours.value, sprint_id: f.sprint_id.value === "" ? null : Number(f.sprint_id.value) };
          if (entry) body.backlog_item_id = f.backlog_item_id.value;
          await requestJSON(entry ? `${base}/effort/${encodeURIComponent(entry.id)}` : `${storyPath}/effort`, { method: entry ? "PUT" : "POST", headers: { "Content-Type": "application/json" }, body: exactHoursJSON(body, ["hours"]) });
          saved = true;
          if (!current()) return;
          await onChanged();
          if (current() && entryOnly) {
            content.replaceChildren(); open.hidden = false; message.textContent = "Horas trabajadas guardadas.";
          } else if (current()) { busy = false; await load(); }
        } catch (error) { if (current()) message.textContent = saved ? `Esfuerzo guardado. Recargá antes de repetir: ${error.message}` : error.message; }
        finally { busy = false; if (!saved) updateSave(); }
      });
      if (entryOnly) {
        const actions = element("div"); actions.className = "actions"; cancel.className = "secondary-button";
        actions.append(save, cancel); form.append(actions);
      } else { form.append(save); if (entry) form.append(cancel); }
      if (blocked) message.textContent = "Fecha de fin desconocida: altas y ediciones de esfuerzo bloqueadas. Podés consultar y eliminar con confirmación.";
      else if (!members.length) message.textContent = "Registrá un miembro del proyecto para cargar esfuerzo.";
      updateSave(); formContainer.replaceChildren(form);
    }
    const formContainer = element("div");
    async function load() {
      const request = ++loadID; open.disabled = true; message.textContent = "Consultando esfuerzo…";
      try {
        const [detail, records, formContext, memberList, storyList, freshProject] = await Promise.all([
          requestJSON(storyPath), requestJSON(`${storyPath}/effort`), requestJSON(`${storyPath}/effort-context`),
          requestJSON(`${base}/members`), requestJSON(`${base}/backlog`), requestJSON(base),
        ]);
        if (!current() || request !== loadID) return;
        context = formContext; members = memberList; stories = storyList; currentProject = freshProject;
        if (entryOnly) {
          message.textContent = "";
          content.replaceChildren(formContainer); renderForm(); open.hidden = true;
          return;
        }
        message.textContent = records.length ? `${records.length} registros activos.` : "Sin registros de esfuerzo.";
        const comparison = element("section"); comparison.className = "effort-comparison";
        comparison.setAttribute("aria-label", "Comparación de esfuerzo");
        comparison.append(element("h4", "Comparación de esfuerzo"));
        const values = element("dl");
        values.append(element("dt", "Horas estimadas"), element("dd", `${detail.estimated_hours ?? "Pendiente de regularización"}${detail.estimated_hours == null ? "" : " h"}`));
        values.append(element("dt", "Horas reales"), element("dd", `${detail.actual_hours} h`));
        content.replaceChildren(comparison, values);
        const list = element("ul");
        for (const entry of records) {
          const row = element("li");
          row.append(element("p", `${members.find(m => m.id === entry.member_id)?.name ?? entry.member_id} (${entry.member_id}) · ${entry.work_date} · ${entry.activity} · ${entry.hours} h · Sprint ${entry.sprint_id ?? "sin asociación"}`));
          const edit = element("button", "Editar esfuerzo"); edit.type = "button";
          edit.disabled = currentProject.status === "Finalizado" && currentProject.end_date == null;
          edit.addEventListener("click", () => renderForm(entry));
          const remove = element("button", "Eliminar esfuerzo"); remove.type = "button";
          remove.addEventListener("click", async () => {
            if (!current() || busy || !confirmDelete(`¿Eliminar el esfuerzo ${entry.id} de ${entry.member_id}, ${entry.work_date}, ${entry.hours} horas? No se guardarán borradores.`)) { remove.focus(); return; }
            busy = true; remove.disabled = true;
            try {
              await requestJSON(`${base}/effort/${encodeURIComponent(entry.id)}/delete`, { method: "POST", headers: { "Content-Type": "application/json" }, body: '{"confirmed":true}' });
              if (current()) { await onChanged(); if (current()) { busy = false; await load(); } }
            } catch (error) { if (current()) message.textContent = error.message; }
            finally { busy = false; remove.disabled = false; }
          });
          row.append(edit, remove); list.append(row);
        }
        content.append(list, formContainer); renderForm();
      } catch (error) { if (current() && request === loadID) message.textContent = error.message; }
      finally { if (current() && request === loadID) open.disabled = false; }
    }
    open.addEventListener("click", load);
    return host;
  }
  return { createWidget, invalidate: () => { ++epoch; } };
}
