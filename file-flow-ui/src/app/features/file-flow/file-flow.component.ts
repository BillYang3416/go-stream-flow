import { Component } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';

@Component({
  selector: 'app-file-flow',
  templateUrl: './file-flow.component.html',
  styleUrls: ['./file-flow.component.scss'],
})
export class FileFlowComponent {
  private selectedFile!: File;

  emailRecipient = new FormControl('', [Validators.required, Validators.email]);

  fileName = '';

  constructor(private apiSvc: ApiService) {}

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
  }

  onSubmit() {
    if (!this.fileName) {
      alert('Please select a file');
      return;
    }

    const formData = new FormData();
    formData.append('file', this.selectedFile);
    formData.append('emailRecipient', this.emailRecipient.value!);

    this.apiSvc.post(ApiPrefix.USER_UPLOADED_FILES, '', formData).subscribe({
      next: (_) => {
        this.onReset();
      },
      error: (err) => {
        alert('Error uploading file' + err.message);
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
