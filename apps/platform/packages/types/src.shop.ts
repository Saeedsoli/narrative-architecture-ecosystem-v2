// packages/types/src.shop.ts

export interface Product {
  id: string;
  title: string;
  description: string;
  price: number;
  imageUrl: string;
}

export interface CartItem extends Product {
  quantity: number;
}