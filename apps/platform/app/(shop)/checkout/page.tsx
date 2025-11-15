// apps/platform/app/(shop)/checkout/page.tsx

'use client';

import { useCartStore } from '@/lib/store/cart-store';
import { CheckoutForm } from '@/components/shop/checkout-form';
import { useMutation } from '@tanstack/react-query';
import { createOrder } from '@/lib/api/shop';

export default function CheckoutPage() {
  const { items, total, clearCart } = useCartStore();

  const mutation = useMutation({
    mutationFn: createOrder,
    onSuccess: (data) => {
      window.location.href = data.paymentUrl;
      clearCart();
    },
    onError: () => {
      alert('خطا در ایجاد سفارش. لطفاً دوباره تلاش کنید.');
    },
  });

  const handleCheckout = (formData) => {
    const orderData = {
      items: items.map(item => ({ productId: item.id, quantity: item.quantity })),
      shippingInfo: formData,
    };
    mutation.mutate(orderData);
  };

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">نهایی کردن خرید</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
        <div>
          <h2 className="text-xl font-semibold mb-4">اطلاعات ارسال</h2>
          <CheckoutForm onSubmit={handleCheckout} isSubmitting={mutation.isLoading} />
        </div>
        <div className="p-6 bg-gray-50 rounded-lg">
          <h2 className="text-xl font-semibold mb-4">خلاصه سفارش</h2>
          {items.map(item => (
            <div key={item.id} className="flex justify-between py-2 border-b">
              <span>{item.title} x {item.quantity}</span>
              <span>{(item.price * item.quantity).toLocaleString('fa-IR')}</span>
            </div>
          ))}
          <div className="flex justify-between font-bold text-lg mt-4">
            <span>جمع کل:</span>
            <span>{total.toLocaleString('fa-IR')} تومان</span>
          </div>
        </div>
      </div>
    </div>
  );
}