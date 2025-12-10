import { Injectable } from '@nestjs/common';
import { FacebookAdsService } from './facebook-ads.service';
import { GoogleAdsService } from './google-ads.service';

@Injectable()
export class AudienceSyncService {
  constructor(
    private readonly facebookAdsService: FacebookAdsService,
    private readonly googleAdsService: GoogleAdsService,
  ) {}

  /**
   * Synchronise une audience vers toutes les plateformes connectées
   */
  async syncToAllPlatforms(merchantId: string, audienceId: string, userIds: string[]): Promise<void> {
    // TODO: Récupérer les comptes connectés pour ce marchand
    // TODO: Synchroniser vers chaque plateforme
    
    // Exemple:
    // const accounts = await this.getConnectedAccounts(merchantId);
    // for (const account of accounts) {
    //   if (account.platform === 'facebook') {
    //     await this.facebookAdsService.syncAudience(account.id, audienceId, userIds);
    //   } else if (account.platform === 'google') {
    //     await this.googleAdsService.syncAudience(account.id, audienceId, userIds);
    //   }
    // }
    
    console.log(`Synchronisation audience ${audienceId} pour merchant ${merchantId}`);
  }
}

