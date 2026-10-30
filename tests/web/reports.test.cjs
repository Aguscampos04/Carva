const {test}=require('node:test');const assert=require('node:assert/strict');
const fs=require('node:fs'),vm=require('node:vm'),path=require('node:path');
const {Element}=require('./app-harness.cjs');
const sample={scope:'project',sprints:[{sprint:{id:1,goal:'Objetivo histórico',closed:true}}],sections:[{title:'Métricas',lines:['Planificados: No disponible','Completados: 0','<img src=x>']},{title:'Defectos',lines:[]}]};
function setup(requestJSON,requestPDF=async()=>({pdf:true})){
 const c=vm.createContext({document:{createElement:tag=>new Element(tag)},encodeURIComponent});
 vm.runInContext(fs.readFileSync(path.join(__dirname,'../../web/static/js/reports.js'),'utf8'),c);
 const e=Object.fromEntries(['select','refresh','download','message','content'].map(k=>[k,new Element()]));
 const saved=[];return {...e,saved,ui:c.createReports({...e,requestJSON,requestPDF,savePDF:(blob,name)=>saved.push({blob,name})})};
}
const all=e=>[e,...e.children.flatMap(all)];
test('US24: show server sections safely, scope selection and download PDF',async()=>{
 const calls=[];const e=setup(async url=>{calls.push(url);return sample},async url=>{calls.push(url);return 'PDF'});
 await e.ui.showProject({id:'p/id'});
 assert(all(e.content).some(n=>n.textContent==='<img src=x>'));
 assert(all(e.content).some(n=>n.textContent==='Sin registros.'));
 assert.equal(e.select.children[1].textContent,'Sprint 1 — Objetivo histórico (cerrado)');
 e.select.value='1';await e.select.listeners.change();await e.download.listeners.click();
 assert(calls.includes('/api/projects/p%2Fid/reports/sprints/1'));
 assert(calls.includes('/api/projects/p%2Fid/reports/sprints/1.pdf'));
 assert.equal(e.saved[0].blob,'PDF');assert.equal(e.saved[0].name,'carva-sprint-1.pdf');
});
test('US24: errors clear stale results, can retry and block duplicate PDF requests',async()=>{
 let fail=false,resolve;const e=setup(async()=>{if(fail)throw new Error('No disponible');return sample},()=>new Promise(r=>{resolve=r}));
 await e.ui.showProject({id:'p'});const pending=e.download.listeners.click();await e.download.listeners.click();
 assert.equal(e.download.disabled,true);resolve('PDF');await pending;assert.equal(e.saved.length,1);
 fail=true;await e.refresh.listeners.click();assert.equal(e.content.children.length,0);assert.equal(e.download.disabled,true);
 assert.match(e.message.textContent,/No disponible/);fail=false;await e.refresh.listeners.click();assert.equal(e.download.disabled,false);
});
test('US24: late JSON and PDF responses cannot populate another route',async()=>{
 let resolveRead,resolvePDF;const e=setup(url=>url.includes('/old/')?new Promise(r=>{resolveRead=r}):Promise.resolve(sample),()=>new Promise(r=>{resolvePDF=r}));
 const old=e.ui.showProject({id:'old'});await e.ui.showProject({id:'new'});resolveRead(sample);await old;
 const pdf=e.download.listeners.click();e.ui.invalidate();resolvePDF('PDF');await pdf;
 assert.equal(e.saved.length,0);assert.equal(e.content.children.length,0);
});
test('US24: failed PDF is visible and does not save a false file',async()=>{
 const e=setup(async()=>sample,async()=>{throw new Error('Error al generar PDF')});await e.ui.showProject({id:'p'});
 await e.download.listeners.click();assert.equal(e.saved.length,0);assert.match(e.message.textContent,/Error al generar PDF/);assert.equal(e.download.disabled,false);
});
