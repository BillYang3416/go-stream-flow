import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { FileFlowComponent } from './file-flow.component';

const routes: Routes = [{ path: '', component: FileFlowComponent }];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class FileFlowRoutingModule {}
