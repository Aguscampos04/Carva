// US-07..11: UI for existing REST contracts; all transitions are decided in Go.
function createSprintWorkspace({requestJSON, canLeave, onChanged}) {
  const find = selector => document.querySelector(selector);
  const message = find('#sprints-message'), list = find('#sprints-list'), detail = find('#sprint-detail');
  let project = null, generation = 0, selection = 0, selected = '', closed = false, mutating = false;
  let sprints = [], backlog = [];
  let storyTab = 'process';
  const node = (tag, text, className) => { const n=document.createElement(tag); if(text!=null)n.textContent=text;if(className)n.className=className;return n; };
  const action = (text, callback, className='secondary-button') => {const n=node('button',text,className);n.type='button';n.addEventListener('click',callback);return n;};
  const base = () => `/api/projects/${encodeURIComponent(project.id)}`;
  const body = value => ({method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(value)});
  const current = context => context === generation && project;
  function explain(error) {
    if (error.status === 409) return `No se pudo realizar la operación. Si hay historias incompletas, creá un Sprint posterior abierto antes de cerrar. Un Sprint cerrado no admite nuevas asignaciones ni finalizaciones. (${error.message})`;
    return error.message;
  }
  async function load(prefix='') {
    if (!project) return;
    const context=generation;
    message.textContent=prefix+'Consultando Sprints…';
    try {
      const [result,items]=await Promise.all([requestJSON(`${base()}/metrics/story-points`),requestJSON(`${base()}/backlog`)]);
      if(!current(context))return;
      sprints=result.sprints;backlog=items;render();
      message.textContent=prefix+(sprints.length?`${sprints.length} ${sprints.length === 1 ? 'Sprint' : 'Sprints'} en el proyecto.`:'Todavía no hay Sprints. Creá el primero para empezar a planificar.');
      if(selected && sprints.some(s=>String(s.id)===selected))await open(selected, true);
    }catch(error){if(current(context))message.textContent=prefix+explain(error);}
  }
  function hideDetails() {
    const previous = selected;
    selected = ''; ++selection; detail.hidden = true; detail.replaceChildren(); render(previous);
  }
  function render(focusID = '') {
    list.replaceChildren();
    const visible=sprints.filter(s=>s.closed===closed);
    if(!visible.length){list.append(node('p',closed?'No hay Sprints cerrados. Acá aparecerá el historial de tus entregas.':'No hay Sprints abiertos.','empty-state'));return;}
    for(const sprint of visible){
      const card=node('article',null,'sprint-card');card.setAttribute('aria-current',String(String(sprint.id)===selected));
      const expanded = String(sprint.id) === selected;
      const toggle = action(expanded ? 'Ocultar detalles ↑' : 'Ver Sprint →', () => {
        if (mutating || !canLeave()) return;
        if (String(sprint.id) === selected) { hideDetails(); return; }
        return open(String(sprint.id));
      });
      toggle.setAttribute('aria-expanded', String(expanded)); toggle.setAttribute('aria-controls', 'sprint-detail');
      card.append(node('span',sprint.closed?'Cerrado · Historial':'Abierto · Planificación','badge'),node('h2',`Sprint ${sprint.id}`),node('p',sprint.goal,'muted'),toggle);list.append(card);
      if (focusID === String(sprint.id)) toggle.focus({preventScroll:true});
    }
  }
  async function mutate(url, options, success, after=()=>{}) {
    if(mutating||!project)return;
    mutating=true; const context=generation;let saved=false;
    const position = {left: window.scrollX || 0, top: window.scrollY || 0};
    message.textContent='Guardando…';
    const controls = [...detail.querySelectorAll('button')].map(button => [button, button.disabled]);
    for(const [button] of controls)button.disabled=true;
    try {
      await requestJSON(url,options); saved=true;
      if(!current(context))return;
      after(); await onChanged();
      if(!current(context))return;
      await load(success+' ');
    }catch(error){if(current(context))message.textContent=(saved?'La operación se guardó. Actualizá antes de repetir. ':'')+explain(error);}
    finally{if(current(context)){
      mutating=false;for(const [button, disabled] of controls)button.disabled=disabled;
      if (typeof window.scrollTo === 'function') window.scrollTo({...position, behavior:'instant'});
    }}
  }
  async function open(id, refreshing = false) {
    if(!project)return;
    const context=generation, request=++selection;
    const expandedCriteria = new Set([...detail.querySelectorAll('details')].filter(n => n.open).map(n => n.dataset.disclosureKey));
    if (selected !== String(id)) storyTab = 'process';
    selected=String(id);render();detail.hidden=false;
    if (!refreshing) detail.replaceChildren(node('p','Consultando planificación…'));
    const url=`${base()}/sprints/${encodeURIComponent(id)}`;
    try {
      const [sprint,items,members]=await Promise.all([requestJSON(url),requestJSON(`${url}/backlog-items`),requestJSON(`${base()}/members`)]);
      if(!current(context)||request!==selection)return;
      const viewRequest = async (target, options) => {
        if (!current(context) || request !== selection) throw new Error('La vista cambió.');
        const value = await requestJSON(target, options);
        if (!current(context) || request !== selection) throw new Error('La vista cambió.');
        return value;
      };
      const heading = node('div', null, 'sprint-detail-heading');
      heading.append(node('h2',`Sprint ${sprint.id} · ${sprint.goal}`), action('Ocultar detalles', () => { if (!mutating && canLeave()) hideDetails(); }));
      detail.replaceChildren(node('span',sprint.closed?'HISTORIAL DEL SPRINT':'PLANIFICACIÓN','step-label'),heading);
      if(sprint.closed)detail.append(node('p','Resultado registrado al cierre. Las horas estimadas vigentes pueden actualizarse; los demás datos planificados y resultados históricos se conservan.','muted'));
      if(!sprint.closed){
        const form=node('form',null,'sprint-assign'),label=node('label','Historia del Backlog'),select=node('select');select.id='assign-story';label.htmlFor=select.id;
        const available=backlog.filter(item=>!items.some(planned=>planned.backlog_item_id===item.id));
        const placeholder=node('option',available.length?'Seleccioná una historia':'No hay historias disponibles');placeholder.value='';select.append(placeholder);select.required=true;
        for(const story of available){const option=node('option',story.title);option.value=story.id;select.append(option);}
        const submit=node('button','Asignar al Sprint');submit.type='submit';submit.disabled=!available.length;
        form.append(label,select,submit);detail.append(form);
        form.addEventListener('submit',async event=>{event.preventDefault();if(!select.value||mutating)return;await mutate(`${url}/backlog-items`,body({backlog_item_id:select.value}),'Historia asignada.');});
      }
      if(!items.length)detail.append(node('p','Todavía no hay historias planificadas.','muted'));
      const storyTabs = node('div', null, 'scope-controls sprint-story-tabs');
      storyTabs.setAttribute('aria-label', 'Estado de las historias del Sprint');
      const pendingPane = node('div'), finishedPane = node('div');
      pendingPane.id = `sprint-${id}-process`; finishedPane.id = `sprint-${id}-finished`;
      const pendingCount = items.filter(item => !item.acceptance_criteria_verified).length;
      const pendingTab = action(`En proceso (${pendingCount})`, () => changeStoryTab('process'));
      const finishedTab = action(`Finalizadas (${items.length - pendingCount})`, () => changeStoryTab('finished'));
      function updateStoryTabs() {
        pendingPane.hidden = storyTab !== 'process'; finishedPane.hidden = storyTab !== 'finished';
        pendingTab.setAttribute('aria-pressed', String(storyTab === 'process')); finishedTab.setAttribute('aria-pressed', String(storyTab === 'finished'));
      }
      function changeStoryTab(tab) { if (mutating || !canLeave()) return; storyTab = tab; updateStoryTabs(); }
      pendingTab.setAttribute('aria-controls', pendingPane.id); finishedTab.setAttribute('aria-controls', finishedPane.id);
      storyTabs.append(pendingTab, finishedTab); updateStoryTabs(); detail.append(storyTabs, pendingPane, finishedPane);
      if (!pendingCount) pendingPane.append(node('p', sprint.closed ? 'No hay historias incompletas al cierre.' : 'No hay historias en proceso. Consultá la pestaña Finalizadas.', 'muted'));
      if (pendingCount === items.length) finishedPane.append(node('p', 'Todavía no hay historias finalizadas.', 'muted'));
      for(const item of CarvaWorkspace.pendingFirst(items)){
        const row=node('article',null,'planned-item');
        const storyDetails = node('details', null, 'sprint-story'); storyDetails.dataset.disclosureKey = `story-${item.backlog_item_id}`;
        storyDetails.open = refreshing && expandedCriteria.has(storyDetails.dataset.disclosureKey);
        storyDetails.append(node('summary', item.title));
        const content = node('div', null, 'sprint-story-content'); storyDetails.append(content); row.append(storyDetails);
        content.append(node('span',sprint.closed?(item.acceptance_criteria_verified?'Completada al cierre':'Incompleta al cierre'):item.status,'badge'),node('p',item.description,'muted'));
        const values=node('p',`${CarvaWorkspace.points(item)} · Horas estimadas vigentes: ${item.estimated_hours ?? 'Pendiente de regularización'}`,'muted');content.append(values);
        const disclosure = node('section', null, 'story-criteria');
        disclosure.append(node('h4', 'Criterios de aceptación'));
        const criteria=node('ul',null,'criteria-list');for(const criterion of item.acceptance_criteria)criteria.append(node('li',criterion));disclosure.append(criteria);content.append(disclosure);
        const assignee = node('section', null, 'story-assignee'); content.append(assignee);
        let assigneeLoaded = false;
        async function loadAssignee() {
          if (assigneeLoaded) return;
          assigneeLoaded = true; assignee.replaceChildren(node('p', 'Consultando responsable…', 'muted'));
          try {
            const target = `${base()}/backlog/${encodeURIComponent(item.backlog_item_id)}/assignee`;
            const value = await viewRequest(target);
            const team = Array.isArray(members) ? members : [];
            assignee.replaceChildren();
            if (sprint.closed) {
              const member = team.find(m => m.id === value.member_id);
              assignee.append(node('p', `Miembro responsable: ${member ? member.name : 'Sin responsable'}`, 'muted')); return;
            }
            const form = node('form', null, 'assignee-form'), label = node('label','Miembro responsable'), select = node('select');
            select.id = `assignee-${id}-${item.backlog_item_id}`; label.htmlFor = select.id;
            const none = node('option','Sin responsable'); none.value = ''; select.append(none);
            for (const member of team) { const homonym = team.filter(m => m.name === member.name).length > 1; const option = node('option', homonym ? `${member.name} · ${member.id}` : member.name); option.value = member.id; select.append(option); }
            select.value = value.member_id || '';
            const save = node('button', 'Guardar responsable'); save.type = 'submit';
            const controls = node('div', null, 'assignee-controls'); controls.append(select, save);
            const status = node('p'); status.setAttribute('role','status'); form.append(label, controls, status); assignee.append(form);
            if (!team.length) assignee.append(node('p','Registrá miembros en Equipo para asignar un responsable.','muted'));
            form.addEventListener('submit', async event => {
              event.preventDefault(); if (save.disabled) return; save.disabled = true;
              try { await viewRequest(target, {...body({member_id:select.value || null}),method:'PUT'}); status.textContent='Responsable guardado.'; }
              catch(error) { status.textContent=error.message; } finally {save.disabled=false;}
            });
          } catch(error) {
            assigneeLoaded = false;
            assignee.replaceChildren(node('p',error.message,'form-error'), action('Reintentar responsable',loadAssignee));
          }
        }
        const effort = node('section', null, 'sprint-hours-entry');
        let effortLoaded = false;
        function loadEffortEntry() {
          if (effortLoaded) return;
          effortLoaded = true;
          const widget = createEffort({requestJSON:viewRequest,entryOnly:true,confirmDiscard:text=>window.confirm(text),confirmDelete:text=>window.confirm(text),onChanged:async()=>{
            await onChanged(); if(current(context)&&request===selection)message.textContent='Esfuerzo actualizado dentro del Sprint.';
          }}).createWidget(project, {id:item.backlog_item_id}, sprint.id);
          effort.append(widget);
        }
        storyDetails.addEventListener('toggle', () => { if (storyDetails.open) { loadAssignee(); loadEffortEntry(); } });
        if (storyDetails.open) { loadAssignee(); loadEffortEntry(); }
        content.append(effort);
        if(!sprint.closed&&!item.acceptance_criteria_verified){
          const label=node('label'),check=node('input');check.type='checkbox';label.append(check,node('span','Confirmo que verifiqué todos los criterios de aceptación de esta historia.'));
          const complete=action('Registrar como completada',()=>{
            if(!check.checked||mutating||!canLeave())return;
            mutate(`${url}/backlog-items/${encodeURIComponent(item.backlog_item_id)}/completion`,body({acceptance_criteria_verified:true}),'Historia completada. Disponible en la pestaña Finalizadas.');
          });complete.disabled=true;check.addEventListener('change',()=>{complete.disabled=!check.checked;});content.append(label,complete);
        }(item.acceptance_criteria_verified ? finishedPane : pendingPane).append(row);
      }
      if(!sprint.closed){const close=action('Cerrar Sprint',()=>{
        if(mutating||!canLeave())return;
        if(!window.confirm('¿Cerrar este Sprint? Las historias incompletas se trasladarán al siguiente Sprint abierto disponible según el servidor. Si no existe y hay incompletas, el cierre se rechazará. El historial se conservará.'))return;
        mutate(`${url}/close`,{method:'POST'},'Sprint cerrado.',()=>{closed=true;setTabs();});
      },'danger-button');detail.append(close);}
    }catch(error){if(current(context)&&request===selection){detail.replaceChildren(node('p',explain(error),'form-error'));}}
  }
  function setTabs(){find('#show-open-sprints').setAttribute('aria-pressed',String(!closed));find('#show-closed-sprints').setAttribute('aria-pressed',String(closed));}
  for(const [id,value]of [['#show-open-sprints',false],['#show-closed-sprints',true]])find(id).addEventListener('click',()=>{if(mutating||!canLeave())return;closed=value;selected='';++selection;detail.hidden=true;setTabs();render();});
  find('#refresh-sprints').addEventListener('click',()=>{if(!mutating&&canLeave())load();});
  find('#new-sprint').addEventListener('click',()=>{if(!mutating&&canLeave()){find('#sprint-create-panel').hidden=false;find('#sprint-goal').focus();}});
  find('#cancel-sprint').addEventListener('click',()=>{if(!mutating&&canLeave()){find('#create-sprint-form').reset();find('#sprint-create-panel').hidden=true;find('#new-sprint').focus();}});
  find('#create-sprint-form').addEventListener('submit',async event=>{
    event.preventDefault();if(mutating||!project)return;const form=event.currentTarget,submit=form.querySelector('button[type=submit]');submit.disabled=true;
    await mutate(`${base()}/sprints`,body({goal:find('#sprint-goal').value}),'Sprint creado.',()=>{form.reset();find('#sprint-create-panel').hidden=true;closed=false;setTabs();});submit.disabled=false;
  });
  return {
    async showProject(value){++generation;++selection;project=value;selected='';closed=false;mutating=false;setTabs();detail.hidden=true;list.replaceChildren();find('#new-sprint').disabled=value.status==='Finalizado';await load();},
    invalidate(){++generation;++selection;project=null;selected='';mutating=false;list.replaceChildren();detail.replaceChildren();detail.hidden=true;},
  };
}
