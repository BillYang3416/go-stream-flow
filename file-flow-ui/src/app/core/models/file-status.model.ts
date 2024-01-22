export interface FileStatus {
  id: number;
  name: string;
  size: number;
  createdAt: string;
  emailSent: boolean;
  emailSentAt: string;
  emailRecipent: string;
  errorMessage: string;
}
