# Utilisation de Scale-Test CLI

## Présentation

`scale-test` est un client en ligne de commande pour l’API Scale-Test. Il permet de créer, récupérer et supprimer des runs de test de charge.

## Authentification

Fournissez votre clé API avec :

- `--api-key <API_KEY>`
- ou `SCALE_TEST_API_KEY`

Le CLI utilise d’abord `--api-key`. Si ce flag n’est pas fourni, il lit `SCALE_TEST_API_KEY`.

## URL de l’API

L’URL de base de l’API peut être remplacée avec :

- `--base-url <URL>`
- ou `SCALE_TEST_BASE_URL`

Valeur par défaut : `https://scale-test.com/api/v1`

## Commandes

### `run create`

Créer un nouveau run.

Exemples :

```bash
./scale-test --api-key <API_KEY> run create --scenario-id 123
```

```bash
./scale-test run create --file scenario.yaml
```

```bash
./scale-test run create --scenario-id 123 --wait --poll-interval 2s
```

### `run get`

Récupérer les détails d’un run.

```bash
./scale-test run get <RUN_UUID>
```

### `run delete`

Supprimer un run.

```bash
./scale-test run delete <RUN_UUID>
```

## Exemple YAML

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

## Notes

- Le JSON est écrit sur `stdout`
- Les messages de progression sont écrits sur `stderr` avec `--wait`
