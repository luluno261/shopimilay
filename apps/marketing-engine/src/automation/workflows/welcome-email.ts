import { Injectable } from '@nestjs/common';
import { AutomationService } from '../automation.service';

@Injectable()
export class WelcomeEmailWorkflow {
  constructor(private readonly automationService: AutomationService) {}

  /**
   * Déclenche la séquence d'email de bienvenue
   */
  async handleNewUser(userId: string, userData: any): Promise<void> {
    console.log(`Envoi email de bienvenue pour utilisateur ${userId}`);

    // Email de bienvenue immédiat
    await this.automationService.sendEmail(
      userId,
      'welcome',
      {
        userName: userData.name || userData.email,
        ...userData,
        delay: 0,
      }
    );

    // Email de découverte après 3 jours
    setTimeout(async () => {
      await this.automationService.sendEmail(
        userId,
        'welcome-followup',
        {
          userName: userData.name || userData.email,
          ...userData,
          delay: 3 * 24 * 60 * 60 * 1000, // 3 jours
        }
      );
    }, 3 * 24 * 60 * 60 * 1000);
  }
}

