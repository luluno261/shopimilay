import { Injectable } from '@nestjs/common';
import { AutomationService } from './automation.service';

@Injectable()
export class AutomationFlows {
  constructor(private readonly automationService: AutomationService) {}

  /**
   * Séquence d'abandon de panier
   */
  async handleCartAbandonment(userId: string, cartData: any): Promise<void> {
    console.log(`Traitement abandon de panier pour utilisateur ${userId}`);
    
    // 1. Email de rappel après 1 heure
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'cart-abandoned', {
        cartUrl: cartData.cartUrl || '#',
        items: cartData.items || [],
      });
    }, 60 * 60 * 1000); // 1 heure

    // 2. Email avec réduction après 24h
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'cart-abandoned', {
        cartUrl: cartData.cartUrl || '#',
        discount: '10%',
      });
    }, 24 * 60 * 60 * 1000); // 24 heures

    // 3. Dernière chance après 72h
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'cart-abandoned', {
        cartUrl: cartData.cartUrl || '#',
        discount: '15%',
        lastChance: true,
      });
    }, 72 * 60 * 60 * 1000); // 72 heures
  }

  /**
   * Séquence de bienvenue pour nouveaux utilisateurs
   */
  async handleWelcomeSequence(userId: string): Promise<void> {
    console.log(`Déclenchement séquence de bienvenue pour utilisateur ${userId}`);
    
    // 1. Email de bienvenue immédiat
    await this.automationService.sendEmail(userId, 'welcome', {
      name: 'Cher client',
    });

    // 2. Email avec guide après 1 jour
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'welcome', {
        type: 'guide',
      });
    }, 24 * 60 * 60 * 1000);

    // 3. Email avec offre spéciale après 3 jours
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'welcome', {
        type: 'offer',
        discount: '10%',
      });
    }, 3 * 24 * 60 * 60 * 1000);
  }

  /**
   * Séquence post-achat
   */
  async handlePostPurchaseSequence(userId: string, orderId: string): Promise<void> {
    console.log(`Déclenchement séquence post-achat pour commande ${orderId}`);
    
    // 1. Email de confirmation immédiat
    await this.automationService.sendEmail(userId, 'order-confirmation', {
      orderId,
    });

    // 2. Email de suivi après livraison (simulé après 2 jours)
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'order-confirmation', {
        orderId,
        type: 'shipping',
      });
    }, 2 * 24 * 60 * 60 * 1000);

    // 3. Demande d'avis après 7 jours
    setTimeout(async () => {
      await this.automationService.sendEmail(userId, 'order-confirmation', {
        orderId,
        type: 'review',
      });
    }, 7 * 24 * 60 * 60 * 1000);
  }

  /**
   * Séquence de réactivation pour utilisateurs inactifs
   */
  async handleReactivationSequence(userId: string, lastActivityDate: Date): Promise<void> {
    console.log(`Déclenchement séquence de réactivation pour utilisateur ${userId}`);
    
    const daysSinceActivity = Math.floor(
      (Date.now() - lastActivityDate.getTime()) / (24 * 60 * 60 * 1000)
    );

    if (daysSinceActivity >= 30 && daysSinceActivity < 60) {
      // Email de rappel après 30 jours
      await this.automationService.sendEmail(userId, 'welcome', {
        type: 'reactivation',
        message: 'Nous vous avons manqué !',
      });
    } else if (daysSinceActivity >= 60 && daysSinceActivity < 90) {
      // Email avec offre spéciale après 60 jours
      await this.automationService.sendEmail(userId, 'welcome', {
        type: 'reactivation',
        discount: '20%',
      });
    } else if (daysSinceActivity >= 90) {
      // Email de dernière chance après 90 jours
      await this.automationService.sendEmail(userId, 'welcome', {
        type: 'reactivation',
        discount: '25%',
        lastChance: true,
      });
    }
  }
}

