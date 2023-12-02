import { Injectable } from '@angular/core';
import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
  HttpParams,
} from '@angular/common/http';
import { throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from 'src/enviroments/environment.dev';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  apiSource = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get<T>(
    apiPrefix: ApiPrefix,
    apiName: string,
    oriParams?: { [key: string]: any }
  ) {
    let params = new HttpParams({ fromObject: oriParams });

    return this.http
      .get<T>(this.getUrl(apiPrefix, apiName), {
        params,
      })
      .pipe(catchError((e) => this.handleError(e)));
  }

  post<T>(
    apiPrefix: ApiPrefix,
    apiName: string,
    body: { [key: string]: any } | FormData,
    headers?: HttpHeaders | { [header: string]: string | string[] }
  ) {
    return this.http.post<T>(this.getUrl(apiPrefix, apiName), body, {
      headers,
    });
  }

  patch<T>(
    apiPrefix: ApiPrefix,
    apiName: string,
    body: { [key: string]: any } | FormData,
    headers?: HttpHeaders | { [header: string]: string | string[] }
  ) {
    return this.http.patch<T>(this.getUrl(apiPrefix, apiName), body, {
      headers,
    });
  }

  delete<T>(
    apiPrefix: ApiPrefix,
    apiName: string,
    headers?: HttpHeaders | { [header: string]: string | string[] }
  ) {
    return this.http.delete<T>(this.getUrl(apiPrefix, apiName), { headers });
  }

  put<T>(
    apiPrefix: ApiPrefix,
    apiName: string,
    body: { [key: string]: any } | FormData,
    headers?: HttpHeaders | { [header: string]: string | string[] }
  ) {
    return this.http.put<T>(this.getUrl(apiPrefix, apiName), body, { headers });
  }

  private handleError(err: HttpErrorResponse) {
    return throwError(err);
  }

  private getUrl(apiPrefix: string, apiName: string): string {
    const basicUrl = `${this.apiSource}/${apiPrefix}`;
    return apiName ? `${basicUrl}/${apiName}` : basicUrl;
  }
}

export enum ApiPrefix {
  AUTH = 'auth',
  USER = 'user-profiles',
}
