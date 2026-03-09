export interface Product {
  sku: string;
  name: string;
  quantity: number;
  price: number;
}

export interface AdjustmentResult {
  sku: string;
  previousQuantity: number;
  adjustment: number;
  newQuantity: number;
}

export class Inventory {
  private products: Map<string, Product> = new Map();

  addProduct(product: Product): void {
    this.products.set(product.sku, { ...product });
  }

  getProduct(sku: string): Product | undefined {
    const product = this.products.get(sku);
    return product ? { ...product } : undefined;
  }

  listProducts(): Product[] {
    return Array.from(this.products.values()).map(p => ({ ...p }));
  }

  adjustQuantity(sku: string, adjustment: number): AdjustmentResult {
    const product = this.products.get(sku);
    if (!product) {
      throw new Error(`Product ${sku} not found`);
    }

    const previousQuantity = product.quantity;
    product.quantity += adjustment;

    return {
      sku: product.sku,
      previousQuantity,
      adjustment,
      newQuantity: product.quantity,
    };
  }

  getTotalValue(): number {
    let total = 0;
    for (const product of this.products.values()) {
      total += product.quantity * product.price;
    }
    return Math.round(total * 100) / 100;
  }
}
