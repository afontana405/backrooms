// ─── Shared Helpers ──────────────────────────────────────────────────────────
// Used across all game JS files. Defined once here, available globally in the IIFE.
function escH(s){return s?s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;'):'';}
function shuffle(a){var b=a.slice();for(var i=b.length-1;i>0;i--){var j=Math.floor(Math.random()*(i+1));var t=b[i];b[i]=b[j];b[j]=t;}return b;}

// ─── Flight Recorder (client ring buffer) ────────────────────────────────────
window.__flight=window.__flight||[];window.__flightCap=500;window.__flightSeq=window.__flightSeq||0;
window.flightRecord=function(c,l,m,d){
  window.__flight.push({id:++window.__flightSeq,time:new Date().toISOString(),
    src:'client',cat:c,level:l,msg:m,detail:d||''});
  if(window.__flight.length>window.__flightCap)window.__flight.shift();
};

// ─── Global Event Instrumentation ────────────────────────────────────────────
// Delegated listeners capture all interactive events for flight recorder.
(function(){
  function elId(el){return el.getAttribute('data-id')||el.id||el.tagName.toLowerCase();}
  function closest(el,sel){while(el&&el!==document){if(el.matches&&el.matches(sel))return el;el=el.parentElement;}return null;}

  // Click — buttons, links, switches, any [on\:click], any [onclick]
  document.addEventListener('click',function(e){
    var btn=closest(e.target,'.cs-button,button,[on\\:click],[onclick],a[href],.cs-switch');
    if(!btn)return;
    var label=btn.textContent.trim().substring(0,60);
    var action=btn.getAttribute('on:click')||btn.getAttribute('onclick')||'';
    var tag=btn.tagName.toLowerCase();
    if(btn.classList.contains('cs-switch')){
      var inp=btn.querySelector('input');
      flightRecord('event',0,'click:switch '+elId(btn),'checked='+(inp?inp.checked:'?'));
    } else if(tag==='a'){
      flightRecord('event',0,'click:link '+elId(btn),'href='+(btn.getAttribute('href')||''));
    } else {
      flightRecord('event',0,'click:button '+elId(btn),(action?'action='+action+' ':'')+'label='+label);
    }
  },true);

  // Change — selects, checkboxes, radios, switches, inputs
  document.addEventListener('change',function(e){
    var el=e.target;
    var tag=el.tagName.toLowerCase();
    var type=el.type||'';
    var sw=closest(el,'.cs-switch');
    if(sw){return;} // already captured by click
    if(tag==='select'||type==='checkbox'||type==='radio'){
      flightRecord('event',0,'change:'+type+' '+elId(el),'value='+(el.type==='checkbox'?el.checked:el.value));
    } else if(tag==='input'||tag==='textarea'){
      flightRecord('event',0,'change:input '+elId(el),'name='+(el.name||''));
    }
  },true);

  // Submit
  document.addEventListener('submit',function(e){
    var form=e.target;
    flightRecord('event',0,'form:submit '+elId(form),'');
  },true);

  // Keydown — Enter and Escape on inputs
  document.addEventListener('keydown',function(e){
    if(e.key==='Enter'&&e.target.tagName==='INPUT'){
      var action=e.target.getAttribute('data-on-enter')||e.target.getAttribute('on:enter')||'';
      if(action) flightRecord('event',0,'key:enter '+elId(e.target),'action='+action);
    }
    if(e.key==='Escape'){
      flightRecord('event',0,'key:escape '+elId(e.target),'');
    }
  },true);

  // Modal/drawer open/close — observe attribute changes on overlay elements
  var overlayObs=new MutationObserver(function(muts){
    for(var i=0;i<muts.length;i++){
      var m=muts[i];
      if(m.type==='attributes'&&(m.attributeName==='class'||m.attributeName==='style')){
        var el=m.target;
        if(el.classList.contains('cs-modal')){
          var open=el.classList.contains('cs-modal--open');
          flightRecord('event',0,(open?'modal:open ':'modal:close ')+elId(el),'');
        }
        if(el.classList.contains('cs-drawer')){
          var open=el.classList.contains('cs-drawer--open');
          flightRecord('event',0,(open?'drawer:open ':'drawer:close ')+elId(el),'');
        }
      }
    }
  });
  document.addEventListener('DOMContentLoaded',function(){
    document.querySelectorAll('.cs-modal,.cs-drawer').forEach(function(el){
      overlayObs.observe(el,{attributes:true,attributeFilter:['class','style']});
    });
  });
  // Also observe newly added modals/drawers after partial nav
  new MutationObserver(function(){
    document.querySelectorAll('.cs-modal,.cs-drawer').forEach(function(el){
      if(!el.__flightObs){el.__flightObs=true;overlayObs.observe(el,{attributes:true,attributeFilter:['class','style']});}
    });
  }).observe(document.body||document.documentElement,{childList:true,subtree:true});
})();

