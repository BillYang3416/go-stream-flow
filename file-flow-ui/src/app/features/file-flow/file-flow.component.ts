import { Component } from '@angular/core';

@Component({
  selector: 'app-file-flow',
  templateUrl: './file-flow.component.html',
  styleUrls: ['./file-flow.component.scss'],
})
export class FileFlowComponent {
  fileName = '';

  onFileSelected(event: any): void {
    const file: File = event.target.files[0];
    if (file) {
      this.fileName = file.name;
    }
  }
}
