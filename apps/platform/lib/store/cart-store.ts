// apps/platform/lib/store/cart-store.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { CartItem, Product } from '@narrative-arch/types';

interface CartState {
  items: CartItem[];
  total: number;
  addItem: (product: Product) => void;
  removeItem: (productId: string) => void;
  updateQuantity: (productId: string, quantity: number) => void;
  clearCart: () => void;
}

const calculateTotal = (items: CartItem[]) => {
  return items.reduce((acc, item) => acc + item.price * item.quantity, 0);
};

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      total: 0,
      addItem: (product) => {
        const currentItems = get().items;
        const existingItem = currentItems.find((item) => item.id === product.id);

        let newItems;
        if (existingItem) {
          newItems = currentItems.map((item) =>
            item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
          );
        } else {
          newItems = [...currentItems, { ...product, quantity: 1 }];
        }
        
        set({
          items: newItems,
          total: calculateTotal(newItems),
        });
      },
      removeItem: (productId) => {
        const newItems = get().items.filter((item) => item.id !== productId);
        set({
          items: newItems,
          total: calculateTotal(newItems),
        });
      },
      updateQuantity: (productId, quantity) => {
        if (quantity < 1) {
          get().removeItem(productId);
          return;
        }
        const newItems = get().items.map((item) =>
          item.id === productId ? { ...item, quantity } : item
        );
        set({
          items: newItems,
          total: calculateTotal(newItems),
        });
      },
      clearCart: () => {
        set({ items: [], total: 0 });
      },
    }),
    {
      name: 'cart-storage', // نام کلید در localStorage
    }
  )
);