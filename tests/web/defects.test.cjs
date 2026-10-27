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
test('US18: register and list a defect with no Sprint by default',async()=>{
 const api=fakeAPI(),e=setup(api.requestJSON);
 await button(e.widget,'Registrar y consultar defectos').listeners.click();
 const f=fields(e.widget);
 assert.equal(f.detection_sprint_id.value,'');assert.equal(f.resolution_sprint_id.value,'');
 f.description.value='No guarda cambios';f.severity.value='Alta';f.status.value='Abierto';
 await all(e.widget).find(n=>n.tagName==='form').listeners.submit({preventDefault(){}});
 const saved=JSON.parse(api.calls.find(c=>c.method==='POST').body);
 assert.equal(saved.backlog_item_id,undefined);
 assert.equal(saved.detection_sprint_id,null);assert.equal(saved.resolution_sprint_id,null);
 assert.equal(api.records.length,1);
});

test('US18: preserve a known detection Sprint and display approved defect details as text',async()=>{
 const api=fakeAPI(),e=setup(api.requestJSON);
 await button(e.widget,'Registrar y consultar defectos').listeners.click();
 const f=fields(e.widget);f.description.value='<img>';f.severity.value='Crítica';f.status.value='En progreso';f.detection_sprint_id.value='12';
 await all(e.widget).find(n=>n.tagName==='form').listeners.submit({preventDefault(){}});
 const saved=JSON.parse(api.calls.find(c=>c.method==='POST').body);
 assert.equal(saved.detection_sprint_id,12);assert.equal(saved.resolution_sprint_id,null);
 assert(all(e.widget).some(n=>n.textContent.includes('<img>')));
});
