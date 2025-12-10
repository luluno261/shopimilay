import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { FacebookAdsService } from './facebook-ads.service';
import { GoogleAdsService } from './google-ads.service';
import { AudienceSyncService } from './audience-sync.service';

@Controller('ads')
export class AdsController {
  constructor(
    private readonly facebookAdsService: FacebookAdsService,
    private readonly googleAdsService: GoogleAdsService,
    private readonly audienceSyncService: AudienceSyncService,
  ) {}

  @Get('accounts')
  async getAccounts() {
    // TODO: Récupérer depuis la base de données
    return { accounts: [] };
  }

  @Post('accounts/facebook/connect')
  async connectFacebook(@Body() body: any) {
    const accountId = await this.facebookAdsService.connectAccount(
      body.merchant_id,
      body.access_token,
    );
    return { account_id: accountId };
  }

  @Post('accounts/google/connect')
  async connectGoogle(@Body() body: any) {
    const accountId = await this.googleAdsService.connectAccount(
      body.merchant_id,
      body.access_token,
    );
    return { account_id: accountId };
  }

  @Post('accounts/:id/sync')
  async syncAudience(@Param('id') accountId: string, @Body() body: any) {
    await this.audienceSyncService.syncToAllPlatforms(
      body.merchant_id,
      body.audience_id,
      body.user_ids,
    );
    return { message: 'Synchronisation démarrée' };
  }
}

