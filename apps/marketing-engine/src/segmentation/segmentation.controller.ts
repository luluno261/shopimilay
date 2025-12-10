import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { SegmentationService } from './segmentation.service';
import { AudienceBuilder } from './audience-builder';

@Controller('segmentation')
export class SegmentationController {
  constructor(
    private readonly segmentationService: SegmentationService,
    private readonly audienceBuilder: AudienceBuilder,
  ) {}

  @Get()
  async getAudiences() {
    // TODO: Récupérer depuis la base de données
    return { audiences: [] };
  }

  @Post()
  async createAudience(@Body() body: any) {
    const audienceId = await this.segmentationService.createAudience(
      body.merchant_id,
      body.name,
      body.rules,
    );
    return { audience_id: audienceId };
  }

  @Post(':id/evaluate')
  async evaluateAudience(@Param('id') id: string) {
    const userIds = await this.segmentationService.evaluateAudience(id);
    return { user_ids: userIds, count: userIds.length };
  }
}

