const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');const vm=require('node:vm');const path=require('node:path');
class Element {
 constructor(tag='div'){this.tagName=tag;this.children=[];this.listeners={};this.value='';this.textContent='';this.disabled=false;}
 append(...nodes){this.children.push(...nodes)} replaceChildren(...nodes){this.children=[...nodes]}
 addEventListener(event,fn){this.listeners[event]=fn} setAttribute(){} focus(){this.focused=true}
}
function setup(requestJSON,confirmDelete=()=>true){
 const c=vm.createContext({document:{createElement:tag=>new Element(tag)},encodeURIComponent});
 vm.runInContext(fs.readFileSync(path.join(__dirname,'../../web/static/js/effort.js'),'utf8'),c);
 let changes=0;const ui=c.createEffort({requestJSON,confirmDelete,onChanged:()=>{changes++}});
 const widget=ui.createWidget({id:'p'}, {id:'a'});
 return {ui,widget,c,changed:()=>changes};
}
function all(e){return [e,...e.children.flatMap(all)]}
function button(e,text){return all(e).find(n=>n.tagName==='button' && n.textContent===text)}
function fields(e){return Object.fromEntries(all(e).filter(n=>n.name).map(n=>[n.name,n]))}
function fakeAPI({blocked=false,records=[]}={}){
 const calls=[];
 const api=async(url,options={})=>{
  calls.push({url,...options});
  if(options.method)return {};
  if(url.endsWith('/effort-context'))return {today:'2026-10-26',sprints:[{id:1,closed:true}]};
  if(url.endsWith('/effort'))return records;
  if(url.endsWith('/members'))return [{id:'m',name:'<b>Ana</b>'},{id:'n',name:'<b>Ana</b>'}];
  if(url.endsWith('/backlog'))return [{id:'a',title:'A'},{id:'b',title:'B'}];
  if(url.endsWith('/projects/p'))return {id:'p',status:blocked?'Finalizado':'En progreso',end_date:null};
  return {id:'a',estimated_hours:null,actual_hours:2};
 };return {api,calls};
}
test('US16 exact decimal payloads do not round cents or huge estimates',()=>{
 const {c}=setup(async()=>{});
 assert.equal(c.exactHoursJSON({estimated_hours:'92233720368547758.07',title:'A'},['estimated_hours']),'{"estimated_hours":92233720368547758.07,"title":"A"}');
 assert.throws(()=>c.exactHoursJSON({hours:'1.001'},['hours']));
});
test('BDD 05/07/11/13: detail, missing estimate, server date, homonym UUID and effort submit',async()=>{
 const api=fakeAPI();const e=setup(api.api);await button(e.widget,'Registrar y consultar esfuerzo').listeners.click();
 assert(all(e.widget).some(n=>n.textContent.includes('Pendiente de regularización')));
 let f=fields(e.widget);assert.equal(f.work_date.value,'2026-10-26');assert.equal(f.member_id.children[1].textContent,'<b>Ana</b> — m');
 f.member_id.value='n';f.activity.value='Review';f.hours.value='1.25';
 await all(e.widget).find(n=>n.tagName==='form').listeners.submit({preventDefault(){}});
 const post=api.calls.find(c=>c.method==='POST');assert.equal(JSON.parse(post.body).member_id,'n');assert.equal(JSON.parse(post.body).hours,1.25);assert.equal(e.changed(),1);
});
test('BDD 24/25/43/44: unknown end blocks writes, keeps consultation and confirmed delete',async()=>{
 const api=fakeAPI({blocked:true,records:[{id:'e',member_id:'m',activity:'Review',hours:2,work_date:'2026-10-20',backlog_item_id:'a',sprint_id:null}]});
 const e=setup(api.api,()=>false);await button(e.widget,'Registrar y consultar esfuerzo').listeners.click();
 assert.equal(button(e.widget,'Guardar esfuerzo').disabled,true);assert.equal(button(e.widget,'Editar esfuerzo').disabled,true);
 await button(e.widget,'Eliminar esfuerzo').listeners.click();assert(!api.calls.some(c=>c.method));
 const yes=setup(api.api,()=>true);await button(yes.widget,'Registrar y consultar esfuerzo').listeners.click();await button(yes.widget,'Eliminar esfuerzo').listeners.click();
 assert(api.calls.some(c=>c.url.endsWith('/delete') && JSON.parse(c.body).confirmed===true));
});
test('BDD 26/46: retain compatible Sprint, require explicit incompatible resolution',async()=>{
 const api=fakeAPI({records:[{id:'e',member_id:'m',activity:'Review',hours:2,work_date:'2026-10-20',backlog_item_id:'a',sprint_id:1}]});let compatible=true;
 const e=setup(async(url,opts)=>url.includes('/b/effort-context')?{today:'2026-10-26',sprints:compatible?[{id:1,closed:true}]:[{id:2,closed:false}]}:api.api(url,opts));
 await button(e.widget,'Registrar y consultar esfuerzo').listeners.click();await button(e.widget,'Editar esfuerzo').listeners.click();
 let f=fields(e.widget);f.backlog_item_id.value='b';await f.backlog_item_id.listeners.change();assert.equal(f.sprint_id.value,'1');
 compatible=false;await f.backlog_item_id.listeners.change();assert.equal(f.sprint_id.value,'__resolve__');assert.equal(button(e.widget,'Guardar esfuerzo').disabled,true);
 f.sprint_id.value='';f.sprint_id.listeners.change();assert.equal(button(e.widget,'Guardar esfuerzo').disabled,false);
});
test('BDD 39: stale responses cannot populate changed project',async()=>{
 let resolve;const pending=new Promise(r=>resolve=r);const e=setup(()=>pending);
 const loading=button(e.widget,'Registrar y consultar esfuerzo').listeners.click();e.ui.invalidate();resolve({});await loading;
 assert.equal(all(e.widget).some(n=>n.tagName==='form'),false);
});
