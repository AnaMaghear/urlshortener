import { Injectable, computed, signal } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';

import { ApiService } from './api.service';
import { AnalyticsResponse, ShortenPayload, ShortenResponse } from '../models/url.models';

@Injectable({
  providedIn: 'root'
})
export class LinkStoreService {
  constructor(private readonly api: ApiService) {}

  readonly submitting = signal(false);
  readonly shortResult = signal<ShortenResponse | null>(null);
  readonly errorMessage = signal<string | null>(null);

  readonly analyticsLoading = signal(false);
  readonly analyticsResult = signal<AnalyticsResponse | null>(null);
  readonly analyticsError = signal<string | null>(null);

  readonly qrLoading = signal(false);
  readonly qrDataUrl = signal<string | null>(null);
  readonly qrError = signal<string | null>(null);

  readonly countryStats = computed(() => {
    const data = this.analyticsResult();
    if (!data) {
      return [];
    }
    return Object.entries(data.clicks_by_country || {}).sort((a, b) => b[1] - a[1]);
  });

  shorten(payload: ShortenPayload): void {
    this.submitting.set(true);
    this.errorMessage.set(null);

    this.api.shortenUrl(payload).subscribe({
      next: (response) => {
        this.shortResult.set(response);
        this.submitting.set(false);
        this.fetchAnalytics(response.code);
        this.fetchQr(response.code);
      },
      error: (error) => {
        this.submitting.set(false);
        this.errorMessage.set(this.buildErrorMessage(error));
      }
    });
  }

  fetchAnalytics(code: string): void {
    this.analyticsLoading.set(true);
    this.analyticsError.set(null);

    this.api.getAnalytics(code).subscribe({
      next: (response) => {
        this.analyticsResult.set(response);
        this.analyticsLoading.set(false);
      },
      error: (error) => {
        this.analyticsResult.set(null);
        this.analyticsLoading.set(false);
        this.analyticsError.set(this.buildErrorMessage(error));
      }
    });
  }

  fetchQr(code: string): void {
    this.qrLoading.set(true);
    this.qrError.set(null);

    this.api.getQr(code).subscribe({
      next: (blob) => {
        const reader = new FileReader();
        reader.onloadend = () => {
          this.qrDataUrl.set(reader.result as string);
          this.qrLoading.set(false);
        };
        reader.onerror = () => {
          this.qrLoading.set(false);
          this.qrError.set('Failed to decode QR code image.');
        };
        reader.readAsDataURL(blob);
      },
      error: (error) => {
        this.qrLoading.set(false);
        this.qrDataUrl.set(null);
        this.qrError.set(this.buildErrorMessage(error));
      }
    });
  }

  private buildErrorMessage(error: unknown): string {
    if (error instanceof HttpErrorResponse) {
      if (typeof error.error === 'string' && error.error.trim().length) {
        return error.error;
      }
      if (error.error?.message) {
        return error.error.message;
      }
      if (error.status === 0) {
        return 'Unable to reach the API. Is the server running?';
      }
      return `Request failed with status ${error.status}.`;
    }
    return 'An unexpected error occurred.';
  }
}
