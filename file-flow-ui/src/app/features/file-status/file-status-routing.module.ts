import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { FileStatusComponent } from './file-status.component';

const routes: Routes = [{ path: '', component: FileStatusComponent }];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class FileStatusRoutingModule {}
