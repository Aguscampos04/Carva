// US-24: the server supplies report content and calculations; the UI only presents it.
function createReports({select, refresh, download, message, content, requestJSON,
  requestPDF = async url => {
    const response = await fetch(url);
    if (!response.ok) {
      let detail;
      try { detail = (await response.json()).error; } catch { /* Generic error below. */ }
      throw new Error(detail || 'No se pudo generar el PDF.');
    }
    if (!response.headers.get('Content-Type')?.startsWith('application/pdf')) throw new Error('La respuesta no es un PDF válido.');
    return response.blob();
  },
  savePDF = (blob, name) => {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a'); link.href = url; link.download = name;
    document.body.append(link); link.click(); link.remove();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  },
}) {
  let generation = 0, selection = 0, projectID = '', loaded = false, downloading = false;
  const element = (tag, text) => {const node = document.createElement(tag); if (text != null) node.textContent = text; return node;};
  const path = () => `/api/projects/${encodeURIComponent(projectID)}/reports/${select.value ? `sprints/${encodeURIComponent(select.value)}` : 'project'}`;
  function invalidate() {
    ++generation; ++selection; projectID = ''; loaded = false; downloading = false;
    content.replaceChildren(); select.replaceChildren(); select.value = '';
    select.disabled = true; download.disabled = true; message.textContent = '';
  }
  function render(report) {
    content.replaceChildren();
    for (const section of report.sections) {
      const panel = element('section'); panel.className = 'panel report-section';
      panel.append(element('h2', section.title));
      for (const text of section.lines?.length ? section.lines : ['Sin registros.']) panel.append(element('p', text));
      content.append(panel);
    }
  }
  async function load(populate = false) {
    if (!projectID) return;
    const version = generation, request = ++selection;
    const current = () => version === generation && request === selection;
    content.replaceChildren(); loaded = false; download.disabled = true;
    message.textContent = 'Generando reporte…';
    try {
      const report = await requestJSON(path());
      if (!current()) return;
      if (populate) {
        select.replaceChildren();
        const project = element('option', 'Proyecto completo'); project.value = ''; select.append(project);
        for (const {sprint} of report.sprints) {
          const option = element('option', `Sprint ${sprint.id} — ${sprint.goal}${sprint.closed ? ' (cerrado)' : ' (abierto)'}`);
          option.value = String(sprint.id); select.append(option);
        }
        select.value = ''; select.disabled = false;
      }
      render(report); loaded = true; download.disabled = downloading;
      message.textContent = 'Reporte disponible. Datos consultados al generarlo.';
    } catch (error) { if (current()) message.textContent = error.message; }
  }
  async function showProject(project) {
    invalidate(); if (!project) return;
    projectID = project.id; await load(true);
  }
  select.addEventListener('change', () => load());
  refresh.addEventListener('click', () => load(select.disabled));
  download.addEventListener('click', async () => {
    if (!projectID || !loaded || downloading) return;
    const version = generation, request = selection, sprint = select.value;
    downloading = true; download.disabled = true; message.textContent = 'Preparando PDF…';
    try {
      const blob = await requestPDF(`${path()}.pdf`);
      if (version !== generation || request !== selection) return;
      savePDF(blob, sprint ? `carva-sprint-${sprint}.pdf` : 'carva-proyecto.pdf');
      message.textContent = 'PDF generado. Descarga iniciada.';
    } catch (error) {
      if (version === generation && request === selection) message.textContent = error.message;
    } finally {
      if (version === generation) { downloading = false; download.disabled = !loaded; }
    }
  });
  invalidate();
  return {showProject, invalidate};
}
