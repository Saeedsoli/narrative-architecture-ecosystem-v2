// apps/platform/lib/api/shop.ts

import { apiClient } from './client';
import type { Product } from '@narrative-arch/types';

export const listProducts = async (): Promise<Product[]> => {
  const { data } = await apiClient.get('/products');
  return data;
};

interface OrderItemData {
  productId: string;
  quantity: number;
}
interface ShippingInfo {
  fullName: string;
  phone: string;
  address: string;
  postalCode: string;
}
interface CreateOrderRequest {
  items: OrderItemData[];
  shippingInfo: ShippingInfo;
}

interface CreateOrderResponse {
  orderId: string;
  paymentUrl: string;
}

export const createOrder = async (orderData: CreateOrderRequest): Promise<CreateOrderResponse> => {
  const { data } = await apiClient.post('/orders', orderData);
  return data;
};