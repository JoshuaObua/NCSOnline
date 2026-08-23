import{P as i}from"./index-BZnBENim.js";function x(){return i.get("/api/v1/admin/dashboard")}function S(t={}){return i.get("/api/v1/admin/applications",{params:t})}function _(t){return i.get(`/api/v1/admin/applications/${t}`)}function j(t,n,p=""){return i.post(`/api/v1/admin/applications/${t}/${n}`,{notes:p})}function L(t){return i.post(`/api/v1/admin/applications/${t}/verify-payment`)}function R(t,n=""){return i.post(`/api/v1/admin/applications/${t}/reject-payment`,{notes:n})}function U(t={}){return i.get("/api/v1/applications",{params:t})}function r(t){return i.get(`/api/v1/applications/${t}`)}function D(t={}){return i.get("/api/v1/transactions",{params:t})}function l(t){return String(t||"").toLowerCase().replaceAll("_"," ").replaceAll("-"," ").replace(/\b\w/g,n=>n.toUpperCase())}function g(t){if(!t)return{};if(typeof t=="object")return t;try{return JSON.parse(t)}catch{return{}}}function a(t){return String(t??"").replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function A(t){const p=String(t||"").split("?")[0].split("#")[0];return decodeURIComponent(p.split("/").pop()||"Uploaded file")}function f(t){if(Array.isArray(t))return t.map(f).join(", ")||"-";if(t&&typeof t=="object")return JSON.stringify(t);const n=String(t??"").trim();return n?/(\.pdf|\.png|\.jpe?g|\.webp|\/uploads\/|\/media\/)/i.test(n)?A(n):n:"-"}function $(t){if(!t)return"-";try{return new Intl.DateTimeFormat("en-UG",{dateStyle:"medium",timeStyle:"short"}).format(new Date(t))}catch{return String(t)}}function O(t,n={}){const p=n.source||(t==null?void 0:t.source)||"custom",m=(t==null?void 0:t.template_title)||(t==null?void 0:t.title)||l((t==null?void 0:t.application_type)||(t==null?void 0:t.form_type)||"Application"),s=(t==null?void 0:t.submission_reference)||(t==null?void 0:t.application_reference)||(t==null?void 0:t.reference)||"Pending",o=g((t==null?void 0:t.answers)??(t==null?void 0:t.form_data)),u=Object.entries(o).map(([h,b])=>`
    <tr><th>${a(l(h))}</th><td>${a(f(b))}</td></tr>
  `).join(""),c=`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>${a(m)} - ${a(s)}</title>
  <style>
    body{font-family:Arial,sans-serif;margin:32px;color:#1f2937}
    h1{margin:0 0 4px;font-size:24px}
    .meta{margin:0 0 22px;color:#6b7280}
    table{width:100%;border-collapse:collapse}
    th,td{padding:10px 12px;border:1px solid #e5e7eb;text-align:left;vertical-align:top}
    th{width:32%;background:#f9fafb}
  </style>
</head>
<body>
  <h1>${a(m)}</h1>
  <p class="meta">Reference: ${a(s)} | Source: ${a(l(p))} | Downloaded: ${a($(new Date))}</p>
  <table>
    <tbody>
      <tr><th>Applicant</th><td>${a((t==null?void 0:t.applicant_name)||(t==null?void 0:t.applicant_email)||"Portal user")}</td></tr>
      <tr><th>Status</th><td>${a(l(t==null?void 0:t.status))}</td></tr>
      <tr><th>Payment Status</th><td>${a(l((t==null?void 0:t.payment_status)||"Not required"))}</td></tr>
      ${u||'<tr><td colspan="2">No response data is available.</td></tr>'}
    </tbody>
  </table>
</body>
</html>`,y=new Blob([c],{type:"text/html;charset=utf-8"}),e=document.createElement("a");e.href=URL.createObjectURL(y),e.download=`${String(s).replace(/[^a-z0-9-]+/gi,"_")}_application.html`,document.body.appendChild(e),e.click(),URL.revokeObjectURL(e.href),e.remove()}export{S as a,D as b,_ as c,O as d,r as e,L as f,x as g,R as h,j as i,U as l};
