import { Module } from '@nestjs/common';
import { FacebookAdsService } from './facebook-ads.service';
import { GoogleAdsService } from './google-ads.service';
import { AudienceSyncService } from './audience-sync.service';

@Module({
  providers: [FacebookAdsService, GoogleAdsService, AudienceSyncService],
  exports: [FacebookAdsService, GoogleAdsService, AudienceSyncService],
})
export class AdsModule {}

