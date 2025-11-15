// apps/platform/components/shop/cart-item.tsx

import Image from 'next/image';
import type { CartItem as CartItemType } from '@/lib/store/cart-store';
import { Input } from '@/packages/ui/src/input';
import { Button } from '@/packages/ui/src/button';

interface CartItemProps {
  item: CartItemType;
  onRemove: (productId: string) => void;
  onUpdateQuantity: (productId: string, quantity: number) => void;
}

export function CartItem({ item, onRemove, onUpdateQuantity }: CartItemProps) {
  return (
    <div className="flex items-center justify-between p-4 border-b">
      <div className="flex items-center gap-4">
        <Image src={item.imageUrl} alt={item.title} width={64} height={64} className="rounded-md object-cover" />
        <div>
          <h4 className="font-semibold">{item.title}</h4>
          <p className="text-sm text-muted-foreground">{item.price.toLocaleString('fa-IR')} تومان</p>
        </div>
      </div>
      <div className="flex items-center gap-4">
        <Input
          type="number"
          min="1"
          value={item.quantity}
          onChange={(e) => onUpdateQuantity(item.id, parseInt(e.target.value))}
          className="w-16 text-center"
        />
        <Button variant="destructive" size="sm" onClick={() => onRemove(item.id)}>
          حذف
        </Button>
      </div>
    </div>
  );
}