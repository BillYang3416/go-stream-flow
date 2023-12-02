import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileFlowComponent } from './file-flow.component';
import { FileFlowRoutingModule } from './file-flow-routing.module';

@NgModule({
  declarations: [FileFlowComponent],
  imports: [CommonModule, FileFlowRoutingModule],
})
export class FileFlowModule {}
