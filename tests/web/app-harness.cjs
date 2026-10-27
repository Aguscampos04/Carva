const fs=require('node:fs'),vm=require('node:vm'),path=require('node:path');
class Element {
  constructor(tag='div'){this.tagName=tag;this.children=[];this.listeners={};this.value='';this.textContent='';this.hidden=false;this.disabled=false;this.dataset={};this.attrs={};this.isConnected=true;this.classList={add(){},remove(){},toggle(){return true;}};}
  append(...nodes){this.children.push(...nodes);}
  replaceChildren(...nodes){this.children=nodes;}
  addEventListener(event,cb){this.listeners[event]=cb;}
  setAttribute(key,value){this.attrs[key]=value;}
  removeAttribute(key){delete this.attrs[key];}
  getAttribute(key){return this.attrs[key];}
  focus(){}
  reset(){}
  closest(){return null;}
  querySelectorAll(selector){return this.children.flatMap(n=>[...(n.tagName===selector?[n]:[]),...n.querySelectorAll(selector)]);}
  querySelector(selector){if(selector.includes('button')||selector.includes('submit')){let b=this.querySelectorAll('button')[0];if(!b){b=new Element('button');this.append(b);}return b;}return this.querySelectorAll(selector)[0]||null;}
}
// The application modules are injected separately so tests can control async reads.
function createApp(fetch){
  const elements=new Map(),get=id=>{if(!elements.has(id))elements.set(id,new Element());return elements.get(id);};
  const pages=['inicio','resumen','backlog','sprints','equipo','datos'].map(view=>{const e=new Element();e.dataset.page=view;return e;});
  const document={activeElement:null,listeners:{},querySelector:get,querySelectorAll:selector=>selector==='[data-page]'?pages:[],createElement:tag=>new Element(tag),addEventListener(event,cb){this.listeners[event]=cb;}};
  const location={hash:'#/'},window={confirm:()=>true,addEventListener(){}};
  const c=vm.createContext({document,window,location,history:{pushState(a,b,hash){location.hash=hash;}},fetch,encodeURIComponent,URL,Intl,console,
    createMembersSection:()=>({showProject(){}}),createMetricsSection:()=>({showProject(){}}),createSprintWorkspace:()=>({invalidate(){},showProject(){}})});
  for(const name of ['workspace','projects','effort','app'])vm.runInContext(fs.readFileSync(path.join(__dirname,`../../web/static/js/${name}.js`),'utf8'),c);
  return {context:c,get,document,location,window,pages,run:code=>vm.runInContext(code,c)};
}
module.exports={Element,createApp};


