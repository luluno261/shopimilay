import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Configuration CORS
  app.enableCors();
  
  // Préfixe global pour les routes API
  app.setGlobalPrefix('api/v1');
  
  const port = process.env.PORT || 8083;
  await app.listen(port);
  
  console.log(`Marketing Engine démarré sur le port ${port}`);
}

bootstrap();

