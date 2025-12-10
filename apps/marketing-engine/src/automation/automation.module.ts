import { Module } from '@nestjs/common';
import { AutomationService } from './automation.service';
import { AutomationFlows } from './automation.flows';
import { EmailService } from './email.service';
import { AutomationController } from './automation.controller';

@Module({
  controllers: [AutomationController],
  providers: [AutomationService, AutomationFlows, EmailService],
  exports: [AutomationService, AutomationFlows, EmailService],
})
export class AutomationModule {}

