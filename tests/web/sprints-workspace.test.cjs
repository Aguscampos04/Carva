const {test}=require('node:test'),assert=require('node:assert/strict'),fs=require('node:fs'),vm=require('node:vm'),path=require('node:path');
const {Element}=require('./app-harness.cjs');
const all=n=>[n,...n.children.flatMap(all)];
function setup(api){
  const elements=new Map(),get=id=>{if(!elements.has(id))elements.set(id,new Element());return elements.get(id);};
  const calls=[];
  const scrolls=[];
  const effortWidgets=[];
  const c=vm.createContext({document:{querySelector:get,createElement:tag=>new Element(tag)},encodeURIComponent,window:{confirm:()=>true,scrollX:0,scrollY:650,scrollTo:value=>scrolls.push(value)},createEffort:config=>({createWidget(project,story,sprintID){effortWidgets.push({project,story,sprintID,config});const widget=new Element(),open=new Element('button');open.click=()=>{};widget.append(open);return widget;}})});
  for(const name of ['workspace','sprints'])vm.runInContext(fs.readFileSync(path.join(__dirname,`../../web/static/js/${name}.js`),'utf8'),c);
  const ui=c.createSprintWorkspace({requestJSON:async(url,options={})=>{calls.push({url,...options});return api(url,options);},canLeave:()=>true,onChanged:async()=>{}});
  return {get,ui,calls,scrolls,effortWidgets};
}
const story={id:'story',backlog_item_id:'story',title:'Demo',status:'Pendiente',acceptance_criteria:['Verified'],estimated_hours:'1.00',acceptance_criteria_verified:false};
const sprint={id:1,goal:'MVP',closed:false};
function api(url){if(url.endsWith('/metrics/story-points'))return {sprints:[sprint]};if(url.endsWith('/backlog'))return [story];if(url.endsWith('/backlog-items'))return [story];return sprint;}
test('UX-03: completion requires explicit verification and duplicate submission is blocked',async()=>{
  let resolveWrite;const e=setup((url,options)=>options.method?new Promise(resolve=>{resolveWrite=resolve;}):api(url));
  await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  const detail=e.get('#sprint-detail'),complete=all(detail).find(n=>n.textContent==='Registrar como completada'),check=all(detail).find(n=>n.type==='checkbox');
  assert.equal(complete.disabled,true);complete.listeners.click();assert.equal(e.calls.filter(c=>c.method).length,0);
  check.checked=true;check.listeners.change();assert.equal(complete.disabled,false);
  complete.listeners.click();complete.listeners.click();
  assert.equal(e.calls.filter(c=>c.method).length,1);
  assert.equal(e.calls.at(-1).body,'{"acceptance_criteria_verified":true}');
  resolveWrite({...story,acceptance_criteria_verified:true});await new Promise(resolve=>setImmediate(resolve));
});
test('UX-03: rejected close stays visible and makes no automatic second request',async()=>{
  const e=setup((url,options)=>{if(options.method){const error=new Error('next Sprint required');error.status=409;throw error;}return api(url);});
  await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  all(e.get('#sprint-detail')).find(n=>n.textContent==='Cerrar Sprint').listeners.click();
  await new Promise(resolve=>setImmediate(resolve));
  assert.match(e.get('#sprints-message').textContent,/Sprint posterior abierto/);
  assert.equal(e.calls.filter(c=>c.method).length,1);
  assert.equal(e.get('#sprint-detail').hidden,false);
});
test('UX-04: old project reads cannot populate the new Sprint workspace',async()=>{
  let resolveOld;const e=setup(url=>url.includes('/old/')?new Promise(resolve=>{if(url.endsWith('/backlog'))resolve([]);else resolveOld=resolve;}):url.endsWith('/backlog')?[]:{sprints:[]});
  const old=e.ui.showProject({id:'old',status:'Planificado'});
  await e.ui.showProject({id:'new',status:'Planificado'});
  resolveOld({sprints:[sprint]});await old;
  assert.match(e.get('#sprints-list').children[0].textContent,/No hay Sprints abiertos/);
});

test('Sprint detail toggles closed and can be opened again; criteria start folded', async()=>{
  const e=setup(api);await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  assert.equal(e.get('#sprint-detail').querySelector('details').open,false);
  const toggle=all(e.get('#sprints-list')).find(n=>n.textContent==='Ocultar detalles ↑');
  assert.equal(toggle.attrs['aria-expanded'],'true');
  toggle.listeners.click();assert.equal(e.get('#sprint-detail').hidden,true);
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  assert.equal(e.get('#sprint-detail').hidden,false);
  assert.equal(all(e.get('#sprint-detail')).some(n=>n.tagName==='a'),false);
  assert(all(e.get('#sprint-detail')).some(n=>n.textContent==='Criterios de aceptación'));
  assert.equal(all(e.get('#sprint-detail')).filter(n=>n.tagName==='details').length,1);
});

test('Completion preserves expanded criteria and scroll; completed stories go last',async()=>{
  let finished=false;
  const e=setup((url,options)=>{
    if(options.method){finished=true;return {...story,acceptance_criteria_verified:true};}
    if(url.endsWith('/backlog-items'))return [{...story,acceptance_criteria_verified:finished,status:finished?'Completada':'Pendiente'},{...story,id:'next',backlog_item_id:'next',title:'Next'}];
    return api(url);
  });
  await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  for(const d of e.get('#sprint-detail').querySelectorAll('details')) if(d.dataset.disclosureKey!=='effort-story')d.open=true;
  const check=all(e.get('#sprint-detail')).find(n=>n.type==='checkbox');check.checked=true;check.listeners.change();
  all(e.get('#sprint-detail')).find(n=>n.textContent==='Registrar como completada').listeners.click();
  await new Promise(resolve=>setImmediate(resolve));
  const pending=all(e.get('#sprint-detail')).find(n=>n.id==='sprint-1-process');
  const finishedPane=all(e.get('#sprint-detail')).find(n=>n.id==='sprint-1-finished');
  assert.equal(pending.querySelector('summary').textContent,'Next');
  assert.equal(finishedPane.querySelector('summary').textContent,'Demo');
  assert.equal(finishedPane.querySelector('details').open,true);
  assert.equal(finishedPane.hidden,true);
  all(e.get('#sprint-detail')).find(n=>n.textContent==='Finalizadas (1)').listeners.click();
  assert.equal(finishedPane.hidden,false);assert.equal(pending.hidden,true);
  assert.equal(e.scrolls.at(-1).top,650);
});

test('Sprint story shows title only while folded and mounts effort in place with Sprint context',async()=>{
  const e=setup(api);await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  const detail=e.get('#sprint-detail'),storyDetails=detail.querySelector('details');
  assert.equal(storyDetails.children[0].tagName,'summary');assert.equal(storyDetails.children[0].textContent,'Demo');assert.equal(storyDetails.open,false);
  assert.equal(e.effortWidgets.length,0);
  storyDetails.open=true;storyDetails.listeners.toggle();
  assert.equal(e.effortWidgets.length,1);assert.equal(e.effortWidgets[0].sprintID,1);assert.equal(e.effortWidgets[0].story.id,'story');
  assert.equal(e.effortWidgets[0].config.entryOnly,true);
  assert.equal(e.get('#sprint-detail').hidden,false);
});

test('Responsible member is saved on the HU endpoint and stale widgets cannot write',async()=>{
  const e=setup((url,options)=>{
    if(url.endsWith('/members'))return [{id:'member',name:'Ana'}];
    if(url.endsWith('/assignee'))return {member_id:null};return api(url);
  });
  await e.ui.showProject({id:'p',status:'En progreso'});
  await all(e.get('#sprints-list')).find(n=>n.textContent==='Ver Sprint →').listeners.click();
  const storyDetails=e.get('#sprint-detail').querySelector('details');storyDetails.open=true;storyDetails.listeners.toggle();
  await new Promise(resolve=>setImmediate(resolve));
  const select=all(storyDetails).find(n=>n.id?.startsWith('assignee-'));select.value='member';
  const form=all(storyDetails).find(n=>n.tagName==='form');
  assert.equal(form.className,'assignee-form');
  assert.equal(form.children[1].children[0],select);
  assert.equal(form.children[1].children[1].textContent,'Guardar responsable');
  assert.equal(all(storyDetails).some(n=>n.textContent==='Responsable actual de la HU'),false);
  await form.listeners.submit({preventDefault(){}});
  assert.equal(e.calls.at(-1).url,'/api/projects/p/backlog/story/assignee');
  assert.equal(e.calls.at(-1).body,'{"member_id":"member"}');
  const count=e.calls.length;e.ui.invalidate();await form.listeners.submit({preventDefault(){}});
  assert.equal(e.calls.length,count);
});
