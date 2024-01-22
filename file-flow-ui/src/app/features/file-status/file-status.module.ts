import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileStatusComponent } from './file-status.component';
import { FileStatusRoutingModule } from './file-status-routing.module';

@NgModule({
  declarations: [FileStatusComponent],
  imports: [CommonModule, FileStatusRoutingModule],
})
export class FileStatusModule {}
