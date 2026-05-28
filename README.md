# Reservation de salles — MVP Go

Application Go composée d'un backend REST et d'une CLI TUI pour réserver des salles sur des créneaux fixes d'une heure.

## Stack

- API REST : `gin`
- PostgreSQL : `pgx/v5`
- Auth : JWT
- TUI : `tview`

## Comptes de démonstration

- `alice@example.com` / `password`
- `bob@example.com` / `password`

## Variables d'environnement serveur

- `DATABASE_URL` : URL de connexion PostgreSQL
- `JWT_SECRET` : secret de signature des tokens
- `PORT` : port HTTP du serveur (`8080` par défaut)
- Le serveur charge automatiquement un fichier `.env` à la racine si présent

## Lancer l'application en local

### Prérequis

- Go installé
- Docker installé et lancé localement

### 1. Démarrer PostgreSQL avec Docker

Le projet inclut maintenant un `docker-compose.yml` prêt à l'emploi avec :

- base : `reservation_salles`
- utilisateur : `postgres`
- mot de passe : `postgres`

Lancer PostgreSQL :

```bash
make db-up
```

Vérifier les logs si besoin :

```bash
make db-logs
```

Arrêter PostgreSQL :

```bash
make db-down
```

Si tu préfères sans Makefile, utilise :

```bash
docker compose up -d postgres
```

Si ton Docker expose encore l'ancien binaire, utilise plutôt :

```bash
docker-compose up -d postgres
```

### 2. Démarrer le serveur API

```bash
go mod tidy
make run-server
```

Un fichier `.env.example` est fourni. Pour initialiser ton environnement local :

```bash
cp .env.example .env
```

Au premier lancement, le serveur crée automatiquement les tables et injecte les données de démonstration.

### 3. Démarrer la CLI TUI

Dans un second terminal :

```bash
make run-cli
```

### 4. Se connecter dans l'application

Comptes de démo :

- `alice@example.com` / `password`
- `bob@example.com` / `password`

Par défaut, la CLI pointe vers `http://localhost:8080`.

### 5. Parcours de test rapide

1. Ouvrir `Login`
2. Se connecter avec `alice@example.com`
3. Aller dans `Salles`
4. Choisir une salle et une date
5. Réserver un créneau libre
6. Vérifier la réservation dans `Mes réservations`
7. Annuler la réservation pour tester le flux complet

## Fichiers locaux CLI

- `config.json` : `~/.config/reservation-salles/config.json`
- `state.json` : `~/.config/reservation-salles/state.json`

## Message prêt à envoyer pour les tests

```text
Salut ! Tu peux tester l'app comme ça :

1. Dans le projet :
   - `make db-up`
2. Puis :
   - `go mod tidy`
   - `cp .env.example .env`
   - `make run-server`
3. Dans un autre terminal :
   - `make run-cli`
4. Connecte-toi avec :
   - `alice@example.com` / `password`

Ensuite tu peux tester :
- Login
- liste des salles
- réservation d'un créneau
- affichage de "Mes réservations"
- annulation d'une réservation

Si besoin pour diagnostiquer PostgreSQL :
- `make db-logs`
```

## Endpoints principaux

- `POST /auth/login`
- `GET /auth/me`
- `GET /rooms`
- `GET /rooms/:roomID/availability?date=YYYY-MM-DD`
- `GET /reservations/me`
- `POST /reservations`
- `DELETE /reservations/:id`
- `GET/PUT /me/config`
- `GET/PUT /me/state`

## Build et release

```bash
make build
make release-snapshot
make release-gh
```

`make build` produit :

- serveur Linux
- CLI Linux
- CLI macOS
- CLI Windows

## Hébergement gratuit

Le serveur est prêt pour un déploiement gratuit sur Render ou Alwaysdata : il lit `PORT`, `DATABASE_URL` et `JWT_SECRET` depuis l'environnement.
