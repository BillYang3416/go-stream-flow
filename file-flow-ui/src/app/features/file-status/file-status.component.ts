import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ViewChild,
} from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTableDataSource } from '@angular/material/table';
import { GetFileStatusResponse } from 'src/app/core/models/GetFileStatusResponse';
import { FileStatus } from 'src/app/core/models/file-status.model';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';

@Component({
  selector: 'app-file-status',
  templateUrl: './file-status.component.html',
  styleUrls: ['./file-status.component.scss'],
})
export class FileStatusComponent implements AfterViewInit {
  displayedColumns: string[] = [
    'name',
    'size',
    'createdAt',
    'emailSent',
    'emailRecipient',
    'errorMessage',
  ];
  dataSource = new MatTableDataSource<FileStatus>([]);

  resultsLength = 0;
  isLoadingResults = false;
  limit = 10;
  private lastID = 0;

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  constructor(
    private apiSvc: ApiService,
    private snackBar: MatSnackBar,
    private changeDetectorRef: ChangeDetectorRef
  ) {}

  ngAfterViewInit() {
    this.toggleSpinner();
    this.apiSvc
      .get<GetFileStatusResponse>(ApiPrefix.USER_UPLOADED_FILES, '', {
        lastID: this.lastID,
        limit: this.limit,
      })
      .subscribe({
        next: (data) => {
          this.resultsLength = data.totalRecords;
          this.dataSource.data = data.files;

          if (data.files && data.files.length > 0) {
            this.lastID = data.files[data.files.length - 1].id;
          }
          this.toggleSpinner();
        },
        error: (err) => {
          this.toggleSpinner();
          this.snackBar.open('Failed to load data.', 'OK', { duration: 2000 });
        },
      });
  }

  private toggleSpinner() {
    this.isLoadingResults = !this.isLoadingResults;
    this.changeDetectorRef.detectChanges();
  }
}
