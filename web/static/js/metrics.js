function createMetricsSection({ section, message, projectMetrics, sprintSelect, sprintMetrics, requestJSON }) {
  let projectID = "";
  let generation = 0;
  let selection = 0;

  function clear() {
    projectMetrics.replaceChildren();
    sprintMetrics.replaceChildren();
    sprintSelect.replaceChildren();
    sprintSelect.disabled = true;
  }

  function showValues(target, values) {
    target.replaceChildren();
    const rows = [
      ["Story Points planificados", values.planned_story_points],
      ["Story Points completados", values.completed_story_points],
      ["Velocidad del equipo", values.team_velocity],
      ["Historias completadas (%)", values.completion_percentage],
    ];
    for (const [label, value] of rows) {
      const term = document.createElement("dt");
      term.textContent = label;
      const detail = document.createElement("dd");
      detail.textContent = value === null || value === undefined
        ? "No disponible"
        : label.endsWith("(%)")
          ? `${new Intl.NumberFormat("es-AR", { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)} %`
          : new Intl.NumberFormat("es-AR", { maximumFractionDigits: 2 }).format(value);
      target.append(term, detail);
    }
  }

  async function showProject(id) {
    projectID = id;
    const request = ++generation;
    ++selection;
    clear();
    section.hidden = !id;
    message.textContent = id ? "Consultando métricas del proyecto…" : "";
    if (!id) return;
    try {
      const result = await requestJSON(`/api/projects/${encodeURIComponent(id)}/metrics/story-points`);
      if (request !== generation) return;
      showValues(projectMetrics, result.metrics);
      const placeholder = document.createElement("option");
      placeholder.value = "";
      placeholder.textContent = result.sprints.length ? "Seleccioná un Sprint" : "No hay Sprints";
      sprintSelect.append(placeholder);
      for (const sprint of result.sprints) {
        const option = document.createElement("option");
        option.value = String(sprint.id);
        option.textContent = `Sprint ${sprint.id} — ${sprint.goal}${sprint.closed ? " (cerrado)" : " (abierto)"}`;
        sprintSelect.append(option);
      }
      sprintSelect.disabled = result.sprints.length === 0;
      message.textContent = "Métricas calculadas con la planificación y finalizaciones registradas.";
    } catch (error) {
      if (request === generation) message.textContent = error.message;
    }
  }

  sprintSelect.addEventListener("change", async () => {
    const request = ++selection;
    const context = generation;
    const id = sprintSelect.value;
    sprintMetrics.replaceChildren();
    if (!projectID || !id) return;
    message.textContent = "Consultando métricas del Sprint…";
    try {
      const result = await requestJSON(`/api/projects/${encodeURIComponent(projectID)}/sprints/${encodeURIComponent(id)}/metrics/story-points`);
      if (request !== selection || context !== generation) return;
      showValues(sprintMetrics, result);
      message.textContent = "Métricas del Sprint actualizadas.";
    } catch (error) {
      if (request === selection && context === generation) message.textContent = error.message;
    }
  });

  return { showProject };
}
