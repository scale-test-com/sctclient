# Scale-Test CLI

Client en ligne de commande en Go pour piloter l API Scale-Test.

## Fonctionnalites

- Creation d un run de test de charge
- Recuperation des details d un run
- Suppression d un run
- Authentification par API key via argument ou variable d environnement
- Option d attente active jusqu a l etat final du run

## Prerequis

- Go 1.22+
- Une API key Scale-Test valide

## Installation

### Depuis le code source

```bash
go build -o scale-test .
```

### Depuis une release GitHub

Telechargez l archive correspondant a votre OS depuis la page Releases.

## Configuration

Les parametres globaux peuvent etre passes en flag ou en variables d environnement.

- API key:
  - flag: --api-key
  - env: SCALE_TEST_API_KEY
- URL API:
  - flag: --base-url
  - env: SCALE_TEST_BASE_URL
  - valeur par defaut: https://scale-test.com/api/v1

Priorite: flag > variable d environnement > valeur par defaut.

## Commandes

### Aide

```bash
./scale-test --help
./scale-test run --help
```

### Creer un run avec un scenario existant

```bash
./scale-test --api-key <API_KEY> run create --scenario-id 123
```

### Creer un run a partir d un fichier YAML

```bash
./scale-test run create --file scenario.yaml
```

### Creer un run et attendre la fin

```bash
./scale-test run create --scenario-id 123 --wait --poll-interval 2s
```

### Recuperer un run

```bash
./scale-test run get <RUN_UUID>
```

### Supprimer un run

```bash
./scale-test run delete <RUN_UUID>
```

## Exemple de scenario YAML

```yaml
name: Test API
request_timeout: 5s
req_target_curve:
  - elapsed_time: 0s
    req_per_sec: 10
  - elapsed_time: 30s
    req_per_sec: 50
  - elapsed_time: 60s
    req_per_sec: 10
operations:
  - uri: https://api.example.com/
    method: GET
```

## Output

- Les resultats JSON sont ecrits sur stdout
- Les messages de progression (mode --wait) sont ecrits sur stderr

## Release automatisee sur tag

Le workflow GitHub Actions publie automatiquement une release et les binaires Linux, macOS, Windows quand un tag est pousse.

Exemple:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Le workflow est defini dans .github/workflows/release.yml.
