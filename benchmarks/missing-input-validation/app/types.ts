/**
 * Type definitions for the inventory API.
 */

export interface AdjustRequest {
  quantity: number;
  reason?: string;
}

export interface ErrorResponse {
  error: string;
}

export interface HealthResponse {
  status: string;
}

export interface ListResponse {
  products: ProductResponse[];
  count: number;
}

export interface ProductResponse {
  sku: string;
  name: string;
  quantity: number;
  price: number;
}
