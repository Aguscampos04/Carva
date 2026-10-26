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
    const finalize = document.createElement("button");
    finalize.type = "button";
    finalize.textContent = "Finalizar ronda y revelar votos";
    finalize.hidden = true;
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
    const result = document.createElement("section");
    result.className = "planning-poker-result";
    result.hidden = true;
    const resultMessage = document.createElement("p");
    resultMessage.setAttribute("aria-live", "polite");
    const resultVotes = document.createElement("ul");
    result.append(resultMessage, resultVotes);
    section.append(title, message, start, form, finalize, result);

    let round = null;
    let saving = false;
    let finalizing = false;
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

    async function loadFinalizedResult() {
      if (!round || !memberSelect.value) return;
      const path = `${base}/rounds/${encodeURIComponent(round.id)}/result?member_id=${encodeURIComponent(memberSelect.value)}`;
      const [roundResult, allMembers] = await Promise.all([
        requestJSON(path),
        requestJSON(`/api/projects/${encodeURIComponent(projectID)}/members`),
      ]);
      const memberNames = new Map(allMembers.map(member => [member.id, member.name]));
      resultVotes.replaceChildren();
      for (const vote of roundResult.votes) {
        const item = document.createElement("li");
        const name = memberNames.get(vote.member_id) || vote.member_id;
        item.textContent = `${name}: ${vote.story_points} Story Points`;
        resultVotes.append(item);
      }
      resultMessage.textContent = roundResult.has_differences
        ? "Hay diferencias entre las estimaciones."
        : "No hay diferencias entre las estimaciones.";
      result.hidden = false;
    }

    async function load() {
      start.hidden = true;
      form.hidden = true;
      finalize.hidden = true;
      result.hidden = true;
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
          const participants = await requestJSON(`/api/projects/${encodeURIComponent(projectID)}/members`);
          memberSelect.replaceChildren();
          for (const member of participants.filter(item => round.participants.includes(item.id))) {
            const option = document.createElement("option");
            option.value = member.id;
            option.textContent = `${member.name} · ${member.id}`;
            memberSelect.append(option);
          }
          if (!memberSelect.value && memberSelect.children.length > 0) memberSelect.value = memberSelect.children[0].value;
          message.textContent = "La ronda está finalizada. Los votos fueron revelados.";
          await loadFinalizedResult();
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
        finalize.hidden = participants.length === 0;
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

    finalize.addEventListener("click", async () => {
      if (finalizing || !round || !memberSelect.value) return;
      finalizing = true;
      finalize.disabled = true;
      try {
        const path = `${base}/rounds/${encodeURIComponent(round.id)}/finalize?member_id=${encodeURIComponent(memberSelect.value)}`;
        await requestJSON(path, { method: "POST" });
        await load();
      } catch (error) {
        message.textContent = error.message;
      } finally {
        finalizing = false;
        finalize.disabled = false;
      }
    });

    memberSelect.addEventListener("change", () => {
      if (round?.status === "finalized") {
        return loadFinalizedResult().catch(error => { message.textContent = error.message; });
      } else {
        return loadOwnVote();
      }
    });
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
