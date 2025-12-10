import { Injectable } from '@nestjs/common';
import { AutomationService } from '../automation.service';

@Injectable()
export class AbandonedCartWorkflow {
  constructor(private readonly automationService: AutomationService) {}

  /**
   * Déclenche la séquence d'abandon de panier
   */
  async handleCartAbandonment(userId: string, cartData: any): Promise<void> {
    console.log(`Traitement abandon de panier pour utilisateur ${userId}`);

    // Email immédiat (0h)
    await this.automationService.sendEmail(
      userId,
      'cart-abandoned',
      {
        ...cartData,
        delay: 0,
      }
    );

    // Email de rappel après 24h
    setTimeout(async () => {
      await this.automationService.sendEmail(
        userId,
        'cart-abandoned-reminder',
        {
          ...cartData,
          delay: 24 * 60 * 60 * 1000, // 24h en millisecondes
        }
      );
    }, 24 * 60 * 60 * 1000);

    // Email final après 72h avec code de réduction
    setTimeout(async () => {
      await this.automationService.sendEmail(
        userId,
        'cart-abandoned-final',
        {
          ...cartData,
          discountCode: this.generateDiscountCode(),
          delay: 72 * 60 * 60 * 1000, // 72h
        }
      );
    }, 72 * 60 * 60 * 1000);
  }

  /**
   * Génère un code de réduction pour inciter à finaliser l'achat
   */
  private generateDiscountCode(): string {
    // Générer un code aléatoire (à implémenter avec génération réelle)
    return 'CART' + Math.random().toString(36).substring(2, 8).toUpperCase();
  }
}

