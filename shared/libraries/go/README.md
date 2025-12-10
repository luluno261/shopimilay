# Bibliothèques Go partagées

Bibliothèques communes utilisées par les microservices Go.

## Packages

### validation

Fonctions de validation communes.

```go
import "github.com/omnisphere/shared/libraries/go/validation"

err := validation.Validate(myStruct)
```

### utils

Fonctions utilitaires communes.

```go
import "github.com/omnisphere/shared/libraries/go/utils"

id := utils.GenerateID()
```

## Utilisation

Pour utiliser ces bibliothèques dans un service Go, ajoutez-les comme dépendance locale :

```go
require github.com/omnisphere/shared/libraries/go v0.0.0

replace github.com/omnisphere/shared/libraries/go => ../../shared/libraries/go
```

