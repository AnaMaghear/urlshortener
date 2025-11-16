import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

import { environment } from '../../environments/environment';
import { AnalyticsResponse, ShortenPayload, ShortenResponse } from '../models/url.models';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private readonly apiUrl = environment.apiUrl;

  constructor(private readonly http: HttpClient) {}

  shortenUrl(payload: ShortenPayload): Observable<ShortenResponse> {
    return this.http.post<ShortenResponse>(`${this.apiUrl}/shorten`, payload);
  }

  getAnalytics(code: string): Observable<AnalyticsResponse> {
    const params = new HttpParams().set('code', code);
    return this.http.get<AnalyticsResponse>(`${this.apiUrl}/analytics`, { params });
  }

  getQr(code: string): Observable<Blob> {
    const params = new HttpParams().set('code', code);
    return this.http.get(`${this.apiUrl}/qr`, {
      params,
      responseType: 'blob'
    });
  }
}
