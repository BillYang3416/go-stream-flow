import { Component } from '@angular/core';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';

@Component({
  selector: 'app-file-flow',
  templateUrl: './file-flow.component.html',
  styleUrls: ['./file-flow.component.scss'],
})
export class FileFlowComponent {
  private selectedFile!: File;

  email = '';

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
    this.email = '';
    this.fileName = '';
    this.selectedFile = undefined!;
  }

  onSubmit() {
    if (!this.email || !this.fileName) {
      alert('Please enter an email and select a file');
      return;
    }

    const formData = new FormData();
    formData.append('file', this.selectedFile);
    formData.append('email', this.email);

    this.apiSvc.post(ApiPrefix.FILE_FLOW, 'upload', formData).subscribe();
  }
}
