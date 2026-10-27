const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const path=require('node:path');
class Element {
  constructor(tag='div'){this.tagName=tag;this.children=[];this.listeners={};this.attributes={};this.textContent='';this.value='';}
  append(...nodes){this.children.push(...nodes)}
  replaceChildren(...nodes){this.children=nodes}
  addEventListener(name,fn){this.listeners[name]=fn}
  setAttribute(key,value){this.attributes[key]=String(value)}
}
const nodes=e=>[e,...e.children.flatMap(nodes)];
const text=e=>nodes(e).map(n=>n.textContent).join(' ');
const points={planned_story_points:10,completed_story_points:0,team_velocity:null,completion_percentage:0};
const hours={estimated_hours:'5.75',actual_hours:'3.75',deviation:-2};
const defects={detected_defects:0,resolved_defects:null};
const project={id:'p',name:'<img src=x>',status:'En progreso'};
function setup(requestJSON){
  const fields=['section','message','projectMetrics','sprintSelect','sprintMetrics','projectHours','sprintHours','projectDefects','sprintDefects','refresh'];
  const e=Object.fromEntries(fields.map(key=>[key,new Element()]));
  const dashboard=Object.fromEntries(['section','identity','projectCharts','sprintCharts','sprintTitle'].map(k=>[k,new Element()]));
  const c=vm.createContext({document:{createElement:tag=>new Element(tag),createElementNS:(_,tag)=>new Element(tag)},Intl,encodeURIComponent});
  vm.runInContext(fs.readFileSync(path.join(__dirname,'../../web/static/js/metrics.js'),'utf8'),c);
  return {...e,dashboard,ui:c.createMetricsSection({...e,dashboard,requestJSON})};
}
function response(url){
  const metric=url.endsWith('/hours')?hours:url.endsWith('/defects')?defects:points;
  return url.includes('/sprints/')?{...metric}:{metrics:{...metric},sprints:[{id:3,goal:'Original goal',closed:true}]};
}
test('US23-01/02/03/04: show official identity, source values, null, zero and signed deviation',async()=>{
  const calls=[];const e=setup(async url=>{calls.push(url);return response(url)});
  await e.ui.showProject('p',project);
  assert.equal(e.dashboard.section.hidden,false);
  assert.match(text(e.dashboard.identity),/En progreso/);
  assert.match(text(e.dashboard.identity),/<img src=x>/);
  assert.equal(nodes(e.dashboard.identity).some(n=>n.tagName==='img'),false);
  const all=nodes(e.dashboard.projectCharts);
  assert.equal(all.filter(n=>n.tagName==='figure').length,5);
  assert.match(text(e.dashboard.projectCharts),/No disponible/);
  const rows=all.filter(n=>n.attributes['data-metric']);
  const unavailable=rows.find(n=>n.attributes['data-metric']==='team_velocity');
  assert.equal(nodes(unavailable).some(n=>n.tagName==='rect'),false);
  assert.match(text(unavailable),/No disponible/);
  const zero=rows.find(n=>n.attributes['data-metric']==='detected_defects');
  assert.match(text(zero),/0/);assert.equal(nodes(zero).find(n=>n.tagName==='rect').attributes.width,'0');
  const deviation=rows.find(n=>n.attributes['data-metric']==='deviation');
  assert.match(text(deviation),/-2/);
  const rect=nodes(deviation).find(n=>n.tagName==='rect');
  const axis=nodes(deviation).find(n=>n.tagName==='line');
  assert(Number(rect.attributes.x)<Number(axis.attributes.x1));
  assert.equal(calls.length,3);assert(calls.every(url=>url.includes('/metrics/')));
  e.sprintSelect.value='3';await e.sprintSelect.listeners.change();
  assert.match(text(e.dashboard.sprintTitle),/Original goal/);
  assert.equal(nodes(e.dashboard.sprintCharts).filter(n=>n.tagName==='figure').length,5);
});
test('US23-05/07: clear every Sprint result on change, discard late responses and retry errors',async()=>{
  let fail=false,pendingResolve;
  const e=setup(async url=>{
    if(url.includes('/sprints/4/'))return new Promise(resolve=>{pendingResolve??=[];pendingResolve.push(()=>resolve(response(url)))});
    if(fail)throw new Error('Sin conexión');return response(url);
  });
  await e.ui.showProject('p',project);e.sprintSelect.value='3';await e.sprintSelect.listeners.change();
  e.sprintSelect.value='4';const pending=e.sprintSelect.listeners.change();
  for(const target of [e.sprintMetrics,e.sprintHours,e.sprintDefects,e.dashboard.sprintCharts])assert.equal(target.children.length,0);
  e.sprintSelect.value='';await e.sprintSelect.listeners.change();pendingResolve.forEach(fn=>fn());await pending;
  assert.equal(e.dashboard.sprintCharts.children.length,0);
  fail=true;await e.ui.showProject('p',project);
  assert.match(e.message.textContent,/Sin conexión/);assert.equal(e.dashboard.projectCharts.children.length,0);
  fail=false;await e.refresh.listeners.click();assert.equal(nodes(e.dashboard.projectCharts).filter(n=>n.tagName==='figure').length,5);
});
test('US23-05/06: hide absent project, keep empty Sprint history and ignore old project results',async()=>{
  const resolves=[];const e=setup(url=>url.includes('/old/')?new Promise(r=>resolves.push(()=>r(response(url)))):Promise.resolve({...response(url),sprints:[]}));
  const old=e.ui.showProject('old',{...project,id:'old'});await e.ui.showProject('p',project);
  resolves.forEach(fn=>fn());await old;
  assert.match(text(e.dashboard.identity),/En progreso/);assert.equal(e.sprintSelect.disabled,true);
  assert.equal(e.dashboard.sprintCharts.children.length,0);
  await e.ui.showProject('');assert.equal(e.dashboard.section.hidden,true);assert.equal(e.dashboard.projectCharts.children.length,0);
});
