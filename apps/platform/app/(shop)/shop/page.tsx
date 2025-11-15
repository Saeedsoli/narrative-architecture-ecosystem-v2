// apps/platform/app/(shop)/shop/page.tsx

'use client';

import { useQuery } from '@tanstack/react-query';
import { listProducts } from '@/lib/api/shop';
import { useCartStore } from '@/lib/store/cart-store';
import { ProductCard } from '@/components/shop/product-card';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

export default function ShopPage() {
  const { data: products, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: listProducts,
  });
  
  const { addItem } = useCartStore();

  if (isLoading) return <LoadingSpinner />;

  return (
    <div className="max-w-7xl mx-auto py-12 px-4">
      <h1 className="text-4xl font-bold mb-8">فروشگاه</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {products?.map((product) => (
          <ProductCard key={product.id} product={product} onAddToCart={() => addItem(product)} />
        ))}
      </div>
    </div>
  );
}