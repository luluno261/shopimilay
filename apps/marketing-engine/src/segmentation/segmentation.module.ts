import { Module } from '@nestjs/common';
import { SegmentationService } from './segmentation.service';
import { AudienceBuilder } from './audience-builder';

@Module({
  providers: [SegmentationService, AudienceBuilder],
  exports: [SegmentationService, AudienceBuilder],
})
export class SegmentationModule {}

