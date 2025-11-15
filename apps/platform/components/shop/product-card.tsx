// apps/platform/components/shop/product-card.tsx

import Image from 'next/image';
import type { Product } from '@narrative-arch/types';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/packages/ui/src/card';
import { Button } from '@/packages/ui/src/button';

interface ProductCardProps {
  product: Product;
  onAddToCart: (product: Product) => void;
}

export function ProductCard({ product, onAddToCart }: ProductCardProps) {
  return (
    <Card>
      <CardHeader>
        <div className="relative aspect-[4/3] w-full">
          <Image
            src={product.imageUrl}
            alt={product.title}
            fill
            className="object-cover rounded-t-lg"
          />
        </div>
      </CardHeader>
      <CardContent>
        <CardTitle className="mb-2">{product.title}</CardTitle>
        <p className="text-sm text-muted-foreground">{product.description}</p>
      </CardContent>
      <CardFooter>
        <span className="font-bold text-lg">{product.price.toLocaleString('fa-IR')} تومان</span>
        <Button onClick={() => onAddToCart(product)}>
          افزودن به سبد
        </Button>
      </CardFooter>
    </Card>
  );
}