const form = document.querySelector("#create-project-form");
const nameInput = document.querySelector("#project-name");
const message = document.querySelector("#message");
const details = document.querySelector("#project-details");
const refreshButton = document.querySelector("#refresh-project");
const backlogSection = document.querySelector("#backlog-section");
const backlogForm = document.querySelector("#create-backlog-form");
const backlogMessage = document.querySelector("#backlog-message");
const backlogItems = document.querySelector("#backlog-items");
let currentProjectID = "";
let projectRequest = 0;
let backlogRequest = 0;
const sprintHistory = createSprintHistory({
  form: document.querySelector("#sprint-history-form"),
  projectInput: document.querySelector("#history-project-id"),
  select: document.querySelector("#closed-sprint-select"),
  message: document.querySelector("#sprint-history-message"),
  goal: document.querySelector("#historical-sprint-goal"),
  details: document.querySelector("#historical-sprint-items"),
  requestJSON,
});
const projectEditor = createProjectEditor({
  form: document.querySelector("#edit-project-form"),
  nameInput: document.querySelector("#edit-project-name"),
  startInput: document.querySelector("#edit-project-start"),
  endOutput: document.querySelector("#project-end-result"),
  finalizeButton: document.querySelector("#finalize-project"),
  saveButton: document.querySelector("#save-project"),
  message: document.querySelector("#project-edit-message"),
  requestJSON,
  confirmFinalization: text => window.confirm(text),
  onSaved: showProject,
  onMutationStart: () => { ++projectRequest; },
  onBusyChange: value => { refreshButton.disabled = value; },
});
const membersUI = createMembersSection({
  section: document.querySelector("#members-section"),
  form: document.querySelector("#create-member-form"),
  nameInput: document.querySelector("#member-name"),
  message: document.querySelector("#members-message"),
  list: document.querySelector("#members-list"),
  refresh: document.querySelector("#refresh-members"),
  submit: document.querySelector("#register-member"),
  requestJSON,
});
const effortUI = createEffort({requestJSON,confirmDelete: text => window.confirm(text),onChanged: () => loadBacklog()});
const defectsUI = createDefects({requestJSON});
const planningPokerUI = createPlanningPoker({ requestJSON, onEstimateRegistered: () => loadBacklog() });

function showProject(project) {
  ++projectRequest;
  effortUI.invalidate();
  defectsUI.invalidate();
  currentProjectID = project.id;
  document.querySelector("#project-name-result").textContent = project.name;
  document.querySelector("#project-id-result").textContent = project.id;
  document.querySelector("#project-status-result").textContent = project.status;
  document.querySelector("#project-start-result").textContent = project.start_date ?? "Sin registrar";
  projectEditor.showProject(project);
  details.hidden = false;
  refreshButton.hidden = false;
  backlogSection.hidden = false;
  message.textContent = "Proyecto disponible.";
  membersUI.showProject(project.id);
  sprintHistory.showProject(project.id);
  loadBacklog();
}

async function requestJSON(url, options = {}) {
  const response = await fetch(url, options);
  const result = response.status === 204 ? null : parseHoursJSON(await response.text());
  if (!response.ok) throw new Error(result.error || "No se pudo completar la operación.");
  return result;
}

function backlogPath(projectID = currentProjectID) {
  return `/api/projects/${encodeURIComponent(projectID)}/backlog`;
}

async function loadBacklog() {
  if (!currentProjectID) return;
  const selectedProject=currentProjectID;
  const request=++backlogRequest;
  effortUI.invalidate();
  defectsUI.invalidate();
  backlogMessage.textContent = "Consultando Product Backlog…";
  backlogItems.replaceChildren();
  try {
    const items = await requestJSON(backlogPath());
    if (currentProjectID !== selectedProject || request !== backlogRequest) return;
    renderBacklog(items);
    backlogMessage.textContent = items.length === 0 ? "Todavía no hay User Stories." : `${items.length} User Stories.`;
  } catch (error) {
    if (currentProjectID === selectedProject && request === backlogRequest) backlogMessage.textContent = error.message;
  }
}

function renderBacklog(items) {
  backlogItems.replaceChildren();
  for (const item of items) backlogItems.append(createBacklogCard(item));
}

function createBacklogCard(item) {
  const card = document.createElement("article");
  card.className = "backlog-card";
  const title = document.createElement("h3");
  title.textContent = item.title;
  const description = document.createElement("p");
  description.textContent = item.description;
  const detailsList = document.createElement("dl");
  addDetail(detailsList, "Prioridad", item.priority);
  addDetail(detailsList, "Estado", item.status);
  if (item.story_points !== null) {
    addDetail(detailsList, "Story Points", item.story_points);
  } else if (item.provisional_story_points !== null) {
    addDetail(detailsList, "Story Points provisionales", item.provisional_story_points);
  } else {
    addDetail(detailsList, "Story Points", "Sin estimar");
  }
  addDetail(detailsList, "Horas estimadas", item.estimated_hours ?? "Pendiente de regularización");
  addDetail(detailsList, "Horas reales", item.actual_hours);
  addDetail(detailsList, "Identificador", item.id);

  const criteriaTitle = document.createElement("h4");
  criteriaTitle.textContent = "Criterios de aceptación";
  const criteria = document.createElement("ul");
  for (const criterion of item.acceptance_criteria) {
    const row = document.createElement("li");
    row.textContent = criterion;
    criteria.append(row);
  }

  const editButton = document.createElement("button");
  editButton.type = "button";
  editButton.textContent = "Editar";
  editButton.addEventListener("click", () => {
    card.replaceChildren(title, description, detailsList, criteriaTitle, criteria, createEditForm(item));
  });
  card.append(title, description, detailsList, criteriaTitle, criteria, editButton, planningPokerUI.createWidget(currentProjectID, item), effortUI.createWidget({id:currentProjectID}, item), defectsUI.createWidget({id:currentProjectID}, item));
  return card;
}

function addDetail(list, label, value) {
  const term = document.createElement("dt");
  term.textContent = label;
  const detail = document.createElement("dd");
  detail.textContent = String(value);
  list.append(term, detail);
}

function createEditForm(item) {
  const editForm = document.createElement("form");
  editForm.className = "backlog-edit-form";
  const title = addFormField(editForm, item, "title", "Título", "input", item.title);
  title.required = true;
  const description = addFormField(editForm, item, "description", "Descripción", "textarea", item.description);
  description.required = true;
  const priority = addFormField(editForm, item, "priority", "Prioridad", "select", item.priority);
  for (const value of ["Alta", "Media", "Baja"]) {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = value;
    priority.append(option);
  }
  const estimated = addFormField(editForm, item, "estimated_hours", "Horas estimadas", "input", item.estimated_hours ?? "");
  estimated.type="number"; estimated.min="0.01"; estimated.step="0.01"; estimated.required=true;
  const points = addFormField(editForm, item, "provisional_story_points", "Story Points provisionales", "input", item.provisional_story_points ?? "");
  points.type = "number";
  points.min = "0";
  points.step = "1";
  points.readOnly = item.story_points !== null;
  const criteria = addFormField(editForm, item, "acceptance_criteria", "Criterios de aceptación (uno por línea)", "textarea", item.acceptance_criteria.join("\n"));
  criteria.required = true;
  const error = document.createElement("p");
  error.className = "form-error";
  error.setAttribute("aria-live", "polite");
  const save = document.createElement("button");
  save.type = "submit";
  save.textContent = "Guardar cambios";
  const cancel = document.createElement("button");
  cancel.type = "button";
  cancel.className = "secondary-button";
  cancel.textContent = "Cancelar";
  cancel.addEventListener("click", () => loadBacklog());
  editForm.append(error, save, cancel);
  editForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    try {
      await requestJSON(`${backlogPath()}/${encodeURIComponent(item.id)}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: exactHoursJSON(readBacklogFields(editForm),["estimated_hours"]),
      });
      await loadBacklog();
    } catch (cause) {
      error.textContent = cause.message;
    }
  });
  return editForm;
}

function addFormField(formElement, item, name, labelText, tagName, value) {
  const id = `edit-${item.id}-${name}`;
  const label = document.createElement("label");
  label.htmlFor = id;
  label.textContent = labelText;
  const field = document.createElement(tagName);
  field.id = id;
  field.name = name;
  field.value = value;
  formElement.append(label, field);
  return field;
}

function readBacklogFields(formElement) {
  const data = new FormData(formElement);
  const pointText = data.get("provisional_story_points");
  return {
    estimated_hours: data.get("estimated_hours"),
    title: data.get("title"),
    description: data.get("description"),
    priority: data.get("priority"),
    provisional_story_points: pointText === "" ? null : Number(pointText),
    acceptance_criteria: String(data.get("acceptance_criteria"))
      .split(/\r?\n/)
      .filter((criterion) => criterion.trim() !== ""),
  };
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  message.textContent = "Guardando proyecto…";
  details.hidden = true;
  refreshButton.hidden = true;
  currentProjectID = "";
  const request = ++projectRequest;
  projectEditor.showProject(null);
  backlogSection.hidden = true;
  membersUI.showProject("");
  sprintHistory.showProject("");

  try {
    const result = await requestJSON("/api/projects", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: nameInput.value }),
    });
    if (request === projectRequest) showProject(result);
  } catch (error) {
    if (request === projectRequest) message.textContent = error.message;
  }
});

refreshButton.addEventListener("click", async () => {
  if (refreshButton.disabled || !currentProjectID) return;
  const id = currentProjectID;
  const request = ++projectRequest;
  message.textContent = "Consultando proyecto…";
  try {
    const project = await requestJSON(`/api/projects/${encodeURIComponent(id)}`);
    if (request === projectRequest && id === currentProjectID) showProject(project);
  } catch (error) {
    if (request === projectRequest && id === currentProjectID) message.textContent = error.message;
  }
});

backlogForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const error = document.querySelector("#backlog-message");
  error.textContent = "Guardando User Story…";
  try {
    await requestJSON(backlogPath(), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: exactHoursJSON(readBacklogFields(backlogForm),["estimated_hours"]),
    });
    backlogForm.reset();
    document.querySelector("#backlog-priority-input").value = "Alta";
    await loadBacklog();
  } catch (cause) {
    error.textContent = cause.message;
  }
});
