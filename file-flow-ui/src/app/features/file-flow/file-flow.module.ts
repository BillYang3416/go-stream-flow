import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileFlowComponent } from './file-flow.component';
import { FileFlowRoutingModule } from './file-flow-routing.module';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
@NgModule({
  declarations: [FileFlowComponent],
  imports: [
    CommonModule,
    FileFlowRoutingModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
  ],
})
export class FileFlowModule {}
