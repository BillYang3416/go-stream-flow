import { FileStatus } from './file-status.model';

export interface GetFileStatusResponse {
  files: FileStatus[];
  totalRecords: number;
}
