const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');const vm=require('node:vm');const path=require('node:path');
class Element {
 constructor(tag='div'){this.tagName=tag;this.children=[];this.listeners={};this.value='';this.textContent='';this.disabled=false;}
 append(...nodes){this.children.push(...nodes)} replaceChildren(...nodes){this.children=[...nodes]}
 addEventListener(event,fn){this.listeners[event]=fn} setAttribute(){} reset(){for(const node of this.children)if(node.name)node.value=''}
}
function setup(requestJSON){
 const c=vm.createContext({document:{createElement:tag=>new Element(tag)},encodeURIComponent});
 vm.runInContext(fs.readFileSync(path.join(__dirname,'../../web/static/js/defects.js'),'utf8'),c);
 const ui=c.createDefects({requestJSON});
 const widget=ui.createWidget({id:'p'},{id:'a'});
 return {ui,widget,c};
}
function all(e){return [e,...e.children.flatMap(all)]}
function button(e,text){return all(e).find(n=>n.tagName==='button'&&n.textContent===text)}
function fields(e){return Object.fromEntries(all(e).filter(n=>n.name).map(n=>[n.name,n]))}
function fakeAPI(){
 const records=[],calls=[];
 const requestJSON=async(url,options={})=>{
  calls.push({url,...options});
  if(options.method==='POST'){const data=JSON.parse(options.body);records.push({id:'d1',project_id:'p',backlog_item_id:'a',...data});return records.at(-1)}
  return records;
 };
 return {requestJSON,records,calls};
}

const original={id:'d1',description:'Original',severity:'Alta',status:'Abierto',detection_sprint_id:1,resolution_sprint_id:null};
async function editFixture(patch){
 const calls=[];
 const e=setup(async(url,options={})=>{calls.push({url,...options});return options.method==='PATCH'?patch(JSON.parse(options.body)):[{...original}];});
 await button(e.widget,'Registrar y consultar defectos').listeners.click();
 button(e.widget,'Editar defecto').listeners.click();
 const edit=all(e.widget).filter(n=>n.tagName==='form').at(-1);
 return {...e,calls,edit,f:fields(edit)};
}
test('US19: prefill and cancel without changing the defect',async()=>{
 const e=await editFixture(()=>assert.fail('unexpected write'));
 assert.deepEqual(Object.keys(e.f).sort(),['detection_sprint_id','resolution_sprint_id','status']);
 assert.equal(e.f.status.value,'Abierto');assert.equal(e.f.detection_sprint_id.value,'1');
 button(e.edit,'Cancelar').listeners.click();assert.equal(e.calls.length,1);
 assert.equal(all(e.widget).filter(n=>n.tagName==='form').length,1);
});
test('US19: send only changed fields and refresh from the confirmed response',async()=>{
 const e=await editFixture(change=>({...original,...change}));
 e.f.status.value='Resuelto';e.f.resolution_sprint_id.value='2';
 await e.edit.listeners.submit({preventDefault(){}});
 assert.deepEqual(JSON.parse(e.calls[1].body),{status:'Resuelto',resolution_sprint_id:2});
 assert.equal(e.calls[1].url,'/api/projects/p/backlog/a/defects/d1');
 assert(all(e.widget).some(n=>n.textContent==='Resuelto'));
});
test('US19: clear an association explicitly without clearing omitted fields',async()=>{
 const e=await editFixture(change=>({...original,...change}));e.f.detection_sprint_id.value='';
 await e.edit.listeners.submit({preventDefault(){}});
 assert.deepEqual(JSON.parse(e.calls[1].body),{detection_sprint_id:null});
});
test('US19: block duplicate submissions, retain draft on error and ignore obsolete responses',async()=>{
 let reject;const e=await editFixture(()=>new Promise((_,r)=>{reject=r}));
 e.f.status.value='Resuelto';const first=e.edit.listeners.submit({preventDefault(){}});
 await e.edit.listeners.submit({preventDefault(){}});assert.equal(e.calls.length,2);
 reject(new Error('No disponible'));await first;
 assert.equal(e.f.status.value,'Resuelto');assert(all(e.widget).some(n=>n.textContent==='No disponible'));
 const before=all(e.widget).map(n=>n.textContent).join();
 e.ui.invalidate();await e.edit.listeners.submit({preventDefault(){}});
 assert.equal(e.calls.length,2);assert.equal(all(e.widget).map(n=>n.textContent).join(),before);
});
test('US19: a late saved response cannot change an invalidated widget',async()=>{
 let resolve;const e=await editFixture(()=>new Promise(r=>{resolve=r}));
 e.f.status.value='Resuelto';const pending=e.edit.listeners.submit({preventDefault(){}});
 e.ui.invalidate();const before=all(e.widget).map(n=>n.textContent).join();
 resolve({...original,status:'Resuelto'});await pending;
 assert.equal(all(e.widget).map(n=>n.textContent).join(),before);
});
