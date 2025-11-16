export interface ShortenPayload {
  url: string;
  custom?: string;
  expires_at?: string;
}

export interface ShortenResponse {
  short_url: string;
  code: string;
}

export interface AnalyticsResponse {
  code: string;
  original_url: string;
  total_clicks: number;
  unique_ips: number;
  clicks_by_country: Record<string, number>;
}
