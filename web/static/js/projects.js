// US-02: the server owns dates and lifecycle rules; this controller owns the form.
function createProjectEditor({ form, nameInput, startInput, endOutput, finalizeButton, saveButton, message, requestJSON, confirmFinalization, onSaved, onMutationStart = () => {}, onBusyChange = () => {} }) {
  let project = null;
  let generation = 0;
  let submitting = false;

  function setSubmitting(value) {
    submitting = value;
    saveButton.disabled = value;
    finalizeButton.disabled = value;
    nameInput.disabled = value;
    startInput.disabled = value;
    onBusyChange(value);
  }

  async function mutate(suffix, method, body) {
    if (!project || submitting) return;
    const context = generation;
    const id = project.id;
    onMutationStart();
    setSubmitting(true);
    message.textContent = "Guardando proyecto…";
    try {
      const saved = await requestJSON(`/api/projects/${encodeURIComponent(id)}${suffix}`, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (context !== generation) return;
      onSaved(saved);
    } catch (error) {
      if (context === generation) message.textContent = `${error.message} Consultá de nuevo antes de reintentar si se perdió la respuesta.`;
    } finally {
      if (context === generation) setSubmitting(false);
    }
  }

  form.addEventListener("submit", event => {
    event.preventDefault();
    return mutate("", "PATCH", { name: nameInput.value, start_date: startInput.value || null });
  });
  finalizeButton.addEventListener("click", () => {
    if (!project || submitting || project.status === "Finalizado") return;
    const accepted = confirmFinalization(`¿Finalizar el proyecto «${project.name}»? Se usarán los datos guardados; los cambios del formulario sin guardar no se aplican. El estado y la fecha de finalización quedarán fijos. Sus Sprints se conservan.`);
    if (!accepted) return;
    return mutate("/finalize", "POST", { confirmed: true });
  });

  return {
    showProject(value) {
      project = value;
      ++generation;
      setSubmitting(false);
      form.hidden = !value;
      finalizeButton.hidden = !value || value.status === "Finalizado";
      nameInput.value = value?.name ?? "";
      startInput.value = value?.start_date ?? "";
      endOutput.textContent = value?.end_date ?? "Sin registrar";
      message.textContent = "";
    },
  };
}
