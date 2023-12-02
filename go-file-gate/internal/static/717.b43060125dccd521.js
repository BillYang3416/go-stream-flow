"use strict";(self.webpackChunkfile_flow_ui=self.webpackChunkfile_flow_ui||[]).push([[717],{6717:(p,s,a)=>{a.r(s),a.d(s,{HomeModule:()=>u});var c=a(6814),r=a(8642),t=a(5879),n=a(5195);const l=[{path:"",component:(()=>{class e{static#t=this.\u0275fac=function(o){return new(o||e)};static#e=this.\u0275cmp=t.Xpm({type:e,selectors:[["app-home"]],decls:46,vars:0,consts:[[1,"home-container"]],template:function(o,h){1&o&&(t.TgZ(0,"div",0)(1,"mat-card")(2,"mat-card-header")(3,"mat-card-title"),t._uU(4,"GoStreamFlow Components"),t.qZA()(),t.TgZ(5,"mat-card-content")(6,"h2"),t._uU(7,"1. Frontend App - FileFlowUI"),t.qZA(),t.TgZ(8,"p"),t._uU(9," A sleek Angular application where users upload files and enter email information. It interacts with the backend to initiate the file and email processing workflow. "),t.qZA(),t.TgZ(10,"h2"),t._uU(11,"2. First Golang Application - GoFileGate"),t.qZA(),t.TgZ(12,"p"),t._uU(13," This app acts as the primary gateway for processing user uploads. It receives files and email data from FileFlowUI, stores the records, and then sends this information to a RabbitMQ broker. "),t.qZA(),t.TgZ(14,"h2"),t._uU(15,"3. Second Golang Application - GoMailDistributor"),t.qZA(),t.TgZ(16,"p"),t._uU(17," This application is responsible for consuming messages from the RabbitMQ broker. It sends the file to the specified email address and posts success or failure messages to another queue on the broker. "),t.qZA()()(),t.TgZ(18,"mat-card")(19,"mat-card-header")(20,"mat-card-title"),t._uU(21,"Workflow"),t.qZA()(),t.TgZ(22,"mat-card-content")(23,"p"),t._uU(24," The user uploads a file and inputs email details on FileFlowUI. FileFlowUI sends this data to GoFileGate, which records it and forwards it to RabbitMQ. GoMailDistributor picks up the message, processes it, and updates the job status on another RabbitMQ queue. GoFileGate receives the status update and modifies the record accordingly. The user is then able to see the status of their file upload and email dispatch on FileFlowUI. "),t.qZA()()(),t.TgZ(25,"mat-card")(26,"mat-card-header")(27,"mat-card-title"),t._uU(28,"Why These Names?"),t.qZA()(),t.TgZ(29,"mat-card-content")(30,"p")(31,"strong"),t._uU(32,"GoStreamFlow:"),t.qZA(),t._uU(33," Encapsulates the core functionality of the project \u2014 streaming and processing data using Golang. "),t.qZA(),t.TgZ(34,"p")(35,"strong"),t._uU(36,"FileFlowUI:"),t.qZA(),t._uU(37," Reflects the user-facing aspect of the project, emphasizing file uploads and user interactions. "),t.qZA(),t.TgZ(38,"p")(39,"strong"),t._uU(40,"GoFileGate:"),t.qZA(),t._uU(41," Suggests a gateway role, handling initial data intake and queuing. "),t.qZA(),t.TgZ(42,"p")(43,"strong"),t._uU(44,"GoMailDistributor:"),t.qZA(),t._uU(45," Accurately describes the app's role in distributing emails and handling message queue responses. "),t.qZA()()()())},dependencies:[n.a8,n.dn,n.dk,n.n5],styles:[".home-container[_ngcontent-%COMP%]{display:flex;flex-direction:column;gap:20px;padding:20px}.home-container[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]{box-shadow:0 4px 8px #0003}.home-container[_ngcontent-%COMP%]   mat-card-header[_ngcontent-%COMP%]{background-color:#f5f5f5}.home-container[_ngcontent-%COMP%]   mat-card-title[_ngcontent-%COMP%]{color:#333}.home-container[_ngcontent-%COMP%]   mat-card-content[_ngcontent-%COMP%]{padding:20px}h2[_ngcontent-%COMP%]{font-size:1.2em;color:#005cbf}p[_ngcontent-%COMP%]{font-size:1em;line-height:1.6;color:#666}"]})}return e})()}];let d=(()=>{class e{static#t=this.\u0275fac=function(o){return new(o||e)};static#e=this.\u0275mod=t.oAB({type:e});static#o=this.\u0275inj=t.cJS({imports:[r.Bz.forChild(l),r.Bz]})}return e})(),u=(()=>{class e{static#t=this.\u0275fac=function(o){return new(o||e)};static#e=this.\u0275mod=t.oAB({type:e});static#o=this.\u0275inj=t.cJS({imports:[c.ez,d,n.QW]})}return e})()}}]);