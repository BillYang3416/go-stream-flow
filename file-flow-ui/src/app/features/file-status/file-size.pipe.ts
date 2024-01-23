import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'fileSize',
})
export class FileSizePipe implements PipeTransform {
  transform(bytes: number, precision: number = 2): string {
    if (isNaN(parseFloat(String(bytes))) || !isFinite(bytes)) return '?';

    let unit = 0;
    while (bytes >= 1024) {
      bytes /= 1024;
      unit++;
    }

    return (
      bytes.toFixed(+precision) +
      ' ' +
      ['B', 'KB', 'MB', 'GB', 'TB', 'PB'][unit]
    );
  }
}
