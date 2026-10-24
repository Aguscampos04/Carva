// US-12: members are selected by their project UUID; open rounds show only the
// selected member's own vote.
function createPlanningPoker({ requestJSON }) {
  function createWidget(projectID, story) {
    const section = document.createElement("section");
    section.className = "planning-poker";
    const title = document.createElement("h4");
    title.textContent = "Planning Poker";
    const message = document.createElement("p");
    message.setAttribute("aria-live", "polite");
    const start = document.createElement("button");
    start.type = "button";
    start.textContent = "Iniciar ronda de estimación";
    start.hidden = true;
    const form = document.createElement("form");
    form.hidden = true;

    const memberLabel = document.createElement("label");
    memberLabel.textContent = "Miembro";
    const memberSelect = document.createElement("select");
    memberSelect.required = true;
    memberLabel.append(memberSelect);

    const pointsLabel = document.createElement("label");
    pointsLabel.textContent = "Tu estimación en Story Points";
    const pointsInput = document.createElement("input");
    pointsInput.type = "number";
    pointsInput.min = "0";
    pointsInput.step = "1";
    pointsInput.required = true;
    pointsLabel.append(pointsInput);

    const submit = document.createElement("button");
    submit.type = "submit";
    submit.textContent = "Registrar o corregir mi voto";
    form.append(memberLabel, pointsLabel, submit);
    section.append(title, message, start, form);

    let round = null;
    let saving = false;
    const base = `/api/projects/${encodeURIComponent(projectID)}/backlog/${encodeURIComponent(story.id)}/estimation`;

    async function loadOwnVote() {
      if (!round || !memberSelect.value) return;
      const path = `${base}/rounds/${encodeURIComponent(round.id)}/votes?member_id=${encodeURIComponent(memberSelect.value)}`;
      try {
        const votes = await requestJSON(path);
        pointsInput.value = votes[0] ? String(votes[0].story_points) : "";
        message.textContent = votes.length ? "Tu voto está guardado y oculto para los demás." : "Todavía no registraste tu voto.";
      } catch (error) {
        message.textContent = error.message;
      }
    }

    async function load() {
      start.hidden = true;
      form.hidden = true;
      message.textContent = "Consultando ronda…";
      try {
        const rounds = await requestJSON(`${base}/rounds`);
        round = rounds.at(-1) ?? null;
        if (!round) {
          message.textContent = "Todavía no hay una ronda de estimación.";
          start.hidden = false;
          return;
        }
        if (round.status !== "open") {
          message.textContent = "La ronda está finalizada.";
          return;
        }
        const allMembers = await requestJSON(`/api/projects/${encodeURIComponent(projectID)}/members`);
        const participants = allMembers.filter(member => round.participants.includes(member.id));
        memberSelect.replaceChildren();
        for (const member of participants) {
          const option = document.createElement("option");
          option.value = member.id;
          option.textContent = `${member.name} · ${member.id}`;
          memberSelect.append(option);
        }
        if (!memberSelect.value && participants.length > 0) memberSelect.value = participants[0].id;
        form.hidden = participants.length === 0;
        message.textContent = participants.length ? "Los votos ajenos permanecen ocultos hasta finalizar la ronda." : "La ronda no tiene miembros participantes.";
        await loadOwnVote();
      } catch (error) {
        message.textContent = error.message;
      }
    }

    start.addEventListener("click", async () => {
      start.disabled = true;
      try {
        await requestJSON(`${base}/rounds`, { method: "POST" });
        await load();
      } catch (error) {
        message.textContent = error.message;
      } finally {
        start.disabled = false;
      }
    });

    memberSelect.addEventListener("change", loadOwnVote);
    form.addEventListener("submit", async event => {
      event.preventDefault();
      if (saving || !round || !memberSelect.value || !pointsInput.validity.valid) return;
      saving = true;
      submit.disabled = true;
      try {
        await requestJSON(`${base}/rounds/${encodeURIComponent(round.id)}/votes/${encodeURIComponent(memberSelect.value)}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ story_points: Number(pointsInput.value) }),
        });
        await loadOwnVote();
      } catch (error) {
        message.textContent = error.message;
      } finally {
        saving = false;
        submit.disabled = false;
      }
    });

    load();
    return section;
  }

  return { createWidget };
}
