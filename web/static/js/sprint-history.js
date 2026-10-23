// Read-only US-11 UI; all history is obtained from persisted server snapshots.
function createSprintHistory({ form, projectInput, select, message, details, goal, requestJSON }) {
  let projectID = "";
  let generation = 0;
  let selection = 0;

  function clearDetails() {
    goal.textContent = "";
    details.replaceChildren();
  }

  function path(id = projectID) {
    return `/api/projects/${encodeURIComponent(id)}/sprints`;
  }

  async function loadProject(id) {
    projectID = id;
    projectInput.value = id;
    const request = ++generation;
    ++selection;
    select.value = "";
    select.replaceChildren();
    select.disabled = true;
    clearDetails();
    message.textContent = id ? "Consultando Sprints cerrados…" : "Ingresá el identificador de un proyecto.";
    if (!id) return;
    try {
      const sprints = await requestJSON(`${path(id)}/closed`);
      if (request !== generation) return;
      const placeholder = document.createElement("option");
      placeholder.value = "";
      placeholder.textContent = "Seleccioná un Sprint cerrado";
      select.append(placeholder);
      for (const sprint of sprints) {
        const option = document.createElement("option");
        option.value = String(sprint.id);
        option.textContent = `Sprint ${sprint.id} — ${sprint.goal}`;
        select.append(option);
      }
      select.disabled = sprints.length === 0;
      message.textContent = sprints.length === 0 ? "No hay Sprints cerrados en este proyecto." : `${sprints.length} Sprints cerrados.`;
    } catch (error) {
      if (request === generation) message.textContent = error.message;
    }
  }

  function addText(parent, tag, text) {
    const node = document.createElement(tag);
    node.textContent = text;
    parent.append(node);
    return node;
  }

  function renderSnapshots(items) {
    if (items.length === 0) addText(details, "p", "Sin historias planificadas.");
    for (const item of items) {
      const card = document.createElement("article");
      card.className = "backlog-card";
      addText(card, "h3", item.title);
      addText(card, "p", item.description);
      addText(card, "p", `Historia: ${item.backlog_item_id}`);
      addText(card, "p", `Prioridad planificada: ${item.priority}`);
      addText(card, "p", `Estado planificado: ${item.planned_status}`);
      addText(card, "p", `Estado registrado: ${item.status}`);
      // This is only a label for the server's persisted verification flag.
      addText(card, "p", `Resultado al cierre: ${item.acceptance_criteria_verified ? "Completada" : "Incompleta"}`);
      if (item.story_points != null) addText(card, "p", `Story Points: ${item.story_points}`);
      else if (item.provisional_story_points != null) addText(card, "p", `Story Points provisionales: ${item.provisional_story_points}`);
      else addText(card, "p", "Story Points: Sin estimar");
      addText(card, "h4", "Criterios de aceptación planificados");
      const criteria = document.createElement("ul");
      for (const criterion of item.acceptance_criteria) addText(criteria, "li", criterion);
      card.append(criteria);
      details.append(card);
    }
  }

  select.addEventListener("change", async () => {
    const request = ++selection;
    const context = generation;
    const id = select.value;
    clearDetails();
    message.textContent = "Seleccioná un Sprint cerrado.";
    if (!projectID || !id) return;
    message.textContent = "Consultando historial…";
    const base = `${path()}/${encodeURIComponent(id)}`;
    try {
      const [sprint, items] = await Promise.all([requestJSON(base), requestJSON(`${base}/backlog-items`)]);
      if (request !== selection || context !== generation) return;
      goal.textContent = `Sprint ${sprint.id} — ${sprint.goal}`;
      renderSnapshots(items);
      message.textContent = "Historial registrado al cierre. Consulta de solo lectura.";
    } catch (error) {
      if (request === selection && context === generation) message.textContent = error.message;
    }
  });

  form.addEventListener("submit", event => {
    event.preventDefault();
    return loadProject(projectInput.value.trim());
  });
  return { showProject: loadProject };
}
