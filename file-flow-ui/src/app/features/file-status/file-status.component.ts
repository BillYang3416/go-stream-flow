import { HttpErrorResponse } from '@angular/common/http';
import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ViewChild,
} from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
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

  pageIndex = 0;
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
    this.getFiles(this.lastID, this.limit);
  }

  onPageChange(event: PageEvent) {
    const previousPageIndex = event.previousPageIndex || 0;
    if (event.pageIndex === 0) {
      this.lastID = 0;
    } else if (
      event.pageIndex ===
      Math.ceil(this.resultsLength / event.pageSize) - 1
    ) {
      this.lastID = event.pageIndex * this.limit;
    } else if (event.pageIndex > previousPageIndex) {
      this.lastID = this.lastID + this.limit;
    } else if (event.pageIndex < previousPageIndex) {
      this.lastID = this.lastID - this.limit;
    }

    this.toggleSpinner();
    this.getFiles(this.lastID, this.limit);
    this.pageIndex = event.pageIndex;
  }

  private toggleSpinner() {
    this.isLoadingResults = !this.isLoadingResults;
    this.changeDetectorRef.detectChanges();
  }

  private getFiles(lastID: number, limit: number) {
    this.apiSvc
      .get<GetFileStatusResponse>(ApiPrefix.USER_UPLOADED_FILES, '', {
        lastID: lastID,
        limit: limit,
      })
      .subscribe({
        next: (data) => {
          this.resultsLength = data.totalRecords;
          this.dataSource.data = data.files;

          this.toggleSpinner();
        },
        error: (err: HttpErrorResponse) => {
          this.toggleSpinner();
          this.snackBar.open(err.error.message, 'OK', { duration: 2000 });
        },
      });
  }
}
