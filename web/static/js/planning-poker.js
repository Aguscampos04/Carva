// US-12: members are selected by their project UUID; open rounds show only the
// selected member's own vote.
function createPlanningPoker({ requestJSON, onEstimateRegistered = () => {} }) {
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
    const resultHeading = document.createElement("h5");
    result.append(resultHeading);
    const resultMessage = document.createElement("p");
    resultMessage.setAttribute("aria-live", "polite");
    const resultVotes = document.createElement("ul");
    const resultHistory = document.createElement("div");
    resultHistory.className = "planning-poker-history";
    const agreement = document.createElement("button");
    agreement.type = "button";
    agreement.textContent = "Registrar estimación acordada";
    agreement.hidden = true;
    result.append(resultMessage, resultVotes, resultHistory, agreement);
    section.append(title, message, start, form, finalize, result);

    let round = null;
    let rounds = [];
    let saving = false;
    let finalizing = false;
    let registeringAgreement = false;
    const registeredAgreementRoundIDs = new Set();
    const base = `/api/projects/${encodeURIComponent(projectID)}/backlog/${encodeURIComponent(story.id)}/estimation`;

    async function registerAgreement(targetRound, button, messageElement) {
      if (registeringAgreement || !memberSelect.value) return;
      registeringAgreement = true;
      button.disabled = true;
      try {
        const path = `${base}/rounds/${encodeURIComponent(targetRound.id)}/agreement?member_id=${encodeURIComponent(memberSelect.value)}`;
        const registered = await requestJSON(path, { method: "POST" });
        registeredAgreementRoundIDs.add(targetRound.id);
        messageElement.textContent = `No hay diferencias entre las estimaciones. Se registraron ${registered.story_points} Story Points acordados.`;
        button.hidden = true;
        await onEstimateRegistered({ storyID: story.id, storyPoints: registered.story_points });
      } catch (error) {
        messageElement.textContent = error.message;
      } finally {
        registeringAgreement = false;
        button.disabled = false;
      }
    }

    function addAgreementButton(container, targetRound, roundResult, messageElement) {
      if (roundResult.has_differences) return;
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = "Registrar estimación acordada";
      button.hidden = registeredAgreementRoundIDs.has(targetRound.id);
      button.addEventListener("click", () => registerAgreement(targetRound, button, messageElement));
      container.append(button);
    }

    async function loadOwnVote() {
      if (!round || !memberSelect.value) return;
      const path = `${base}/rounds/${encodeURIComponent(round.id)}/votes?member_id=${encodeURIComponent(memberSelect.value)}`;
      try {
        const votes = await requestJSON(path);
        pointsInput.value = votes[0] ? String(votes[0].story_points) : "";
        message.textContent = votes.length
          ? "Tu voto está guardado y oculto para los demás."
          : "Todavía no registraste tu voto. Los votos ajenos permanecen ocultos hasta finalizar la ronda.";
      } catch (error) {
        message.textContent = error.message;
      }
    }

    async function loadFinalizedResult(targetRound, memberNames, messageElement, votesElement) {
      const path = `${base}/rounds/${encodeURIComponent(targetRound.id)}/result?member_id=${encodeURIComponent(memberSelect.value)}`;
      const roundResult = await requestJSON(path);
      votesElement.replaceChildren();
      for (const vote of roundResult.votes) {
        const item = document.createElement("li");
        const name = memberNames.get(vote.member_id) || vote.member_id;
        item.textContent = `${name}: ${vote.story_points} Story Points`;
        votesElement.append(item);
      }
      messageElement.textContent = roundResult.has_differences
        ? "Hay diferencias entre las estimaciones."
        : "No hay diferencias entre las estimaciones.";
      agreement.hidden = roundResult.has_differences || registeredAgreementRoundIDs.has(targetRound.id);
    }

    async function loadFinalizedHistory(excludedRoundID, memberNames) {
      resultHistory.replaceChildren();
      const finalizedRounds = rounds.filter(item => item.status === "finalized" && item.id !== excludedRoundID);
      for (const previousRound of finalizedRounds) {
        if (!previousRound.participants.includes(memberSelect.value)) continue;
        const path = `${base}/rounds/${encodeURIComponent(previousRound.id)}/result?member_id=${encodeURIComponent(memberSelect.value)}`;
        const previousResult = await requestJSON(path);
        const card = document.createElement("article");
        const heading = document.createElement("h5");
        heading.textContent = `Ronda ${previousRound.number}`;
        const summary = document.createElement("p");
        summary.textContent = previousResult.has_differences
          ? "Hay diferencias entre las estimaciones."
          : "No hay diferencias entre las estimaciones.";
        const votes = document.createElement("ul");
        for (const vote of previousResult.votes) {
          const item = document.createElement("li");
          const name = memberNames.get(vote.member_id) || vote.member_id;
          item.textContent = `${name}: ${vote.story_points} Story Points`;
          votes.append(item);
        }
        card.append(heading, summary, votes);
        addAgreementButton(card, previousRound, previousResult, summary);
        resultHistory.append(card);
      }
      return resultHistory.children.length > 0;
    }

    async function load() {
      start.hidden = true;
      form.hidden = true;
      finalize.hidden = true;
      result.hidden = true;
      agreement.hidden = true;
      resultHeading.textContent = "";
      resultMessage.textContent = "";
      resultVotes.replaceChildren();
      resultHistory.replaceChildren();
      message.textContent = "Consultando rondas…";
      try {
        rounds = await requestJSON(`${base}/rounds`);
        round = rounds.at(-1) ?? null;
        if (!round) {
          message.textContent = "Todavía no hay una ronda de estimación.";
          start.textContent = "Iniciar ronda de estimación";
          start.hidden = false;
          return;
        }
        const allMembers = await requestJSON(`/api/projects/${encodeURIComponent(projectID)}/members`);
        const memberNames = new Map(allMembers.map(member => [member.id, member.name]));
        const participants = allMembers.filter(member => round.participants.includes(member.id));
        memberSelect.replaceChildren();
        for (const member of participants) {
          const option = document.createElement("option");
          option.value = member.id;
          option.textContent = `${member.name} · ${member.id}`;
          memberSelect.append(option);
        }
        if (!memberSelect.value && participants.length > 0) memberSelect.value = participants[0].id;

        if (round.status === "finalized") {
          start.textContent = "Iniciar nueva ronda";
          start.hidden = false;
          message.textContent = "La ronda está finalizada. Los votos fueron revelados.";
          if (memberSelect.value) {
            resultHeading.textContent = `Ronda ${round.number}`;
            await loadFinalizedResult(round, memberNames, resultMessage, resultVotes);
            result.hidden = false;
          }
        } else {
          form.hidden = participants.length === 0;
          finalize.hidden = participants.length === 0;
          message.textContent = participants.length ? "Los votos ajenos permanecen ocultos hasta finalizar la ronda." : "La ronda no tiene miembros participantes.";
          if (memberSelect.value) await loadOwnVote();
          resultHeading.textContent = "Rondas anteriores";
          resultMessage.textContent = "Resultados de rondas anteriores";
        }

        const hasHistory = await loadFinalizedHistory(round.status === "finalized" ? round.id : null, memberNames);
        if (round.status !== "finalized") result.hidden = !hasHistory;
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

    agreement.addEventListener("click", async () => {
      if (!round) return;
      await registerAgreement(round, agreement, resultMessage);
    });

    memberSelect.addEventListener("change", () => {
      return load().catch(error => { message.textContent = error.message; });
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
