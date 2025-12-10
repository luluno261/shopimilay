import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { EventsModule } from './events/events.module';
import { AutomationModule } from './automation/automation.module';
import { SegmentationModule } from './segmentation/segmentation.module';
import { CaptureModule } from './capture/capture.module';
import { AdsModule } from './ads/ads.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    EventsModule,
    AutomationModule,
    SegmentationModule,
    CaptureModule,
    AdsModule,
  ],
})
export class AppModule {}

