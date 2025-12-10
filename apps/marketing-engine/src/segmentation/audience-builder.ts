import { Injectable } from '@nestjs/common';
import { SegmentationService } from './segmentation.service';

@Injectable()
export class AudienceBuilder {
  constructor(private readonly segmentationService: SegmentationService) {}

  /**
   * Construit une audience avec des règles complexes
   */
  async buildComplexAudience(merchantId: string, rules: any[]): Promise<string[]> {
    // TODO: Implémenter la logique de construction d'audience complexe
    // Support des opérateurs AND, OR, NOT
    return await this.segmentationService.buildAudience(merchantId, rules);
  }

  /**
   * Exemples de règles prédéfinies
   */
  getPredefinedRules() {
    return {
      purchasedProduct: (productId: string) => ({
        type: 'purchased_product',
        product_id: productId,
      }),
      notVisitedSince: (days: number) => ({
        type: 'last_visit',
        days,
        operator: 'greater_than',
      }),
      ltvGreaterThan: (value: number) => ({
        type: 'ltv',
        value,
        operator: 'greater_than',
      }),
      cartAbandoned: () => ({
        type: 'cart_abandoned',
      }),
    };
  }
}

