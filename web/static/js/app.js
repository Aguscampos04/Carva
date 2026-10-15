const form = document.querySelector("#create-project-form");
const nameInput = document.querySelector("#project-name");
const message = document.querySelector("#message");
const details = document.querySelector("#project-details");
const refreshButton = document.querySelector("#refresh-project");
let currentProjectID = "";

function showProject(project) {
  currentProjectID = project.id;
  document.querySelector("#project-name-result").textContent = project.name;
  document.querySelector("#project-id-result").textContent = project.id;
  document.querySelector("#project-status-result").textContent = project.status;
  details.hidden = false;
  refreshButton.hidden = false;
  message.textContent = "Proyecto guardado y consultado.";
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  message.textContent = "Guardando proyecto…";
  details.hidden = true;
  refreshButton.hidden = true;

  try {
    const response = await fetch("/api/projects", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: nameInput.value }),
    });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || "No se pudo crear el proyecto.");
    showProject(result);
  } catch (error) {
    message.textContent = error.message;
  }
});

refreshButton.addEventListener("click", async () => {
  message.textContent = "Consultando proyecto…";
  try {
    const response = await fetch(`/api/projects/${encodeURIComponent(currentProjectID)}`);
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || "No se pudo consultar el proyecto.");
    showProject(result);
  } catch (error) {
    message.textContent = error.message;
  }
});
