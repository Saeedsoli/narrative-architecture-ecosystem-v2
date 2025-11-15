// apps/platform/app/(shop)/cart/page.tsx

'use client';

import Link from 'next/link';
import { useCartStore } from '@/lib/store/cart-store';
import { CartItem } from '@/components/shop/cart-item';
import { Button } from '@/packages/ui/src/button';

export default function CartPage() {
  const { items, removeItem, updateQuantity, total } = useCartStore();

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">سبد خرید شما</h1>
      {items.length > 0 ? (
        <>
          <div className="border rounded-lg">
            {items.map((item) => (
              <CartItem
                key={item.id}
                item={item}
                onRemove={removeItem}
                onUpdateQuantity={updateQuantity}
              />
            ))}
          </div>
          <div className="mt-8 flex justify-between items-center">
            <span className="text-xl font-bold">جمع کل: {total.toLocaleString('fa-IR')} تومان</span>
            <Button asChild size="lg">
              <Link href="/checkout">ادامه و پرداخت</Link>
            </Button>
          </div>
        </>
      ) : (
        <p>سبد خرید شما خالی است.</p>
      )}
    </div>
  );
}