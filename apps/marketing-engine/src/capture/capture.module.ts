import { Module } from '@nestjs/common';
import { PopupBuilderService } from './popup-builder.service';

@Module({
  providers: [PopupBuilderService],
  exports: [PopupBuilderService],
})
export class CaptureModule {}

