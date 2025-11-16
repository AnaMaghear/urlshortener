import { Component, computed, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { ShortenPayload } from './models/url.models';
import { LinkStoreService } from './services/link-store.service';
import { environment } from '../environments/environment';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatToolbarModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  private readonly urlPattern = /^(https?:\/\/).+/i;
  private readonly baseShortUrl = environment.apiUrl.replace(/\/$/, '');

  readonly shortenForm = this.fb.group({
    url: ['', [Validators.required, Validators.pattern(this.urlPattern)]],
    custom: ['', [Validators.pattern(/^$|^[a-zA-Z0-9_-]{3,16}$/)]],
    expiresAt: ['']
  });

  readonly analyticsForm = this.fb.group({
    code: ['', Validators.required]
  });

  readonly submitting = this.linkStore.submitting;
  readonly analyticsLoading = this.linkStore.analyticsLoading;
  readonly qrLoading = this.linkStore.qrLoading;

  readonly shortResult = this.linkStore.shortResult;
  readonly analyticsResult = this.linkStore.analyticsResult;
  readonly qrDataUrl = this.linkStore.qrDataUrl;
  readonly displayedShortLink = computed(() => {
    const shortResult = this.shortResult();
    if (shortResult) {
      return {
        url: shortResult.short_url,
        code: shortResult.code
      };
    }

    const analytics = this.analyticsResult();
    if (analytics) {
      return {
        url: `${this.baseShortUrl}/${analytics.code}`,
        code: analytics.code
      };
    }

    return null;
  });

  readonly errorMessage = this.linkStore.errorMessage;
  readonly analyticsError = this.linkStore.analyticsError;
  readonly qrError = this.linkStore.qrError;

  readonly countryStats = this.linkStore.countryStats;

  constructor(
    private readonly fb: FormBuilder,
    private readonly linkStore: LinkStoreService
  ) {
    effect(() => {
      const result = this.shortResult();
      if (result) {
        this.analyticsForm.patchValue({ code: result.code }, { emitEvent: false });
      }
    });
  }

  onShortenSubmit(): void {
    if (this.shortenForm.invalid) {
      this.shortenForm.markAllAsTouched();
      return;
    }

    const { url, custom, expiresAt } = this.shortenForm.value;
    if (!url) {
      return;
    }

    const payload: ShortenPayload = {
      url: url.trim()
    };

    if (custom) {
      payload.custom = custom.trim();
    }

    if (expiresAt) {
      payload.expires_at = new Date(expiresAt).toISOString();
    }

    this.linkStore.shorten(payload);
  }

  onAnalyticsSubmit(): void {
    if (this.analyticsForm.invalid) {
      this.analyticsForm.markAllAsTouched();
      return;
    }
    const code = this.analyticsForm.value.code?.trim();
    if (!code) {
      return;
    }

    this.linkStore.fetchAnalytics(code);
    this.linkStore.fetchQr(code);
  }

  copyShortUrl(): void {
    const shortUrl = this.displayedShortLink()?.url;
    if (!shortUrl || !navigator.clipboard) {
      return;
    }
    navigator.clipboard.writeText(shortUrl).catch(() => {
      this.errorMessage.set('Unable to copy the URL. Please copy it manually.');
    });
  }

  get urlControl() {
    return this.shortenForm.get('url');
  }

  get customControl() {
    return this.shortenForm.get('custom');
  }

  get expiresAtControl() {
    return this.shortenForm.get('expiresAt');
  }

  get analyticsCodeControl() {
    return this.analyticsForm.get('code');
  }

  get expiresSelection(): Date | null {
    const value = this.expiresAtControl?.value;
    if (!value) {
      return null;
    }
    return new Date(value);
  }
}
