import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileStatusComponent } from './file-status.component';
import { FileStatusRoutingModule } from './file-status-routing.module';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { FileSizePipe } from './file-size.pipe';

@NgModule({
  declarations: [FileStatusComponent, FileSizePipe],
  imports: [
    CommonModule,
    FileStatusRoutingModule,
    MatCardModule,
    MatTableModule,
    MatPaginatorModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatIconModule,
  ],
})
export class FileStatusModule {}
