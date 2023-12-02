import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileFlowComponent } from './file-flow.component';
import { FileFlowRoutingModule } from './file-flow-routing.module';

import { ReactiveFormsModule } from '@angular/forms';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBarModule } from '@angular/material/snack-bar';

@NgModule({
  declarations: [FileFlowComponent],
  imports: [
    CommonModule,
    FileFlowRoutingModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
  ],
})
export class FileFlowModule {}
