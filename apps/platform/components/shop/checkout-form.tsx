// apps/platform/components/shop/checkout-form.tsx

'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Input } from '@/packages/ui/src/input';
import { Button } from '@/packages/ui/src/button';

const checkoutSchema = z.object({
  fullName: z.string().min(3, 'نام کامل الزامی است'),
  phone: z.string().regex(/^09[0-9]{9}$/, 'شماره موبایل نامعتبر است (مثال: 09123456789)'),
  address: z.string().min(10, 'آدرس کامل الزامی است'),
  postalCode: z.string().regex(/^[0-9]{10}$/, 'کد پستی باید ۱۰ رقم باشد'),
});

type CheckoutFormValues = z.infer<typeof checkoutSchema>;

interface CheckoutFormProps {
  onSubmit: (data: CheckoutFormValues) => void;
  isSubmitting: boolean;
}

export function CheckoutForm({ onSubmit, isSubmitting }: CheckoutFormProps) {
  const { register, handleSubmit, formState: { errors } } = useForm<CheckoutFormValues>({
    resolver: zodResolver(checkoutSchema),
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label htmlFor="fullName">نام و نام خانوادگی</label>
        <Input id="fullName" {...register('fullName')} className="mt-2" />
        {errors.fullName && <p className="text-destructive text-sm mt-1">{errors.fullName.message}</p>}
      </div>
      <div>
        <label htmlFor="phone">شماره موبایل</label>
        <Input id="phone" {...register('phone')} className="mt-2" />
        {errors.phone && <p className="text-destructive text-sm mt-1">{errors.phone.message}</p>}
      </div>
      <div>
        <label htmlFor="address">آدرس کامل</label>
        <textarea id="address" {...register('address')} rows={3} className="w-full mt-2 p-2 border rounded-md bg-background" />
        {errors.address && <p className="text-destructive text-sm mt-1">{errors.address.message}</p>}
      </div>
      <div>
        <label htmlFor="postalCode">کد پستی</label>
        <Input id="postalCode" {...register('postalCode')} className="mt-2" />
        {errors.postalCode && <p className="text-destructive text-sm mt-1">{errors.postalCode.message}</p>}
      </div>
      <Button type="submit" disabled={isSubmitting} className="w-full" size="lg">
        {isSubmitting ? 'در حال انتقال به درگاه...' : 'پرداخت و نهایی کردن خرید'}
      </Button>
    </form>
  );
}