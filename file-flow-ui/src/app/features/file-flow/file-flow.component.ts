import { HttpErrorResponse } from '@angular/common/http';
import { Component, ElementRef, ViewChild } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';

@Component({
  selector: 'app-file-flow',
  templateUrl: './file-flow.component.html',
  styleUrls: ['./file-flow.component.scss'],
})
export class FileFlowComponent {
  @ViewChild('fileUpload') fileUploadEl!: ElementRef;

  private selectedFile!: File;

  emailRecipient = new FormControl('', [Validators.required, Validators.email]);

  fileName = '';

  constructor(private apiSvc: ApiService, private snackBar: MatSnackBar) {}

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      this.fileName = file.name;
      this.selectedFile = file;
    }
  }

  onReset() {
    this.emailRecipient.reset();
    this.fileName = '';
    this.selectedFile = undefined!;
    this.fileUploadEl.nativeElement.value = '';
  }

  onSubmit() {
    if (this.emailRecipient.invalid) {
      this.snackBar.open('Please enter a valid email', 'OK');
      return;
    }

    if (!this.fileName) {
      this.snackBar.open('Please select a file', 'OK');
      return;
    }

    const formData = new FormData();
    formData.append('file', this.selectedFile);
    formData.append('emailRecipient', this.emailRecipient.value!);

    this.apiSvc.post(ApiPrefix.USER_UPLOADED_FILES, '', formData).subscribe({
      next: (_) => {
        this.onReset();
      },
      error: (err: HttpErrorResponse) => {
        this.snackBar.open(err.error.message);
      },
    });
  }

  getEmailErrorMessage() {
    if (this.emailRecipient.hasError('required')) {
      return 'You must enter a value';
    }

    return this.emailRecipient.hasError('email') ? 'Not a valid email' : '';
  }
}
