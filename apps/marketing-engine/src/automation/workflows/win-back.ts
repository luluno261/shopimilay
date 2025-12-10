import { Injectable } from '@nestjs/common';
import { AutomationService } from '../automation.service';

@Injectable()
export class WinBackWorkflow {
  constructor(private readonly automationService: AutomationService) {}

  /**
   * Déclenche la séquence win-back pour un client inactif
   */
  async handleInactiveCustomer(userId: string, customerData: any): Promise<void> {
    const daysInactive = customerData.daysInactive || 30;
    
    console.log(`Traitement win-back pour utilisateur ${userId} (inactif depuis ${daysInactive} jours)`);

    // Email de rappel avec offre spéciale
    await this.automationService.sendEmail(
      userId,
      'win-back',
      {
        userName: customerData.name || customerData.email,
        daysInactive,
        discountCode: this.generateWinBackDiscount(),
        lastPurchaseDate: customerData.lastPurchaseDate,
        ...customerData,
      }
    );
  }

  /**
   * Génère un code de réduction pour win-back
   */
  private generateWinBackDiscount(): string {
    // Générer un code avec un pourcentage de réduction attractif
    return 'WELCOMEBACK' + Math.random().toString(36).substring(2, 6).toUpperCase();
  }

  /**
   * Vérifie si un client est inactif (pas d'achat depuis X jours)
   */
  async checkCustomerInactivity(userId: string, thresholdDays: number = 30): Promise<boolean> {
    // TODO: Implémenter la vérification depuis la base de données
    // Vérifier la date du dernier achat
    return false;
  }
}

