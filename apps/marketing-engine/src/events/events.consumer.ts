import { Injectable, OnModuleInit } from '@nestjs/common';
import { Kafka, Consumer, EachMessagePayload } from 'kafkajs';
import { EventsService } from './events.service';

@Injectable()
export class EventsConsumer implements OnModuleInit {
  private consumer: Consumer;
  private kafka: Kafka;

  constructor(private readonly eventsService: EventsService) {
    this.kafka = new Kafka({
      clientId: 'marketing-engine',
      brokers: [process.env.KAFKA_BROKERS || 'localhost:9092'],
    });

    this.consumer = this.kafka.consumer({ groupId: 'marketing-engine-group' });
  }

  async onModuleInit() {
    await this.consumer.connect();
    
    // S'abonner aux topics Kafka
    const topics = [
      'user.events',
      'order.events',
      'product.events',
      'cart.events',
    ];

    for (const topic of topics) {
      await this.consumer.subscribe({ topic, fromBeginning: false });
    }

    // Consommer les messages
    await this.consumer.run({
      eachMessage: async (payload: EachMessagePayload) => {
        await this.handleMessage(payload);
      },
    });

    console.log('EventsConsumer initialisé et connecté à Kafka');
  }

  private async handleMessage(payload: EachMessagePayload) {
    const { topic, partition, message } = payload;
    
    try {
      const event = JSON.parse(message.value?.toString() || '{}');
      
      console.log(`Événement reçu - Topic: ${topic}, Partition: ${partition}`);
      
      // Traiter l'événement selon son type
      await this.eventsService.processEvent(topic, event);
      
    } catch (error) {
      console.error(`Erreur lors du traitement de l'événement:`, error);
      // TODO: Implémenter un système de retry ou dead letter queue
    }
  }

  async onModuleDestroy() {
    await this.consumer.disconnect();
  }
}

