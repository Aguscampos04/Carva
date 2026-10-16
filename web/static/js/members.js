// US-03: the backend owns validation and identity. This controller owns UI state.
function createMembersSection({ section, form, nameInput, message, list, refresh, submit, requestJSON }) {
  let projectID = "";
  let generation = 0;
  let listRequest = 0;
  let submitting = false;

  function membersPath(id) {
    return `/api/projects/${encodeURIComponent(id)}/members`;
  }

  function renderMembers(items) {
    list.replaceChildren();
    for (const member of items) {
      const row = document.createElement("li");
      const name = document.createElement("span");
      name.className = "member-name";
      name.textContent = member.name;
      const id = document.createElement("small");
      id.className = "member-id";
      id.textContent = `UUID: ${member.id}`;
      row.append(name, id);
      list.append(row);
    }
  }

  async function loadMembers(prefix = "") {
    if (!projectID) return;
    const context = generation;
    const request = ++listRequest;
    message.textContent = `${prefix}Consultando miembros…`;
    list.replaceChildren();
    try {
      const items = await requestJSON(membersPath(projectID));
      if (context !== generation || request !== listRequest) return;
      renderMembers(items);
      message.textContent = prefix + (items.length === 0 ? "Todavía no hay miembros." : `${items.length} miembros registrados.`);
    } catch (error) {
      if (context === generation && request === listRequest) message.textContent = prefix + error.message;
    }
  }

  refresh.addEventListener("click", () => loadMembers());
  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    if (!projectID || submitting) return;
    const context = generation;
    const target = projectID;
    submitting = true;
    submit.disabled = true;
    // Invalidate any pending list so it cannot replace the saving/error message.
    ++listRequest;
    message.textContent = "Registrando miembro…";
    try {
      await requestJSON(membersPath(target), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: nameInput.value }),
      });
      if (context !== generation) return;
      form.reset();
      await loadMembers("Miembro registrado. ");
    } catch (error) {
      if (context === generation) message.textContent = error.message;
    } finally {
      if (context === generation) {
        submitting = false;
        submit.disabled = false;
      }
    }
  });

  return {
    async showProject(id) {
      projectID = id;
      ++generation;
      ++listRequest;
      submitting = false;
      submit.disabled = false;
      form.reset();
      list.replaceChildren();
      message.textContent = "";
      section.hidden = !id;
      if (id) await loadMembers();
    },
  };
}
