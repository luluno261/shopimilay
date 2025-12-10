import { Module } from '@nestjs/common';
import { EventsConsumer } from './events.consumer';
import { EventsService } from './events.service';
import { AutomationModule } from '../automation/automation.module';

@Module({
  imports: [AutomationModule],
  providers: [EventsConsumer, EventsService],
  exports: [EventsService],
})
export class EventsModule {}

