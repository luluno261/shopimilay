import Joi from 'joi';

/**
 * Valide une adresse email
 */
export function validateEmail(email: string): boolean {
  const schema = Joi.string().email().required();
  const { error } = schema.validate(email);
  return !error;
}

/**
 * Valide un ID
 */
export function validateID(id: string): boolean {
  const schema = Joi.string().uuid().required();
  const { error } = schema.validate(id);
  return !error;
}

/**
 * Génère un ID unique
 */
export function generateID(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

