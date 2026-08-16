# Mimic

<p align="center">
  <strong>Serveur de Mock d'API Intelligent et Moteur de Contrats Haute Performance en Go.</strong>
</p>

<p align="center">
  <a href="README.md">Read in English</a> •
  <a href="#demarrage-rapide">Demarrage Rapide</a> •
  <a href="#fonctionnalites">Fonctionnalites</a> •
  <a href="#sources-supportees">Sources Supportees</a> •
  <a href="#commandes-cli">Commandes CLI</a>
</p>

---

## Presentation

**Mimic** est un outil CLI ultra-rapide (< 5 ms de temps de demarrage) sans dependance externe ecrit en Go, transformant vos contrats d'API en serveurs de mock HTTP realistes, interactifs et persistants en memoire vive pour les equipes frontend, mobile (Flutter, React Native, iOS, Android) et QA.

---

## Fonctionnalites

- **Latence HTTP Ultra-Faible** : Propulse par les goroutines Go, servant plus de 100 000 requetes/seconde sans configuration.
- **Ingestion Universelle de Contrats** : Lit nativement les formats OpenAPI 3.0/3.1, Swagger 2.0, Collections Postman v2.1, fichiers de routes Laravel `routes/api.php`, et fichiers declaratifs `mimic.yaml`.
- **Persistance CRUD en Memoire (Stateful)** : Les requetes `POST`, `PUT`, `PATCH`, et `DELETE` modifient reellement les collections en memoire vive. Un appel `GET /resource/:id` retrouve les elements crees.
- **Moteur Faker Intelligent** : Genere des noms, emails, UUIDs, dates ISO et donnees localisees (numeros de telephone Senegal `+221`, Francs CFA `currency_xof`).
- **Simulation Reseau & Chaos** :
  - Latence artificielle reglable (`--latency=200ms`)
  - Simulation de gigue/instabilite (`--jitter=50ms`)
  - Injection de pannes 500 aleatoires (`--chaos=0.05` pour 5% d'erreurs)
- **Support CORS Integre** : Accessible immediatement depuis un navigateur web ou une application mobile.

---

## Demarrage Rapide

### 1. Generer un Contrat d'Exemple
```bash
mimic init
```
Cree un fichier `mimic.yaml` pret a l'emploi avec CRUD utilisateurs, authentification et paiements mobiles.

### 2. Demarrer le Serveur Mock
```bash
mimic serve -f mimic.yaml --port=8080 --latency=200ms
```

### 3. Tester les Endpoints
```bash
# Lister les utilisateurs
curl http://localhost:8080/api/v1/users

# Creer un utilisateur
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Mamadou Diop", "email": "mamadou@example.com"}'

# Recuperer l'utilisateur cree
curl http://localhost:8080/api/v1/users/1
```

---

## Sources Supportees

| Format | Option / Fichier | Exemple |
| :--- | :--- | :--- |
| **YAML / JSON Declaratif** | `-f mimic.yaml` | Format natif avec mapping CRUD et Faker |
| **OpenAPI 3.0 / 3.1** | `-f openapi.yaml` | Standard international de l'industrie |
| **Swagger 2.0** | `-f swagger.json` | Schemas Swagger standards |
| **Collection Postman v2.1** | `-f collection.json` | Export direct depuis Postman |
| **Routes Laravel** | `-f routes/api.php` | Analyse AST directe du fichier de routes |

---

## Commandes CLI

```text
Utilisation:
  mimic [commande]

Commandes Disponibles:
  init        Genere un fichier template mimic.yaml
  serve       Demarre le serveur mock HTTP a partir d'un contrat d'API
  validate    Verifie la syntaxe et la validite d'un contrat d'API
  version     Affiche la version de Mimic
```

---

## Licence

Distribue sous licence libre [MIT](LICENSE).
