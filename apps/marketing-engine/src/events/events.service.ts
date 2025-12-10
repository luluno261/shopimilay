import { Injectable } from '@nestjs/common';
import { AutomationFlows } from '../automation/automation.flows';
import { AutomationService } from '../automation/automation.service';

@Injectable()
export class EventsService {
  constructor(
    private readonly automationFlows: AutomationFlows,
    private readonly automationService: AutomationService,
  ) {}

  /**
   * Traite un événement reçu depuis Kafka
   */
  async processEvent(topic: string, event: any): Promise<void> {
    console.log(`Traitement événement depuis ${topic}:`, event);
    
    switch (topic) {
      case 'user.events':
        await this.handleUserEvent(event);
        break;
      case 'order.events':
        await this.handleOrderEvent(event);
        break;
      case 'product.events':
        await this.handleProductEvent(event);
        break;
      case 'cart.events':
        await this.handleCartEvent(event);
        break;
      default:
        console.warn(`Topic non géré: ${topic}`);
    }
  }

  private async handleUserEvent(event: any): Promise<void> {
    console.log('Traitement événement utilisateur:', event);
    
    if (event.type === 'user.created' || event.type === 'user.registered') {
      // Déclencher la séquence de bienvenue
      await this.automationFlows.handleWelcomeSequence(event.user_id);
    }
  }

  private async handleOrderEvent(event: any): Promise<void> {
    console.log('Traitement événement commande:', event);
    
    if (event.type === 'order.created' || event.type === 'order.paid') {
      // Déclencher la séquence post-achat
      await this.automationFlows.handlePostPurchaseSequence(
        event.user_id,
        event.order_id
      );
    }
  }

  private async handleProductEvent(event: any): Promise<void> {
    console.log('Traitement événement produit:', event);
    // Enrichir le profil utilisateur avec les interactions produits
  }

  private async handleCartEvent(event: any): Promise<void> {
    console.log('Traitement événement panier:', event);
    
    if (event.type === 'cart_abandoned') {
      // Déclencher la séquence d'abandon de panier
      await this.automationFlows.handleCartAbandonment(
        event.user_id,
        {
          cartUrl: event.cart_url,
          items: event.items || [],
        }
      );
    }
  }
}

