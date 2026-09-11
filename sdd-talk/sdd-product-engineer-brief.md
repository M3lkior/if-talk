# Brief de session — Talk "Le développeur Product Engineer / développement piloté par les spécifications"

> **Document de handoff.** Contexte, plan validé, sources et travail restant, pour reprendre la session dans Claude Code.
> Généré le 10/09/2026. Recherches web effectuées à cette date (les frameworks SDD évoluent vite : revérifier les versions et les chiffres avant de figer les slides).

---

## 0. Contexte et objectif

### La demande
Produire un **slide deck** pour animer une intervention de **45 minutes** devant un **club de développeurs**, en introduction d'une demi-journée intitulée :

> **Le développeur "Product Engineer" / développement piloté par les spécifications**
> Spec-driven development, DDD, structuration des spécifications et évolution du rôle du développeur avec les outils d'IA générative.

### Les 5 questions auxquelles le deck doit répondre
1. Introduire le **nouveau métier de product engineer** : comment il a été créé, et pourquoi.
2. **Que faut-il présenter avant** d'introduire le SDD (les prérequis de compréhension) ?
3. Comment bien présenter **l'intérêt** du SDD ?
4. Présenter un **overview des principaux frameworks** du marché : mode plan vs OpenSpec / SpecKit / Superpowers / BMAD.
5. Montrer en quoi le SDD permet de **découpler et d'accélérer réellement** le développement humain et de feature.

### Le fil rouge (thèse unique du talk)

> **L'IA a rendu le code abondant. Le goulot d'étranglement s'est déplacé en amont (l'intention) et en aval (la vérification). Le SDD est la discipline qui réoutille ces deux extrémités — et c'est ce déplacement qui crée le "Product Engineer".**

Tout slide qui ne sert pas cette phrase est candidat à la coupe. Piège n°1 à éviter : le catalogue d'outils.

### Le pont conceptuel à ne jamais perdre
C'est la charnière qui fait tenir la demi-journée en **un** sujet et pas deux :

> Le product engineer a besoin d'un artefact pour porter l'intention à travers la frontière humain/machine. Cet artefact, c'est la spec.
> **Sans SDD, un product engineer n'est qu'un dev surchargé à qui on a retiré son PM.**

Ce pont apparaît deux fois dans le deck : slide S5 (pose) et slide S30 (résolution). **À écrire en premier, avant tous les autres slides.** Si le lien ne se formule pas en une phrase, réduire la partie product engineer à 4 min et la déplacer en fin de talk.

---

## 1. Recherche : généalogie du rôle "Product Engineer"

Le point contre-intuitif et la force de l'accroche : **le product engineer n'a pas été créé par l'IA.** Il existait avant, pour résoudre un problème d'organisation. L'IA l'a rendu inévitable. Raconté dans cet ordre, on passe de "encore un buzzword" à "un mouvement de fond avec 6 ans d'historique".

### Vague 0 — Ce que le rôle défait (1930-2015)
Une ligne suffit, pour montrer que la spécialisation est une construction historique récente : le product management a ses racines dans les années 1930 quand Procter & Gamble a créé le rôle de "brand man" ; Hewlett-Packard l'a adapté à la tech dans les années 1960 ; Microsoft a suivi avec ses "program managers" dans les années 1980 ; au début des années 2000, l'essor du SaaS a consacré le PM comme pièce maîtresse, en même temps que montaient les équipes pluridisciplinaires.

**Message :** la séparation PM / dev / designer est une réponse à une contrainte, pas une loi de la nature.

### Vague 1 — La naissance réelle du rôle (2018-2024), sans IA
Le rôle est né dans des scale-ups produit, pour une raison purement organisationnelle : **le coût des handoffs dépasse leur bénéfice quand l'équipe est petite et que la boucle de feedback doit être courte.**

Sources primaires :

- **PostHog** — la référence documentaire, avec un *Product Engineer Handbook* public. Leur définition : un développeur directement responsable de l'amélioration du produit, ayant un avis sur la roadmap, capable de prendre seul des décisions de design, faisant du support utilisateur directement, moins attaché aux fonctionnalités qu'il a ajoutées par le passé, conscient de la place de son travail dans la stratégie, et motivé avant tout par l'impact et les résultats. Modèle organisationnel associé : une petite équipe efficace doit avoir un seul leader, pouvoir livrer et décider de façon autonome avec un minimum de dépendances, et posséder sa propre mission, ses objectifs long terme, ses métriques clés et ses clients cibles. **Chiffre à mettre en slide : 187 personnes organisées en 47 équipes de 4,2 personnes en moyenne.**
- **incident.io** — résume les product engineers comme des développeurs qui se soucient davantage des résultats et de l'impact que de l'implémentation exacte ou des outils utilisés.
- **Jean-Michel Lemieux (ex-VP Engineering, Shopify)** — des ingénieurs qui ont soif d'utiliser la technologie pour court-circuiter les problèmes humains et utilisateurs.
- **Product Engineer Manifesto** (productengineer.org / github.com/anttiviljami/product-engineer-manifesto, par Viljami Kuosmanen, ~2024) — c'est la responsabilité des builders de chercher d'abord à comprendre le problème avant de plonger dans les solutions, et de s'occuper des domaines design, technique et business en prenant une part active dans chacun.
- **En France, le rôle existait dès 2024** — Le Talent Club a publié une interview d'une product engineer chez Flowie, arrivée par la voie PM. Son avertissement, à reprendre tel quel : beaucoup de gens qui veulent faire du produit veulent faire de la stratégie et prendre des décisions, alors qu'au quotidien c'est surtout de la delivery ; il ne faut pas devenir product engineer pour faire de la stratégie, ce ne sera pas la mission principale.

**Pourquoi le rôle a été créé, en une phrase de slide :** pour supprimer la traduction. Chaque passage de main entre celui qui comprend le problème et celui qui écrit le code coûte du temps, de l'information et de l'intention.

### Vague 2 — L'IA fait passer le rôle de niche culturelle à norme (2025-2026)

**Le mécanisme, pivot de l'accroche :** la spécialisation était une réponse à la complexité ; quand l'IA absorbe une partie de la complexité d'exécution, la division du travail perd sa justification. L'IA élimine la distance entre les disciplines — un fondateur construit un prototype fonctionnel en quelques heures, un développeur dessine des interfaces, un designer génère du code prêt pour la production.

Formulations qui marchent en présentation :

- **Le "one-person band"** — le Chief AI Officer de Pendo décrit le glissement : traditionnellement un ingénieur livre, un product manager parle aux utilisateurs pour voir si c'était viable, et la boucle continue ; en se rapprochant du product engineer, on obtient une seule personne qui livre, itère, collecte le feedback et itère de nouveau, au lieu de répartir ça entre engineering, produit et design. Et parce qu'ils passent moins de temps à écrire du code grâce à Codex, Claude Code ou l'agent de leur choix, ce rôle se concentre sur le travail à plus forte valeur : la résolution créative de problèmes et les échanges avec les utilisateurs.
- **L'argument économique (le plus solide)** — la ressource rare n'est plus l'exécution, c'est le jugement : savoir quoi construire, et savoir si ce qu'on a construit est bon. La rémunération monte pour ceux qui construisent concrètement et baisse pour ceux qui ne font que coordonner. Le rôle se situe à l'intersection de quatre disciplines : ingénierie logicielle/IA, product management, UX/UI et jugement sur la donnée.
- **L'argument du goulot (le meilleur pont vers le SDD)** — le temps de développement est rarement le vrai frein à la vélocité produit ; les goulots sont en amont, dans des exigences floues, des décisions produit faibles et des frictions entre l'ingénierie et le reste du métier. Une IA qui rend le code moins cher ne répare rien de ça, elle l'expose, parce que la partie du process qui absorbait le mou avance désormais plus vite que tout ce qui l'entoure. Conséquence : l'ingénieur le plus utile n'est pas celui qui produit le plus de code, c'est celui qui opère bien sur tout le chemin allant d'un problème métier à un changement livré et fonctionnel.
- **Lecture du marché français** — SFEIR positionne le product engineer comme le rôle central de l'ère 10x, qui ne se définit plus par une spécialité technique unique mais par sa capacité à coordonner toute la chaîne de production logicielle augmentée par l'IA ; le développeur passe de créateur de code à coordinateur de production, les compétences techniques restant essentielles pour le discernement et la validation alors que le code devient une commodité.
- **Donnée de marché** (à manier avec précaution : petits volumes, Angleterre uniquement) — sur les 6 mois précédant le 5 janvier 2026, 72 offres permanentes portaient "Product Engineer" dans le titre, contre 20 sur la même période de 2025, soit 0,14 % des offres permanentes contre 0,041 %, et un bond de 129 places au classement des intitulés. Directionnel, pas probant : le dire.

### Les nuances à porter (pour ne pas faire un talk de hype)

- Tous les ingénieurs ne peuvent pas devenir de bons product engineers : le rôle exige profondeur technique, intuition produit, communication, conscience business et forte autogestion. Le recrutement devient plus difficile car il faut évaluer les candidats au-delà de la capacité à coder, et les organisations doivent soit sélectionner beaucoup plus durement, soit investir lourdement pour faire évoluer leurs ingénieurs — les deux voies coûtent nettement plus qu'une structure d'ingénierie classique.
- L'un des pièges opérationnels les plus dangereux : déléguer l'autorité produit trop tôt, sans supervision suffisante ni maturité organisationnelle.
- **Ce n'est pas la mort du PM.** Le manifeste s'adresse directement aux designers dans une lettre ouverte pour désamorcer l'idée d'un empiètement, et le rôle reste compatible avec un PM qui garde du recul sur la priorisation globale.
- **Piège franco-français** : sur les job boards FR, "product engineer" désigne aussi un ingénieur produit R&D classique — bureaux d'études, automobile, aéronautique, électronique, prototypage. Quelqu'un dans la salle qui a cherché le terme sur un site d'emploi n'a pas trouvé ce métier. Bonne accroche, et vraie information.

---

## 2. Le plan de slides validé — 45 minutes, 33 slides

| Partie | Durée | Slides |
|---|---|---|
| 0. Accroche : un nouveau métier, et son problème | 7 min | S1–S5 |
| 1. Les prérequis avant le SDD | 8 min | S6–S10 |
| 2. Le SDD et son intérêt | 10 min | S11–S17 |
| 3. Panorama des frameworks | 10 min | S18–S24 |
| 4. Découplage et accélération | 7 min | S25–S29 |
| 5. Retour au product engineer + sortie | 3 min | S30–S33 |

---

### Partie 0 — Accroche : un nouveau métier, et son problème (7 min)

**S1 — Titre + toi.** Qui tu es, pourquoi tu parles de ça, la promesse de la demi-journée.

**S2 — "Product engineer" : le métier que le marché a inventé deux fois.** Deux colonnes. À gauche 2018-2024, né sans IA, pour supprimer le coût des handoffs (PostHog 187 personnes / 47 équipes / 4,2 en moyenne ; incident.io ; Lemieux ; le manifeste de 2024). À droite 2025-2026, généralisé par l'IA. Punchline : *ce rôle n'a pas été créé par l'IA, il a été révélé par elle.*

**S3 — Pourquoi maintenant : le mécanisme.** La spécialisation était une réponse à la complexité. L'IA absorbe une partie de la complexité d'exécution, donc la division du travail perd sa raison d'être : l'IA élimine la distance entre les disciplines. Illustrer avec le "one-person band" de Pendo. Puis l'argument économique : la ressource rare n'est plus l'exécution mais le jugement.

**S4 — Le paradoxe de productivité.** Le personnage est posé, montrer son problème :
- Télémétrie sur plus de 10 000 développeurs : +21 % de tâches complétées, +98 % de PR mergées à l'échelle individuelle, métriques de delivery organisationnelles plates.
- Analyse LinearB de 8,1 millions de PR : le code généré par IA attend 4,6× plus longtemps sa première revue.
- METR : oscillation de 37 points sur 12 mois, d'un ralentissement de 19 % à une accélération de 18 % — la courbe en J.
- Conclusion DORA : l'IA ne répare pas une équipe, elle amplifie ce qui est déjà là.

**S5 — Le diagnostic, et la promesse du talk.** Le temps de développement est rarement le vrai frein ; les goulots sont en amont dans les exigences floues et les décisions produit faibles, et une IA qui rend le code moins cher ne répare rien de ça, elle l'expose. Donc : *le product engineer a besoin d'un artefact pour porter l'intention jusqu'à la machine. Cet artefact, c'est la spec. Sans SDD, un product engineer n'est qu'un dev surchargé à qui on a retiré son PM.* Annoncer le plan en 4 temps.

---

### Partie 1 — Les prérequis avant d'introduire le SDD (8 min)

**S6 — L'échelle de délégation.** Autocomplétion → chat → agent → agent autonome. Le secteur passe de pratiques assistées comme le vibe coding, où l'assistant accélère un développeur isolé, vers l'Agentic Software Engineering où des agents autonomes reçoivent des tâches au niveau de l'objectif. Faire lever la main : *qui est où aujourd'hui ?* → calibre la salle pour la suite.

**S7 — Comment marche vraiment un agent (le slide le plus utile du talk).** Pas de mémoire persistante, seulement une fenêtre de contexte. Tout le SDD découle de là : puisque les LLM n'ont pas de mémoire persistante, il faut créer des artefacts durables et versionnés — `constitution.md`, `spec.md`, `plan.md` — pour externaliser l'état du projet. Introduire le terme **context engineering**. Si ce slide passe, le reste coule.

**S8 — Les modes de défaillance du vibe coding.** Concret, ils l'ont tous vécu : l'agent produit mille lignes de code qui ont l'air correctes et qu'on n'a pas demandées ; on voulait un petit fix, on récupère un refactor, une logique "améliorée", trois nouveaux fichiers et des tests qui passent parce qu'ils ne testent rien. Trois causes racines : intention qui ne vit que dans l'historique de chat, dérive, output non vérifiable faute de critères d'acceptation.

**S9 — Rien de tout ça n'est nouveau (slide anti-objection).** Frise : user story + critères d'acceptation, Gherkin/BDD, TDD, contract-first OpenAPI, ADR, DDD. Message : *on a toujours écrit des specs ; ce qui change, c'est qu'elles doivent être lisibles par une machine et versionnées.* Glisser ici le glossaire en encart plutôt qu'en slide séparé (spec / plan / task / constitution / gate / subagent / drift) → gain d'une minute.

**S10 — Rappel DDD, et pourquoi il compte PLUS qu'avant.** Langage ubiquitaire, bounded contexts, context map, agrégats, invariants. L'argument : un nommage précis et cohérent affûte les prompts et empêche les agents de confondre des concepts distincts ; les agents ont accès à tout le dépôt et cherchent les noms qu'on leur donne, ils peuvent naïvement essayer d'unifier des entités similaires, donc il faut expliciter qu'elles sont séparées pour une raison. Trois usages directs : langage ubiquitaire = désambiguïsation de l'agent, bounded contexts = périmètre de contexte et de parallélisation, invariants = critères vérifiables. Mise en garde qui compte : la valeur du modèle de domaine réside dans la compréhension partagée qu'il construit dans l'équipe, pas dans l'artefact — une spec générée par IA que personne ne lit ne résout rien.

---

### Partie 2 — Le SDD et son intérêt (10 min)

**S11 — Définition.** Le SDD traite les spécifications comme des contrats exécutables desquels les agents dérivent le code, en empêchant la dérive architecturale par une application automatisée plutôt que par de la documentation passive. La formule : *the spec is the prompt.*

**S12 — Les 3 niveaux d'ambition (Böckeler).** Le cadre analytique le plus cité du domaine, à ne pas sauter :

| Niveau | Principe | Statut |
|---|---|---|
| **Spec-first** | La spec guide l'implémentation puis est abandonnée ; le code reste l'artefact maintenu | Point d'entrée pragmatique, là où sont la plupart des équipes qui démarrent |
| **Spec-anchored** | La spec persiste comme contrat vivant, versionnée et mise à jour | La cible réaliste 2026 |
| **Spec-as-source** | La spec EST la source ; le code est généré, jetable, régénérable | Vision long terme, encore largement expérimentale |

Dire explicitement à quel niveau on propose de jouer → évite 80 % des malentendus dans la salle.

**S13 — La boucle canonique.** intention → clarification → spec (quoi/pourquoi + critères) → plan (comment) → tasks atomiques → implémentation → vérification → réconciliation spec↔code. Gates humains en rouge. **Ce diagramme revient 3 fois dans le deck comme repère de progression.**

**S14 — Anatomie d'une bonne spec.** Une bonne spec fixe les résultats attendus, les limites de périmètre, les contraintes, les décisions antérieures, la découpe en tâches et les critères de vérification — et l'agent comble tout ce que la spec laisse ouvert. Cette dernière proposition est la punchline : *tout ce que tu n'écris pas, l'agent l'invente.* Ajouter un contre-exemple (spec floue) et un exemple réel du repo.

**S15 — Pourquoi ça marche mécaniquement.** (1) externalise l'état que le LLM n'a pas ; (2) déplace la revue en amont — 200 lignes de markdown au lieu de 2000 lignes de diff ; (3) rend l'output vérifiable ; (4) rend le travail reprenable entre sessions et transférable entre agents ; (5) permet le parallélisme. Appui : le SDD attrape les violations architecturales et la dérive de contrat d'API que les tests unitaires ne peuvent structurellement pas détecter, et il scale sur des agents parallèles en séparant le rôle qui implémente de celui qui vérifie.

**S16 — SDD vs TDD vs BDD vs cycle en V.** Tableau compact. Traiter frontalement la critique Beck/Fowler : Fowler soutient qu'une spécification utile dépend de l'apprentissage acquis pendant le développement, et que la clé de l'usage plein de l'IA est d'accélérer les boucles de feedback. La nuance : Beck s'oppose à l'écriture de la spécification complète avant l'implémentation ; le niveau 2 traite la spec comme un document vivant tout au long de l'implémentation ; la critique vise surtout le niveau 3 et tout workflow qui gèle les hypothèses.

**S17 — Le coût, honnêtement (slide de crédibilité).** Trois coûts : cérémonie mal dimensionnée, tokens, dérive spec↔code. L'anecdote parfaite : Böckeler a lancé Kiro sur une petite correction de bug et a obtenu quatre user stories et seize critères d'acceptation ; le SDD passe mal à l'échelle vers le bas — pour un null check d'une ligne, on n'a pas à toucher à `/constitution`. Sur la dérive : aucun de ces outils ne réconcilie automatiquement, il faut le déclencher et relire la sortie, et c'est cette charge que la plupart des équipes esquivent jusqu'à ce que leurs specs aient six mois de retard.

---

### Partie 3 — Panorama des frameworks (10 min)

**S18 — La carte par couches + le niveau 0.** Deux choses sur un slide. D'abord pourquoi c'est confus : les outils opèrent à des couches différentes — définition des artefacts d'exigence, conversion en graphe de tâches, exécution du code, intégration IDE ; une cartographie communautaire recensait plus de 30 outils début 2026. Ensuite la baseline : **le mode plan** (Claude Code, Cursor) — plan éphémère, aucun artefact persisté, aucun gate, zéro installation. Excellent pour une tâche de 30 min, insuffisant dès qu'il y a plusieurs sessions, plusieurs agents ou une revue par un tiers. **C'est le témoin de comparaison sur les 4 slides suivants.**

**S19 — GitHub Spec Kit.** CLI Python, environ 129 000 étoiles et 38 intégrations d'agents ; quatre commandes portent le workflow : `/speckit.specify` capture le contexte métier et les critères de succès, `/speckit.plan` traduit en décisions d'architecture, `/speckit.tasks` décompose en unités testables. Différenciateurs : une constitution définie une fois pour le projet dont chaque spec hérite, et des templates qui marquent les inconnues en `NEEDS CLARIFICATION` plutôt que de deviner. Faiblesses : pas d'étape de revue de code intégrée, et changer de direction implique de relancer les commandes concernées, chacune régénérant tout son document. Profil : greenfield, gates forts, verbeux.

**S20 — OpenSpec (Fission-AI).** L'anti-cérémonie. Empreinte la plus légère : on écrit des **delta specs**, uniquement ce qui change. Installation en 5 minutes contre 30, pas de Python, sortie d'environ 250 lignes contre 800, conçu pour les bases de code existantes, 3 commandes IA contre 8, mise à jour de n'importe quel artefact à tout moment sans phase gates rigides. Workflow OPSX : `explore`, `propose`, `apply`, `sync`, `archive`. Contrepartie : pas de gates de revue entre phases.

**S21 — Superpowers (Jesse Vincent / Prime Radiant).** Le cas le plus intéressant pédagogiquement : ce n'est pas un dépôt de specs mais **une méthodologie livrée comme skills**. Plugin qui impose un workflow structuré avant qu'une seule ligne de code ne soit écrite : brainstormer d'abord, isoler sa branche, écrire un plan détaillé, exécuter — chaque étape conditionnant la suivante. Le plan découpe le travail en tâches de 2 à 5 minutes avec chemins de fichiers exacts, commandes exactes et code complet ; des subagents implémentent chaque tâche avec une revue en deux étapes. Pratiques imposées : cycles TDD red-green-refactor où les tests doivent échouer avant l'implémentation, méthodologie de debug en quatre phases exigeant l'investigation de la cause racine avant tout correctif, sessions de brainstorming socratique. **L'angle à souligner devant des devs :** c'est une réponse directe au problème de discipline — quelqu'un écrit une bonne skill specify-plan-implement, l'utilise une semaine, puis retourne discrètement au prompting non structuré dès qu'une deadline approche. Anecdote : Jesse Vincent, créateur de RT, contributeur de Perl 5, auteur du client mail K-9, a quasiment cessé de coder lui-même.

**S22 — BMAD-METHOD.** Le poids lourd. Simule une équipe agile complète via des personas nommés : Analyst, Product Manager, Architect, UX Designer, Scrum Master, Dev, et BMad Master comme orchestrateur ; l'agent Scrum Master crée des fichiers de story détaillés portant le contexte architectural, les guidelines d'implémentation et les critères de test pour l'agent Dev. En v6 : intelligence adaptative à l'échelle qui ajuste la profondeur de planification du bugfix au système d'entreprise, 19 agents spécialisés, trois pistes — Quick Flow avec tech spec seule, BMad Method avec PRD + architecture + UX, Enterprise pour la conformité.

**S23 — Le tableau de décision.** Le chiffre qui fait rire la salle : sur un même build de dashboard CRM, la même tâche a pris **12 minutes avec OpenSpec, 90 minutes avec Spec Kit et 5 h 30 avec BMAD**. Puis la grille : pour la plupart des équipes travaillant sur du code existant, OpenSpec offre le meilleur équilibre vitesse/flexibilité ; pour les nouveaux projets avec des rôles clairs, Spec Kit apporte structure et documentation ; pour la complexité d'échelle entreprise, BMAD gère l'orchestration multi-agents. L'honnêteté qui fait la différence : BMAD est excellent, mais aussi coûteux et disproportionné pour la plupart du travail hebdomadaire d'ingénierie. *(Mentionner en une ligne, sans slide : Kiro, Tessl, cc-sdd, Antigravity, et la convergence AGENTS.md / Agent Skills comme lingua franca émergente.)*

**S24 — Ce sur quoi ils sont tous d'accord (slide de sortie de partie).** Quatre primitives universelles : **règles/constitution, spec, plan, tasks + des gates humains**. Malgré des approches différentes, tous ces frameworks s'accordent sur un point : l'humain reste dans la boucle, mais pas pour tout. Conclusion libératrice : *tu peux commencer demain avec trois fichiers markdown et un gate ; le framework est une optimisation, pas un prérequis.*

---

### Partie 4 — Découplage et accélération réelle (7 min)

**S25 — Les 5 découplages.**
1. **Le quoi / le comment** — la spec devient l'interface entre jugement humain et exécution machine.
2. **La conception / l'exécution dans le temps** — tu spécifies maintenant, les agents exécutent pendant que tu fais autre chose. Ton temps de clavier n'est plus le facteur limitant.
3. **La revue / le code** — revoir l'intention (200 lignes lisibles) au lieu du diff (2000 lignes qu'on n'a pas écrites).
4. **Les développeurs entre eux** — specs comme contrats + bounded contexts = agents parallèles sur des tranches indépendantes, sans collision (worktrees, subagents).
5. **L'humain / la session** — spec et plan comme mémoire durable : reprise après crash de session, transfert entre agents, onboarding.

**S26 — Pourquoi le découplage n°3 est LE gain.** Le rapport DORA 2026 montre que l'IA n'élimine pas les goulots, elle déplace souvent le problème à l'étape suivante : si l'équipe augmente sa production de code mais continue à revoir les changements de la même façon, avec un contexte limité et peu de signaux de risque clairs, une partie de la vitesse gagnée en développement se perd en validation. Le SDD attaque ce point précis en donnant au reviewer l'intention explicite contre laquelle juger.

**S27 — Les chiffres d'accélération, avec leurs limites.** Le meilleur point de données à l'échelle d'une équipe : **Mercari**, place de marché japonaise d'environ 22 millions d'utilisateurs mensuels, a publié sa méthodologie interne "Agent Spec-Driven Development" en décembre 2025 ; après six mois, elle rapportait un gain de vitesse de **150 %** sur sa baseline traditionnelle et de **80 %** sur le prompting IA en format libre. Compléments : GitHub rapporte que les équipes utilisant Spec Kit livrent avec environ un ordre de grandeur moins de cycles de régénération from scratch ; AWS documente des cas où des fonctionnalités de 40 heures ont été livrées en moins de 8 heures de temps humain quand elles étaient d'abord rédigées comme specs.

**S28 — Où le SDD n'accélère PAS.** Indispensable pour la crédibilité. DORA 2026 montre des gains de 35 à 40 % sur les tâches simples mais moins de 10 % sur du code legacy complexe. Donc : petits changements, spikes exploratoires, domaine encore flou, legacy sans tests. Message : *le SDD est un investissement dont le retour dépend de la taille du changement.*

**S29 — Quoi mesurer.** Pas les lignes de code. Lead time jusqu'à la prod, temps de revue, taux de rework, change failure rate, et surtout **la part du rework due à une intention mal comprise** — c'est la métrique que le SDD prétend améliorer. À citer : seules 7,3 % des équipes ont un taux de rework sous 2 %, ce qui révèle une taxe cachée sur la productivité que la génération de code par IA peut facilement aggraver.

---

### Partie 5 — Retour au product engineer et sortie (3 min)

**S30 — Boucle fermée.** Reprendre le slide S5. Le product engineer avait besoin d'un artefact pour porter l'intention ; on vient de passer 35 minutes à décrire cet artefact et son outillage. La nouvelle pile de compétences : cadrage de problème, modélisation de domaine, écriture de critères vérifiables, arbitrage, design de la vérification. Moins de vitesse de frappe.

**S31 — Les risques, honnêtement.** Tous les ingénieurs ne peuvent pas devenir de bons product engineers — le rôle exige profondeur technique, intuition produit, communication, conscience business et forte autogestion, et l'erreur la plus dangereuse est de déléguer l'autorité produit trop tôt sans supervision ni maturité organisationnelle. Ajouter : atrophie des compétences chez les juniors, capacité de revue, responsabilité en cas d'incident. Rappeler l'avertissement de la product engineer française : ne devenez pas product engineer pour faire de la stratégie, au quotidien c'est surtout de la delivery.

**S32 — "Lundi matin".** 5 actions par coût croissant :
1. Écrire un `AGENTS.md` / `constitution.md` pour un repo.
2. Utiliser le mode plan systématiquement pendant une semaine.
3. Faire une vraie feature avec OpenSpec.
4. Versionner les specs dans git à côté du code.
5. Instaurer un seul gate humain — validation de spec avant implémentation — et mesurer le temps de revue avant/après.

**S33 — Transition vers la suite de la demi-journée.** Ce qui sera fait en atelier, et les 3 questions laissées ouvertes.

---

## 3. Ressources

### Le rôle de product engineer

| Ressource | Pourquoi |
|---|---|
| **PostHog, Product Engineer Handbook** — posthog.com/product-engineer (pages *what-is-a-product-engineer*, *traits*, *culture*) | La source primaire la plus documentée, avec des chiffres d'organisation réels |
| **Product Engineer Manifesto** — productengineer.org + github.com/anttiviljami/product-engineer-manifesto | La formalisation de 2024, pré-IA. Preuve que le rôle n'a pas été inventé par l'IA |
| **Hones** — hones.fr/ressources/product-engineer-nouveau-role-equipe-tech (Matthieu Sénéchal) | L'angle recrutement français. ⚠️ Le site bloque l'accès automatisé : à lire manuellement et à intégrer |
| **Le Talent Club**, *Product Engineer, c'est quoi ce "nouveau" métier ?* — letalentclub.substack.com | Le témoignage FR de 2024 et le meilleur avertissement anti-fantasme |
| **CIO**, *The rise of the product engineer: How AI is reshaping modern tech teams* | Les risques organisationnels et de recrutement |
| **Turing College**, *The Rise of the AI Product Engineer* | L'argument économique : jugement rare, comp qui bascule |
| **hitechnology.io**, *The Rise of the Product Engineer* (étude 100+ leaders) | L'argument du goulot en amont — le pont vers le SDD |
| **SFEIR**, page concept Product Engineer — sfeir.com/concepts/product-engineer | La lecture du marché français / conseil |
| **itjobswatch.co.uk** (Product Engineer, England) | La seule donnée de marché chiffrée, avec ses limites |

### Incontournables SDD (à lire en premier, dans cet ordre)

| Ressource | Pourquoi |
|---|---|
| Böckeler, *Understanding Spec-Driven Development: Kiro, spec-kit, and Tessl* — martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html | La taxonomie des 3 niveaux. Fondation de la partie 2 |
| Podcast Thoughtworks, *What is spec-driven development* (Böckeler + Laura Tacho, AWS) | Le meilleur cadrage du flou terminologique |
| DORA 2026, *The ROI of AI-assisted Software Development* (Google) | Source de chiffres pour les parties 0 et 4 |
| github/spec-kit | Lire les **templates**, pas juste le README |
| github.com/Fission-AI/OpenSpec — `docs/opsx.md` et `docs/commands.md` | Le contraste de cérémonie avec Spec Kit |
| claude.com/plugins/superpowers + obra/superpowers-marketplace | Lire 2-3 fichiers de skills : meilleure démo de "la méthode comme artefact" |
| github.com/bmad-code-org/BMAD-METHOD (v6) | Survoler les workflows et la notion de tracks |

### Argumentation et nuances

- **arXiv 2609.00252** — *Spec-Driven Development for Agentic Software Engineering: Harnessing Human-Agent Teamwork* (août 2026). Le cadrage académique du paradoxe de productivité. Précaution à citer : le travail est présenté comme un premier pas vers un consensus académique-industriel plutôt qu'une théorie validée.
- **arXiv 2605.01160** — *The Productivity-Reliability Paradox: Specification-Driven Governance*. Contient une section DDD/DDT par tiers de délégation, directement exploitable pour S10.
- **martinfowler.com/articles/structured-prompt-driven** — SPDD, variante utile pour montrer que le débat bouge encore.
- **threedots.tech/post/ddd-and-ai-coding** — *Domain-Driven Design matters more when AI writes your code*. Le meilleur article pour la brique DDD.
- **gitnation** — talk *From Prompt Spaghetti to Bounded Contexts: DDD for Agentic Codebases* (AI Coding Summit 2026).
- **François Zaninotto**, *Waterfall Strikes Back* — la critique à connaître pour tenir les questions.
- **ranthebuilder.cloud**, *I Tested Three Spec-Driven AI Tools* — comparatif hands-on avec numéros de version, le plus honnête trouvé.
- **reenbit.com/bmad-vs-spec-kit-vs-openspec** — la grille de décision et les chiffres de durée.
- **github.com/ianhxu/agentic-engineering-field-study** (fichier `04-spec-driven-development.md`) — revue de littérature déjà faite, avec le cas Mercari.
- **medium.com/@damien.gouron** — *Spec-Driven Development : le guide pratique*, à partager après le talk (source FR).

---

## 4. Travail restant

### À produire soi-même (le plus important, et le plus long)
1. **Une spec réelle de son propre code, en deux versions** : floue vs complète, avec le résultat produit par l'agent dans les deux cas. C'est le slide dont on se souviendra. → alimente S14.
2. **Une capture du même petit bugfix** passé au mode plan, à OpenSpec et à Spec Kit, pour illustrer le problème de dimensionnement de la cérémonie. → alimente S18 et S23.
3. **Un `AGENTS.md` / `constitution.md` d'exemple**, court, lisible en 20 secondes à l'écran. → alimente S32.
4. **Un diagramme unique de la boucle canonique** avec les gates humains, réutilisé 3 fois dans le deck. → S13.

### Décisions ouvertes
- **Format de sortie du deck** : PowerPoint (.pptx), Google Slides, ou markdown type Slidev / Marp. Non tranché.
- **Calibrage de l'audience** : le plan suppose des devs seniors à l'aise avec les agents de codage mais pas avec le SDD. Si le club est plus mixte (juniors, peu d'usage d'agents), transférer 3 minutes de la partie 3 (en faire un tableau unique au lieu de 4 slides) vers les parties 1 et 2.
- **Intégration de l'article Hones** une fois lu manuellement.
- **Revérifier avant de figer** : versions et nombres d'étoiles des frameworks, chiffres DORA/Mercari, données itjobswatch. Le domaine bouge vite.

### Arbitrages retenus
- **Le risque n°1 du plan : deux sujets au lieu d'un.** L'accroche de 7 min sur le product engineer ne tient que si S5 et S30 font vraiment leur travail de charnière. **Écrire ces deux slides en premier.** Si le lien ne se formule pas en une phrase, réduire la partie product engineer à 4 min et la déplacer en fin de talk.
- **Pas de démo live.** Avec ce volume de concepts, une démo qui plante coûte 8 minutes et la crédibilité. Enregistrer un asciinema de 90 secondes, ou garder la démo pour l'atelier de l'après-midi.
- **Les 5 questions à préparer** (elles tomberont) :
  1. *C'est le cycle en V déguisé ?* → S12 + S16.
  2. *Ça coûte combien en tokens ?* → S17.
  3. *Et le legacy ?* → S28 + OpenSpec.
  4. *Qui maintient les specs quand elles divergent du code ?* → S17, et assumer que c'est le point faible non résolu du domaine.
  5. *Donc on supprime les PM ?* → S31, répondre non et citer la lettre ouverte du manifeste aux designers.

---

## 5. Prompt de reprise pour Claude Code

À coller au démarrage de la session :

```
Lis le fichier sdd-product-engineer-brief.md à la racine.

C'est le brief complet d'un talk de 45 min sur le Spec-Driven Development
agentique et le rôle de Product Engineer, destiné à un club de développeurs.
Le plan de slides (33 slides, 6 parties) est validé — ne le restructure pas
sans me demander.

Objectif de cette session : produire le deck.

Avant de commencer, pose-moi les questions nécessaires sur :
- le format de sortie voulu (pptx / Slidev / Marp / autre)
- le niveau de l'audience, pour arbitrer le calibrage décrit en section 4
- ce que je veux que tu rédiges toi-même vs ce que je garde à ma main

Puis attaque par les slides S5 et S30 : ce sont les charnières du talk,
et la section 0 explique pourquoi elles conditionnent tout le reste.
```
