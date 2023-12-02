import { Component } from '@angular/core';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';

@Component({
  selector: 'app-file-flow',
  templateUrl: './file-flow.component.html',
  styleUrls: ['./file-flow.component.scss'],
})
export class FileFlowComponent {
  private selectedFile!: File;

  emailRecipient = '';

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
    this.emailRecipient = '';
    this.fileName = '';
    this.selectedFile = undefined!;
  }

  onSubmit() {
    if (!this.emailRecipient || !this.fileName) {
      alert('Please enter an email and select a file');
      return;
    }

    const formData = new FormData();
    formData.append('file', this.selectedFile);
    formData.append('emailRecipient', this.emailRecipient);

    this.apiSvc.post(ApiPrefix.USER_UPLOADED_FILES, '', formData).subscribe({
      next: (_) => {
        this.onReset();
      },
      error: (err) => {
        alert('Error uploading file' + err.message);
      },
    });
  }
}
