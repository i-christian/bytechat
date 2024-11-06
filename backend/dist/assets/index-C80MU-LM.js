(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const r of document.querySelectorAll('link[rel="modulepreload"]'))s(r);new MutationObserver(r=>{for(const o of r)if(o.type==="childList")for(const l of o.addedNodes)l.tagName==="LINK"&&l.rel==="modulepreload"&&s(l)}).observe(document,{childList:!0,subtree:!0});function n(r){const o={};return r.integrity&&(o.integrity=r.integrity),r.referrerPolicy&&(o.referrerPolicy=r.referrerPolicy),r.crossOrigin==="use-credentials"?o.credentials="include":r.crossOrigin==="anonymous"?o.credentials="omit":o.credentials="same-origin",o}function s(r){if(r.ep)return;r.ep=!0;const o=n(r);fetch(r.href,o)}})();const st=(e,t)=>e===t,xe=Symbol("solid-proxy"),ot=typeof Proxy=="function",it=Symbol("solid-track"),ie={equals:st};let Fe=Ge;const U=1,le=2,Ue={owned:null,cleanups:null,context:null,owner:null};var x=null;let me=null,lt=null,_=null,P=null,N=null,fe=0;function Y(e,t){const n=_,s=x,r=e.length===0,o=t===void 0?s:t,l=r?Ue:{owned:null,cleanups:null,context:o?o.context:null,owner:o},i=r?e:()=>e(()=>j(()=>Z(l)));x=l,_=null;try{return B(i,!0)}finally{_=n,x=s}}function C(e,t){t=t?Object.assign({},ie,t):ie;const n={value:e,observers:null,observerSlots:null,comparator:t.equals||void 0},s=r=>(typeof r=="function"&&(r=r(n.value)),Ke(n,r));return[We.bind(n),s]}function O(e,t,n){const s=Pe(e,t,!1,U);z(s)}function at(e,t,n){Fe=ht;const s=Pe(e,t,!1,U);s.user=!0,N?N.push(s):z(s)}function A(e,t,n){n=n?Object.assign({},ie,n):ie;const s=Pe(e,t,!0,0);return s.observers=null,s.observerSlots=null,s.comparator=n.equals||void 0,z(s),We.bind(s)}function ct(e){return B(e,!1)}function j(e){if(_===null)return e();const t=_;_=null;try{return e()}finally{_=t}}function Ae(e,t,n){const s=Array.isArray(e);let r,o=n&&n.defer;return l=>{let i;if(s){i=Array(e.length);for(let u=0;u<e.length;u++)i[u]=e[u]()}else i=e();if(o)return o=!1,l;const a=j(()=>t(i,r,l));return r=i,a}}function ke(e){at(()=>j(e))}function Ce(e){return x===null||(x.cleanups===null?x.cleanups=[e]:x.cleanups.push(e)),e}function Be(){return x}function De(e,t){const n=x,s=_;x=e,_=null;try{return B(t,!0)}catch(r){Le(r)}finally{x=n,_=s}}function ut(e){const t=_,n=x;return Promise.resolve().then(()=>{_=t,x=n;let s;return B(e,!1),_=x=null,s?s.done:void 0})}function qe(e,t){const n=Symbol("context");return{id:n,Provider:pt(n),defaultValue:e}}function Me(e){let t;return x&&x.context&&(t=x.context[e.id])!==void 0?t:e.defaultValue}function Ee(e){const t=A(e),n=A(()=>ve(t()));return n.toArray=()=>{const s=n();return Array.isArray(s)?s:s!=null?[s]:[]},n}function We(){if(this.sources&&this.state)if(this.state===U)z(this);else{const e=P;P=null,B(()=>ce(this),!1),P=e}if(_){const e=this.observers?this.observers.length:0;_.sources?(_.sources.push(this),_.sourceSlots.push(e)):(_.sources=[this],_.sourceSlots=[e]),this.observers?(this.observers.push(_),this.observerSlots.push(_.sources.length-1)):(this.observers=[_],this.observerSlots=[_.sources.length-1])}return this.value}function Ke(e,t,n){let s=e.value;return(!e.comparator||!e.comparator(s,t))&&(e.value=t,e.observers&&e.observers.length&&B(()=>{for(let r=0;r<e.observers.length;r+=1){const o=e.observers[r],l=me&&me.running;l&&me.disposed.has(o),(l?!o.tState:!o.state)&&(o.pure?P.push(o):N.push(o),o.observers&&Ve(o)),l||(o.state=U)}if(P.length>1e6)throw P=[],new Error},!1)),t}function z(e){if(!e.fn)return;Z(e);const t=fe;ft(e,e.value,t)}function ft(e,t,n){let s;const r=x,o=_;_=x=e;try{s=e.fn(t)}catch(l){return e.pure&&(e.state=U,e.owned&&e.owned.forEach(Z),e.owned=null),e.updatedAt=n+1,Le(l)}finally{_=o,x=r}(!e.updatedAt||e.updatedAt<=n)&&(e.updatedAt!=null&&"observers"in e?Ke(e,s):e.value=s,e.updatedAt=n)}function Pe(e,t,n,s=U,r){const o={fn:e,state:s,updatedAt:null,owned:null,sources:null,sourceSlots:null,cleanups:null,value:t,owner:x,context:x?x.context:null,pure:n};return x===null||x!==Ue&&(x.owned?x.owned.push(o):x.owned=[o]),o}function ae(e){if(e.state===0)return;if(e.state===le)return ce(e);if(e.suspense&&j(e.suspense.inFallback))return e.suspense.effects.push(e);const t=[e];for(;(e=e.owner)&&(!e.updatedAt||e.updatedAt<fe);)e.state&&t.push(e);for(let n=t.length-1;n>=0;n--)if(e=t[n],e.state===U)z(e);else if(e.state===le){const s=P;P=null,B(()=>ce(e,t[0]),!1),P=s}}function B(e,t){if(P)return e();let n=!1;t||(P=[]),N?n=!0:N=[],fe++;try{const s=e();return dt(n),s}catch(s){n||(N=null),P=null,Le(s)}}function dt(e){if(P&&(Ge(P),P=null),e)return;const t=N;N=null,t.length&&B(()=>Fe(t),!1)}function Ge(e){for(let t=0;t<e.length;t++)ae(e[t])}function ht(e){let t,n=0;for(t=0;t<e.length;t++){const s=e[t];s.user?e[n++]=s:ae(s)}for(t=0;t<n;t++)ae(e[t])}function ce(e,t){e.state=0;for(let n=0;n<e.sources.length;n+=1){const s=e.sources[n];if(s.sources){const r=s.state;r===U?s!==t&&(!s.updatedAt||s.updatedAt<fe)&&ae(s):r===le&&ce(s,t)}}}function Ve(e){for(let t=0;t<e.observers.length;t+=1){const n=e.observers[t];n.state||(n.state=le,n.pure?P.push(n):N.push(n),n.observers&&Ve(n))}}function Z(e){let t;if(e.sources)for(;e.sources.length;){const n=e.sources.pop(),s=e.sourceSlots.pop(),r=n.observers;if(r&&r.length){const o=r.pop(),l=n.observerSlots.pop();s<r.length&&(o.sourceSlots[l]=s,r[s]=o,n.observerSlots[s]=l)}}if(e.tOwned){for(t=e.tOwned.length-1;t>=0;t--)Z(e.tOwned[t]);delete e.tOwned}if(e.owned){for(t=e.owned.length-1;t>=0;t--)Z(e.owned[t]);e.owned=null}if(e.cleanups){for(t=e.cleanups.length-1;t>=0;t--)e.cleanups[t]();e.cleanups=null}e.state=0}function gt(e){return e instanceof Error?e:new Error(typeof e=="string"?e:"Unknown error",{cause:e})}function Le(e,t=x){throw gt(e)}function ve(e){if(typeof e=="function"&&!e.length)return ve(e());if(Array.isArray(e)){const t=[];for(let n=0;n<e.length;n++){const s=ve(e[n]);Array.isArray(s)?t.push.apply(t,s):t.push(s)}return t}return e}function pt(e,t){return function(s){let r;return O(()=>r=j(()=>(x.context={...x.context,[e]:s.value},Ee(()=>s.children))),void 0),r}}const mt=Symbol("fallback");function Te(e){for(let t=0;t<e.length;t++)e[t]()}function yt(e,t,n={}){let s=[],r=[],o=[],l=0,i=t.length>1?[]:null;return Ce(()=>Te(o)),()=>{let a=e()||[],u=a.length,f,c;return a[it],j(()=>{let h,w,d,g,m,p,S,$,E;if(u===0)l!==0&&(Te(o),o=[],s=[],r=[],l=0,i&&(i=[])),n.fallback&&(s=[mt],r[0]=Y(I=>(o[0]=I,n.fallback())),l=1);else if(l===0){for(r=new Array(u),c=0;c<u;c++)s[c]=a[c],r[c]=Y(y);l=u}else{for(d=new Array(u),g=new Array(u),i&&(m=new Array(u)),p=0,S=Math.min(l,u);p<S&&s[p]===a[p];p++);for(S=l-1,$=u-1;S>=p&&$>=p&&s[S]===a[$];S--,$--)d[$]=r[S],g[$]=o[S],i&&(m[$]=i[S]);for(h=new Map,w=new Array($+1),c=$;c>=p;c--)E=a[c],f=h.get(E),w[c]=f===void 0?-1:f,h.set(E,c);for(f=p;f<=S;f++)E=s[f],c=h.get(E),c!==void 0&&c!==-1?(d[c]=r[f],g[c]=o[f],i&&(m[c]=i[f]),c=w[c],h.set(E,c)):o[f]();for(c=p;c<u;c++)c in d?(r[c]=d[c],o[c]=g[c],i&&(i[c]=m[c],i[c](c))):r[c]=Y(y);r=r.slice(0,l=u),s=a.slice(0)}return r});function y(h){if(o[c]=h,i){const[w,d]=C(c);return i[c]=d,t(a[c],w)}return t(a[c])}}}function k(e,t){return j(()=>e(t||{}))}function se(){return!0}const bt={get(e,t,n){return t===xe?n:e.get(t)},has(e,t){return t===xe?!0:e.has(t)},set:se,deleteProperty:se,getOwnPropertyDescriptor(e,t){return{configurable:!0,enumerable:!0,get(){return e.get(t)},set:se,deleteProperty:se}},ownKeys(e){return e.keys()}};function ye(e){return(e=typeof e=="function"?e():e)?e:{}}function wt(){for(let e=0,t=this.length;e<t;++e){const n=this[e]();if(n!==void 0)return n}}function xt(...e){let t=!1;for(let l=0;l<e.length;l++){const i=e[l];t=t||!!i&&xe in i,e[l]=typeof i=="function"?(t=!0,A(i)):i}if(ot&&t)return new Proxy({get(l){for(let i=e.length-1;i>=0;i--){const a=ye(e[i])[l];if(a!==void 0)return a}},has(l){for(let i=e.length-1;i>=0;i--)if(l in ye(e[i]))return!0;return!1},keys(){const l=[];for(let i=0;i<e.length;i++)l.push(...Object.keys(ye(e[i])));return[...new Set(l)]}},bt);const n={},s=Object.create(null);for(let l=e.length-1;l>=0;l--){const i=e[l];if(!i)continue;const a=Object.getOwnPropertyNames(i);for(let u=a.length-1;u>=0;u--){const f=a[u];if(f==="__proto__"||f==="constructor")continue;const c=Object.getOwnPropertyDescriptor(i,f);if(!s[f])s[f]=c.get?{enumerable:!0,configurable:!0,get:wt.bind(n[f]=[c.get.bind(i)])}:c.value!==void 0?c:void 0;else{const y=n[f];y&&(c.get?y.push(c.get.bind(i)):c.value!==void 0&&y.push(()=>c.value))}}}const r={},o=Object.keys(s);for(let l=o.length-1;l>=0;l--){const i=o[l],a=s[i];a&&a.get?Object.defineProperty(r,i,a):r[i]=a?a.value:void 0}return r}const vt=e=>`Stale read from <${e}>.`;function $t(e){const t="fallback"in e&&{fallback:()=>e.fallback};return A(yt(()=>e.each,e.children,t||void 0))}function He(e){const t=e.keyed,n=A(()=>e.when,void 0,{equals:(s,r)=>t?s===r:!s==!r});return A(()=>{const s=n();if(s){const r=e.children;return typeof r=="function"&&r.length>0?j(()=>r(t?s:()=>{if(!j(n))throw vt("Show");return e.when})):r}return e.fallback},void 0,void 0)}function St(e,t,n){let s=n.length,r=t.length,o=s,l=0,i=0,a=t[r-1].nextSibling,u=null;for(;l<r||i<o;){if(t[l]===n[i]){l++,i++;continue}for(;t[r-1]===n[o-1];)r--,o--;if(r===l){const f=o<s?i?n[i-1].nextSibling:n[o-i]:a;for(;i<o;)e.insertBefore(n[i++],f)}else if(o===i)for(;l<r;)(!u||!u.has(t[l]))&&t[l].remove(),l++;else if(t[l]===n[o-1]&&n[i]===t[r-1]){const f=t[--r].nextSibling;e.insertBefore(n[i++],t[l++].nextSibling),e.insertBefore(n[--o],f),t[r]=n[o]}else{if(!u){u=new Map;let c=i;for(;c<o;)u.set(n[c],c++)}const f=u.get(t[l]);if(f!=null)if(i<f&&f<o){let c=l,y=1,h;for(;++c<r&&c<o&&!((h=u.get(t[c]))==null||h!==f+y);)y++;if(y>f-i){const w=t[l];for(;i<f;)e.insertBefore(n[i++],w)}else e.replaceChild(n[i++],t[l++])}else l++;else t[l++].remove()}}}const je="_$DX_DELEGATE";function _t(e,t,n,s={}){let r;return Y(o=>{r=o,t===document?e():R(t,e(),t.firstChild?null:void 0,n)},s.owner),()=>{r(),t.textContent=""}}function T(e,t,n){let s;const r=()=>{const l=document.createElement("template");return l.innerHTML=e,l.content.firstChild},o=()=>(s||(s=r())).cloneNode(!0);return o.cloneNode=o,o}function V(e,t=window.document){const n=t[je]||(t[je]=new Set);for(let s=0,r=e.length;s<r;s++){const o=e[s];n.has(o)||(n.add(o),t.addEventListener(o,At))}}function $e(e,t,n){n==null?e.removeAttribute(t):e.setAttribute(t,n)}function R(e,t,n,s){if(n!==void 0&&!s&&(s=[]),typeof t!="function")return ue(e,t,s,n);O(r=>ue(e,t(),r,n),s)}function At(e){let t=e.target;const n=`$$${e.type}`,s=e.target,r=e.currentTarget,o=a=>Object.defineProperty(e,"target",{configurable:!0,value:a}),l=()=>{const a=t[n];if(a&&!t.disabled){const u=t[`${n}Data`];if(u!==void 0?a.call(t,u,e):a.call(t,e),e.cancelBubble)return}return t.host&&typeof t.host!="string"&&!t.host._$host&&t.contains(e.target)&&o(t.host),!0},i=()=>{for(;l()&&(t=t._$host||t.parentNode||t.host););};if(Object.defineProperty(e,"currentTarget",{configurable:!0,get(){return t||document}}),e.composedPath){const a=e.composedPath();o(a[0]);for(let u=0;u<a.length-2&&(t=a[u],!!l());u++){if(t._$host){t=t._$host,i();break}if(t.parentNode===r)break}}else i();o(s)}function ue(e,t,n,s,r){for(;typeof n=="function";)n=n();if(t===n)return n;const o=typeof t,l=s!==void 0;if(e=l&&n[0]&&n[0].parentNode||e,o==="string"||o==="number"){if(o==="number"&&(t=t.toString(),t===n))return n;if(l){let i=n[0];i&&i.nodeType===3?i.data!==t&&(i.data=t):i=document.createTextNode(t),n=G(e,n,s,i)}else n!==""&&typeof n=="string"?n=e.firstChild.data=t:n=e.textContent=t}else if(t==null||o==="boolean")n=G(e,n,s);else{if(o==="function")return O(()=>{let i=t();for(;typeof i=="function";)i=i();n=ue(e,i,n,s)}),()=>n;if(Array.isArray(t)){const i=[],a=n&&Array.isArray(n);if(Se(i,t,n,r))return O(()=>n=ue(e,i,n,s,!0)),()=>n;if(i.length===0){if(n=G(e,n,s),l)return n}else a?n.length===0?Ie(e,i,s):St(e,n,i):(n&&G(e),Ie(e,i));n=i}else if(t.nodeType){if(Array.isArray(n)){if(l)return n=G(e,n,s,t);G(e,n,null,t)}else n==null||n===""||!e.firstChild?e.appendChild(t):e.replaceChild(t,e.firstChild);n=t}}return n}function Se(e,t,n,s){let r=!1;for(let o=0,l=t.length;o<l;o++){let i=t[o],a=n&&n[e.length],u;if(!(i==null||i===!0||i===!1))if((u=typeof i)=="object"&&i.nodeType)e.push(i);else if(Array.isArray(i))r=Se(e,i,a)||r;else if(u==="function")if(s){for(;typeof i=="function";)i=i();r=Se(e,Array.isArray(i)?i:[i],Array.isArray(a)?a:[a])||r}else e.push(i),r=!0;else{const f=String(i);a&&a.nodeType===3&&a.data===f?e.push(a):e.push(document.createTextNode(f))}}return r}function Ie(e,t,n=null){for(let s=0,r=t.length;s<r;s++)e.insertBefore(t[s],n)}function G(e,t,n,s){if(n===void 0)return e.textContent="";const r=s||document.createTextNode("");if(t.length){let o=!1;for(let l=t.length-1;l>=0;l--){const i=t[l];if(r!==i){const a=i.parentNode===e;!o&&!l?a?e.replaceChild(r,i):e.insertBefore(r,n):a&&i.remove()}else o=!0}}else e.insertBefore(r,n);return[r]}const kt=!1;function Je(){let e=new Set;function t(r){return e.add(r),()=>e.delete(r)}let n=!1;function s(r,o){if(n)return!(n=!1);const l={to:r,options:o,defaultPrevented:!1,preventDefault:()=>l.defaultPrevented=!0};for(const i of e)i.listener({...l,from:i.location,retry:a=>{a&&(n=!0),i.navigate(r,{...o,resolve:!1})}});return!l.defaultPrevented}return{subscribe:t,confirm:s}}let _e;function Re(){(!window.history.state||window.history.state._depth==null)&&window.history.replaceState({...window.history.state,_depth:window.history.length-1},""),_e=window.history.state._depth}Re();function Ct(e){return{...e,_depth:window.history.state&&window.history.state._depth}}function Et(e,t){let n=!1;return()=>{const s=_e;Re();const r=s==null?null:_e-s;if(n){n=!1;return}r&&t(r)?(n=!0,window.history.go(-r)):e()}}const Pt=/^(?:[a-z0-9]+:)?\/\//i,Lt=/^\/+|(\/)\/+$/g,Xe="http://sr";function Q(e,t=!1){const n=e.replace(Lt,"$1");return n?t||/^[?#]/.test(n)?n:"/"+n:""}function oe(e,t,n){if(Pt.test(t))return;const s=Q(e),r=n&&Q(n);let o="";return!r||t.startsWith("/")?o=s:r.toLowerCase().indexOf(s.toLowerCase())!==0?o=s+r:o=r,(o||"/")+Q(t,!o)}function Rt(e,t){if(e==null)throw new Error(t);return e}function Ot(e,t){return Q(e).replace(/\/*(\*.*)?$/g,"")+Q(t)}function Ye(e){const t={};return e.searchParams.forEach((n,s)=>{s in t?Array.isArray(t[s])?t[s].push(n):t[s]=[t[s],n]:t[s]=n}),t}function Tt(e,t,n){const[s,r]=e.split("/*",2),o=s.split("/").filter(Boolean),l=o.length;return i=>{const a=i.split("/").filter(Boolean),u=a.length-l;if(u<0||u>0&&r===void 0&&!t)return null;const f={path:l?"":"/",params:{}},c=y=>n===void 0?void 0:n[y];for(let y=0;y<l;y++){const h=o[y],w=h[0]===":",d=w?a[y]:a[y].toLowerCase(),g=w?h.slice(1):h.toLowerCase();if(w&&be(d,c(g)))f.params[g]=d;else if(w||!be(d,g))return null;f.path+=`/${d}`}if(r){const y=u?a.slice(-u).join("/"):"";if(be(y,c(r)))f.params[r]=y;else return null}return f}}function be(e,t){const n=s=>s===e;return t===void 0?!0:typeof t=="string"?n(t):typeof t=="function"?t(e):Array.isArray(t)?t.some(n):t instanceof RegExp?t.test(e):!1}function jt(e){const[t,n]=e.pattern.split("/*",2),s=t.split("/").filter(Boolean);return s.reduce((r,o)=>r+(o.startsWith(":")?2:3),s.length-(n===void 0?0:1))}function Qe(e){const t=new Map,n=Be();return new Proxy({},{get(s,r){return t.has(r)||De(n,()=>t.set(r,A(()=>e()[r]))),t.get(r)()},getOwnPropertyDescriptor(){return{enumerable:!0,configurable:!0}},ownKeys(){return Reflect.ownKeys(e())}})}function Ze(e){let t=/(\/?\:[^\/]+)\?/.exec(e);if(!t)return[e];let n=e.slice(0,t.index),s=e.slice(t.index+t[0].length);const r=[n,n+=t[1]];for(;t=/^(\/\:[^\/]+)\?/.exec(s);)r.push(n+=t[1]),s=s.slice(t[0].length);return Ze(s).reduce((o,l)=>[...o,...r.map(i=>i+l)],[])}const It=100,ze=qe(),et=qe(),tt=()=>Rt(Me(ze),"<A> and 'use' router primitives can be only used inside a Route."),H=()=>tt().navigatorFactory(),Nt=()=>tt().location;function Ft(e,t=""){const{component:n,preload:s,load:r,children:o,info:l}=e,i=!o||Array.isArray(o)&&!o.length,a={key:e,component:n,preload:s||r,info:l};return nt(e.path).reduce((u,f)=>{for(const c of Ze(f)){const y=Ot(t,c);let h=i?y:y.split("/*",1)[0];h=h.split("/").map(w=>w.startsWith(":")||w.startsWith("*")?w:encodeURIComponent(w)).join("/"),u.push({...a,originalPath:f,pattern:h,matcher:Tt(h,!i,e.matchFilters)})}return u},[])}function Ut(e,t=0){return{routes:e,score:jt(e[e.length-1])*1e4-t,matcher(n){const s=[];for(let r=e.length-1;r>=0;r--){const o=e[r],l=o.matcher(n);if(!l)return null;s.unshift({...l,route:o})}return s}}}function nt(e){return Array.isArray(e)?e:[e]}function rt(e,t="",n=[],s=[]){const r=nt(e);for(let o=0,l=r.length;o<l;o++){const i=r[o];if(i&&typeof i=="object"){i.hasOwnProperty("path")||(i.path="");const a=Ft(i,t);for(const u of a){n.push(u);const f=Array.isArray(i.children)&&i.children.length===0;if(i.children&&!f)rt(i.children,u.pattern,n,s);else{const c=Ut([...n],s.length);s.push(c)}n.pop()}}}return n.length?s:s.sort((o,l)=>l.score-o.score)}function we(e,t){for(let n=0,s=e.length;n<s;n++){const r=e[n].matcher(t);if(r)return r}return[]}function Bt(e,t,n){const s=new URL(Xe),r=A(f=>{const c=e();try{return new URL(c,s)}catch{return console.error(`Invalid path ${c}`),f}},s,{equals:(f,c)=>f.href===c.href}),o=A(()=>r().pathname),l=A(()=>r().search,!0),i=A(()=>r().hash),a=()=>"",u=Ae(l,()=>Ye(r()));return{get pathname(){return o()},get search(){return l()},get hash(){return i()},get state(){return t()},get key(){return a()},query:n?n(u):Qe(u)}}let M;function Dt(){return M}function qt(e,t,n,s={}){const{signal:[r,o],utils:l={}}=e,i=l.parsePath||(v=>v),a=l.renderPath||(v=>v),u=l.beforeLeave||Je(),f=oe("",s.base||"");if(f===void 0)throw new Error(`${f} is not a valid base path`);f&&!r().value&&o({value:f,replace:!0,scroll:!1});const[c,y]=C(!1);let h;const w=(v,b)=>{b.value===d()&&b.state===m()||(h===void 0&&y(!0),M=v,h=b,ut(()=>{h===b&&(g(h.value),p(h.state),E[1]([]))}).finally(()=>{h===b&&ct(()=>{M=void 0,v==="navigate"&&he(h),y(!1),h=void 0})}))},[d,g]=C(r().value),[m,p]=C(r().state),S=Bt(d,m,l.queryWrapper),$=[],E=C([]),I=A(()=>typeof s.transformUrl=="function"?we(t(),s.transformUrl(S.pathname)):we(t(),S.pathname)),ee=()=>{const v=I(),b={};for(let L=0;L<v.length;L++)Object.assign(b,v[L].params);return b},te=l.paramsWrapper?l.paramsWrapper(ee,t):Qe(ee),D={pattern:f,path:()=>f,outlet:()=>null,resolvePath(v){return oe(f,v)}};return O(Ae(r,v=>w("native",v),{defer:!0})),{base:D,location:S,params:te,isRouting:c,renderPath:a,parsePath:i,navigatorFactory:de,matches:I,beforeLeave:u,preloadRoute:W,singleFlight:s.singleFlight===void 0?!0:s.singleFlight,submissions:E};function ne(v,b,L){j(()=>{if(typeof b=="number"){b&&(l.go?l.go(b):console.warn("Router integration does not support relative routing"));return}const q=!b||b[0]==="?",{replace:ge,resolve:K,scroll:pe,state:J}={replace:!1,resolve:!q,scroll:!0,...L},re=K?v.resolvePath(b):oe(q&&S.pathname||"",b);if(re===void 0)throw new Error(`Path '${b}' is not a routable path`);if($.length>=It)throw new Error("Too many redirects");const Oe=d();(re!==Oe||J!==m())&&(kt||u.confirm(re,L)&&($.push({value:Oe,replace:ge,scroll:pe,state:m()}),w("navigate",{value:re,state:J})))})}function de(v){return v=v||Me(et)||D,(b,L)=>ne(v,b,L)}function he(v){const b=$[0];b&&(o({...v,replace:b.replace,scroll:b.scroll}),$.length=0)}function W(v,b){const L=we(t(),v.pathname),q=M;M="preload";for(let ge in L){const{route:K,params:pe}=L[ge];K.component&&K.component.preload&&K.component.preload();const{preload:J}=K;b&&J&&De(n(),()=>J({params:pe,location:{pathname:v.pathname,search:v.search,hash:v.hash,query:Ye(v),state:null,key:""},intent:"preload"}))}M=q}}function Mt(e,t,n,s){const{base:r,location:o,params:l}=e,{pattern:i,component:a,preload:u}=s().route,f=A(()=>s().path);a&&a.preload&&a.preload();const c=u?u({params:l,location:o,intent:M||"initial"}):void 0;return{parent:t,pattern:i,path:f,outlet:()=>a?k(a,{params:l,location:o,data:c,get children(){return n()}}):n(),resolvePath(h){return oe(r.path(),h,f())}}}const Wt=e=>t=>{const{base:n}=t,s=Ee(()=>t.children),r=A(()=>rt(s(),t.base||""));let o;const l=qt(e,r,()=>o,{base:n,singleFlight:t.singleFlight,transformUrl:t.transformUrl});return e.create&&e.create(l),k(ze.Provider,{value:l,get children(){return k(Kt,{routerState:l,get root(){return t.root},get preload(){return t.rootPreload||t.rootLoad},get children(){return[A(()=>(o=Be())&&null),k(Gt,{routerState:l,get branches(){return r()}})]}})}})};function Kt(e){const t=e.routerState.location,n=e.routerState.params,s=A(()=>e.preload&&j(()=>{e.preload({params:n,location:t,intent:Dt()||"initial"})}));return k(He,{get when(){return e.root},keyed:!0,get fallback(){return e.children},children:r=>k(r,{params:n,location:t,get data(){return s()},get children(){return e.children}})})}function Gt(e){const t=[];let n;const s=A(Ae(e.routerState.matches,(r,o,l)=>{let i=o&&r.length===o.length;const a=[];for(let u=0,f=r.length;u<f;u++){const c=o&&o[u],y=r[u];l&&c&&y.route.key===c.route.key?a[u]=l[u]:(i=!1,t[u]&&t[u](),Y(h=>{t[u]=h,a[u]=Mt(e.routerState,a[u-1]||e.routerState.base,Ne(()=>s()[u+1]),()=>e.routerState.matches()[u])}))}return t.splice(r.length).forEach(u=>u()),l&&i?l:(n=a[0],a)}));return Ne(()=>s()&&n)()}const Ne=e=>()=>k(He,{get when(){return e()},keyed:!0,children:t=>k(et.Provider,{value:t,get children(){return t.outlet()}})}),X=e=>{const t=Ee(()=>e.children);return xt(e,{get children(){return t()}})};function Vt([e,t],n,s){return[e,s?r=>t(s(r)):t]}function Ht(e){let t=!1;const n=r=>typeof r=="string"?{value:r}:r,s=Vt(C(n(e.get()),{equals:(r,o)=>r.value===o.value&&r.state===o.state}),void 0,r=>(!t&&e.set(r),r));return e.init&&Ce(e.init((r=e.get())=>{t=!0,s[1](n(r)),t=!1})),Wt({signal:s,create:e.create,utils:e.utils})}function Jt(e,t,n){return e.addEventListener(t,n),()=>e.removeEventListener(t,n)}function Xt(e,t){const n=e&&document.getElementById(e);n?n.scrollIntoView():t&&window.scrollTo(0,0)}const Yt=new Map;function Qt(e=!0,t=!1,n="/_server",s){return r=>{const o=r.base.path(),l=r.navigatorFactory(r.base);let i,a;function u(d){return d.namespaceURI==="http://www.w3.org/2000/svg"}function f(d){if(d.defaultPrevented||d.button!==0||d.metaKey||d.altKey||d.ctrlKey||d.shiftKey)return;const g=d.composedPath().find(I=>I instanceof Node&&I.nodeName.toUpperCase()==="A");if(!g||t&&!g.hasAttribute("link"))return;const m=u(g),p=m?g.href.baseVal:g.href;if((m?g.target.baseVal:g.target)||!p&&!g.hasAttribute("state"))return;const $=(g.getAttribute("rel")||"").split(/\s+/);if(g.hasAttribute("download")||$&&$.includes("external"))return;const E=m?new URL(p,document.baseURI):new URL(p);if(!(E.origin!==window.location.origin||o&&E.pathname&&!E.pathname.toLowerCase().startsWith(o.toLowerCase())))return[g,E]}function c(d){const g=f(d);if(!g)return;const[m,p]=g,S=r.parsePath(p.pathname+p.search+p.hash),$=m.getAttribute("state");d.preventDefault(),l(S,{resolve:!1,replace:m.hasAttribute("replace"),scroll:!m.hasAttribute("noscroll"),state:$?JSON.parse($):void 0})}function y(d){const g=f(d);if(!g)return;const[m,p]=g;s&&(p.pathname=s(p.pathname)),r.preloadRoute(p,m.getAttribute("preload")!=="false")}function h(d){clearTimeout(i);const g=f(d);if(!g)return a=null;const[m,p]=g;a!==m&&(s&&(p.pathname=s(p.pathname)),i=setTimeout(()=>{r.preloadRoute(p,m.getAttribute("preload")!=="false"),a=m},20))}function w(d){if(d.defaultPrevented)return;let g=d.submitter&&d.submitter.hasAttribute("formaction")?d.submitter.getAttribute("formaction"):d.target.getAttribute("action");if(!g)return;if(!g.startsWith("https://action/")){const p=new URL(g,Xe);if(g=r.parsePath(p.pathname+p.search),!g.startsWith(n))return}if(d.target.method.toUpperCase()!=="POST")throw new Error("Only POST forms are supported for Actions");const m=Yt.get(g);if(m){d.preventDefault();const p=new FormData(d.target,d.submitter);m.call({r,f:d.target},d.target.enctype==="multipart/form-data"?p:new URLSearchParams(p))}}V(["click","submit"]),document.addEventListener("click",c),e&&(document.addEventListener("mousemove",h,{passive:!0}),document.addEventListener("focusin",y,{passive:!0}),document.addEventListener("touchstart",y,{passive:!0})),document.addEventListener("submit",w),Ce(()=>{document.removeEventListener("click",c),e&&(document.removeEventListener("mousemove",h),document.removeEventListener("focusin",y),document.removeEventListener("touchstart",y)),document.removeEventListener("submit",w)})}}function Zt(e){const t=()=>{const s=window.location.pathname.replace(/^\/+/,"/")+window.location.search,r=window.history.state&&window.history.state._depth&&Object.keys(window.history.state).length===1?void 0:window.history.state;return{value:s+window.location.hash,state:r}},n=Je();return Ht({get:t,set({value:s,replace:r,scroll:o,state:l}){r?window.history.replaceState(Ct(l),"",s):window.history.pushState(l,"",s),Xt(decodeURIComponent(window.location.hash.slice(1)),o),Re()},init:s=>Jt(window,"popstate",Et(s,r=>{if(r&&r<0)return!n.confirm(r);{const o=t();return!n.confirm(o.value,{state:o.state})}})),create:Qt(e.preload,e.explicitLinks,e.actionBase,e.transformUrl),utils:{go:s=>window.history.go(s),beforeLeave:n}})(e)}function zt(e){const t=H(),n=Nt(),{href:s,state:r}=e,o=typeof s=="function"?s({navigate:t,location:n}):s;return t(o,{replace:!0,state:r}),null}var en=T('<div class="flex h-screen"><aside class="w-1/4 bg-gray-200 p-4 border-r border-gray-300"><h2 class="text-xl font-semibold mb-4">Available Rooms</h2><button class="mb-4 p-2 bg-red-500 text-white rounded-lg">Logout</button><ul class=space-y-2></ul></aside><main class="flex-1 flex flex-col p-4">'),tn=T("<p class=text-red-500>"),nn=T('<li class="p-2 bg-white rounded-md shadow cursor-pointer hover:bg-gray-100"><h3 class="text-lg font-medium"></h3><p class="text-sm text-gray-600">'),rn=T('<h2 class="text-2xl font-bold">'),sn=T('<p class="text-sm text-gray-500 mb-4">'),on=T('<div class="flex-1 bg-gray-100 rounded-lg p-4 overflow-y-auto"><p>Chat messages for <!> will appear here...'),ln=T('<div class=mt-4><input type=text placeholder="Type your message here..."class="w-full p-2 border rounded-lg"><button class="mt-2 p-2 bg-blue-600 text-white rounded-lg">Send'),an=T('<div class="flex items-center justify-center h-full"><p class=text-gray-500>Select a room to start chatting');const cn=()=>{const[e,t]=C([]),[n,s]=C(null),[r,o]=C(""),l=H(),i=async()=>{try{const f=await fetch(`//${window.location.host}/api/rooms`,{method:"GET",credentials:"include"});if(f.status===403){F(!1),l("/login");return}if(!f.ok)throw new Error("Failed to load rooms");const c=await f.json();t(c)}catch{o("Failed to load rooms. Please try again later.")}},a=async f=>{try{const c=await fetch(`//${window.location.host}/api/rooms`,{method:"POST",credentials:"include",headers:{"Content-Type":"application/json"},body:JSON.stringify({room_id:f})});if(c.status===403){sessionStorage.removeItem("isLoggedIn"),F(!1),l("/login");return}if(!c.ok)throw new Error("Failed to join room");const y=e().find(h=>h.id===f);s(y||null)}catch{o("Could not join the room. Try again.")}},u=async()=>{try{(await fetch(`//${window.location.host}/api/auth/logout`,{method:"GET",credentials:"include"})).status!==500?(F(!1),l("/")):console.log("Internal server error")}catch(f){console.error("Error logging out:",f)}};return ke(i),(()=>{var f=en(),c=f.firstChild,y=c.firstChild,h=y.nextSibling,w=h.nextSibling,d=c.nextSibling;return R(c,(()=>{var g=A(()=>!!r());return()=>g()&&(()=>{var m=tn();return R(m,r),m})()})(),h),h.$$click=u,R(w,k($t,{get each(){return e()},children:g=>(()=>{var m=nn(),p=m.firstChild,S=p.nextSibling;return m.$$click=()=>a(g.id),R(p,()=>g.name),R(S,()=>g.description),m})()})),R(d,(()=>{var g=A(()=>!!n());return()=>g()?[(()=>{var m=rn();return R(m,()=>n().name),m})(),(()=>{var m=sn();return R(m,()=>n().description),m})(),(()=>{var m=on(),p=m.firstChild,S=p.firstChild,$=S.nextSibling;return $.nextSibling,R(p,()=>n().name,$),m})(),ln()]:an()})()),f})()};V(["click"]);var un=T(`<form class="px-5 min-h-screen flex flex-col items-center justify-center bg-gray-100"><div class="flex flex-col bg-white shadow-md px-4 sm:px-6 md:px-8 lg:px-10 py-8 rounded-md w-50 max-w-md"><h1 class="lg:text-2xl text-xl text-center">Log In</h1><fieldset class=mt-10><label for=email class="text-xs tracking-wide text-gray-600">Email Address:</label><div class="relative mb-4"><input type=email name=email class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"placeholder="Enter your email"></div><label for=password class="text-xs tracking-wide text-gray-600">Password:</label><div class="relative mb-4"><input name=password class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"placeholder="Enter your password"><button class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm"type=button></button></div><div class="flex w-full"><button type=submit class="flex mt-2 items-center justify-center focus:outline-none text-white text-sm sm:text-base bg-black hover:bg-slate-950 rounded-md py-2 w-full transition duration-150 ease-in"><span class="mr-2 uppercase">Log In &rarr;</span></button></div></fieldset></div><div class="flex justify-center items-center mt-6"><span class="inline-flex items-center text-gray-700 font-medium text-xs text-center">Don't have an account?<a href=/register class="text-xs ml-2 text-black font-semibold">Register here`);const fn=()=>{const[e,t]=C(""),[n,s]=C(""),[r,o]=C(!1),l=H();ke(async()=>{try{(await fetch(`//${window.location.host}/api/check`,{method:"GET",credentials:"include"})).ok&&(F(!0),l("/"))}catch(u){console.error("Auth check error:",u)}});const i=async u=>{u.preventDefault();const f=`//${window.location.host}/api/auth/login`;try{(await fetch(f,{method:"POST",mode:"cors",credentials:"include",headers:{"Content-Type":"application/json"},body:JSON.stringify({email:e(),password:n()})})).ok?(F(!0),localStorage.setItem("isLoggedIn","true"),l("/")):alert("Login failed")}catch(c){console.error("Login error:",c),alert("An error occurred while logging in.")}},a=u=>{u.preventDefault(),o(!r())};return(()=>{var u=un(),f=u.firstChild,c=f.firstChild,y=c.nextSibling,h=y.firstChild,w=h.nextSibling,d=w.firstChild,g=w.nextSibling,m=g.nextSibling,p=m.firstChild,S=p.nextSibling;return u.addEventListener("submit",i),d.$$input=$=>t($.target.value),p.$$input=$=>s($.target.value),S.$$click=a,R(S,()=>r()?"Hide":"Show"),O(()=>$e(p,"type",r()?"text":"password")),O(()=>d.value=e()),O(()=>p.value=n()),u})()};V(["input","click"]);var dn=T('<form class="px-5 min-h-screen flex flex-col items-center justify-center bg-gray-100"><div class="flex flex-col bg-white shadow-md px-4 sm:px-6 md:px-8 lg:px-10 py-8 rounded-md w-50 max-w-md"><h1 class="lg:text-2xl text-xl text-center">Register</h1><fieldset class=mt-10><label for=name class="text-xs tracking-wide text-gray-600">Name:</label><div class="relative mb-4"><input type=text name=name class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"required placeholder="Enter your name"></div><label for=email class="text-xs tracking-wide text-gray-600">Email:</label><div class="relative mb-4"><input type=email name=email class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"required placeholder="Enter your email"></div><label for=password class="text-xs tracking-wide text-gray-600">Password:</label><div class="relative mb-4"><input name=password class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"required><button class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm"type=button></button></div><label for=confirm class="text-xs tracking-wide text-gray-600">Confirm Password:</label><div class="relative mb-4"><input name=confirm class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"required><button class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm"type=button></button></div><div class="flex w-full"><button type=submit class="flex mt-2 items-center justify-center focus:outline-none text-white text-sm sm:text-base bg-black hover:bg-slate-950 rounded-md py-2 w-full transition duration-150 ease-in"><span class="mr-2 uppercase">Sign Up &rarr;</span></button></div></fieldset></div><div class="flex justify-center items-center mt-6"><span class="inline-flex items-center text-gray-700 font-medium text-xs text-center">You have an account?<a href=/login class="text-xs ml-2 text-black font-semibold">Login here');function hn(){const[e,t]=C(""),[n,s]=C(""),[r,o]=C(""),[l,i]=C(""),[a,u]=C(!1),f=H(),c=async h=>{h.preventDefault();const w=`//${window.location.host}/api/auth/register`;try{await fetch(w,{method:"POST",mode:"cors",headers:{"Content-Type":"application/json"},body:JSON.stringify({name:e(),email:n(),password:r()})}),f("/login")}catch(d){console.error(`Error: ${d}`)}},y=h=>{h.preventDefault(),u(!a())};return(()=>{var h=dn(),w=h.firstChild,d=w.firstChild,g=d.nextSibling,m=g.firstChild,p=m.nextSibling,S=p.firstChild,$=p.nextSibling,E=$.nextSibling,I=E.firstChild,ee=E.nextSibling,te=ee.nextSibling,D=te.firstChild,ne=D.nextSibling,de=te.nextSibling,he=de.nextSibling,W=he.firstChild,v=W.nextSibling;return h.addEventListener("submit",c),S.$$input=b=>t(b.target.value),I.$$input=b=>s(b.target.value),D.$$input=b=>o(b.target.value),ne.$$click=y,R(ne,()=>a()?"Hide":"Show"),W.$$input=b=>i(b.target.value),v.$$click=y,R(v,()=>a()?"Hide":"Show"),O(b=>{var L=a()?"text":"password",q=a()?"text":"password";return L!==b.e&&$e(D,"type",b.e=L),q!==b.t&&$e(W,"type",b.t=q),b},{e:void 0,t:void 0}),O(()=>S.value=e()),O(()=>I.value=n()),O(()=>D.value=r()),O(()=>W.value=l()),h})()}V(["input","click"]);var gn=T('<section class=bg-inherit><div class="flex sm:flex-wrap items-center min-h-screen px-6 py-12 mx-auto"><div><p class="text-sm font-medium text-red-500 dark:text-red-400">403 Forbidden</p><h1 class="mt-3 text-2xl font-semibold text-gray-800 dark:text-white md:text-3xl">Access Denied</h1><p class="mt-4 text-gray-500 dark:text-gray-400">Sorry, you do not have the necessary permissions to view this page.</p><div class="flex items-center mt-6 gap-x-3"><button class="flex items-center justify-center w-1/2 px-5 py-2 text-sm text-gray-700 transition-colors duration-200 bg-white border rounded-lg gap-x-2 sm:w-auto dark:hover:bg-gray-800 dark:bg-gray-900 hover:bg-gray-100 dark:text-gray-200 dark:border-gray-700"><svg xmlns=http://www.w3.org/2000/svg fill=none viewBox="0 0 24 24"stroke-width=1.5 stroke=currentColor class="w-5 h-5 rtl:rotate-180"><path stroke-linecap=round stroke-linejoin=round d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18"></path></svg><span>Go home</span></button><button class="w-1/2 px-5 py-2 text-sm text-gray-700 transition-colors duration-200 bg-white border rounded-lg sm:w-auto dark:bg-gray-900 hover:bg-gray-100 dark:text-gray-200 dark:border-gray-700">Go Back');const pn=()=>{const e=H();return(()=>{var t=gn(),n=t.firstChild,s=n.firstChild,r=s.firstChild,o=r.nextSibling,l=o.nextSibling,i=l.nextSibling,a=i.firstChild,u=a.nextSibling;return a.$$click=()=>e("/"),u.$$click=()=>e(-1),t})()};V(["click"]);var mn=T(`<section><div class="flex items-center min-h-screen px-6 py-12 mx-auto"><div><p class="text-sm font-medium text-blue-500 dark:text-blue-400">404 error</p><h1 class="mt-3 text-2xl font-semibold text-gray-800 dark:text-white md:text-3xl">We can’t find that page</h1><p class="mt-4 text-gray-500 dark:text-gray-400">Sorry, the page you are looking for doesn't exist or has been moved.</p><div class="flex items-center mt-6 gap-x-3"><button class="flex items-center justify-center w-1/2 px-5 py-2 text-sm text-gray-700 transition-colors duration-200 bg-white border rounded-lg gap-x-2 sm:w-auto dark:hover:bg-gray-800 dark:bg-gray-900 hover:bg-gray-100 dark:text-gray-200 dark:border-gray-700"><svg xmlns=http://www.w3.org/2000/svg fill=none viewBox="0 0 24 24"stroke-width=1.5 stroke=currentColor class="w-5 h-5 rtl:rotate-180"><path stroke-linecap=round stroke-linejoin=round d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18"></path></svg><span>Go home`);const yn=()=>{const e=H();return(()=>{var t=mn(),n=t.firstChild,s=n.firstChild,r=s.firstChild,o=r.nextSibling,l=o.nextSibling,i=l.nextSibling,a=i.firstChild;return a.$$click=()=>e("/"),t})()};V(["click"]);const[bn,F]=C(!1),wn=()=>(ke(async()=>{try{(await fetch(`//${window.location.host}/api/check`,{method:"GET",credentials:"include"})).ok?(F(!0),localStorage.setItem("isLoggedIn","true")):(F(!1),localStorage.removeItem("isLoggedIn"))}catch(e){console.error("Auth check error:",e),F(!1)}}),k(Zt,{get children(){return[k(X,{path:"/",component:()=>bn()?k(cn,{}):k(zt,{href:"/login"})}),k(X,{path:"/register",component:hn}),k(X,{path:"/login",component:fn}),k(X,{path:"/403",component:pn}),k(X,{path:"*",component:yn})]}}));_t(()=>k(wn,{}),document.getElementById("root"));
