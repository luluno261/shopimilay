# Bibliothèques Node.js partagées

Bibliothèques communes utilisées par les services Node.js.

## Installation

```bash
npm install
npm run build
```

## Utilisation

```typescript
import { validateEmail, generateID } from '@omnisphere/shared-libraries';

const isValid = validateEmail('user@example.com');
const id = generateID();
```

## Packages

- `validation` - Fonctions de validation
- `utils` - Fonctions utilitaires

