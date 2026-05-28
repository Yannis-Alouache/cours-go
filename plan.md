# Plan MVP — réservation de salles

## État actuel
- Dépôt quasi vide : seul `enonce.md` est présent.
- Aucun code Go, aucune structure CLI/API, aucune configuration de build ou de release.
- Le plan part donc d'un démarrage from scratch avec un périmètre volontairement réduit.

## Approche retenue
- Monorepo Go simple avec deux binaires : `cmd/server` et `cmd/cli`.
- Réservations limitées à des créneaux fixes d'1h sur une journée donnée pour éviter un calendrier complexe.
- Salles et comptes de démonstration préchargés en base pour éviter un CRUD admin.
- Configuration locale et état local stockés en JSON dans le dossier utilisateur, synchronisables via l'API.

## Todos
- Poser le squelette minimal du projet.
- Mettre en place PostgreSQL et les données de base.
- Créer l'API REST d'authentification.
- Créer l'API REST métier de réservation.
- Ajouter les endpoints d'import/export de config et d'état.
- Construire le socle CLI : fichiers locaux + client HTTP.
- Brancher la TUI minimale pour la démo.
- Finaliser erreurs, builds, déploiement et release.

## 1. Périmètre MVP

- **Ce qu'on fait**
  - Authentification simple par email/mot de passe avec token JWT.
  - Liste des salles préchargées en base.
  - Consultation des créneaux libres d'une salle pour une date donnée.
  - Réservation d'un créneau.
  - Consultation de mes réservations.
  - Annulation d'une réservation.
  - TUI avec écran d'accueil, login, paramètres, écrans métier.
  - Stockage local `config.json` et `state.json`.
  - Import/export de la configuration et de l'état entre CLI et serveur.
  - Cross-compilation, release GitHub et hébergement gratuit du serveur.

- **Ce qu'on ne fait pas**
  - CRUD admin complet des salles ou des utilisateurs.
  - Inscription publique, reset mot de passe, rôles avancés.
  - Notifications, récurrence, participants, validation managériale.
  - Calendrier avancé semaine/mois, drag-and-drop, gestion de fuseaux.
  - Mode offline, synchronisation complexe, cache riche.

## 2. Stack minimale recommandée

- **TUI : `tview`**
  - Formulaires, listes, tables et navigation multi-pages très rapides à monter.
  - Suffisant pour une interface colorée et agréable sans la complexité d'une architecture Bubble Tea.

- **API REST : `gin`**
  - Recommandé dans l'énoncé.
  - Routing, binding JSON et middleware simples pour aller vite.

- **PostgreSQL : `pgx/v5` + `pgxpool`**
  - Requêtes SQL directes, sans ORM.
  - Très bien pour un schéma court et explicable à l'oral.

- **Auth : `golang-jwt/jwt/v5`**
  - Plus simple qu'une session serveur pour une CLI.
  - Le token peut être stocké dans `state.json`.

- **Config/état local : `os.UserConfigDir()` + `encoding/json`**
  - Pas de dépendance inutile.
  - Cross-plateforme propre pour Windows, Linux et macOS.

- **Build/release : `Makefile` + `go build` croisé + `gh release create`**
  - Plus simple qu'une usine à gaz de CI/CD.
  - Suffisant pour satisfaire cross-compilation et GitHub Releases.

## 3. Modèle de données minimal

- **`users`**
  - `id`
  - `email` unique
  - `password_hash`
  - `created_at`

- **`rooms`**
  - `id`
  - `name` unique

- **`reservations`**
  - `id`
  - `user_id`
  - `room_id`
  - `day` (date)
  - `start_time`
  - `end_time`
  - `created_at`
  - contrainte unique sur `(room_id, day, start_time, end_time)`

- **`client_sync`**
  - `user_id` unique
  - `config_json` JSONB
  - `state_json` JSONB
  - `updated_at`

**Note simplificatrice :** les créneaux disponibles ne sont pas stockés en table. Ils sont calculés à partir d'une liste fixe d'heures ouvrées, par exemple `08:00-18:00` avec pas de `1h`.

## 4. Endpoints REST minimaux

- **Auth**
  - `POST /auth/login`
  - `GET /auth/me`

- **Salles**
  - `GET /rooms`

- **Créneaux**
  - `GET /rooms/:roomID/availability?date=YYYY-MM-DD`

- **Réservations**
  - `GET /reservations/me`
  - `POST /reservations`
  - `DELETE /reservations/:id`

- **Import/export config/state**
  - `GET /me/config`
  - `PUT /me/config`
  - `GET /me/state`
  - `PUT /me/state`

## 5. Écrans TUI minimaux

- **Accueil**
  - Menu principal.
  - Si non connecté : `Login`, `Paramètres`, `Quitter`.
  - Si connecté : `Salles`, `Mes réservations`, `Paramètres`, `Quitter`.

- **Login**
  - Formulaire `email` + `mot de passe`.
  - Message d'erreur lisible si échec.
  - Sauvegarde du token dans `state.json` si succès.

- **Paramètres**
  - URL de l'API.
  - Choix de thème/couleurs simple.
  - Actions `Importer config`, `Exporter config`, `Importer état`, `Exporter état`.

- **Écrans métier**
  - `Salles` : liste des salles.
  - `Disponibilités` : choix d'une date, affichage des créneaux libres de la salle sélectionnée, action `Réserver`.
  - `Mes réservations` : liste des réservations de l'utilisateur, action `Annuler`.

## 6. Plan d'implémentation en 6 à 10 étapes maximum

1. **Poser le squelette minimal**
   - **But** : créer une base simple avec deux binaires et des dossiers cohérents.
   - **Composants/fichiers concernés** : `go.mod`, `cmd/server/main.go`, `cmd/cli/main.go`, `internal/config/`, `Makefile`.
   - **Dépendances éventuelles** : aucune.
   - **Résultat concret attendu** : le projet compile avec un serveur vide et une CLI vide.

2. **Mettre en place PostgreSQL et les données de démonstration**
   - **But** : créer le schéma minimal et précharger salles + comptes.
   - **Composants/fichiers concernés** : `db/migrations/001_init.sql`, `db/seed.sql`, `internal/db/postgres.go`.
   - **Dépendances éventuelles** : étape 1.
   - **Résultat concret attendu** : base prête avec tables minimales, salles de test et au moins un utilisateur de démo.

3. **Livrer l'authentification API**
   - **But** : permettre le login et protéger les routes privées.
   - **Composants/fichiers concernés** : `internal/server/router.go`, `internal/server/auth.go`, `internal/server/middleware.go`.
   - **Dépendances éventuelles** : étape 2.
   - **Résultat concret attendu** : `POST /auth/login` renvoie un token, `GET /auth/me` fonctionne avec Bearer token.

4. **Livrer l'API métier de réservation**
   - **But** : couvrir le coeur fonctionnel demandé par l'énoncé.
   - **Composants/fichiers concernés** : `internal/server/rooms.go`, `internal/server/reservations.go`, logique de calcul des créneaux.
   - **Dépendances éventuelles** : étape 3.
   - **Résultat concret attendu** : lister les salles, voir les créneaux libres, réserver, lister mes réservations, annuler.

5. **Ajouter l'import/export de configuration et d'état**
   - **But** : couvrir le critère JSON local <-> serveur.
   - **Composants/fichiers concernés** : `internal/server/client_sync.go`, accès table `client_sync`.
   - **Dépendances éventuelles** : étape 3.
   - **Résultat concret attendu** : les endpoints `GET/PUT` config/state stockent et restituent du JSONB par utilisateur.

6. **Construire le socle CLI hors interface**
   - **But** : gérer les fichiers locaux et les appels HTTP avant de brancher la TUI.
   - **Composants/fichiers concernés** : `internal/local/config.go`, `internal/local/state.go`, `internal/client/api.go`, `cmd/cli/main.go`.
   - **Dépendances éventuelles** : étapes 3 à 5.
   - **Résultat concret attendu** : la CLI sait lire/écrire `config.json` et `state.json`, se connecter et appeler l'API.

7. **Brancher la TUI minimale de démo**
   - **But** : assembler les écrans demandés dans un parcours fluide.
   - **Composants/fichiers concernés** : `internal/tui/app.go`, `home.go`, `login.go`, `settings.go`, `rooms.go`, `reservations.go`.
   - **Dépendances éventuelles** : étape 6.
   - **Résultat concret attendu** : démo complète dans le terminal du login à l'annulation d'une réservation.

8. **Durcir les erreurs et préparer la livraison**
   - **But** : fiabiliser la démo et satisfaire les points hors fonctionnalités.
   - **Composants/fichiers concernés** : handlers API, vues TUI, `Makefile`, scripts de build, `README.md`.
   - **Dépendances éventuelles** : étape 7.
   - **Résultat concret attendu** : messages d'erreur propres, build Linux pour le serveur, builds Windows/Linux/macOS pour la CLI, serveur gratuit déployé, release GitHub prête.

## 7. Tableau court : critère de l'énoncé -> solution prévue

| Critère de l'énoncé | Solution prévue |
| --- | --- |
| Écran d'accueil TUI | Menu principal avec navigation selon état connecté/non connecté |
| Écran d'authentification | Formulaire login email/mot de passe |
| Écran de paramètres | Thème + URL API + import/export config/state |
| Écrans métier TUI | Salles, disponibilités, mes réservations |
| Stockage local config JSON | `config.json` dans `os.UserConfigDir()` |
| Stockage local état JSON | `state.json` dans `os.UserConfigDir()` |
| Import/export config et état | Endpoints `GET/PUT /me/config` et `GET/PUT /me/state` |
| API REST auth | `POST /auth/login`, `GET /auth/me` |
| Endpoints de données | `/rooms`, `/availability`, `/reservations` |
| PostgreSQL | 4 tables : `users`, `rooms`, `reservations`, `client_sync` |
| Gestion des erreurs | JSON d'erreur côté API + messages visibles côté TUI |
| Cross-compilation serveur Linux | cible `make build-server-linux` |
| Cross-compilation CLI multi-OS | cibles `make build-cli-linux`, `darwin`, `windows` |
| GitHub Releases | `gh release create` avec binaires `dist/` |
| Hébergement gratuit | déploiement simple sur Render ou Alwaysdata |

## 8. Ordre de priorité

1. Backend : base + auth + liste des salles.
2. CLI : login + lecture des salles.
3. Réservation : voir les créneaux libres puis réserver.
4. Mes réservations : afficher puis annuler.
5. Paramètres + import/export config/state.
6. Finition démo : erreurs, build croisé, hébergement, release.

## 9. Éléments exclus

- CRUD admin des salles.
- Création de compte depuis la CLI.
- Gestion de plusieurs sites, bâtiments ou équipements.
- Réservations récurrentes.
- Invitations/participants.
- Notifications mail ou push.
- Tableau de bord analytics.
- Docker, microservices, CI/CD complète, optimisation prématurée.
