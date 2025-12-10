import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { PopupBuilderService } from './popup-builder.service';

@Controller('capture')
export class CaptureController {
  constructor(private readonly popupBuilderService: PopupBuilderService) {}

  @Get('popups')
  async getPopups() {
    // TODO: Récupérer depuis la base de données avec merchant_id
    return { popups: [] };
  }

  @Post('popups')
  async createPopup(@Body() body: any) {
    const popupId = await this.popupBuilderService.createPopup(body.merchant_id, body);
    return { popup_id: popupId };
  }

  @Get('popups/:id/code')
  async getPopupCode(@Param('id') id: string) {
    const code = this.popupBuilderService.generatePopupCode(id);
    return { code };
  }
}

