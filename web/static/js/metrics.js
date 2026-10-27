// Presentation only: all metric values, including availability, come from Go.
function renderMetricCharts(target, points, hours, defects) {
  const element = (tag, text) => {
    const node = document.createElement(tag);
    if (text !== undefined) node.textContent = text;
    return node;
  };
  const svgElement = (tag, attributes) => {
    const node = document.createElementNS("http://www.w3.org/2000/svg", tag);
    for (const [key, value] of Object.entries(attributes)) node.setAttribute(key, value);
    return node;
  };
  const groups = [
    { title: "Story Points", unit: "SP", values: points, fields: [["planned_story_points", "Planificados"], ["completed_story_points", "Completados"], ["team_velocity", "Velocidad"]] },
    { title: "Avance", unit: "%", values: points, limit: 100, fields: [["completion_percentage", "Historias completadas"]] },
    { title: "Esfuerzo", unit: "h", values: hours, fields: [["estimated_hours", "Horas estimadas"], ["actual_hours", "Horas reales"]] },
    { title: "Desviación", unit: "h", values: hours, signed: true, fields: [["deviation", "Reales menos estimadas"]] },
    { title: "Defectos", unit: "defectos", values: defects, fields: [["detected_defects", "Detectados"], ["resolved_defects", "Resueltos"]] },
  ];
  target.replaceChildren();
  for (const group of groups) {
    const figure = element("figure"); figure.className = "metric-chart";
    figure.append(element("figcaption", group.title));
    const number = key => {
      const value = group.values?.[key];
      return value === null || value === undefined || !Number.isFinite(Number(value)) ? null : Number(value);
    };
    const max = group.limit ?? Math.max(1, ...group.fields.map(([key]) => Math.abs(number(key) ?? 0)));
    const min = group.signed ? -max : 0;
    const position = value => 8 + ((value - min) / (max - min)) * 384;
    const format = value => new Intl.NumberFormat("es-AR", { maximumFractionDigits: 2 }).format(value);
    figure.append(element("p", `Escala: ${format(min)} a ${format(max)} ${group.unit}`));
    for (const [key, label] of group.fields) {
      const row = element("div"); row.setAttribute("data-metric", key);
      const value = number(key);
      row.append(element("p", `${label}: ${value === null ? "No disponible" : `${format(group.values[key])} ${group.unit}`}`));
      if (value !== null) {
        const svg = svgElement("svg", { viewBox: "0 0 400 26", "aria-hidden": "true", focusable: "false" });
        const zero = position(0), end = position(value);
        svg.append(svgElement("line", { x1: zero, y1: 0, x2: zero, y2: 26, class: "metric-axis" }));
        svg.append(svgElement("rect", { x: Math.min(zero, end), y: 4, width: Math.abs(end - zero), height: 18, class: value < 0 ? "metric-bar negative" : "metric-bar" }));
        row.append(svg);
      }
      figure.append(row);
    }
    target.append(figure);
  }
}

function createMetricsSection({ section, message, projectMetrics, sprintSelect, sprintMetrics, projectHours, sprintHours, projectDefects, sprintDefects, requestJSON, dashboard, refresh }) {
  let projectID = "";
  let generation = 0;
  let selection = 0;
  let currentProject = null;
  let sprints = [];

  function clearSprint() {
    sprintMetrics.replaceChildren();
    sprintHours.replaceChildren();
    sprintDefects.replaceChildren();
    if (dashboard) {
      dashboard.sprintCharts.replaceChildren();
      dashboard.sprintTitle.textContent = "Seleccioná un Sprint para ver sus gráficos.";
    }
  }

  function clear() {
    projectMetrics.replaceChildren();
    clearSprint();
    projectHours.replaceChildren();
    projectDefects.replaceChildren();
    sprintSelect.replaceChildren();
    sprintSelect.disabled = true;
    sprints = [];
    if (dashboard) dashboard.projectCharts.replaceChildren();
  }

  function showHours(target, values) {
    const rows = [
      ["Horas estimadas", values.estimated_hours],
      ["Horas reales", values.actual_hours],
      ["Desviación", values.deviation],
    ].map(([label, value]) => [label, value === null || value === undefined
      ? "No disponible"
      : `${new Intl.NumberFormat("es-AR", { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)} h`]);
    renderMetricRows(target, rows);
  }

  function showDefectValues(target, values) {
    const rows = [
      ["Defectos detectados", values.detected_defects],
      ["Defectos resueltos", values.resolved_defects],
    ].map(([label, value]) => [label, value === null || value === undefined
      ? "No disponible"
      : new Intl.NumberFormat("es-AR", { maximumFractionDigits: 0 }).format(value)]);
    renderMetricRows(target, rows);
  }

  function showValues(target, values) {
    const rows = [
      ["Story Points planificados", values.planned_story_points],
      ["Story Points completados", values.completed_story_points],
      ["Velocidad del equipo", values.team_velocity],
      ["Historias completadas (%)", values.completion_percentage],
    ].map(([label, value]) => [label, value === null || value === undefined
      ? "No disponible"
      : label.endsWith("(%)")
        ? `${new Intl.NumberFormat("es-AR", { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)} %`
        : new Intl.NumberFormat("es-AR", { maximumFractionDigits: 2 }).format(value)]);
    renderMetricRows(target, rows);
  }

  function renderMetricRows(target, rows) {
    target.replaceChildren();
    for (const [label, value] of rows) {
      const term = document.createElement("dt");
      term.textContent = label;
      const detail = document.createElement("dd");
      detail.textContent = value;
      target.append(term, detail);
    }
  }

  async function showProject(id, project = null) {
    projectID = id;
    currentProject = project;
    const request = ++generation;
    ++selection;
    clear();
    section.hidden = !id;
    if (dashboard) {
      dashboard.section.hidden = !id;
      dashboard.identity.textContent = project ? `${project.name} — ${project.id} — Estado: ${project.status}` : "";
    }
    message.textContent = id ? "Consultando métricas del proyecto…" : "";
    if (!id) return;
    try {
      const [result, hoursResult, defectsResult] = await Promise.all([
        requestJSON(`/api/projects/${encodeURIComponent(id)}/metrics/story-points`),
        requestJSON(`/api/projects/${encodeURIComponent(id)}/metrics/hours`),
        requestJSON(`/api/projects/${encodeURIComponent(id)}/metrics/defects`),
      ]);
      if (request !== generation) return;
      showValues(projectMetrics, result.metrics);
      showHours(projectHours, hoursResult.metrics);
      showDefectValues(projectDefects, defectsResult.metrics);
      if (dashboard) renderMetricCharts(dashboard.projectCharts, result.metrics, hoursResult.metrics, defectsResult.metrics);
      sprints = result.sprints;
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
      message.textContent = "Métricas calculadas con la planificación, las finalizaciones, el esfuerzo y los defectos registrados.";
    } catch (error) {
      if (request === generation) message.textContent = error.message;
    }
  }

  if (refresh) refresh.addEventListener("click", async () => {
    const selected = sprintSelect.value;
    const project = projectID;
    const request = generation + 1;
    await showProject(projectID, currentProject);
    if (request !== generation || project !== projectID || !sprints.some(sprint => String(sprint.id) === selected)) return;
    sprintSelect.value = selected;
    await loadSprint();
  });

  async function loadSprint() {
    const request = ++selection;
    const context = generation;
    const id = sprintSelect.value;
    clearSprint();
    if (!projectID || !id) { message.textContent = "Seleccioná un Sprint para consultar sus métricas."; return; }
    message.textContent = "Consultando métricas del Sprint…";
    try {
      const [result, hoursResult, defectsResult] = await Promise.all([
        requestJSON(`/api/projects/${encodeURIComponent(projectID)}/sprints/${encodeURIComponent(id)}/metrics/story-points`),
        requestJSON(`/api/projects/${encodeURIComponent(projectID)}/sprints/${encodeURIComponent(id)}/metrics/hours`),
        requestJSON(`/api/projects/${encodeURIComponent(projectID)}/sprints/${encodeURIComponent(id)}/metrics/defects`),
      ]);
      if (request !== selection || context !== generation) return;
      showValues(sprintMetrics, result);
      showHours(sprintHours, hoursResult);
      showDefectValues(sprintDefects, defectsResult);
      if (dashboard) {
        renderMetricCharts(dashboard.sprintCharts, result, hoursResult, defectsResult);
        const selected = sprints.find(sprint => String(sprint.id) === id);
        dashboard.sprintTitle.textContent = selected ? `Sprint ${selected.id} — ${selected.goal}` : `Sprint ${id}`;
      }
      message.textContent = "Métricas del Sprint actualizadas.";
    } catch (error) {
      if (request === selection && context === generation) message.textContent = error.message;
    }
  }
  sprintSelect.addEventListener("change", loadSprint);

  return { showProject };
}
