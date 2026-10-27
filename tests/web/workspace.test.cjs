const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const context = vm.createContext({});
vm.runInContext(fs.readFileSync(require('node:path').join(__dirname, '../../web/static/js/workspace.js'), 'utf8'), context);
const ui = vm.runInContext('CarvaWorkspace', context);
const id = '9544ff7b-81cb-4d80-b5c9-4ab23aa3e051';

test('UX-01: round trip of project, view and selected story; rejects invalid locations', () => {
  const route = ui.parseRoute(ui.route(id, 'backlog', 'story-1', 'effort'));
  assert.equal(route.project, id);
  assert.equal(route.view, 'backlog');
  assert.equal(route.story, 'story-1');
  assert.equal(route.tab, 'effort');
  assert.equal(ui.parseRoute('#/proyecto/x/no-existe'), null);
  assert.equal(ui.parseRoute('#/%E0%A4%A'), null);
  assert.equal(ui.parseRoute('#/proyecto/x/resumen/sobra'), null);
  assert.equal(ui.parseRoute('').view, 'inicio');
});

test('UX-02: accent insensitive combined filters preserve original data and order', () => {
  const stories = [{title:'Gestión de equipo', priority:'Alta', status:'Pendiente'}, {title:'Gestión de Sprint', priority:'Baja', status:'Completada'}];
  const before = JSON.stringify(stories);
  assert.equal(ui.filterStories(stories, 'gestion', 'Alta', 'Pendiente').length, 1);
  assert.equal(ui.filterStories(stories, '', '', 'Completada')[0].title, 'Gestión de Sprint');
  assert.equal(ui.filterStories(stories, 'ausente', '', '').length, 0);
  assert.equal(JSON.stringify(stories), before);
});

test('UX-04: newer navigation invalidates requests including an interrupted request', () => {
  const gate = ui.createGate();
  const a = gate.next();
  const b = gate.next();
  assert.equal(gate.current(a), false);
  assert.equal(gate.current(b), true);
  gate.next();
  assert.equal(gate.current(b), false);
});

test('UX-06: distinguish provisional, agreed zero and missing Story Points', () => {
  assert.equal(ui.points({story_points:0, provisional_story_points:5}), '0 SP acordados');
  assert.equal(ui.points({story_points:null, provisional_story_points:0}), '0 SP provisionales');
  assert.equal(ui.points({story_points:null, provisional_story_points:null}), 'Sin estimar');
});

test('UX-04: full app discards late project navigation responses', async () => {
  const {createApp}=require('./app-harness.cjs');
  let resolveOld;
  const response=data=>({ok:true,status:200,text:async()=>JSON.stringify(data)});
  const e=createApp(url=>url.includes('/old')?new Promise(resolve=>{resolveOld=()=>resolve(response({id:'old',name:'Old'}));}):response({id:'new',name:'New',status:'Planificado'}));
  e.location.hash='#/proyecto/old/datos';const old=e.run('navigate()');
  e.location.hash='#/proyecto/new/datos';await e.run('navigate()');
  resolveOld();await old;
  assert.equal(e.get('#project-name-result').textContent,'New');
});

test('UX-05: cancelling draft discard preserves location and draft', () => {
  const {createApp,Element}=require('./app-harness.cjs');
  const e=createApp(()=>{}),draft=new Element('form');
  e.context.draft=draft;e.run('dirtyForms.add(draft)');e.window.confirm=()=>false;
  e.run("go('#/proyecto/new/datos')");assert.equal(e.location.hash,'#/');assert.equal(e.run('dirtyForms.size'),1);
});

test('Home lists saved projects with direct links and safe titles',async()=>{
  const {createApp}=require('./app-harness.cjs');
  const saved=[{id:'one',name:'<b>Alpha</b>',status:'Planificado',start_date:null},{id:'two',name:'Beta',status:'Finalizado',start_date:'2026-01-01'}];
  const e=createApp(async()=>({ok:true,status:200,text:async()=>JSON.stringify(saved)}));
  await e.run('navigate()');
  const cards=e.get('#projects-list').children;
  assert.equal(cards.length,2);
  assert.equal(cards[0].querySelector('h3').children[0].textContent,'<b>Alpha</b>');
  assert.equal(cards[0].querySelector('h3').children[0].href,'#/proyecto/one/resumen');
  assert.match(e.get('#projects-message').textContent,/2 proyectos/);
});

test('Completed stories sort last without modifying the source or other relative order',()=>{
  const items=[{id:'a',status:'Completada'},{id:'b',status:'Pendiente'},{id:'c',status:'En progreso'},{id:'d',acceptance_criteria_verified:true}];
  assert.equal(ui.pendingFirst(items).map(n=>n.id).join(','),'b,c,a,d');
  assert.equal(items.map(n=>n.id).join(','),'a,b,c,d');
});
