import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { AutomationService } from './automation.service';
import { AutomationFlows } from './automation.flows';

@Controller('automation')
export class AutomationController {
  constructor(
    private readonly automationService: AutomationService,
    private readonly automationFlows: AutomationFlows,
  ) {}

  @Get()
  async getAutomations() {
    // TODO: Récupérer depuis la base de données
    return { automations: [] };
  }

  @Post()
  async createAutomation(@Body() body: any) {
    // TODO: Créer une automation
    return { message: 'Automation créée' };
  }

  @Post(':id/trigger')
  async triggerAutomation(@Param('id') id: string, @Body() body: any) {
    await this.automationService.triggerAutomation(id, body.user_id, body.context);
    return { message: 'Automation déclenchée' };
  }
}

