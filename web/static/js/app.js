// Local UX preview: navigation and orchestration only. Go owns business rules.
const $ = selector => document.querySelector(selector);
const W = CarvaWorkspace;
const routeGate = W.createGate();
let currentProject = null;
let currentRoute = {view:'inicio', project:''};
let activeHash = location.hash || '#/';
let stories = [];
let backlogGeneration = 0;
let homeGeneration = 0;
let busy = 0;
const dirtyForms = new Set();
const names = {inicio:'Inicio', resumen:'Resumen', backlog:'Backlog', sprints:'Sprints', equipo:'Equipo', datos:'Datos del proyecto', reportes:'Reportes'};
const reportsUI = createReports({
  select: $('#report-scope'), refresh: $('#refresh-report'), download: $('#download-report'),
  message: $('#report-message'), content: $('#report-content'), requestJSON,
});

function element(tag, text, className) {
  const node = document.createElement(tag);
  if (text != null) node.textContent = text;
  if (className) node.className = className;
  return node;
}
function button(text, action, className = '') {
  const node = element('button', text, className); node.type = 'button';
  node.addEventListener('click', action); return node;
}
function badge(value) { const node = element('span', value, 'badge'); node.dataset.value = value; return node; }
function notify(text) { $('#app-message').textContent = text; $('#app-message').hidden = !text; }
function clearDrafts() { dirtyForms.clear(); }
function hasDrafts() {
  for (const form of dirtyForms) if (!form.isConnected) dirtyForms.delete(form);
  return dirtyForms.size > 0;
}
function canLeave() {
  if (busy) { notify('Hay una operación en curso. Esperá a que termine antes de cambiar de sección.'); return false; }
  if (!hasDrafts()) return true;
  if (!window.confirm('Tenés cambios sin guardar. ¿Descartarlos y continuar?')) return false;
  for (const form of dirtyForms) if (form.isConnected) form.reset();
  clearDrafts(); return true;
}
function go(hash) {
  if (hash === location.hash) return;
  if (canLeave()) { activeHash = hash; location.hash = hash; }
}
function path(id = currentProject?.id) { return `/api/projects/${encodeURIComponent(id)}`; }
async function requestJSON(url, options = {}) {
  const writing = options.method && options.method !== 'GET';
  const form = document.activeElement?.closest('form');
  if (writing) busy++;
  try {
    const response = await fetch(url, options);
    const text = await response.text();
    let data;
    try { data = text ? parseHoursJSON(text) : null; }
    catch { throw new Error('La respuesta no se pudo leer. Consultá de nuevo antes de repetir la operación.'); }
    if (!response.ok) {
      const error = new Error(data?.error || 'No se pudo completar la operación.'); error.status = response.status; throw error;
    }
    if (writing && form) dirtyForms.delete(form);
    return data;
  } catch (error) {
    if (error instanceof TypeError) throw new Error('No hay conexión con Carva. Si estabas guardando, consultá antes de volver a enviar.');
    throw error;
  } finally { if (writing) busy--; }
}
function scopedRequest(token) {
  return async (url, options) => {
    if (!routeGate.current(token)) throw new Error('La vista cambió.');
    const result = await requestJSON(url, options);
    if (!routeGate.current(token)) throw new Error('La vista cambió.');
    return result;
  };
}
const jsonOptions = (method, body) => ({method, headers:{'Content-Type':'application/json'}, body:JSON.stringify(body)});

const projectEditor = createProjectEditor({
  form:$('#edit-project-form'), nameInput:$('#edit-project-name'), startInput:$('#edit-project-start'),
  endOutput:$('#project-end-result'), finalizeButton:$('#finalize-project'), saveButton:$('#save-project'),
  message:$('#project-edit-message'), requestJSON,
  confirmFinalization:text => window.confirm(text),
  onBusyChange:value => { $('#refresh-project').disabled = value; },
  onSaved:project => {
    dirtyForms.delete($('#edit-project-form')); currentProject = project; updateIdentity(); projectEditor.showProject(project);
    $('#project-edit-message').textContent = 'Proyecto guardado.';
  },
});
const membersUI = createMembersSection({section:$('#members-section'), form:$('#create-member-form'), nameInput:$('#member-name'), message:$('#members-message'), list:$('#members-list'), refresh:$('#refresh-members'), submit:$('#register-member'), requestJSON});
const metricsUI = createMetricsSection({
  section:$('#metrics-section'), message:$('#metrics-message'), projectMetrics:$('#project-metrics'), sprintSelect:$('#metrics-sprint-select'), sprintMetrics:$('#sprint-metrics'), projectHours:$('#project-hours'), sprintHours:$('#sprint-hours'), projectDefects:$('#project-defects'), sprintDefects:$('#sprint-defects'), requestJSON,
  dashboard:{section:$('#metrics-section'), identity:$('#dashboard-identity'), projectCharts:$('#dashboard-project-charts'), sprintCharts:$('#dashboard-sprint-charts'), sprintTitle:$('#dashboard-sprint-title')}, refresh:$('#refresh-metrics'),
});
const sprintsUI = createSprintWorkspace({requestJSON, canLeave, onChanged:async () => {
  currentProject = await requestJSON(path()); updateIdentity();
}});

function updateIdentity() {
  const p = currentProject;
  $('#sidebar-project-name').textContent = p?.name || 'Tu próximo proyecto';
  $('#sidebar-project-status').textContent = p?.status || 'Todo empieza con una idea';
  $('#project-navigation').hidden = !p;
  $('#menu-toggle').hidden = !p;
  for (const link of document.querySelectorAll('[data-view], [data-go]')) {
    const view = link.dataset.view || link.dataset.go;
    link.href = W.route(p?.id, view);
    if (link.dataset.view && view === currentRoute.view) link.setAttribute('aria-current', 'page');
    else link.removeAttribute('aria-current');
  }
  $('#breadcrumb-view').textContent = names[currentRoute.view];
  document.title = p ? `${names[currentRoute.view]} · ${p.name} · Carva` : 'Carva · Tu espacio de trabajo';
  if (!p) return;
  for (const [id, value] of Object.entries({
    'project-name-result':p.name, 'project-id-result':p.id, 'project-status-result':p.status,
    'project-start-result':p.start_date || 'Sin registrar', 'project-end-result':p.end_date || 'Sin registrar',
    'overview-name':p.name, 'overview-status':p.status, 'overview-start':p.start_date || 'Sin registrar', 'overview-end':p.end_date || 'Sin registrar',
  })) $(`#${id}`).textContent = value;
  $('#overview-status').dataset.value = p.status;
}
function showPage(view) {
  for (const section of document.querySelectorAll('[data-page]')) section.hidden = section.dataset.page !== view;
}
async function loadProjects() {
  const request = ++homeGeneration;
  const list = $('#projects-list');
  $('#projects-message').textContent = 'Consultando proyectos…';
  try {
    const projects = await requestJSON('/api/projects');
    if (request !== homeGeneration || currentRoute.view !== 'inicio') return;
    list.replaceChildren();
    $('#projects-message').textContent = projects.length ? `${projects.length} ${projects.length === 1 ? 'proyecto disponible' : 'proyectos disponibles'}. Elegí uno para continuar.` : 'Todavía no hay proyectos. Usá el formulario de arriba para crear el primero.';
    for (const project of projects) {
      const card = element('article', null, 'panel project-card');
      const title = element('a', project.name, 'story-link'); title.href = W.route(project.id);
      const open = element('a', 'Abrir proyecto →', 'text-link'); open.href = title.href;
      card.append(badge(project.status), element('h3'), element('p', `Inicio: ${project.start_date || 'Sin registrar'}`, 'muted'));
      card.querySelector('h3').append(title);
      const identity = element('details', null, 'project-identity'); identity.append(element('summary', 'Identificador'), element('small', project.id));
      card.append(identity, open); list.append(card);
    }
  } catch (error) {
    if (request === homeGeneration && currentRoute.view === 'inicio') $('#projects-message').textContent = `No se pudo cargar la lista. ${error.message} Usá Actualizar para reintentar.`;
  }
}
async function navigate() {
  if (location.hash !== activeHash && !canLeave()) { history.pushState(null, '', activeHash); return; }
  activeHash = location.hash || '#/';
  const token = routeGate.next();
  ++backlogGeneration;
  ++homeGeneration;
  sprintsUI.invalidate();
  reportsUI.invalidate();
  membersUI.showProject(''); metricsUI.showProject(''); projectEditor.showProject(null);
  $('#story-content').replaceChildren(); $('#story-heading').replaceChildren();
  $('#story-create-panel').hidden = true; $('#sprint-create-panel').hidden = true;
  $('#create-backlog-form').reset(); $('#create-sprint-form').reset();
  $('#project-navigation').classList.remove('is-open'); $('#menu-toggle').setAttribute('aria-expanded','false');
  const next = W.parseRoute(activeHash);
  clearDrafts(); notify('');
  showPage('none');
  if (!next) {
    currentProject = null; currentRoute = {view:'inicio', project:''}; updateIdentity(); showPage('inicio');
    notify('Este enlace no corresponde a una sección de Carva. Abrí un proyecto para continuar.'); return;
  }
  currentRoute = next;
  if (!next.project) { currentProject = null; updateIdentity(); showPage('inicio'); await loadProjects(); return; }
  notify('Abriendo proyecto…');
  try {
    const project = await requestJSON(path(next.project));
    if (!routeGate.current(token)) return;
    if (currentProject?.id !== project.id) {
      $('#backlog-search').value = ''; $('#filter-priority').value = ''; $('#filter-status').value = '';
    }
    currentProject = project; updateIdentity(); showPage(next.view); notify('');
    $('#main-content').focus({preventScroll:true});
    if (next.view === 'datos') projectEditor.showProject(project);
    if (next.view === 'equipo') await membersUI.showProject(project.id);
    if (next.view === 'resumen') { setScope(false); await metricsUI.showProject(project.id, project); }
    if (next.view === 'backlog') await loadBacklog(token);
    if (next.view === 'sprints') await sprintsUI.showProject(project);
    if (next.view === 'reportes') await reportsUI.showProject(project);
  } catch (error) {
    if (!routeGate.current(token)) return;
    currentProject = null; updateIdentity(); showPage('inicio');
    notify(error.status === 404 ? 'No encontramos ese proyecto. Revisá el identificador.' : error.message);
  }
}

async function loadBacklog(token) {
  const request = ++backlogGeneration;
  const target = currentProject.id;
  const isCurrent = () => request === backlogGeneration && target === currentProject?.id && (!token || routeGate.current(token));
  $('#backlog-message').textContent = 'Consultando historias…';
  $('#backlog-items').replaceChildren();
  try {
    const result = await requestJSON(`${path(target)}/backlog`);
    if (!isCurrent()) return;
    stories = result; renderStories();
    if (currentRoute.story) renderStory(token);
  } catch (error) { if (isCurrent()) $('#backlog-message').textContent = error.message; }
}
function empty(title, text) { const node=element('div', null, 'empty-state'); node.append(element('h2', title), element('p',text)); return node; }
function renderStories() {
  const selected = W.pendingFirst(W.filterStories(stories, $('#backlog-search').value, $('#filter-priority').value, $('#filter-status').value));
  $('#backlog-list-view').hidden = !!currentRoute.story;
  $('#story-detail').hidden = !currentRoute.story;
  $('#new-story').hidden = !!currentRoute.story;
  $('#backlog-message').textContent = `${selected.length} de ${stories.length} historias`;
  const list = $('#backlog-items'); list.replaceChildren();
  if (!selected.length) { list.append(empty(stories.length ? 'Sin coincidencias' : 'Tu Backlog empieza acá', stories.length ? 'Probá otra búsqueda o cambiá los filtros.' : 'Agregá tu primera historia para empezar a planificar.')); return; }
  for (const item of selected) {
    const row = element('article', null, 'backlog-row');
    const info = element('div'); const link = element('a', item.title, 'story-link'); link.href = W.route(currentProject.id, 'backlog', item.id);
    const meta = element('div',null,'story-meta'); meta.append(badge(item.priority), badge(item.status),element('span',W.points(item)));
    info.append(link,meta);
    const hours = element('div',null,'row-hours');
    hours.append(element('div',`Estimadas: ${item.estimated_hours == null ? 'Pendiente de regularización' : item.estimated_hours + ' h'}`), element('div',`Reales: ${item.actual_hours ?? '0.00'} h`));
    row.append(info,hours); list.append(row);
  }
}
function renderStory(token) {
  const item = stories.find(item => item.id === currentRoute.story);
  const host = $('#story-content'); host.replaceChildren();
  const heading = $('#story-heading'); heading.replaceChildren();
  $('#back-to-backlog').href = W.route(currentProject.id,'backlog');
  if (!item) { heading.append(empty('Historia no encontrada','Volvé al Backlog para seleccionar una historia de este proyecto.')); $('#story-tabs').replaceChildren(); return; }
  heading.append(element('h2',item.title));
  const meta=element('div',null,'story-meta'); meta.append(badge(item.priority),badge(item.status),element('span',W.points(item))); heading.append(meta,element('p',`Identificador: ${item.id}`,'muted'));
  const tabs=$('#story-tabs'); tabs.replaceChildren();
  for (const [key,label] of [['info','Información'],['poker','Planning Poker'],['effort','Esfuerzo'],['defects','Defectos']]) {
    const link=element('a',label); link.href=W.route(currentProject.id,'backlog',item.id,key);
    if (key===currentRoute.tab) link.setAttribute('aria-current','page'); tabs.append(link);
  }
  const request = scopedRequest(token);
  const reloadStory = async () => {
    const latest = await request(`${path()}/backlog/${encodeURIComponent(item.id)}`);
    stories = stories.map(story => story.id === latest.id ? latest : story);
    const meta = $('#story-heading .story-meta');
    meta.replaceChildren(badge(latest.priority), badge(latest.status), element('span', W.points(latest)));
    notify('Datos de la historia actualizados.');
  };
  if (currentRoute.tab === 'poker') host.append(createPlanningPoker({requestJSON:request,onEstimateRegistered:reloadStory}).createWidget(currentProject.id,item));
  else if (currentRoute.tab === 'effort') {
    const widget = createEffort({requestJSON:request,confirmDelete:text => window.confirm(text),onChanged:reloadStory}).createWidget(currentProject,item);
    host.append(widget); const open = widget.querySelector('button'); open.textContent = '↻ Actualizar esfuerzo'; open.click();
  }
  else if (currentRoute.tab === 'defects') {
    const widget = createDefects({requestJSON:request,sprintOptions:async()=> (await request(`${path()}/metrics/story-points`)).sprints}).createWidget(currentProject,item);
    host.append(widget); const open = widget.querySelector('button'); open.textContent = '↻ Actualizar defectos'; open.click();
  }
  else {
    host.append(element('p',item.description,'story-description'));
    const dl=element('dl');
    addDetail(dl,'Story Points',W.points(item)); addDetail(dl,'Horas estimadas',item.estimated_hours == null ? 'Pendiente de regularización' : `${item.estimated_hours} h`); addDetail(dl,'Horas reales',`${item.actual_hours ?? '0.00'} h`); host.append(dl,element('h3','Criterios de aceptación'));
    const criteria=element('ul',null,'criteria-list'); for(const text of item.acceptance_criteria) criteria.append(element('li',text)); host.append(criteria,button('Editar historia',()=> { if(canLeave()) {host.replaceChildren(createEditForm(item,token)); host.querySelector('input').focus();} },'secondary-button'));
  }
}
function addDetail(list,label,value) { list.append(element('dt',label),element('dd',String(value))); }
function field(form, name, label, tag, value) {
  const id=`story-edit-${name}`, input=element(tag); input.id=id; input.name=name;
  const title=element('label',label); title.htmlFor=id; form.append(title,input); input.value=value; return input;
}
function readStory(form) {
  const data=new FormData(form), points=data.get('provisional_story_points');
  return {title:data.get('title'), description:data.get('description'), priority:data.get('priority'), estimated_hours:data.get('estimated_hours'), provisional_story_points:points === '' ? null : Number(points), acceptance_criteria:String(data.get('acceptance_criteria')).split(/\r?\n/).filter(text=>text.trim())};
}
function createEditForm(item,token) {
  const form=element('form'); form.append(element('h3','Editar historia'));
  field(form,'title','Título','input',item.title).required=true;
  field(form,'description','Descripción','textarea',item.description).required=true;
  const priority=field(form,'priority','Prioridad','select',''); for(const value of ['Alta','Media','Baja']) { const option=element('option',value); option.value=value; priority.append(option); } priority.value=item.priority;
  const hours=field(form,'estimated_hours','Horas estimadas','input',item.estimated_hours ?? ''); hours.type='number';hours.min='0.01';hours.step='0.01';hours.required=true;
  const points=field(form,'provisional_story_points','Story Points provisionales','input',item.provisional_story_points ?? ''); points.type='number';points.min='0';points.step='1';points.readOnly=item.story_points!=null;
  field(form,'acceptance_criteria','Criterios de aceptación (uno por línea)','textarea',item.acceptance_criteria.join('\n')).required=true;
  const message=element('p');message.setAttribute('role','status');
  const save=element('button','Guardar cambios');save.type='submit';
  const cancel=button('Cancelar',()=>{if(canLeave()){dirtyForms.delete(form);renderStory(token);}},'secondary-button');
  form.append(message,save,cancel);
  form.addEventListener('submit',async event=>{
    event.preventDefault(); if(save.disabled)return;save.disabled=true;cancel.disabled=true;let saved=false;
    try {
      await scopedRequest(token)(`${path()}/backlog/${encodeURIComponent(item.id)}`,{method:'PUT',headers:{'Content-Type':'application/json'},body:exactHoursJSON(readStory(form),['estimated_hours'])});
      saved=true;dirtyForms.delete(form);await navigate();notify('Historia guardada.');
    }catch(error){message.textContent=(saved?'Historia guardada. Actualizá antes de repetir. ':'')+error.message;}
    finally{save.disabled=false;cancel.disabled=false;}
  });return form;
}

function setScope(sprint) {
  $('#scope-project').setAttribute('aria-pressed',String(!sprint));$('#scope-sprint').setAttribute('aria-pressed',String(sprint));
  $('#project-summary').hidden=sprint;$('#sprint-summary').hidden=!sprint;$('#metrics-sprint-picker').hidden=!sprint;
}
$('#scope-project').addEventListener('click',()=>setScope(false));$('#scope-sprint').addEventListener('click',()=>setScope(true));
$('#menu-toggle').addEventListener('click',()=>{const open=$('#project-navigation').classList.toggle('is-open');$('#menu-toggle').setAttribute('aria-expanded',String(open));});
for (const id of ['backlog-search','filter-priority','filter-status']) $(`#${id}`).addEventListener(id==='backlog-search'?'input':'change',renderStories);
$('#refresh-backlog').addEventListener('click',()=>{if(canLeave())loadBacklog();});
$('#new-story').addEventListener('click',()=>{if(canLeave()){$('#story-create-panel').hidden=false;$('#backlog-title-input').focus();}});
$('#cancel-new-story').addEventListener('click',()=>{if(canLeave()){$('#create-backlog-form').reset();$('#story-create-panel').hidden=true;$('#new-story').focus();}});
$('#create-backlog-form').addEventListener('submit',async event=>{
  event.preventDefault();const form=event.currentTarget,submit=form.querySelector('[type=submit]');if(submit.disabled)return;submit.disabled=true;
  const target=currentProject.id;let saved=false;$('#story-create-message').textContent='Guardando historia…';
  try {
    await requestJSON(`${path(target)}/backlog`,{method:'POST',headers:{'Content-Type':'application/json'},body:exactHoursJSON(readStory(form),['estimated_hours'])});
    saved=true;dirtyForms.delete(form);form.reset();$('#story-create-panel').hidden=true;await loadBacklog();notify('Historia agregada al Backlog.');
  }catch(error){$('#story-create-message').textContent=(saved?'Historia guardada. Actualizá antes de repetir. ':'')+error.message;}
  finally{submit.disabled=false;}
});
$('#open-project-form').addEventListener('submit',event=>{event.preventDefault();dirtyForms.delete(event.currentTarget);go(W.route($('#open-project-id').value.trim()));});
$('#create-project-form').addEventListener('submit',async event=>{
  event.preventDefault();const form=event.currentTarget,submit=form.querySelector('button');if(submit.disabled)return;submit.disabled=true;notify('Creando proyecto…');
  try{const project=await requestJSON('/api/projects',jsonOptions('POST',{name:$('#project-name').value}));dirtyForms.delete(form);form.reset();go(W.route(project.id));}
  catch(error){notify(error.message);}finally{submit.disabled=false;}
});
$('#refresh-project').addEventListener('click',()=>{if(canLeave())navigate();});
$('#refresh-projects').addEventListener('click',loadProjects);
$('#copy-project-id').addEventListener('click',async()=>{try{await navigator.clipboard.writeText(currentProject.id);notify('Identificador copiado.');}catch{notify('No se pudo copiar. Seleccioná el identificador en Información general.');}});
// Only actual forms participate in draft protection; search and selectors do not.
for(const eventName of ['input','change']) document.addEventListener(eventName,event=>{const form=event.target.closest('form');if(form)dirtyForms.add(form);});
document.addEventListener('click',event=>{
  const link=event.target.closest('a[href^="#/"]');if(!link||event.ctrlKey||event.metaKey||event.shiftKey||event.altKey)return;
  event.preventDefault();go(link.getAttribute('href'));
});
window.addEventListener('hashchange',navigate);
window.addEventListener('beforeunload',event=>{if(hasDrafts()||busy){event.preventDefault();event.returnValue='';}});
navigate();
