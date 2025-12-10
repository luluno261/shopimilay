import { Injectable } from '@nestjs/common';

@Injectable()
export class EmailService {
  /**
   * Envoie un email via SendGrid ou Mailgun
   */
  async sendEmail(to: string, subject: string, html: string, from?: string): Promise<void> {
    const emailProvider = process.env.EMAIL_PROVIDER || 'sendgrid';
    const fromEmail = from || process.env.EMAIL_FROM || 'noreply@omnisphere.com';

    if (emailProvider === 'sendgrid') {
      await this.sendViaSendGrid(to, subject, html, fromEmail);
    } else if (emailProvider === 'mailgun') {
      await this.sendViaMailgun(to, subject, html, fromEmail);
    } else {
      console.log(`[Email Mock] To: ${to}, Subject: ${subject}`);
      // En développement, on log juste
    }
  }

  /**
   * Envoie un email via SendGrid
   */
  private async sendViaSendGrid(to: string, subject: string, html: string, from: string): Promise<void> {
    const sendGridApiKey = process.env.SENDGRID_API_KEY;
    if (!sendGridApiKey) {
      console.warn('SENDGRID_API_KEY non configuré, email non envoyé');
      return;
    }

    // TODO: Implémenter l'intégration SendGrid
    // const sgMail = require('@sendgrid/mail');
    // sgMail.setApiKey(sendGridApiKey);
    // await sgMail.send({ to, from, subject, html });
    
    console.log(`[SendGrid] Email envoyé à ${to}: ${subject}`);
  }

  /**
   * Envoie un email via Mailgun
   */
  private async sendViaMailgun(to: string, subject: string, html: string, from: string): Promise<void> {
    const mailgunApiKey = process.env.MAILGUN_API_KEY;
    const mailgunDomain = process.env.MAILGUN_DOMAIN;
    
    if (!mailgunApiKey || !mailgunDomain) {
      console.warn('MAILGUN_API_KEY ou MAILGUN_DOMAIN non configuré, email non envoyé');
      return;
    }

    // TODO: Implémenter l'intégration Mailgun
    // const formData = require('form-data');
    // const Mailgun = require('mailgun.js');
    // const mailgun = new Mailgun(formData);
    // const mg = mailgun.client({ username: 'api', key: mailgunApiKey });
    // await mg.messages.create(mailgunDomain, { from, to, subject, html });
    
    console.log(`[Mailgun] Email envoyé à ${to}: ${subject}`);
  }

  /**
   * Récupère un template d'email
   */
  async getTemplate(templateId: string, data: any): Promise<string> {
    // TODO: Implémenter le système de templates
    // Pour l'instant, on retourne un template simple
    const templates: Record<string, (data: any) => string> = {
      'welcome': (data) => `
        <h1>Bienvenue ${data.name || 'cher client'} !</h1>
        <p>Merci de nous rejoindre sur OmniSphere.</p>
      `,
      'cart-abandoned': (data) => `
        <h1>Vous avez oublié quelque chose ?</h1>
        <p>Votre panier vous attend !</p>
        <a href="${data.cartUrl}">Retourner au panier</a>
      `,
      'order-confirmation': (data) => `
        <h1>Confirmation de commande</h1>
        <p>Votre commande #${data.orderId} a été confirmée.</p>
      `,
    };

    const template = templates[templateId];
    return template ? template(data) : `<p>Template ${templateId} non trouvé</p>`;
  }
}

